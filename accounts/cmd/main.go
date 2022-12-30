package main

import (
	"context"
	"fmt"
	"github.com/arikkfir/greenstar/accounts/internal"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/rs/zerolog/log"
	"github.com/rueian/rueidis"
	"strconv"
)

func main() {
	config := internal.NewConfig()
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

	// Setup Redis Search index
	if resp := redisClient.Do(ctx, redisClient.B().FtDropindex().Index("accounts").Build()); resp.Error() != nil {
		log.Ctx(ctx).Warn().Err(resp.Error()).Msg("Failed to drop index in preparation of recreating it")
	}
	createIndexCommand := redisClient.B().
		FtCreate().Index("accounts").OnJson().Prefix(1).Prefix("account:").
		Schema().
		FieldName("$.id").As("id").Text().
		FieldName("$.displayName").As("name").Text().
		FieldName("$.labels").As("labels").Text().
		Build()
	if resp := redisClient.Do(ctx, createIndexCommand); resp.Error() != nil {
		log.Ctx(ctx).Fatal().Err(resp.Error()).Msg("Failed to create index")
	}

	// Create HTTP server
	server, err := internal.NewServer(*config, redisClient, neo4jDriver)
	if err != nil {
		log.Ctx(ctx).Fatal().Err(err).Msg("Failed to create server")
	}

	// Start server
	if err := server.Run(ctx); err != nil {
		log.Err(err).Msg("HTTP server failed")
	}
}
