package middleware

import (
	"github.com/arikkfir/greenstar/backend/internal/util/observability"
	"github.com/google/uuid"
	"net/http"
)

const RequestIDHeaderName = "X-Request-ID"

func RequestIDMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rid := r.Header.Get(RequestIDHeaderName)
		if rid == "" {
			rid = uuid.New().String()
			r.Header.Set(RequestIDHeaderName, rid)
			w.Header().Set(RequestIDHeaderName, rid)
		}

		next.ServeHTTP(w, r.WithContext(observability.WithRequestID(r.Context(), rid)))
	}
}
