package main

import (
	"context"
	"github.com/arikkfir/greenstar/backend/internal/util/observability"
	"github.com/arikkfir/greenstar/backend/internal/util/tenant"
	"log/slog"
)

type ServerSLogHandler struct {
	slog.Handler
}

func (h ServerSLogHandler) Handle(ctx context.Context, r slog.Record) error {
	r.AddAttrs(slog.String("tenantID", tenant.GetTenantID(ctx)))
	r.AddAttrs(slog.String("rid", observability.GetRequestID(ctx)))
	return h.Handler.Handle(ctx, r)
}
