package main

import (
	"cloud.google.com/go/pubsub"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/99designs/gqlgen/graphql/executor"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	adminGQL "github.com/arik-kfir/greenstar/backend/admin/gql"
	adminResolver "github.com/arik-kfir/greenstar/backend/admin/resolver"
	"github.com/arik-kfir/greenstar/backend/auth"
	operationsGQL "github.com/arik-kfir/greenstar/backend/operations/gql"
	operationsResolver "github.com/arik-kfir/greenstar/backend/operations/resolver"
	publicGQL "github.com/arik-kfir/greenstar/backend/public/gql"
	publicResolver "github.com/arik-kfir/greenstar/backend/public/resolver"
	"github.com/arik-kfir/greenstar/backend/util/ginutil"
	"github.com/arik-kfir/greenstar/backend/util/gqlutil"
	"github.com/arik-kfir/greenstar/backend/util/httputil"
	"github.com/arik-kfir/greenstar/backend/util/neo4jutil"
	"github.com/arik-kfir/greenstar/backend/util/pubsubutil"
	"github.com/arik-kfir/greenstar/backend/util/redisutil"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jessevdk/go-flags"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/rueian/rueidis"
	"github.com/tavsec/gin-healthcheck"
	"github.com/tavsec/gin-healthcheck/config"
	"net/http"
	"os"
	"strconv"
)

func init() {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
}

type Config struct {
	LogLevel    string                       `env:"LOG_LEVEL" value-name:"LEVEL" long:"log-level" description:"Log level" default:"info" enum:"trace,debug,info,warn,error,fatal,panic"`
	DevMode     bool                         `env:"DEV_MODE" long:"dev-mode" description:"Development mode"`
	DevModeID   string                       `env:"DEV_MODE_ID" long:"id" description:"Development mode ID, used for generating developer-specific resources."`
	Redis       redisutil.RedisConfig        `group:"redis" namespace:"redis" env-namespace:"REDIS"`
	Neo4j       neo4jutil.Neo4jConfig        `group:"neo4j" namespace:"neo4j" env-namespace:"NEO4J"`
	GoogleCloud pubsubutil.GoogleCloudConfig `group:"gcp" namespace:"gcp" env-namespace:"GCP"`
	HTTP        httputil.HTTPConfig          `group:"http" namespace:"http" env-namespace:"HTTP"`
	Auth        auth.Config                  `group:"auth" namespace:"auth" env-namespace:"AUTH"`
}

