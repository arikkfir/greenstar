package resolver

import (
	"context"
	"github.com/arikkfir/greenstar/api/model"
	"github.com/arikkfir/greenstar/api/util"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/rs/zerolog/log"
	"github.com/rueian/rueidis"
)

const (
	convertXLSChannelName  = "convert-xls-to-xlsx"
	processXLSXChannelName = "process-xlsx"
)

type Resolver struct {
	Redis rueidis.Client
	Neo4j neo4j.DriverWithContext
}

func (r *Resolver) getNeo4jSession(ctx context.Context, mode neo4j.AccessMode) neo4j.SessionWithContext {
	return r.Neo4j.NewSession(ctx, neo4j.SessionConfig{
		AccessMode:   mode,
		DatabaseName: util.GetTenantSlug(ctx),
		BoltLogger:   &util.Neo4jZerologBoltLogger{Logger: log.Ctx(ctx)},
	})
}

func (r *Resolver) readAccount(node neo4j.Node) *model.Account {
	return &model.Account{
		ID:          node.Props["accountID"].(string),
		DisplayName: node.Props["displayName"].(string),
	}
}

func (r *Resolver) readKeyAndValues(keyAndValueList []map[string]string) []*model.KeyAndValue {
	labels := make([]*model.KeyAndValue, 0)
	for _, m := range keyAndValueList {
		labels = append(labels, &model.KeyAndValue{
			Key:   m["label"],
			Value: m["value"],
		})
	}
	return labels
}
