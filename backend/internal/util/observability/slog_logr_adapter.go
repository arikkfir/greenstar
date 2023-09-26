package observability

import (
	"context"
	"fmt"
	"github.com/go-logr/logr"
	"log/slog"
	"runtime"
	"time"
)

const (
	logrSkipCount = 4
)

type SlogLogrAdapter struct {
	l *slog.Logger
}

func (a *SlogLogrAdapter) Init(_ logr.RuntimeInfo) {}

func (a *SlogLogrAdapter) Enabled(level int) bool {
	return level <= 4
	// Defer decision on whether to log or not to log
	// return true
}

func (a *SlogLogrAdapter) Info(level int, msg string, keysAndValues ...interface{}) {
	if level < 4 {
		var pcs [1]uintptr
		runtime.Callers(logrSkipCount, pcs[:]) // skip [Callers, Info]
		r := slog.NewRecord(time.Now(), slog.LevelInfo, fmt.Sprintf(msg, keysAndValues...), pcs[0])
		_ = a.l.Handler().Handle(context.Background(), r)
	} else if level < 8 {
		var pcs [1]uintptr
		runtime.Callers(logrSkipCount, pcs[:]) // skip [Callers, Info]
		r := slog.NewRecord(time.Now(), slog.LevelDebug, fmt.Sprintf(msg, keysAndValues...), pcs[0])
		_ = a.l.Handler().Handle(context.Background(), r)
	} else {
		var pcs [1]uintptr
		runtime.Callers(logrSkipCount, pcs[:]) // skip [Callers, Info]
		r := slog.NewRecord(time.Now(), LevelTrace, fmt.Sprintf(msg, keysAndValues...), pcs[0])
		_ = a.l.Handler().Handle(context.Background(), r)
	}
}

func (a *SlogLogrAdapter) Error(err error, msg string, keysAndValues ...interface{}) {
	var pcs [1]uintptr
	runtime.Callers(logrSkipCount, pcs[:]) // skip [Callers, Error]
	r := slog.NewRecord(time.Now(), slog.LevelError, fmt.Sprintf(msg, append([]any{"err", err}, keysAndValues...)...), pcs[0])
	_ = a.l.Handler().Handle(context.Background(), r)
}

func (a *SlogLogrAdapter) WithValues(keysAndValues ...interface{}) logr.LogSink {
	return &SlogLogrAdapter{l: a.l.With(keysAndValues...)}
}

func (a *SlogLogrAdapter) WithName(_ string) logr.LogSink {
	return a
}
