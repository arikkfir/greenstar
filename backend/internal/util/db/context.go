package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
)

type txKeyType struct{}

var txKey = txKeyType{}

func WithTransaction(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, txKey, tx)
}

func GetTransaction(ctx context.Context) pgx.Tx {
	if v := ctx.Value(txKey); v == nil {
		panic("no transaction in context")
	} else if tx, ok := v.(pgx.Tx); !ok {
		panic(fmt.Errorf("invalid transaction type: %T", v))
	} else {
		return tx
	}
}
