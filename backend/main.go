package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	adminGQL "github.com/arikkfir/greenstar/backend/admin/gql"
	adminResolver "github.com/arikkfir/greenstar/backend/admin/resolver"
	"github.com/arikkfir/greenstar/backend/auth"
	opsGQL "github.com/arikkfir/greenstar/backend/operations/gql"
	opsResolver "github.com/arikkfir/greenstar/backend/operations/resolver"
	publicGQL "github.com/arikkfir/greenstar/backend/public/gql"
	publicResolver "github.com/arikkfir/greenstar/backend/public/resolver"
	"github.com/arikkfir/greenstar/backend/util/ginutil"
	"github.com/arikkfir/greenstar/backend/util/gqlutil"
	"github.com/arikkfir/greenstar/backend/util/httputil"
	"github.com/arikkfir/greenstar/backend/util/natsutil"
	"github.com/arikkfir/greenstar/backend/util/neo4jutil"
	"github.com/arikkfir/greenstar/backend/util/redisutil"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jessevdk/go-flags"
	"github.com/markbates/goth/providers/google"
	"github.com/nats-io/nats.go"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/rueian/rueidis"
	"github.com/tavsec/gin-healthcheck"
	"github.com/tavsec/gin-healthcheck/config"
	"golang.org/x/oauth2"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

func init() {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
}

type Config struct {
	LogLevel string                `env:"LOG_LEVEL" value-name:"LEVEL" long:"log-level" description:"Log level" default:"info" enum:"trace,debug,info,warn,error,fatal,panic"`
	DevMode  bool                  `env:"DEV_MODE" long:"dev-mode" description:"Development mode"`
	AppURL   string                `env:"APP_URL" value-name:"URL" long:"app-url" description:"URL to redirect to after login or logout" required:"yes"`
	Redis    redisutil.RedisConfig `group:"redis" namespace:"redis" env-namespace:"REDIS"`
	Neo4j    neo4jutil.Neo4jConfig `group:"neo4j" namespace:"neo4j" env-namespace:"NEO4J"`
	NATS     natsutil.NATSConfig   `group:"nats" namespace:"nats" env-namespace:"NATS"`
	HTTP     httputil.HTTPConfig   `group:"http" namespace:"http" env-namespace:"HTTP"`
	Auth     auth.Config           `group:"auth" namespace:"auth" env-namespace:"AUTH"`
}

func (c *Config) IsSecure() bool {
	secure, err := httputil.IsSecure(c.AppURL)
	if err != nil {
		panic(fmt.Errorf("failed to check if URL '%s' is secure: %w", c.AppURL, err))
	}
	return secure
}

