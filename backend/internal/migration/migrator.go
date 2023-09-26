package migration

import (
	"context"
	"embed"
	"fmt"
	"github.com/arikkfir/greenstar/backend/internal/util/lang"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	_ "github.com/lib/pq"
	"log/slog"
)

var (
	//go:embed schema/*.sql
	fs embed.FS
)

func Migrate(ctx context.Context, pool *pgxpool.Pool) (result error) {
	slog.Default().InfoContext(ctx, "Verifying and potentially upgrading database schema")

	sourceDriver, err := iofs.New(fs, "schema")
	if err != nil {
		return fmt.Errorf("failed creating schema migration source driver: %w", err)
	}

	db := stdlib.OpenDBFromPool(pool)
	defer db.Close()

	databaseDriver, err := pgx.WithInstance(db, &pgx.Config{})
	if err != nil {
		return fmt.Errorf("failed creating schema migrator database driver: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", sourceDriver, "pgx", databaseDriver)
	if err != nil {
		return fmt.Errorf("failed creating schema migrator: %w", err)
	}
	defer m.Close()

	if err := m.Up(); lang.IgnoreErrorOfType(err, migrate.ErrNoChange) != nil {
		return fmt.Errorf("failed migrating database schema: %w", err)
	}

	return nil
}
