package observability

import (
	"github.com/go-logr/logr"
	"go.opentelemetry.io/otel"
	"log/slog"
)

func ConfigureOTEL(version string) (func() error, error) {
	otel.SetLogger(logr.New(&slogLogrAdapter{Target: slog.Default()}))
	return func() error { return nil }, nil
}
