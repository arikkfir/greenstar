package observability

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel/trace"
	"runtime"
	"time"
)

func NamedTrace(ctx context.Context, spanName string, kind trace.SpanKind, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	ctxWithSpan, span := Tracer.Start(
		ctx,
		spanName,
		append(opts, trace.WithSpanKind(kind), trace.WithTimestamp(time.Now()))...,
	)
	return ctxWithSpan, span
}

func Trace(ctx context.Context, kind trace.SpanKind) (context.Context, trace.Span) {
	spanName := "<unknown>"

	var pc, callerFile, callerLine, ok = runtime.Caller(1)
	if ok {
		fn := runtime.FuncForPC(pc)
		spanName = fmt.Sprintf("%s [%s:%d]", fn.Name(), callerFile, callerLine)
	}

	return NamedTrace(ctx, spanName, kind)
}
