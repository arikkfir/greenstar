package observability

import (
	"context"
	"fmt"
	"github.com/lmittmann/tint"
	"log/slog"
	"os"
	"time"

	"github.com/go-logr/logr"
	"go.opentelemetry.io/otel"

	"github.com/arikkfir/greenstar/backend/internal/util/version"
)

const (
	LevelTrace slog.Level = -8
)

// LoggingHook is a command hook for configuring logging.
type LoggingHook struct {
	DisableJSONLogging bool   `desc:"Disable JSON logging."`
	LogLevel           string `required:"true" desc:"Log level, must be one of: trace,debug,info,warn,error,fatal,panic"`
}

func (c *LoggingHook) PreRun(_ context.Context) error {

	// Parse level
	level := slog.LevelInfo
	switch c.LogLevel {
	case "trace":
		level = LevelTrace
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn", "warning":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		return fmt.Errorf("invalid log level: %s", c.LogLevel)
	}

	// Create log handler
	var handler slog.Handler
	if c.DisableJSONLogging {
		handler = tint.NewHandler(os.Stdout, &tint.Options{
			AddSource:  true,
			Level:      LevelTrace,
			TimeFormat: time.TimeOnly,
		})
	} else {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: false,
			Level:     level,
		})
	}

	// Create the default logger
	slog.SetDefault(slog.New(handler).With("version", version.Version))

	// Redirect logs sent to OTEL to slog
	otel.SetLogger(logr.New(&SlogLogrAdapter{l: slog.Default()}))

	return nil
}
