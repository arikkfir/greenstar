package main

import (
	"context"
	"github.com/arikkfir/greenstar/common"
	"github.com/arikkfir/greenstar/xlsconverter/internal"
	"github.com/rs/zerolog/log"
)

func main() {
	config := internal.NewConfig()

	// Create a context which will be cancelled on termination
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create the Redis client
	redisClient := common.NewRedisClient(config.Redis)
	defer redisClient.Close()

	// Subscribe to Redis queue
	server := internal.NewServer(*config, redisClient)
	if err := server.Run(ctx); err != nil {
		log.Fatal().Err(err).Msg("Server failed")
	}
}
