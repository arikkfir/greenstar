package main

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/arikkfir/greenstar/api/internal"
	"github.com/arikkfir/greenstar/common"
	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rueian/rueidis"
	"html/template"
	"net/http"
	"strconv"
)

type Config struct {
	General common.GeneralConfig
	Redis   RedisConfig `group:"redis" namespace:"redis" env-namespace:"REDIS"`
	Neo4j   Neo4jConfig `group:"neo4j" namespace:"neo4j" env-namespace:"NEO4J"`
	HTTP    HTTPConfig  `group:"http" namespace:"http" env-namespace:"HTTP"`
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

//go:embed templates
var templates embed.FS

func main() {
	config := Config{}
	common.ReadConfig(&config)
	config.General.Apply()

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
	gqlServer := handler.NewDefaultServer(internal.NewExecutableSchema(internal.Config{Resolvers: &internal.Resolver{
		Redis: redisClient,
		Neo4j: neo4jDriver,
	}}))
	gqlPlaygroundHandler := playground.Handler("GraphQL playground", "/")

	// Initialize Gin framework
	gin.DefaultWriter = log.Logger.Level(zerolog.TraceLevel)
	gin.DefaultErrorWriter = log.Logger.Level(zerolog.ErrorLevel)
	if config.General.DevMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// Setup routes
	router := gin.Default()
	router.SetHTMLTemplate(template.Must(template.ParseFS(templates, "templates/*.html")))
	router.MaxMultipartMemory = 8 << 20
	if err != nil {
		log.Ctx(ctx).Fatal().Err(err).Msg("Failed to initialize Gin router")
	}
	router.GET("/playground", func(c *gin.Context) { gqlPlaygroundHandler(c.Writer, c.Request) })
	// TODO: wrap each individual GraphQL query in a transaction (to ensure consistency between resolvers)
	router.POST("/", func(c *gin.Context) { gqlServer.ServeHTTP(c.Writer, c.Request) })

	// Setup HTTP server
	httpServer := &http.Server{Addr: ":" + strconv.Itoa(config.HTTP.Port), Handler: router}

	// Start the transactions receiver
	go func() {
		processor := &internal.XLSXProcessor{
			Resolver: &internal.Resolver{
				Redis: redisClient,
				Neo4j: neo4jDriver,
			},
		}
		log.Ctx(ctx).Info().Msg("Starting XLSX processor")
		processor.Run(ctx)
	}()

	// Start server
	log.Ctx(ctx).Info().Str("addr", httpServer.Addr).Msg("Starting HTTP server")
	if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Err(err).Msg("HTTP server failed")
	}
}
