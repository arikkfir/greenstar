package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/arikkfir/greenstar/backend/internal/sample"
	"github.com/arikkfir/greenstar/backend/internal/server/resources/account"
	"github.com/arikkfir/greenstar/backend/internal/server/resources/tenant"
	"github.com/arikkfir/greenstar/backend/internal/server/resources/transaction"
	"github.com/arikkfir/greenstar/backend/internal/util/db"
	"github.com/arikkfir/greenstar/backend/internal/util/lang"
	"github.com/arikkfir/greenstar/backend/internal/util/observability"
	"github.com/arikkfir/greenstar/backend/internal/util/version"
	"github.com/iancoleman/strcase"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

var (
	showHelp           bool
	disableJSONLogging bool
	logLevel           slog.Level
)

func main() {

	// Configure CLI flags
	parseFlags()

	// Configure logger
	observability.ConfigureLogging(disableJSONLogging, logLevel, version.Version)

	// Execute
	if err := execute(); err != nil {
		slog.Error("Failed generating sample ACME tenant data", "err", err)
		os.Exit(1)
	}
}

func parseFlags() {
	fs := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	fs.BoolVar(&showHelp, "h", false, "Show help screen")
	fs.BoolVar(&showHelp, "help", false, "Show help screen")
	fs.BoolVar(&disableJSONLogging, "disable-json-logging", false, "Disable JSON logging in favor of human-readable logging.")
	fs.Func("log-level", "Minimum `log level`, must be one of: trace, debug, info, warn, error", func(s string) error {
		switch s {
		case "trace":
			logLevel = observability.LevelTrace
		case "debug":
			logLevel = slog.LevelDebug
		case "info":
			logLevel = slog.LevelInfo
		case "warn", "warning":
			logLevel = slog.LevelWarn
		case "error":
			logLevel = slog.LevelError
		default:
			return fmt.Errorf("invalid log level: %s", s)
		}
		return nil
	})

	if err := fs.Parse(lang.SliceFrom(1, os.Args)); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		fs.Usage()
		os.Exit(2)
	} else if showHelp {
		fs.Usage()
		fs.PrintDefaults()
		os.Exit(0)
	}

	fs.VisitAll(func(f *flag.Flag) {
		if envValue, envSet := os.LookupEnv(strcase.ToScreamingSnake(f.Name)); envSet {
			if err := f.Value.Set(envValue); err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Failed setting flag '%s' from environment variable '%s': %v\n", f.Name, envValue, err)
				fs.Usage()
				os.Exit(2)
			}
		}
	})
}

func execute() (re error) {

	// Configure OTEL
	shutdown, err := observability.ConfigureOTEL(version.Version)
	if err != nil {
		slog.Error("Failed configuring telemetry", "err", err)
	}
	defer func() {
		if err := shutdown(); err != nil {
			slog.Error("Failed shutting down telemetry", "err", err)
			re = errors.Join(re, err)
		}
	}()

	// Setup context
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// PostgreSQL
	pgPool, err := db.NewPostgreSQLPool(ctx)
	if err != nil {
		return fmt.Errorf("failed creating PostgreSQL pool: %w", err)
	}
	defer pgPool.Close()

	// Wait for database to become available
	if err := db.WaitFor(ctx, pgPool); err != nil {
		return err
	}

	// Generate sample data for the ACME tenant
	accountsHandler := &account.HandlerImpl{}
	tenantsHandler := &tenant.HandlerImpl{}
	transactionsHandler := &transaction.HandlerImpl{AccountsHandler: accountsHandler}
	if err := sample.Generate(ctx, pgPool, tenantsHandler, "acme", "A.C.M.E", transactionsHandler, accountsHandler); err != nil {
		return err
	}

	return nil
}
