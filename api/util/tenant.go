package util

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/rs/zerolog/log"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
)

const tenantKey = "###greenstar$tenantSlug"

var (
	tenantRE           = regexp.MustCompile(`([a-zA-Z0-9_-]+)\..*`)
	tenantAllowedRunes = []rune("abcdefghijklmnopqrstuvwxyz")
)

func GetTenantSlug(ctx context.Context) string {
	return ctx.Value(tenantKey).(string)
}

func setTenantSlug(ctx context.Context, slug string) context.Context {
	return context.WithValue(ctx, tenantKey, slug)
}

func AddTenantSlugToContextMiddleware(c *gin.Context) {
	ctx := c.Request.Context()

	matches := tenantRE.FindStringSubmatch(c.Request.Host)
	if matches == nil {
		c.Negotiate(http.StatusBadRequest, gin.Negotiate{
			Offered:  []string{gin.MIMEPlain, gin.MIMEHTML, gin.MIMEJSON, gin.MIMEYAML},
			Data:     errors.New("hostname must start with a tenant, e.g. acme.api.greenstar.com"),
			HTMLName: strconv.Itoa(http.StatusBadRequest),
		})
		return
	}

	tenantSlug := matches[1]
	// TODO: validate tenant exists

	ctxWithTenantSlug := setTenantSlug(ctx, tenantSlug)
	c.Request = c.Request.WithContext(ctxWithTenantSlug)
	c.Next()
}

func NewCreateTenantDBMiddleware(neo4jDriver neo4j.DriverWithContext) func(*gin.Context) {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		session := neo4jDriver.NewSession(ctx, neo4j.SessionConfig{
			AccessMode: neo4j.AccessModeWrite,
			BoltLogger: &Neo4jZerologBoltLogger{Logger: log.Ctx(ctx)},
		})
		defer session.Close(ctx)

		//goland:noinspection SqlNoDataSourceInspection
		const createTenantDBQuery = `CREATE DATABASE $tenant IF NOT EXISTS`
		tenantSlug := GetTenantSlug(ctx)

		log.Ctx(ctx).Debug().Str("tenant", tenantSlug).Msg("Creating tenant Neo4j database")
		_, err := session.Run(c.Request.Context(), createTenantDBQuery, map[string]any{"tenant": tenantSlug})
		if err != nil {
			log.Ctx(ctx).Error().Str("tenant", tenantSlug).Err(err).Msg("Failed creating tenant Neo4j database")
			c.Negotiate(http.StatusInternalServerError, gin.Negotiate{
				Offered:  []string{gin.MIMEPlain, gin.MIMEHTML, gin.MIMEJSON, gin.MIMEYAML},
				Data:     err,
				HTMLName: strconv.Itoa(http.StatusInternalServerError),
			})
			return
		}

		c.Next()
	}
}

func NewDropTenantDBHandler(neo4jDriver neo4j.DriverWithContext) func(*gin.Context) {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		session := neo4jDriver.NewSession(ctx, neo4j.SessionConfig{
			AccessMode: neo4j.AccessModeWrite,
			BoltLogger: &Neo4jZerologBoltLogger{Logger: log.Ctx(ctx)},
		})
		defer session.Close(ctx)

		//goland:noinspection SqlNoDataSourceInspection
		const dropTenantDBQuery = `DROP DATABASE $tenant IF EXISTS`
		tenantSlug := GetTenantSlug(ctx)

		log.Ctx(ctx).Debug().Str("tenant", tenantSlug).Msg("Dropping tenant Neo4j database")
		_, err := session.Run(c.Request.Context(), dropTenantDBQuery, map[string]any{"tenant": tenantSlug})
		if err != nil {
			log.Ctx(ctx).Error().Str("tenant", tenantSlug).Err(err).Msg("Failed dropping tenant Neo4j database")
			c.Negotiate(http.StatusInternalServerError, gin.Negotiate{
				Offered:  []string{gin.MIMEPlain, gin.MIMEHTML, gin.MIMEJSON, gin.MIMEYAML},
				Data:     err,
				HTMLName: strconv.Itoa(http.StatusInternalServerError),
			})
			return
		}

		c.Status(http.StatusOK)
	}
}

func RandomTenant(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = tenantAllowedRunes[rand.Intn(len(tenantAllowedRunes))]
	}
	return string(b)
}
