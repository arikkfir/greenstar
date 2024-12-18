package observability

import (
	"context"
	"fmt"
)

type ridKeyType struct{}

var ridKey = &ridKeyType{}

func GetRequestID(ctx context.Context) string {
	v := ctx.Value(ridKey)
	if v == nil {
		return ""
	} else if rid, ok := v.(string); ok {
		return rid
	} else {
		panic(fmt.Sprintf("unexpected request ID type '%T' encountered: %+v", v, v))
	}
}

func WithRequestID(ctx context.Context, rid string) context.Context {
	return context.WithValue(ctx, ridKey, rid)
}
