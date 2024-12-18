package observability

import (
	"github.com/lmittmann/tint"
	"log/slog"
	"os"
	"time"
)

const (
	LevelTrace slog.Level = -8
)

func ConfigureLogging(disableJSONLogging bool, logLevel slog.Level, version string) {
	// Create log handler
	var handler slog.Handler
	if disableJSONLogging {
		handler = tint.NewHandler(os.Stdout, &tint.Options{
			AddSource:  true,
			Level:      logLevel,
			TimeFormat: time.TimeOnly,
		})
	} else {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
			Level:     logLevel,
		})
	}

	// Configure the default logger
	slog.SetDefault(slog.New(handler).With("version", version))
}
