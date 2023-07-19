package web

import (
	"github.com/rs/cors"
	"net/http"
	"time"
)

func CORSMiddleware(
	allowedOrigins []string,
	allowedMethods []string,
	allowedHeaders []string,
	disableCredentials bool,
	exposeHeaders []string,
	maxAge time.Duration,
	next http.HandlerFunc) http.HandlerFunc {

	c := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   allowedMethods,
		AllowedHeaders:   allowedHeaders,
		ExposedHeaders:   exposeHeaders,
		MaxAge:           int(maxAge.Seconds()),
		AllowCredentials: !disableCredentials,
	})

	return func(w http.ResponseWriter, r *http.Request) {
		c.ServeHTTP(w, r, next)
	}
}
