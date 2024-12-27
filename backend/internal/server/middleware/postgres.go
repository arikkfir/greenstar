package middleware

import (
	"context"
	"github.com/arikkfir/greenstar/backend/internal/server/util"
	"github.com/arikkfir/greenstar/backend/internal/util/db"
	"github.com/arikkfir/greenstar/backend/internal/util/observability"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
)

func PostgresMiddleware(pool *pgxpool.Pool, mode pgx.TxAccessMode, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := observability.GetLogger(r.Context())

		ww := &responseStatusCodeRecorder{ResponseWriter: w, statusCode: http.StatusOK}

		txOptions := pgx.TxOptions{IsoLevel: pgx.Serializable, AccessMode: mode, DeferrableMode: pgx.NotDeferrable}
		tx, err := pool.BeginTx(r.Context(), txOptions)
		if err != nil {
			l.ErrorContext(r.Context(), "Failed to begin transaction", "err", err, "txMode", mode)
			util.ServeError(w, r, err)
			return
		}
		defer tx.Rollback(context.Background()) // fresh context, to ensure rollback succeeds for timed out requests too

		next.ServeHTTP(ww, r.WithContext(db.WithTransaction(r.Context(), tx)))

		if ww.statusCode >= 200 && ww.statusCode < 400 {
			// fresh context, to ensure rollback succeeds for timed out requests too
			if err := tx.Commit(context.Background()); err != nil {
				l.ErrorContext(r.Context(), "Failed to commit transaction", "err", err, "txMode", mode)
			}
		}
	})
}
