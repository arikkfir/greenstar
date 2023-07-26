//go:generate go run github.com/99designs/gqlgen generate --verbose
package main

import (
	"context"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/arik-kfir/greenstar/backend/gql"
	"github.com/arik-kfir/greenstar/backend/resolver"
	"github.com/arik-kfir/greenstar/backend/services"
	"github.com/arik-kfir/greenstar/backend/util"
	"github.com/arik-kfir/greenstar/backend/web"
	descope "github.com/descope/go-sdk/descope/client"
	"github.com/descope/go-sdk/descope/logger"
	"github.com/rs/zerolog/log"
	"github.com/secureworks/errors"
	"net/http"
)

func main() {
	ctx := context.Background()
	log.Ctx(ctx).Info().Msg("Starting Greenstar Admin")

	// Create the Redis client
	redisClient, err := util.CreateRedisClient(cfg.Redis.Host, cfg.Redis.Port, "greenstar-admin", cfg.Redis.TLS)
	if err != nil {
		log.Ctx(ctx).Fatal().Err(err).Msg("Failed to connect to Redis")
	}
	defer redisClient.Close()
	log.Ctx(ctx).Info().Msg("Connected to Redis")

	// Create the Neo4j client
	neo4jDriver, err := util.NewNeo4jDriver(cfg.Neo4j.Host, cfg.Neo4j.Port, cfg.Neo4j.TLS)
	if err != nil {
		log.Ctx(ctx).Fatal().Err(err).Msg("Failed to connect to Neo4j")
	}
	defer neo4jDriver.Close(ctx)
	log.Ctx(ctx).Info().Msg("Connected to Neo4j")

	// Descope
	var ll logger.LogLevel
	switch cfg.Auth.DescopeLogLevel {
	case "none":
		ll = logger.LogNone
	case "debug":
		ll = logger.LogDebugLevel
	default:
		ll = logger.LogInfoLevel
	}
	descopeClient, err := descope.NewWithConfig(&descope.Config{
		ProjectID:     cfg.Auth.DescopeProjectID,
		ManagementKey: cfg.Auth.DescopeManagementKey,
		LogLevel:      ll,
		Logger:        &util.DescopeLogger{Logger: log.Ctx(ctx), LogLevel: ll},
	})
	if err != nil {
		log.Ctx(ctx).Fatal().Err(err).Msg("Failed to setup Descope client")
	}
	log.Ctx(ctx).Info().Msg("Initialized Descope")

	// Setup health check
	hc := web.NewHealthCheck(cfg.HTTP.HealthPort)
	go hc.Start(ctx)
	defer hc.Stop(ctx)
	log.Ctx(ctx).Info().Msg("Initialized health check")

	// Admin routes
	service := services.Service{Neo4j: neo4jDriver, Redis: redisClient, Descope: descopeClient}
	graphResolver := &resolver.Resolver{
		TenantsService:      &services.TenantsService{Service: service},
		AccountsService:     &services.AccountsService{Service: service},
		TransactionsService: &services.TransactionsService{Service: service},
		OperationsService:   &services.OperationsService{Service: service},
	}
	graphSchema := gql.NewExecutableSchema(gql.Config{Resolvers: graphResolver})
	graphServer := handler.NewDefaultServer(graphSchema)
	graphServer.SetErrorPresenter(util.GraphErrorPresenter)
	graphServer.SetRecoverFunc(util.GraphPanicRecoverer)
	graphPlayground := playground.Handler("GraphQL Playground", "/query")

	// Graph routes handler
	graphMux := http.NewServeMux()
	graphMux.Handle("/playground", graphPlayground)
	graphMux.Handle("/query", graphServer)
	graphHandler := func(w http.ResponseWriter, r *http.Request) { graphMux.ServeHTTP(w, r) }

	// Setup HTTP server
	server := web.NewServer(
		cfg.HTTP.Port,
		!cfg.HTTP.DisableAccessLog,
		cfg.HTTP.CORS.AllowedOrigins,
		cfg.HTTP.CORS.AllowMethods,
		cfg.HTTP.CORS.AllowHeaders,
		cfg.HTTP.CORS.DisableCredentials,
		cfg.HTTP.CORS.ExposeHeaders,
		cfg.HTTP.CORS.MaxAge,
		redisClient,
		neo4jDriver,
		descopeClient,
		graphHandler,
	)

	// Start server
	log.Ctx(ctx).Info().Str("addr", server.Addr).Msg("Starting HTTP server")
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Err(err).Msg("HTTP server failed")
	}
}
