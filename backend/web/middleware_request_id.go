package web

import (
	"github.com/google/uuid"
	"net/http"
)

const RequestIDHeaderName = "X-Request-ID"

func RequestIDMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rid := r.Header.Get(RequestIDHeaderName)
		if rid == "" {
			rid = uuid.New().String()
			r.Header.Set(RequestIDHeaderName, rid)
			w.Header().Set(RequestIDHeaderName, rid)
		}
		next(w, r)
	}
}
