package util

import (
	"fmt"
	"github.com/descope/go-sdk/descope/logger"
	"github.com/rs/zerolog"
)

type DescopeLogger struct {
	LogLevel logger.LogLevel
	Logger   *zerolog.Logger
}

func (l *DescopeLogger) Print(v ...interface{}) {
	var e *zerolog.Event
	switch l.LogLevel {
	case logger.LogNone:
		e = l.Logger.Trace()
	case logger.LogDebugLevel:
		e = l.Logger.Debug()
	default:
		e = l.Logger.Info()
	}
	e.Msg(fmt.Sprint(v...))
}