func main() {
	cfg := Config{}

	// Parse configuration
	parser := flags.NewParser(&cfg, flags.HelpFlag|flags.PassDoubleDash)
	if _, err := parser.Parse(); err != nil {
		fmt.Printf("ERROR: %s\n\n", err)
		parser.WriteHelp(os.Stderr)
		os.Exit(1)
	}
	if cfg.DevMode {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		zerolog.DefaultContextLogger = &log.Logger
	}
	if level, err := zerolog.ParseLevel(cfg.LogLevel); err != nil {
		log.Fatal().Err(err).Msg("Failed to parse config")
	} else {
		zerolog.SetGlobalLevel(level)
	}
	zerolog.DefaultContextLogger = &log.Logger

	// Service context
	ctx := context.Background()

	// Create the Redis client
	redisClientOption := rueidis.ClientOption{
		InitAddress:      []string{cfg.Redis.Host + ":" + strconv.Itoa(cfg.Redis.Port)},
		ClientName:       "greenstar-backend",
		BlockingPoolSize: cfg.Redis.PoolSize,
	}
	if cfg.Redis.TLS {
		redisClientOption.TLSConfig = &tls.Config{ServerName: cfg.Redis.Host}
	}
	redisClient, err := rueidis.NewClient(redisClientOption)
	if err != nil {
		log.Ctx(ctx).Fatal().Err(err).Msg("Failed to connect to Redis")
	}
	defer redisClient.Close()

	// Create the Neo4j client
	neo4jURL := fmt.Sprintf("bolt+s://%s:%d", cfg.Neo4j.Host, cfg.Neo4j.Port)
	neo4jDriver, err := neo4j.NewDriverWithContext(neo4jURL, neo4j.NoAuth())
	if err != nil {
		log.Ctx(ctx).Fatal().Err(err).Msg("Failed to connect to Neo4j")
	}
	defer neo4jDriver.Close(ctx)

	// Creates a Pub/Sub client.
	pubSubClient, err := pubsub.NewClient(ctx, cfg.GoogleCloud.ProjectID)
	if err != nil {
		log.Ctx(ctx).Fatal().Err(err).Msg("Failed to create Pub/Sub client")
	}
	defer pubSubClient.Close()
	if cfg.DevMode {
		if cfg.DevModeID == "" {
			log.Ctx(ctx).Fatal().Msg("Dev mode ID must be set when dev mode is enabled")
		}
		subs := map[string]*string{
			cfg.DevModeID + "-admin":      &cfg.GoogleCloud.AdminSubscription,
			cfg.DevModeID + "-operations": &cfg.GoogleCloud.OperationsSubscription,
			cfg.DevModeID + "-public":     &cfg.GoogleCloud.PublicSubscription,
		}
		for topicName, subNameCfg := range subs {

			// Create topic
			topic, err := pubsubutil.CreateDevTopic(ctx, pubSubClient, topicName)
			if err != nil {
				log.Ctx(ctx).Fatal().Err(err).Str("topic", topicName).Msg("Failed to create Pub/Sub topic")
			}

			// Create dead-letter queue topic
			if _, err := pubsubutil.CreateDevTopic(ctx, pubSubClient, topicName+"-dl"); err != nil {
				log.Ctx(ctx).Fatal().Err(err).Str("topic", topicName+"-dl").Msg("Failed to create Pub/Sub topic")
			}

			// Create subscription
			sub, err := pubsubutil.CreateDevSubscription(ctx, pubSubClient, topic, topicName)
			if err != nil {
				log.Ctx(ctx).Fatal().Err(err).Str("subscription", topicName).Msg("Failed to create Pub/Sub subscription")
			} else {
				*subNameCfg = sub.ID()
			}
		}
	}

	// Setup routes
	router := ginutil.NewGin(cfg.DevMode)
	router.Use(redisutil.CreateSetRedisMiddleware(redisClient))
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = cfg.HTTP.CORS.AllowedOrigins
	corsConfig.AddAllowMethods(cfg.HTTP.CORS.AllowMethods...)
	corsConfig.AddAllowHeaders(cfg.HTTP.CORS.AllowHeaders...)
	corsConfig.AddExposeHeaders(cfg.HTTP.CORS.ExposeHeaders...)
	corsConfig.AllowAllOrigins = false
	corsConfig.AllowBrowserExtensions = false
	corsConfig.AllowCredentials = !cfg.HTTP.CORS.DisableCredentials
	corsConfig.AllowFiles = false
	corsConfig.AllowWebSockets = true
	corsConfig.AllowWildcard = true
	corsConfig.MaxAge = cfg.HTTP.CORS.MaxAge
	corsMiddleware := cors.New(corsConfig)
	router.OPTIONS("/*path", corsMiddleware)
	router.Use(corsMiddleware)

	// Admin routes
	adminAPIResolver := &adminResolver.Resolver{Neo4j: neo4jDriver}
	adminAPISchema := adminGQL.NewExecutableSchema(adminGQL.Config{Resolvers: adminAPIResolver})
	adminAPIHandler := handler.NewDefaultServer(adminAPISchema)
	adminAPIHandler.SetErrorPresenter(gqlutil.ErrorPresenter)
	adminAPIHandler.SetRecoverFunc(gqlutil.PanicRecoverer)
	adminAPIPlaygroundHandler := playground.Handler("Admin GraphQL playground", "/api/admin/query")
	router.Group("/api/admin").
		Use(auth.CreateJWTValidationMiddleware(cfg.Auth.Auth0Domain, []string{"https://greenstar.kfirs.com/api/admin"})).
		GET("/playground", func(c *gin.Context) { adminAPIPlaygroundHandler(c.Writer, c.Request) }).
		POST("/query", func(c *gin.Context) { adminAPIHandler.ServeHTTP(c.Writer, c.Request) })

	adminAPIExecutor := executor.New(adminAPISchema)
	adminAPIExecutor.SetErrorPresenter(gqlutil.ErrorPresenter)
	adminAPIExecutor.SetRecoverFunc(gqlutil.PanicRecoverer)
	go func(subscriptionName string, handler func(context.Context, *pubsub.Message)) {
		subscription := pubSubClient.Subscription(subscriptionName)
		if err := subscription.Receive(ctx, handler); err != nil && err != context.Canceled {
			log.Ctx(ctx).Fatal().Err(err).Str("subscription", subscriptionName).Msg("Cloud Pub/Sub subscriber failed")
		}
	}(cfg.GoogleCloud.AdminSubscription, adminAPIResolver.HandleMessage)

	// Operations routes
	operationsAPIResolver := &operationsResolver.Resolver{Redis: redisClient}
	operationsAPISchema := operationsGQL.NewExecutableSchema(operationsGQL.Config{Resolvers: operationsAPIResolver})
	operationsAPIHandler := handler.NewDefaultServer(operationsAPISchema)
	operationsAPIHandler.SetErrorPresenter(gqlutil.ErrorPresenter)
	operationsAPIHandler.SetRecoverFunc(gqlutil.PanicRecoverer)
	operationsAPIPlaygroundHandler := playground.Handler("Operations GraphQL playground", "/api/operations/query")
	router.Group("/api/operations").
		Use(auth.CreateJWTValidationMiddleware(cfg.Auth.Auth0Domain, []string{"https://greenstar.kfirs.com/api/operations"})).
		GET("/playground", func(c *gin.Context) { operationsAPIPlaygroundHandler(c.Writer, c.Request) }).
		POST("/query", func(c *gin.Context) { operationsAPIHandler.ServeHTTP(c.Writer, c.Request) })

	operationsAPIExecutor := executor.New(operationsAPISchema)
	operationsAPIExecutor.SetErrorPresenter(gqlutil.ErrorPresenter)
	operationsAPIExecutor.SetRecoverFunc(gqlutil.PanicRecoverer)
	go func(subscriptionName string, handler func(context.Context, *pubsub.Message)) {
		subscription := pubSubClient.Subscription(subscriptionName)
		if err := subscription.Receive(ctx, handler); err != nil && err != context.Canceled {
			log.Ctx(ctx).Fatal().Err(err).Str("subscription", subscriptionName).Msg("Cloud Pub/Sub subscriber failed")
		}
	}(cfg.GoogleCloud.OperationsSubscription, operationsAPIResolver.HandleMessage)

	// Public API routes
	publicAPIResolver := &publicResolver.Resolver{Neo4j: neo4jDriver, Redis: redisClient}
	publicAPISchema := publicGQL.NewExecutableSchema(publicGQL.Config{Resolvers: publicAPIResolver})
	publicAPIHandler := handler.NewDefaultServer(publicAPISchema)
	publicAPIHandler.SetErrorPresenter(gqlutil.ErrorPresenter)
	publicAPIHandler.SetRecoverFunc(gqlutil.PanicRecoverer)
	publicAPIPlaygroundHandler := playground.Handler("Public GraphQL playground", "/api/public/query")
	router.Group("/api/public").
		Use(auth.CreateJWTValidationMiddleware(cfg.Auth.Auth0Domain, []string{"https://greenstar.kfirs.com/api/public"})).
		GET("/playground", func(c *gin.Context) { publicAPIPlaygroundHandler(c.Writer, c.Request) }).
		POST("/query", func(c *gin.Context) { publicAPIHandler.ServeHTTP(c.Writer, c.Request) })

	publicAPIExecutor := executor.New(publicAPISchema)
	publicAPIExecutor.SetErrorPresenter(gqlutil.ErrorPresenter)
	publicAPIExecutor.SetRecoverFunc(gqlutil.PanicRecoverer)
	go func(subscriptionName string, handler func(context.Context, *pubsub.Message)) {
		subscription := pubSubClient.Subscription(subscriptionName)
		if err := subscription.Receive(ctx, handler); err != nil && err != context.Canceled {
			log.Ctx(ctx).Fatal().Err(err).Str("subscription", subscriptionName).Msg("Cloud Pub/Sub subscriber failed")
		}
	}(cfg.GoogleCloud.PublicSubscription, publicAPIResolver.HandleMessage)

	// Setup generic, public, get-org-ID API
	router.Group("/api/util").
		GET("/v1/organizations/:tenant", auth.CreateGetOrgIDHandler(cfg.Auth))

	// Setup health checks under "/healthz"
	if err := gin_healthcheck.New(router, config.DefaultConfig(), nil); err != nil {
		log.Ctx(ctx).Fatal().Err(err).Msg("Failed to setup health checks")
	}

	// Setup HTTP server
	httpServer := &http.Server{Addr: ":" + strconv.Itoa(cfg.HTTP.Port), Handler: router}

	// Start server
	log.Ctx(ctx).Info().Str("addr", httpServer.Addr).Msg("Starting HTTP server")
	if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Err(err).Msg("HTTP server failed")
	}
}
