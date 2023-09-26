package auth

import (
	"fmt"
	"github.com/descope/go-sdk/descope/logger"
	"log/slog"
)

type DescopeLogger struct {
	LogLevel logger.LogLevel
	Logger   *slog.Logger
}

func (l *DescopeLogger) Print(v ...interface{}) {
	switch l.LogLevel {
	case logger.LogNone:
		// No-op
	case logger.LogDebugLevel:
		l.Logger.Debug(fmt.Sprint(v...))
	default:
		l.Logger.Info(fmt.Sprint(v...))
	}
}
