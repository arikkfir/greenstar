package main

//go:generate go run github.com/99designs/gqlgen generate

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/arikkfir/greenstar/api/gql"
	"github.com/arikkfir/greenstar/api/resolver"
	"github.com/arikkfir/greenstar/api/util"
	"github.com/gin-gonic/gin"
	"github.com/jessevdk/go-flags"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/rueian/rueidis"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"html/template"
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
	LogLevel string      `env:"LOG_LEVEL" value-name:"LEVEL" long:"log-level" description:"Log level" default:"info" enum:"trace,debug,info,warn,error,fatal,panic"`
	DevMode  bool        `env:"DEV_MODE" long:"dev-mode" description:"Development mode"`
	Redis    RedisConfig `group:"redis" namespace:"redis" env-namespace:"REDIS"`
	Neo4j    Neo4jConfig `group:"neo4j" namespace:"neo4j" env-namespace:"NEO4J"`
	HTTP     HTTPConfig  `group:"http" namespace:"http" env-namespace:"HTTP"`
}

type RedisConfig struct {
	Host     string `env:"HOST" value-name:"HOST" long:"host" description:"Redis host name" default:"localhost"`
	Port     int    `env:"PORT" value-name:"PORT" long:"port" description:"Redis port" default:"6379"`
	PoolSize int    `env:"POOL_SIZE" value-name:"POOL_SIZE" long:"pool-size" description:"Redis connection pool size" default:"3"`
}

type Neo4jConfig struct {
	Host string `env:"HOST" value-name:"HOST" long:"host" description:"Neo4j host name" default:"localhost"`
	Port int    `env:"PORT" value-name:"PORT" long:"port" description:"Neo4j port" default:"7687"`
}

type HTTPConfig struct {
	Port int `env:"PORT" value-name:"PORT" long:"port" description:"Port to listen on" default:"8000"`
}

var (
	//go:embed templates
	templates embed.FS
)

func main() {
	// Ensure we have a random seed
	rand.Seed(time.Now().UnixNano())

	config := Config{}
	parser := flags.NewParser(&config, flags.HelpFlag|flags.PassDoubleDash)
	if _, err := parser.Parse(); err != nil {
		fmt.Printf("ERROR: %s\n\n", err)
		parser.WriteHelp(os.Stderr)
		os.Exit(1)
	}
	if config.DevMode {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		zerolog.DefaultContextLogger = &log.Logger
	}
	if level, err := zerolog.ParseLevel(config.LogLevel); err != nil {
		log.Fatal().Err(err).Msg("Failed to parse config")
	} else {
		zerolog.SetGlobalLevel(level)
	}
	zerolog.DefaultContextLogger = &log.Logger

	// Service context
	ctx := context.Background()

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
	gqlResolver := &resolver.Resolver{Redis: redisClient, Neo4j: neo4jDriver}
	gqlSchema := gql.NewExecutableSchema(gql.Config{Resolvers: gqlResolver})
	gqlHandler := handler.NewDefaultServer(gqlSchema)
	gqlHandler.SetErrorPresenter(func(ctx context.Context, e error) *gqlerror.Error {

		if gqlErr, ok := e.(*gqlerror.Error); ok {
			log.Ctx(ctx).Debug().Err(e).Msg("Presenting error")
			return gqlErr
		}

		path := graphql.GetPath(ctx)

		log.Ctx(ctx).
			Error().
			Err(e).
			Interface("gqlPath", path).
			Msg("GraphQL handling error")

		return &gqlerror.Error{
			Message:    e.Error(),
			Path:       path,
			Extensions: map[string]interface{}{"code": "INTERNAL_SERVER_ERROR"},
		}
	})
	gqlHandler.SetRecoverFunc(func(ctx context.Context, err interface{}) error {
		log.Ctx(ctx).Debug().Interface("panic", err).Msg("Recovered from panic")
		if e, ok := err.(error); ok {
			return e
		}
		return fmt.Errorf("internal server error: %v", err)
	})
	gqlPlaygroundHandler := playground.Handler("GraphQL playground", "/")

	// Start the transactions receiver
	go func() {
		log.Ctx(ctx).Info().Msg("Starting XLSX processor")
		gqlResolver.TransactionsXLSXProcessor().Run(ctx)
	}()

	// Initialize Gin framework
	gin.DefaultWriter = log.Logger.Level(zerolog.TraceLevel)
	gin.DefaultErrorWriter = log.Logger.Level(zerolog.ErrorLevel)
	if config.DevMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// Setup routes
	// TODO: add authentication to tenants router group
	router := gin.New()
	router.SetHTMLTemplate(template.Must(template.ParseFS(templates, "templates/*.html")))
	router.MaxMultipartMemory = 8 << 20
	router.Use(util.SetLoggerMiddleware)
	router.Use(util.AccessLogMiddleware)

	tenants := router.Group("/tenants")
	tenants.POST("/", util.NewCreateTenantHandler(neo4jDriver))
	tenants.Group("/:tenant").
		Use(util.AddTenantSlugToContextMiddleware).
		DELETE("/", util.NewDeleteTenantHandler(neo4jDriver)).
		GET("/playground", func(c *gin.Context) { gqlPlaygroundHandler(c.Writer, c.Request) }).
		POST("/query", func(c *gin.Context) { gqlHandler.ServeHTTP(c.Writer, c.Request) })

	// Setup HTTP server
	httpServer := &http.Server{Addr: ":" + strconv.Itoa(config.HTTP.Port), Handler: router}

	// Start server
	log.Ctx(ctx).Info().Str("addr", httpServer.Addr).Msg("Starting HTTP server")
	if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Err(err).Msg("HTTP server failed")
	}
}
