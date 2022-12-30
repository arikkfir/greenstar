package internal

import (
	"context"
	"errors"
	"fmt"
	"github.com/arikkfir/greenstar/common/ginutil"
	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rueian/rueidis"
	"net/http"
	"strconv"
)

type Server struct {
	Config     Config
	Redis      rueidis.Client
	Neo4j      neo4j.DriverWithContext
	HTTPServer *http.Server
}

func NewServer(config Config, Redis rueidis.Client, neo4jDriver neo4j.DriverWithContext) (*Server, error) {
	s := &Server{Config: config, Redis: Redis, Neo4j: neo4jDriver}

	// Initialize Gin framework
	gin.DefaultWriter = log.Logger.Level(zerolog.TraceLevel)
	gin.DefaultErrorWriter = log.Logger.Level(zerolog.ErrorLevel)
	if config.General.DevMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// Setup routes
	router, err := ginutil.NewRouter()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize router: %w", err)
	}
	router.POST("/", s.createAccount)
	router.GET("/", s.getAccounts)
	router.PUT("/:id", s.putAccount)
	router.DELETE("/:id", s.deleteAccount)

	// Setup HTTP server
	s.HTTPServer = &http.Server{Addr: ":" + strconv.Itoa(config.HTTP.Port), Handler: router}
	return s, nil
}

func (s *Server) Run(ctx context.Context) error {
	log.Ctx(ctx).Info().Str("addr", s.HTTPServer.Addr).Msg("Starting HTTP server")
	if err := s.HTTPServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}
