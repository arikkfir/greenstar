package middleware

import (
	"github.com/arikkfir/greenstar/backend/internal/util/tenant"
	"net/http"
)

const TenantIDHeaderName = "X-GreenSTAR-Tenant-ID"

func TenantIDMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tid := r.Header.Get(TenantIDHeaderName)
		if tid != "" {
			next.ServeHTTP(w, r.WithContext(tenant.WithTenantID(r.Context(), tid)))
		} else {
			next.ServeHTTP(w, r)
		}
	}
}
