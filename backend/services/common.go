package services

import (
	"context"
	"github.com/arik-kfir/greenstar/backend/util"
	descope "github.com/descope/go-sdk/descope/client"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/redis/rueidis"
	"github.com/rs/zerolog/log"
)

const GlobalTenantID = "global"

type Service struct {
	Descope *descope.DescopeClient
	Neo4j   neo4j.DriverWithContext
	Redis   rueidis.Client
}

func (s *Service) getNeo4jSessionForSystem(ctx context.Context, mode neo4j.AccessMode) neo4j.SessionWithContext {
	return s.Neo4j.NewSession(ctx, neo4j.SessionConfig{
		AccessMode:   mode,
		DatabaseName: "neo4j",
		BoltLogger:   &util.Neo4jZerologBoltLogger{Logger: log.Ctx(ctx)},
	})
}

func (s *Service) getNeo4jSessionForTenant(ctx context.Context, mode neo4j.AccessMode, tenantID string) neo4j.SessionWithContext {
	return s.Neo4j.NewSession(ctx, neo4j.SessionConfig{
		AccessMode:   mode,
		DatabaseName: tenantID,
		BoltLogger:   &util.Neo4jZerologBoltLogger{Logger: log.Ctx(ctx)},
	})
}
