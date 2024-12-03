package main

import (
	"context"
	"fmt"
	"github.com/arikkfir/command"
	"github.com/arikkfir/greenstar/backend/internal/auth"
	"github.com/arikkfir/greenstar/backend/internal/migration"
	"github.com/arikkfir/greenstar/backend/internal/server"
	"github.com/arikkfir/greenstar/backend/internal/server/resources/account"
	"github.com/arikkfir/greenstar/backend/internal/server/resources/tenant"
	"github.com/arikkfir/greenstar/backend/internal/server/resources/transaction"
	"github.com/arikkfir/greenstar/backend/internal/server/util"
	"github.com/arikkfir/greenstar/backend/internal/util/db"
	"github.com/arikkfir/greenstar/backend/internal/util/lang"
	"github.com/arikkfir/greenstar/backend/internal/util/observability"
	"go.opentelemetry.io/otel/trace"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Action struct {
	Descope                       auth.DescopeConfig
	CurrencyAPIKey                string `required:"true" env:"CURRENCY_API_KEY"`
	GenerateSampleData            bool   `flag:"true"`
	Namespace                     string `flag:"true" env:"POD_NAMESPACE"`
	ExchangeRatesCronJobName      string `flag:"true" env:"EXCHANGE_RATES_CRONJOB_NAME"`
	HistoricalExchangeRatesPeriod string `flag:"true"`
}

func (e *Action) Run(ctx context.Context) error {
	ctx, span := observability.NamedTrace(ctx, "InitJob:Run", trace.SpanKindInternal)
	defer span.End()

	var err error

	ctx = util.ContextWithLogger(ctx, slog.Default())

	// Descope
	descopeClient, err := e.Descope.NewSDK(ctx)
	if err != nil {
		return fmt.Errorf("failed creating Descope client: %w", err)
	}

	// PostgreSQL
	pgPool, err := db.NewPostgreSQLPool(ctx)
	if err != nil {
		return fmt.Errorf("failed creating PostgreSQL pool: %w", err)
	}
	defer pgPool.Close()

	// Wait for PostgreSQL to be ready
	slog.Default().InfoContext(ctx, "Waiting for Postgres connection pool to become available")
	pgPingInterval := 3 * time.Second
	timer := time.NewTimer(2 * time.Minute)
	defer timer.Stop()
	ticker := time.NewTicker(pgPingInterval)
	defer ticker.Stop()
	pgAvailable := false
	for !pgAvailable {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
			return fmt.Errorf("timed out waiting for Postgres connection pool to become available")
		case <-ticker.C:
			if err := pgPool.Ping(ctx); err != nil {
				slog.Default().InfoContext(ctx, "PostgreSQL not yet available", "err", err)
			} else {
				pgAvailable = true
			}
		}
	}
	timer.Stop()
	ticker.Stop()

	// Migrate database schema & generate sample data
	if err := migration.Migrate(ctx, pgPool); err != nil {
		return err
	}

	exchangeRatesManager := migration.NewExchangeRatesManager(e.CurrencyAPIKey, pgPool)

	// Ensure currencies are populated
	if err := exchangeRatesManager.PopulateCurrencies(ctx); err != nil {
		return fmt.Errorf("failed populating currencies: %w", err)
	}

	// Ensure currencies are populated
	historicalExchangeRatesPeriod := migration.HistoricalExchangeRatesPeriod(e.HistoricalExchangeRatesPeriod)
	if err := exchangeRatesManager.PopulateHistoricalRates(ctx, historicalExchangeRatesPeriod); err != nil {
		return fmt.Errorf("failed populating historical exchange rates: %w", err)
	}

	// Create REST server
	accountsHandler := &account.HandlerImpl{}
	appServer := server.Server{
		Pool:                pgPool,
		Descope:             descopeClient,
		AccountsHandler:     accountsHandler,
		TenantsHandler:      &tenant.HandlerImpl{Descope: descopeClient},
		TransactionsHandler: &transaction.HandlerImpl{AccountsHandler: accountsHandler},
	}

	// Generate sample data
	if e.GenerateSampleData {
		if err := appServer.GenerateSampleTenant(ctx, "acme", "A.C.M.E", e.Descope.DescopeBackendAccessKeyToken); err != nil {
			return err
		}
	}

	// Ensure fill-in any missing rates
	if err := exchangeRatesManager.PopulateMissingExchangeRates(ctx); err != nil {
		return fmt.Errorf("failed populating missing exchange rates: %w", err)
	}

	// Unsuspend the exchange rates cron-job if specified
	if e.ExchangeRatesCronJobName != "" {
		if e.Namespace == "" {
			return fmt.Errorf("--namespace is required when --exchange-rates-cronjob-name is specified")
		}

		// Connect to host Kubernetes cluster
		config, err := rest.InClusterConfig()
		if err != nil {
			return fmt.Errorf("failed creating Kubernetes config: %w", err)
		}
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			return fmt.Errorf("failed creating Kubernetes client: %w", err)
		}
		slog.Default().InfoContext(ctx, "Un-suspending the exchange-rates cron-job")
		err = wait.PollUntilContextTimeout(ctx, time.Second, 5*time.Minute, true, func(ctx context.Context) (bool, error) {
			cronJob, err := clientset.BatchV1().CronJobs(e.Namespace).Get(ctx, e.ExchangeRatesCronJobName, metav1.GetOptions{})
			if err != nil {
				return false, fmt.Errorf("failed getting cron job '%s': %w", e.ExchangeRatesCronJobName, err)
			}

			cronJob.Spec.Suspend = lang.PtrOf(false)
			if _, err = clientset.BatchV1().CronJobs(e.Namespace).Update(ctx, cronJob, metav1.UpdateOptions{}); err != nil {
				return false, fmt.Errorf("failed updating cron job '%s': %w", e.ExchangeRatesCronJobName, err)
			}
			return true, nil
		})
		if err != nil {
			return fmt.Errorf("failed un-suspending the exchange-rates cron-job: %w", err)
		}
	}

	return nil
}

func main() {

	// Create command structure
	cmd := command.MustNew(
		filepath.Base(os.Args[0]),
		"GreenSTAR personal accountant backend",
		`GreenSTAR Personal Accountant is your own personal assistant with everything related to your
personal expenses, income and how to balance them.

This is the backend server for the GreenSTAR application.`,
		&Action{
			Descope: auth.DescopeConfig{
				DescopeLogLevel: "none",
			},
			HistoricalExchangeRatesPeriod: string(migration.Latest),
		},
		[]any{
			&observability.LoggingHook{LogLevel: "info"},
			&observability.OTelHook{ServiceName: "greenstar-init-job"},
		},
	)

	// Execute the correct command
	os.Exit(int(command.Execute(&dynamicWriter{}, cmd, os.Args[1:], command.EnvVarsArrayToMap(os.Environ()))))
}

type dynamicWriter struct{}

func (w *dynamicWriter) Write(p []byte) (n int, err error) {
	return os.Stderr.Write(p)
}
