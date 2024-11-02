package util

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
)

type loggerKeyType struct{}

var loggerKey = loggerKeyType{}

func Logger(ctx context.Context) *slog.Logger {
	if l := ctx.Value(loggerKey); l == nil {
		panic(fmt.Errorf("no logger associated with context"))
	} else {
		return l.(*slog.Logger)
	}
}

func ContextWithLogger(ctx context.Context, l *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, l)
}

func RequestWithLogger(r *http.Request, l *slog.Logger) *http.Request {
	contextWithLoggerValue := context.WithValue(r.Context(), loggerKey, l)
	return r.WithContext(contextWithLoggerValue)
}
