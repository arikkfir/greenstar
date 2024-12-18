package util

import (
	"context"
	"github.com/arikkfir/greenstar/backend/internal/util/auth"
	"github.com/arikkfir/greenstar/backend/internal/util/observability"
	"github.com/arikkfir/greenstar/backend/internal/util/tenant"
)

func VerifyPermissions(ctx context.Context, permissions ...auth.Permission) error {
	l := observability.GetLogger(ctx)
	authToken := auth.GetToken(ctx)
	tenantID := tenant.GetTenantID(ctx)
	if !authToken.IsPermittedGlobally(permissions...) {
		if tenantID != "" {
			if !authToken.IsPermittedForTenant(tenantID, permissions...) {
				l.WarnContext(ctx, "Access denied", "permissions", permissions)
				return ErrForbidden
			}
		} else {
			l.WarnContext(ctx, "Access denied", "permissions", permissions)
			return ErrForbidden
		}
	}
	return nil
}
