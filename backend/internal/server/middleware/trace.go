package middleware

import (
	"fmt"
	"github.com/arikkfir/greenstar/backend/internal/util/observability"
	"github.com/arikkfir/greenstar/backend/internal/util/tenant"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"net/http"
)

func TraceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, span := observability.NamedTrace(
			r.Context(),
			fmt.Sprintf("%s %s%s", r.Method, r.URL.Host, r.URL.Path),
			trace.SpanKindServer,
			trace.WithAttributes(
				attribute.String("rid", observability.GetRequestID(r.Context())),
				attribute.String("tid", tenant.GetTenantID(r.Context())),
				attribute.String("proto", r.Proto),
				attribute.String("host", r.Host),
				attribute.String("method", r.Method),
				attribute.String("remoteAddr", r.RemoteAddr),
				attribute.String("requestURI", r.RequestURI),
			),
		)
		defer span.End()

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
