package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
)

type txKeyType struct{}

var txKey = txKeyType{}

func NewContextWithTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, txKey, tx)
}

func TxFromContext(ctx context.Context) pgx.Tx {
	if v := ctx.Value(txKey); v == nil {
		panic("no transaction in context")
	} else if tx, ok := v.(pgx.Tx); !ok {
		panic(fmt.Errorf("invalid transaction type: %T", v))
	} else {
		return tx
	}
}

func NewTxFromContext(ctx context.Context) (pgx.Tx, error) {
	return TxFromContext(ctx).Begin(ctx)
}