func main() {
	cfg := Config{}

	// Ensure we have a random seed
	rand.Seed(time.Now().UnixNano())

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
	redisClient, err := rueidis.NewClient(rueidis.ClientOption{
		InitAddress:      []string{cfg.Redis.Host + ":" + strconv.Itoa(cfg.Redis.Port)},
		ClientName:       "greenstar.operations",
		BlockingPoolSize: cfg.Redis.PoolSize,
	})
	if err != nil {
		log.Ctx(ctx).Fatal().Err(err).Msg("Failed to connect to Redis")
	}
	defer redisClient.Close()

	// Create the Neo4j client
	neo4jURL := fmt.Sprintf("neo4j://%s:%d", cfg.Neo4j.Host, cfg.Neo4j.Port)
	neo4jDriver, err := neo4j.NewDriverWithContext(neo4jURL, neo4j.NoAuth())
	if err != nil {
		log.Ctx(ctx).Fatal().Err(err).Msg("Failed to connect to Neo4j")
	}
	defer neo4jDriver.Close(ctx)

	// Connect to NATS
	natsClient, err := nats.Connect(cfg.NATS.URL, nats.RetryOnFailedConnect(true))
	if err != nil {
		log.Ctx(ctx).Fatal().Err(err).Msg("Failed to connect to NATS")
	}
	defer natsClient.Close()
	jetStreamClient, err := natsClient.JetStream()
	if err != nil {
		log.Ctx(ctx).Fatal().Err(err).Msg("Failed to establish JetStream connection")
	}

	// Listen for operation commands
	opsSub, err := jetStreamClient.Subscribe("OPERATIONS", func(m *nats.Msg) {
		log.Ctx(ctx).Info().Interface("msg", m).Msg("Received a JetStream message")
		if err := m.Ack(); err != nil {
			log.Ctx(ctx).Error().Err(err).Interface("msg", m).Msg("Failed to ack message")
		}
	}, nats.ManualAck(), nats.AckExplicit(), nats.MaxAckPending(10))
	defer opsSub.Unsubscribe()

	// Setup Google OAuth
	var googleOauthConfig = &oauth2.Config{
		RedirectURL:  cfg.Auth.Google.CallbackURL,
		ClientID:     cfg.Auth.Google.ClientID,
		ClientSecret: cfg.Auth.Google.ClientSecret,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	// Authentication manager
	authenticator := &auth.Authenticator{
		Config:              &cfg.Auth,
		SecureCookies:       cfg.IsSecure(),
		DefaultPostLoginURL: cfg.AppURL,
		OAuth:               googleOauthConfig,
	}

	// Setup routes
	router := ginutil.NewGin(cfg.DevMode)
	router.Use(redisutil.CreateSetRedisMiddleware(redisClient))
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:        false,
		AllowOrigins:           cfg.HTTP.CORS.AllowedOrigins,
		AllowMethods:           cfg.HTTP.CORS.AllowMethods,
		AllowHeaders:           cfg.HTTP.CORS.AllowHeaders,
		AllowCredentials:       !cfg.HTTP.CORS.DisableCredentials,
		ExposeHeaders:          cfg.HTTP.CORS.ExposeHeaders,
		MaxAge:                 cfg.HTTP.CORS.MaxAge,
		AllowBrowserExtensions: false,
		AllowFiles:             false,
		AllowWebSockets:        true,
		AllowWildcard:          true,
	}))

	// Authentication routes
	router.Group("/auth").
		GET("/google/login", authenticator.InitiateLogin).
		GET("/google/callback", authenticator.HandleCallback).
		GET("/google/logout", authenticator.HandleLogout).
		GET("/user", authenticator.CreateMiddlewareForFilteringUnauthenticated("greenstar.auth"), authenticator.HandleUserInfo)

	// Operations routes
	opsGraphResolver := &opsResolver.Resolver{Redis: redisClient}
	opsGraphSchema := opsGQL.NewExecutableSchema(opsGQL.Config{Resolvers: opsGraphResolver})
	opsGraphHandler := handler.NewDefaultServer(opsGraphSchema)
	opsGraphHandler.SetErrorPresenter(gqlutil.ErrorPresenter)
	opsGraphHandler.SetRecoverFunc(gqlutil.PanicRecoverer)
	opsGraphPlaygroundHandler := playground.Handler("Operations GraphQL playground", "/operations/query")
	router.Group("/api/operations").
		Use(authenticator.CreateMiddlewareForFilteringUnauthenticated("greenstar.operations")).
		GET("/playground", func(c *gin.Context) { opsGraphPlaygroundHandler(c.Writer, c.Request) }).
		POST("/query", func(c *gin.Context) { opsGraphHandler.ServeHTTP(c.Writer, c.Request) })

	// Admin routes
	adminGraphResolver := &adminResolver.Resolver{Neo4j: neo4jDriver}
	adminGraphSchema := adminGQL.NewExecutableSchema(adminGQL.Config{Resolvers: adminGraphResolver})
	adminGraphHandler := handler.NewDefaultServer(adminGraphSchema)
	adminGraphHandler.SetErrorPresenter(gqlutil.ErrorPresenter)
	adminGraphHandler.SetRecoverFunc(gqlutil.PanicRecoverer)
	adminGraphPlaygroundHandler := playground.Handler("Admin GraphQL playground", "/admin/query")
	router.Group("/api/admin").
		Use(authenticator.CreateMiddlewareForFilteringUnauthenticated("greenstar.admin")).
		GET("/playground", func(c *gin.Context) { adminGraphPlaygroundHandler(c.Writer, c.Request) }).
		POST("/query", func(c *gin.Context) { adminGraphHandler.ServeHTTP(c.Writer, c.Request) })

	// Public API routes
	publicGraphResolver := &publicResolver.Resolver{Neo4j: neo4jDriver}
	publicGraphSchema := publicGQL.NewExecutableSchema(publicGQL.Config{Resolvers: publicGraphResolver})
	publicGraphHandler := handler.NewDefaultServer(publicGraphSchema)
	publicGraphHandler.SetErrorPresenter(gqlutil.ErrorPresenter)
	publicGraphHandler.SetRecoverFunc(gqlutil.PanicRecoverer)
	publicGraphPlaygroundHandler := playground.Handler("Public GraphQL playground", "/public/query")
	router.Group("/api/public").
		Use(authenticator.CreateMiddlewareForFilteringUnauthenticated("greenstar.public")).
		GET("/playground", func(c *gin.Context) { publicGraphPlaygroundHandler(c.Writer, c.Request) }).
		POST("/query", func(c *gin.Context) { publicGraphHandler.ServeHTTP(c.Writer, c.Request) })

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
