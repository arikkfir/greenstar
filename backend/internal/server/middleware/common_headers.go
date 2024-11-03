package middleware

import (
	"github.com/arikkfir/greenstar/backend/internal/util/version"
	"net/http"
)

func CommonHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", "Greenstar")
		w.Header().Set("X-Greenstar-Version", version.Version)
		next.ServeHTTP(w, r)
	})
}
