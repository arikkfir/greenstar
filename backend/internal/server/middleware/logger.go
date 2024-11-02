package middleware

import (
	"github.com/arikkfir/greenstar/backend/internal/server/util"
	"net/http"
)

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, util.RequestWithLogger(r, util.Logger(r.Context())))
	})
}

func PathVariablesLoggerMiddleware(next http.Handler, pathVariableNames ...string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := util.Logger(r.Context())

		var args []any
		for _, name := range pathVariableNames {
			args = append(args, "req:path:"+name, r.PathValue(name))
		}

		if len(args) > 0 {
			next.ServeHTTP(w, util.RequestWithLogger(r, l.With(args...)))
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
