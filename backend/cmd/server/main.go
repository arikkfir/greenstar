package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/arikkfir/greenstar/backend/internal/server"
	"github.com/arikkfir/greenstar/backend/internal/server/middleware"
	"github.com/arikkfir/greenstar/backend/internal/server/resources/account"
	"github.com/arikkfir/greenstar/backend/internal/server/resources/tenant"
	"github.com/arikkfir/greenstar/backend/internal/server/resources/transaction"
	"github.com/arikkfir/greenstar/backend/internal/util/db"
	flagutil "github.com/arikkfir/greenstar/backend/internal/util/flag"
	"github.com/arikkfir/greenstar/backend/internal/util/lang"
	"github.com/arikkfir/greenstar/backend/internal/util/observability"
	"github.com/arikkfir/greenstar/backend/internal/util/version"
	"github.com/iancoleman/strcase"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

var (
	showHelp                    bool
	disableJSONLogging          bool
	logLevel                    slog.Level
	serverPort                  int
	healthPort                  int
	accessLogSuccessfulRequests bool
	accessLogExcludedHeaders    flagutil.CommaSeparatedValue
	accessLogExcludeRemoteAddr  bool
	corsAllowedOrigins          flagutil.CommaSeparatedValue
	corsAllowMethods            = flagutil.CommaSeparatedValue{"GET", "POST", "PUT", "PATCH", "DELETE"}
	corsAllowHeaders            = flagutil.CommaSeparatedValue{"accept", "authorization", "content-type", middleware.RequestIDHeaderName, middleware.TenantIDHeaderName}
	corsDisableCredentials      bool
	corsExposeHeaders           flagutil.CommaSeparatedValue
	corsMaxAge                  time.Duration
)

func main() {

	// Configure CLI flags
	parseFlags()

	// Configure logger
	observability.ConfigureLogging(disableJSONLogging, logLevel, version.Version)
	slog.SetDefault(slog.New(&ServerSLogHandler{Handler: slog.Default().Handler()}))

	// Execute
	if err := execute(); err != nil {
		slog.Error("Failed starting backend server", "err", err)
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
	fs.IntVar(&serverPort, "server-port", 8080, "`Port` for the REST HTTP server")
	fs.IntVar(&healthPort, "health-port", 9000, "`Port` for the health HTTP server")
	fs.BoolVar(&accessLogSuccessfulRequests, "log-successful-requests", false, "Whether to log successful HTTP requests to the access log")
	fs.Var(&accessLogExcludedHeaders, "excluded-accesslog-headers", "Comma-separated list of HTTP `header names` to exclude from the access log")
	fs.BoolVar(&accessLogExcludeRemoteAddr, "exclude-remote-addr", false, "Whether to exclude the remote address from the access log")
	fs.Var(&corsAllowedOrigins, "allowed-origins", "Comma-separated list of allowed `origins` (or origin patterns) that can call this API, correlating to the CORS header 'Access-Control-Allow-Origin'")
	fs.Var(&corsAllowMethods, "allowed-methods", "Comma-separated list of allowed HTTP `methods` that can be called on this API, correlating to the CORS header 'Access-Control-Allow-Methods'")
	fs.Var(&corsAllowHeaders, "allowed-headers", "Comma-separated list of allowed HTTP `headers` that can be used in calls to this API, correlating to the CORS header 'Access-Control-Allow-Headers'")
	fs.BoolVar(&corsDisableCredentials, "disable-credentials", false, "Whether to disable the CORS header 'Access-Control-Allow-Credentials'")                                              // TODO: verify this description
	fs.Var(&corsExposeHeaders, "exposed-headers", "Comma-separated list of HTTP `headers` that can be exposed in this API, correlating to the CORS header 'Access-Control-Expose-Headers'") // TODO: verify this description
	fs.DurationVar(&corsMaxAge, "max-age", 5*time.Second, "The duration that CORS preflight responses may be cached by the client, corresponding to the CORS header 'Access-Control-Max-Age'")

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

	// REST server
	svr := createServer(pgPool)

	// Create REST routes
	serverMux := http.NewServeMux()
	if err := svr.Register(serverMux); err != nil {
		return fmt.Errorf("failed creating HTTP server routes: %w", err)
	}
	mux := withMiddleware(serverMux)

	// Create HTTP servers
	apiServer := &http.Server{
		Addr: ":" + strconv.Itoa(serverPort),
		BaseContext: func(listener net.Listener) context.Context {
			return observability.WithLogger(ctx, slog.Default())
		},
		Handler: mux,
	}
	healthServer := &http.Server{
		Addr:    ":" + strconv.Itoa(healthPort),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNoContent) }),
	}

	// Start the servers
	stop := make(chan string, 100)
	errs := make(chan error, 100)
	startHTTPServer := func(ctx context.Context, stopChan chan string, errChan chan error, name string, server *http.Server) {
		slog.DebugContext(ctx, "Starting HTTP server", "server", name, "addr", server.Addr)
		if err := server.ListenAndServe(); lang.IgnoreErrorOfType(err, http.ErrServerClosed) != nil {
			errChan <- fmt.Errorf("%s server failed: %w", name, err)
			stopChan <- name
		}
	}
	go startHTTPServer(ctx, stop, errs, "health", healthServer)
	go startHTTPServer(ctx, stop, errs, "api", apiServer)

	// Wait for either: one of the HTTP servers to prematurely exit, or an OS interrupt signal
	select {
	case name := <-stop:
		slog.ErrorContext(ctx, "One of the HTTP servers failed", "server", name)
	case <-ctx.Done():
		slog.ErrorContext(ctx, "Interrupt signal received")
	}

	// Gracefully shutdown all HTTP servers
	shutdownCtx := context.Background()
	err = errors.Join(err, healthServer.Shutdown(shutdownCtx))
	err = errors.Join(err, apiServer.Shutdown(shutdownCtx))

	// Close the errors channel & collect all errors that occurred so far
	close(errs)
	for e := range errs {
		err = errors.Join(err, e)
	}
	return err
}

func createServer(pgPool *pgxpool.Pool) server.Server {
	accountsHandler := &account.HandlerImpl{}
	tenantsHandler := &tenant.HandlerImpl{}
	transactionsHandler := &transaction.HandlerImpl{AccountsHandler: accountsHandler}
	return server.Server{
		Pool:                pgPool,
		AccountsHandler:     accountsHandler,
		TenantsHandler:      tenantsHandler,
		TransactionsHandler: transactionsHandler,
	}
}

func withMiddleware(svrMux *http.ServeMux) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/",
		middleware.CommonHeadersMiddleware(
			middleware.PreventCachingMiddleware(
				middleware.RequestIDMiddleware(
					middleware.TenantIDMiddleware(
						middleware.AccessLogMiddleware(accessLogSuccessfulRequests, accessLogExcludeRemoteAddr, accessLogExcludedHeaders,
							middleware.CORSMiddleware(corsAllowedOrigins, corsAllowMethods, corsAllowHeaders, corsDisableCredentials, corsExposeHeaders, corsMaxAge,
								middleware.TraceMiddleware(
									middleware.TokenMiddleware(
										svrMux,
									),
								),
							),
						),
					),
				),
			),
		),
	)
	return mux
}
