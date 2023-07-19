package web

import (
	descope "github.com/descope/go-sdk/descope/client"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/rueian/rueidis"
	"net/http"
	"strconv"
	"time"
)

func NewServer(
	port int,
	accessLog bool,
	allowedOrigins []string,
	allowedMethods []string,
	allowedHeaders []string,
	disableCredentials bool,
	exposeHeaders []string,
	maxAge time.Duration,
	redisClient rueidis.Client,
	neo4jDriver neo4j.DriverWithContext,
	descopeClient *descope.DescopeClient,
	graphHandler http.HandlerFunc) *http.Server {

	addr := ":" + strconv.Itoa(port)

	handler := graphHandler
	handler = DescopeAuthenticationMiddleware(descopeClient, handler)
	handler = Neo4jClientMiddleware(neo4jDriver, handler)
	handler = RedisClientMiddleware(redisClient, handler)
	if accessLog {
		handler = AccessLogMiddleware(handler)
	}
	handler = CORSMiddleware(allowedOrigins, allowedMethods, allowedHeaders, disableCredentials, exposeHeaders, maxAge, handler)
	handler = RequestIDMiddleware(handler)

	return &http.Server{Addr: addr, Handler: handler}
}
