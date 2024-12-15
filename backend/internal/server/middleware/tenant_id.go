package middleware

import (
	"context"
	"fmt"
	"github.com/arikkfir/greenstar/backend/internal/server/util"
	"net/http"
)

const TenantIDHeaderName = "X-GreenSTAR-Tenant-ID"

type tidKeyType struct{}

var tidKey = &tidKeyType{}

func GetTenantID(ctx context.Context) string {
	v := ctx.Value(tidKey)
	if v == nil {
		return ""
	} else if tid, ok := v.(string); ok {
		return tid
	} else {
		panic(fmt.Sprintf("unexpected tenant ID type '%T' encountered: %+v", v, v))
	}
}

func WithTenantID(ctx context.Context, tenantID string) context.Context {
	return context.WithValue(ctx, tidKey, tenantID)
}

func TenantIDMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tid := r.Header.Get(TenantIDHeaderName)
		if tid != "" {
			ctxWithTID := WithTenantID(r.Context(), tid)
			reqWithTID := r.WithContext(ctxWithTID)

			logger := util.Logger(r.Context()).With("tenantID", tid)
			reqWithTIDLogger := util.RequestWithLogger(reqWithTID, logger)

			next.ServeHTTP(w, reqWithTIDLogger)
		} else {
			next.ServeHTTP(w, r)
		}
	}
}
