package resolver

import (
	"cloud.google.com/go/pubsub"
	"context"
	_ "embed"
	"github.com/arik-kfir/greenstar/backend/util/neo4jutil"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/rs/zerolog/log"
	"github.com/rueian/rueidis"
	"text/template"
)

var (
	//go:embed scraper-job.tmpl.yaml
	scraperJobTmplString string
	scraperJobTmpl       *template.Template
)

func init() {
	tmpl, err := template.New("scraper-job").Parse(scraperJobTmplString)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to parse scraper-job.tmpl.yaml")
	} else {
		scraperJobTmpl = tmpl
	}
}

type Resolver struct {
	Redis rueidis.Client
	Neo4j neo4j.DriverWithContext
}

func (r *Resolver) getNeo4jSession(ctx context.Context, mode neo4j.AccessMode) neo4j.SessionWithContext {
	return r.Neo4j.NewSession(ctx, neo4j.SessionConfig{
		AccessMode:   mode,
		DatabaseName: "", // TODO: auth.GetClaims(ctx).Tenant,
		BoltLogger:   &neo4jutil.Neo4jZerologBoltLogger{Logger: log.Ctx(ctx)},
	})
}

func (r *Resolver) HandleMessage(ctx context.Context, msg *pubsub.Message) {
	panic("implement me")
}
