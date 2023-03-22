package resolver

import (
	"context"
	"github.com/arikkfir/greenstar/backend/auth"
	"github.com/arikkfir/greenstar/backend/util/neo4jutil"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/rs/zerolog/log"
	"github.com/rueian/rueidis"
)

const (
	processXLSXChannelName = "process-xlsx"
)

type Resolver struct {
	Redis rueidis.Client
	Neo4j neo4j.DriverWithContext
}

func (r *Resolver) getNeo4jSession(ctx context.Context, mode neo4j.AccessMode) neo4j.SessionWithContext {
	return r.Neo4j.NewSession(ctx, neo4j.SessionConfig{
		AccessMode:   mode,
		DatabaseName: auth.GetSession(ctx).Tenant,
		BoltLogger:   &neo4jutil.Neo4jZerologBoltLogger{Logger: log.Ctx(ctx)},
	})
}
