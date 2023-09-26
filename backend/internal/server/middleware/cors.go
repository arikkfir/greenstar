package middleware

import (
	"github.com/rs/cors"
	"net/http"
	"strings"
	"time"
)

func flattenCORSValues(a []string) []string {
	values := make([]string, 0)
	for _, s := range a {
		tokens := strings.Split(s, ",")
		for _, token := range tokens {
			values = append(values, strings.TrimSpace(token))
		}
	}
	return values
}

func CORSMiddleware(
	allowedOrigins []string,
	allowedMethods []string,
	allowedHeaders []string,
	disableCredentials bool,
	exposeHeaders []string,
	maxAge time.Duration,
	next http.Handler) http.HandlerFunc {

	c := cors.New(cors.Options{
		AllowedOrigins:   flattenCORSValues(allowedOrigins),
		AllowedMethods:   flattenCORSValues(allowedMethods),
		AllowedHeaders:   flattenCORSValues(allowedHeaders),
		ExposedHeaders:   flattenCORSValues(exposeHeaders),
		MaxAge:           int(maxAge.Seconds()),
		AllowCredentials: !disableCredentials,
	})

	handler := func(w http.ResponseWriter, r *http.Request) { next.ServeHTTP(w, r) }

	return func(w http.ResponseWriter, r *http.Request) {
		c.ServeHTTP(w, r, handler)
	}
}
