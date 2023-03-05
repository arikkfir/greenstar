//go:generate go run github.com/99designs/gqlgen generate ./...
package main

import (
	"context"
	"errors"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	auth "github.com/arikkfir/greenstar/backend/auth/pkg"
	"github.com/arikkfir/greenstar/backend/operations/gql"
	"github.com/arikkfir/greenstar/backend/operations/resolver"
	"github.com/arikkfir/greenstar/backend/util/bootutil"
	"github.com/arikkfir/greenstar/backend/util/ginutil"
	"github.com/arikkfir/greenstar/backend/util/gqlutil"
	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
	"github.com/rueian/rueidis"
	"net/http"
	"strconv"
)

func main() {
	config := Config{}
	bootutil.Boot(&config)

	// Service context
	ctx := context.Background()

	// Create the Redis client
	redisClient, err := rueidis.NewClient(rueidis.ClientOption{
		InitAddress:      []string{config.Redis.Host + ":" + strconv.Itoa(config.Redis.Port)},
		ClientName:       "greenstar.operations",
		BlockingPoolSize: config.Redis.PoolSize,
	})
	if err != nil {
		log.Ctx(ctx).Fatal().Err(err).Msg("Failed to connect to Redis")
	}
	defer redisClient.Close()

	// Connect to NATS
	natsClient, err := nats.Connect(config.NATS.URL, nats.RetryOnFailedConnect(true))
	if err != nil {
		log.Ctx(ctx).Fatal().Err(err).Msg("Failed to connect to NATS")
	}
	defer natsClient.Close()
	jetStreamClient, err := natsClient.JetStream()
	if err != nil {
		log.Ctx(ctx).Fatal().Err(err).Msg("Failed to establish JetStream connection")
	}

	// Initialize GraphQL servers
	graphResolver := &resolver.Resolver{Redis: redisClient}
	graphSchema := gql.NewExecutableSchema(gql.Config{Resolvers: graphResolver})
	graphHandler := handler.NewDefaultServer(graphSchema)
	graphHandler.SetErrorPresenter(gqlutil.ErrorPresenter)
	graphHandler.SetRecoverFunc(gqlutil.PanicRecoverer)
	graphPlaygroundHandler := playground.Handler("GraphQL playground", "/query")

	// Listen for operation commands
	sub, err := jetStreamClient.Subscribe("OPERATIONS", func(m *nats.Msg) {
		log.Ctx(ctx).Info().Interface("msg", m).Msg("Received a JetStream message")
		if err := m.Ack(); err != nil {
			log.Ctx(ctx).Error().Err(err).Interface("msg", m).Msg("Failed to ack message")
		}
	}, nats.ManualAck(), nats.AckExplicit(), nats.MaxAckPending(10))
	defer sub.Unsubscribe()

	// Setup routes
	router := ginutil.NewGin(config.DevMode)
	router.Use(auth.CreateVerifyTokenMiddleware("greenstar.auth", "greenstar.operations"))
	router.GET("/playground", func(c *gin.Context) { graphPlaygroundHandler(c.Writer, c.Request) })
	router.POST("/query", func(c *gin.Context) { graphHandler.ServeHTTP(c.Writer, c.Request) })

	// Setup HTTP server
	httpServer := &http.Server{Addr: ":" + strconv.Itoa(config.HTTP.Port), Handler: router}

	// Start server
	log.Ctx(ctx).Info().Str("addr", httpServer.Addr).Msg("Starting HTTP server")
	if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Err(err).Msg("HTTP server failed")
	}
}
