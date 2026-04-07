package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"nutrometra/api/internal/platform/config"
)

// Pool is an exported alias for pgxpool.Pool.
type Pool = pgxpool.Pool

// New creates and validates a PostgreSQL connection pool.
func New(ctx context.Context, cfg config.PostgresConfig) (*Pool, error) {
	pool, err := pgxpool.New(ctx, cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("db: failed to create pool: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("db: ping failed: %w", err)
	}
	return pool, nil
}

// TxFunc is a function executed inside a transaction.
type TxFunc func(ctx context.Context, tx pgx.Tx) error

// RunInTx executes fn inside a transaction. Rolls back on error.
func RunInTx(ctx context.Context, pool *Pool, fn TxFunc) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("db: begin tx: %w", err)
	}
	if err := fn(ctx, tx); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("db: commit tx: %w", err)
	}
	return nil
}
