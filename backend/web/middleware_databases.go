package web

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/rueian/rueidis"
	"net/http"
)

const (
	contextRedisKey = "$$$_redis_client_$$$"
	contextNeo4jKey = "$$$_neo4j_driver_$$$"
)

func RedisClientMiddleware(redisClient rueidis.Client, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(context.WithValue(r.Context(), contextRedisKey, redisClient))
		next(w, r)
	}
}

func Neo4jClientMiddleware(neo4jDriver neo4j.DriverWithContext, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(context.WithValue(r.Context(), contextNeo4jKey, neo4jDriver))
		next(w, r)
	}
}
