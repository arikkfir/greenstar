package web

import (
	descope "github.com/descope/go-sdk/descope/client"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/redis/rueidis"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func flatten(a []string) []string {
	values := make([]string, 0)
	for _, s := range a {
		tokens := strings.Split(s, ",")
		for _, token := range tokens {
			values = append(values, strings.TrimSpace(token))
		}
	}
	return values
}

func NewServer(
	port int,
	accessLog bool,
	accessLogExcludedHeaders []string,
	accessLogExcludeRemoteAddr bool,
	allowedOrigins []string,
	allowedMethods []string,
	allowedHeaders []string,
	disableCredentials bool,
	exposeHeaders []string,
	maxAge time.Duration,
	redisClient rueidis.Client,
	neo4jDriver neo4j.DriverWithContext,
	descopeClient *descope.DescopeClient,
	graphHandler http.HandlerFunc,
) *http.Server {

	allowedOrigins = flatten(allowedOrigins)
	allowedMethods = flatten(allowedMethods)
	allowedHeaders = flatten(allowedHeaders)
	exposeHeaders = flatten(exposeHeaders)

	addr := ":" + strconv.Itoa(port)

	handler := graphHandler
	handler = DescopeAuthenticationMiddleware(descopeClient, handler)
	handler = Neo4jClientMiddleware(neo4jDriver, handler)
	handler = RedisClientMiddleware(redisClient, handler)
	if accessLog {
		handler = AccessLogMiddleware(accessLogExcludeRemoteAddr, accessLogExcludedHeaders, handler)
	}
	handler = CORSMiddleware(allowedOrigins, allowedMethods, allowedHeaders, disableCredentials, exposeHeaders, maxAge, handler)
	handler = RequestIDMiddleware(handler)

	return &http.Server{Addr: addr, Handler: handler}
}
