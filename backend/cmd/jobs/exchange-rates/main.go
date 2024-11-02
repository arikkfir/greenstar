package main

import (
	"context"
	"fmt"
	"github.com/arikkfir/command"
	"github.com/arikkfir/greenstar/backend/internal/migration"
	"github.com/arikkfir/greenstar/backend/internal/server/util"
	"github.com/arikkfir/greenstar/backend/internal/util/db"
	"github.com/arikkfir/greenstar/backend/internal/util/observability"
	"log/slog"
	"os"
	"path/filepath"
)

type Action struct {
	Observability  observability.Config
	CurrencyAPIKey string `required:"true" env:"CURRENCY_API_KEY"`
}

func (e *Action) Run(ctx context.Context) error {
	var err error

	ctx = util.ContextWithLogger(ctx, slog.Default())

	// PostgreSQL
	pgPool, err := db.NewPostgreSQLPool(ctx)
	if err != nil {
		return fmt.Errorf("failed creating PostgreSQL pool: %w", err)
	}
	defer pgPool.Close()

	exchangeRatesManager := migration.NewExchangeRatesManager(e.CurrencyAPIKey, pgPool)

	if err := exchangeRatesManager.PopulateCurrencies(ctx); err != nil {
		return fmt.Errorf("failed populating currencies: %w", err)
	}

	if err := exchangeRatesManager.PopulateMissingExchangeRates(ctx); err != nil {
		return fmt.Errorf("failed populating currencies: %w", err)
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
		&Action{},
		[]any{
			&observability.LoggingHook{LogLevel: "info"},
			&observability.OTelHook{ServiceName: "greenstar-currencies-job"},
		},
	)

	// Execute the correct command
	os.Exit(int(command.Execute(&dynamicWriter{}, cmd, os.Args[1:], command.EnvVarsArrayToMap(os.Environ()))))
}

type dynamicWriter struct{}

func (w *dynamicWriter) Write(p []byte) (n int, err error) {
	return os.Stderr.Write(p)
}
