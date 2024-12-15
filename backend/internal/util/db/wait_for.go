package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"time"
)

const (
	PingInterval = 3 * time.Second
)

func WaitFor(ctx context.Context, pgPool *pgxpool.Pool) error {
	timer := time.NewTimer(2 * time.Minute)
	defer timer.Stop()

	ticker := time.NewTicker(PingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
			return fmt.Errorf("timed out waiting for Postgres connection pool to become available")
		case <-ticker.C:
			if err := pgPool.Ping(ctx); err != nil {
				slog.InfoContext(ctx, "PostgreSQL not yet available", "err", err)
			} else {
				return nil
			}
		}
	}
}
