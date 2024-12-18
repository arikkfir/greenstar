package tenant

import (
	"context"
	"fmt"
)

type tidKeyType struct{}

var tidKey = &tidKeyType{}

func GetTenantID(ctx context.Context) string {
	v := ctx.Value(tidKey)
	if v == nil {
		return ""
	} else if tid, ok := v.(string); ok {
		return tid
	} else {
		panic(fmt.Sprintf("unexpected tenant ID type '%T' encountered: %+v", v, v))
	}
}

func WithTenantID(ctx context.Context, tenantID string) context.Context {
	return context.WithValue(ctx, tidKey, tenantID)
}
