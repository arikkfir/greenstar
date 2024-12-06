package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/arikkfir/greenstar/backend/internal/auth"
	"github.com/arikkfir/greenstar/backend/internal/server"
	"github.com/arikkfir/greenstar/backend/internal/server/middleware"
	"github.com/arikkfir/greenstar/backend/internal/server/resources/account"
	"github.com/arikkfir/greenstar/backend/internal/server/resources/tenant"
	"github.com/arikkfir/greenstar/backend/internal/server/resources/transaction"
	"github.com/arikkfir/greenstar/backend/internal/server/util"
	"github.com/arikkfir/greenstar/backend/internal/util/db"
	"github.com/arikkfir/greenstar/backend/internal/util/lang"
	"github.com/arikkfir/greenstar/backend/internal/util/observability"
	"github.com/descope/go-sdk/descope/client"
	"log/slog"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/arikkfir/command"
)

type ServerConfig struct {
	ServerPort                      int           `flag:"true"`
	HealthPort                      int           `flag:"true"`
	HTTPAccessLogSuccessfulRequests bool          `flag:"true"`
	HTTPAccessLogExcludedHeaders    []string      `flag:"true"`
	HTTPAccessLogExcludeRemoteAddr  bool          `flag:"true"`
	HTTPAllowedOrigins              []string      `flag:"true"`
	HTTPAllowMethods                []string      `flag:"true"`
	HTTPAllowHeaders                []string      `flag:"true"`
	HTTPDisableCredentials          bool          `flag:"true"`
	HTTPExposeHeaders               []string      `flag:"true"`
	HTTPMaxAge                      time.Duration `flag:"true"`
}

type Action struct {
	Descope auth.DescopeConfig
	Server  ServerConfig
}

func (e *Action) Run(ctx context.Context) error {
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

	// Create REST server
	mux := http.NewServeMux()
	accountsHandler := &account.HandlerImpl{}
	appServer := server.Server{
		Pool:                pgPool,
		Descope:             descopeClient,
		AccountsHandler:     accountsHandler,
		TenantsHandler:      &tenant.HandlerImpl{Descope: descopeClient},
		TransactionsHandler: &transaction.HandlerImpl{AccountsHandler: accountsHandler},
	}
	if err := appServer.Register(mux); err != nil {
		return err
	}
	appHTTPServer := e.newServer(ctx, e.Server.HTTPAccessLogSuccessfulRequests, descopeClient, mux)

	// All HTTP servers to start
	servers := map[string]*http.Server{
		"health": e.newHealthCheckServer(),
		"api":    appHTTPServer,
	}

	// Start the servers
	stop := make(chan string, 100)
	errs := make(chan error, 100)
	for name, s := range servers {
		go e.startHTTPServer(ctx, stop, errs, name, s)
	}

	// Wait for either one of the HTTP servers to prematurely exit, or an OS interrupt signal
	select {
	case name := <-stop:
		slog.Default().ErrorContext(ctx, "One of the HTTP servers failed", "server", name)
	case <-ctx.Done():
		slog.Default().ErrorContext(ctx, "Interrupt signal received")
	}

	// Gracefully shutdown all HTTP servers
	for _, s := range servers {
		err = errors.Join(err, s.Shutdown(context.Background()))
	}

	// Close the errors channel & collect all errors that occurred so far
	close(errs)
	for e := range errs {
		err = errors.Join(err, e)
	}
	return err
}

func (e *Action) startHTTPServer(ctx context.Context, stopChan chan string, errChan chan error, name string, server *http.Server) {
	slog.Default().DebugContext(ctx, "Starting HTTP server", "server", name)
	if err := server.ListenAndServe(); lang.IgnoreErrorOfType(err, http.ErrServerClosed) != nil {
		errChan <- fmt.Errorf("%s server failed: %w", name, err)
		stopChan <- name
	}
}

func (e *Action) newHealthCheckServer() *http.Server {
	return &http.Server{
		Addr: ":" + strconv.Itoa(e.Server.HealthPort),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		}),
	}
}

func (e *Action) newServer(ctx context.Context, logSuccessfulRequests bool, descopeClient *client.DescopeClient, routes http.Handler) *http.Server {

	// Define global mux
	mux := http.NewServeMux()
	mux.Handle("/",
		middleware.CommonHeadersMiddleware(
			middleware.PreventCachingMiddleware(
				middleware.RequestIDMiddleware(
					middleware.AccessLogMiddleware(logSuccessfulRequests, e.Server.HTTPAccessLogExcludeRemoteAddr, e.Server.HTTPAccessLogExcludedHeaders,
						middleware.CORSMiddleware(e.Server.HTTPAllowedOrigins, e.Server.HTTPAllowMethods, e.Server.HTTPAllowHeaders, e.Server.HTTPDisableCredentials, e.Server.HTTPExposeHeaders, e.Server.HTTPMaxAge,
							middleware.TraceMiddleware(auth.WithSDKMiddleware(descopeClient,
								auth.AuthenticationMiddleware(
									routes,
								),
							)),
						),
					),
				),
			),
		),
	)

	return &http.Server{
		Addr: ":" + strconv.Itoa(e.Server.ServerPort),
		BaseContext: func(listener net.Listener) context.Context {
			return util.ContextWithLogger(ctx, slog.Default())
		},
		Handler: mux,
	}
}

func main() {

	defaultMaxAge, err := time.ParseDuration("5s")
	if err != nil {
		panic(err)
	}

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
			Server: ServerConfig{
				ServerPort:       8080,
				HealthPort:       9000,
				HTTPAllowHeaders: []string{"accept", "authorization", "content-type", "x-request-id"},
				HTTPAllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
				HTTPMaxAge:       defaultMaxAge,
			},
		},
		[]any{
			&observability.LoggingHook{LogLevel: "info"},
			&observability.OTelHook{ServiceName: "greenstar-backend"},
		},
	)

	// Execute the correct command
	os.Exit(int(command.Execute(&dynamicWriter{}, cmd, os.Args[1:], command.EnvVarsArrayToMap(os.Environ()))))
}

type dynamicWriter struct{}

func (w *dynamicWriter) Write(p []byte) (n int, err error) {
	return os.Stderr.Write(p)
}
