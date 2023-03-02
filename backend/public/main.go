//go:generate go run github.com/99designs/gqlgen generate ./...
package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	auth "github.com/arikkfir/greenstar/backend/auth/pkg"
	"github.com/arikkfir/greenstar/backend/public/gql"
	"github.com/arikkfir/greenstar/backend/public/resolver"
	"github.com/arikkfir/greenstar/backend/util/bootutil"
	"github.com/arikkfir/greenstar/backend/util/ginutil"
	"github.com/arikkfir/greenstar/backend/util/gqlutil"
	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
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

	// Create the Neo4j client
	neo4jURL := fmt.Sprintf("neo4j://%s:%d", config.Neo4j.Host, config.Neo4j.Port)
	neo4jDriver, err := neo4j.NewDriverWithContext(neo4jURL, neo4j.NoAuth())
	if err != nil {
		log.Ctx(ctx).Fatal().Err(err).Msg("Failed to connect to Neo4j")
	}
	defer neo4jDriver.Close(ctx)

	// Create the Redis client
	redisClient, err := rueidis.NewClient(rueidis.ClientOption{
		InitAddress:      []string{config.Redis.Host + ":" + strconv.Itoa(config.Redis.Port)},
		ClientName:       "greenstar-accounts",
		BlockingPoolSize: config.Redis.PoolSize,
	})
	if err != nil {
		log.Ctx(ctx).Fatal().Err(err).Msg("Failed to connect to Redis")
	}
	defer redisClient.Close()

	// Initialize GraphQL servers
	graphResolver := &resolver.Resolver{Redis: redisClient, Neo4j: neo4jDriver}
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

	// Initialize Gin framework
	router := ginutil.NewGin(config.DevMode)
	router.Use(auth.CreateVerifyTokenMiddleware(
		config.HTTP.ClaimsCookieName,
		config.Auth.Google.ClientSecret,
		"greenstar.auth",
		"greenstar.public",
	))
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
