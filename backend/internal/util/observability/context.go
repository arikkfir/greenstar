package observability

import (
	"context"
	"fmt"
	"log/slog"
)

type loggerKeyType struct{}

var loggerKey = loggerKeyType{}

func GetLogger(ctx context.Context) *slog.Logger {
	if l := ctx.Value(loggerKey); l == nil {
		panic(fmt.Errorf("no logger associated with context"))
	} else {
		return l.(*slog.Logger)
	}
}

func WithLogger(ctx context.Context, l *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, l)
}
