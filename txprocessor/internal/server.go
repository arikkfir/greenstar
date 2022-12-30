package internal

import (
	"context"
	"encoding/json"
	"github.com/arikkfir/greenstar/txprocessor/pkg"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"
)

type Server struct {
	Config Config
	Redis  *redis.Client
}

func NewServer(config Config, Redis *redis.Client) *Server {
	return &Server{
		Config: config,
		Redis:  Redis,
	}
}

func (s *Server) Run(ctx context.Context) error {
	logger := *zerolog.Ctx(ctx)
	subscriber := s.Redis.Subscribe(ctx, pkg.InputChannelName)
	for {
		msg, err := subscriber.ReceiveMessage(ctx)
		if err != nil {
			logger.Fatal().Err(err).Msg("Failed to receive message from Redis")
		}

		tx := pkg.Transaction{}
		if err := json.Unmarshal([]byte(msg.Payload), &tx); err != nil {
			logger.Fatal().Err(err).Msg("Failed to unmarshal message")
		}

		logger.Info().Interface("tx", tx).Msg("Received transaction")
	}
}
