package middleware

import (
	authutil "github.com/arikkfir/greenstar/backend/internal/util/auth"
	"net/http"
)

func TokenMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctxWithToken := authutil.WithToken(r.Context(), &authutil.Token{})
		next.ServeHTTP(w, r.WithContext(ctxWithToken))
	}
}
