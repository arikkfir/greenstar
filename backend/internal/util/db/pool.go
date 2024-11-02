package db

import (
	"context"
	"fmt"
	pgxdecimal "github.com/jackc/pgx-shopspring-decimal"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
)

func NewPostgreSQLPool(ctx context.Context) (*pgxpool.Pool, error) {
	slog.Default().InfoContext(ctx, "Creating Postgres connection pool")

	// PostgreSQL (note that standard environment variables are used to configure connection details to Postgres)
	// See: https://www.postgresql.org/docs/current/libpq-envars.html
	pgxpoolConfig, err := pgxpool.ParseConfig("")
	if err != nil {
		return nil, fmt.Errorf("failed parsing pgxpool config: %w", err)
	}
	pgxpoolConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		pgxdecimal.Register(conn.TypeMap())
		return nil
	}
	pgPool, err := pgxpool.NewWithConfig(ctx, pgxpoolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed creating pgxpool: %w", err)
	}

	return pgPool, nil
}
