package queue

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const defaultMaxAttempts = 3

// PostgresQueue implements Queue backed by the background_jobs table.
type PostgresQueue struct {
	pool         *pgxpool.Pool
	handlers     map[string]HandlerFunc
	pollInterval time.Duration
	cancel       context.CancelFunc
	done         chan struct{}
}

// NewPostgresQueue creates a new PostgresQueue.
func NewPostgresQueue(pool *pgxpool.Pool) *PostgresQueue {
	return &PostgresQueue{
		pool:         pool,
		handlers:     make(map[string]HandlerFunc),
		pollInterval: 1 * time.Second,
	}
}

// Enqueue inserts a new job into the background_jobs table.
func (q *PostgresQueue) Enqueue(ctx context.Context, job Job) error {
	maxAttempts := job.MaxAttempts
	if maxAttempts <= 0 {
		maxAttempts = defaultMaxAttempts
	}
	scheduledAt := job.ScheduleAt
	if scheduledAt.IsZero() {
		scheduledAt = time.Now().UTC()
	}
	payload := job.Payload
	if payload == nil {
		payload = json.RawMessage("{}")
	}

	_, err := q.pool.Exec(ctx,
		`INSERT INTO background_jobs (id, tenant_id, job_type, payload_json, max_attempts, scheduled_at)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		uuid.New(), job.TenantID, job.JobType, payload, maxAttempts, scheduledAt,
	)
	if err != nil {
		return fmt.Errorf("queue: enqueue: %w", err)
	}
	return nil
}

// RegisterHandler registers a handler function for a given job type.
func (q *PostgresQueue) RegisterHandler(jobType string, handler HandlerFunc) {
	q.handlers[jobType] = handler
}

// dequeue fetches the next eligible job using SELECT FOR UPDATE SKIP LOCKED.
func (q *PostgresQueue) dequeue(ctx context.Context) (*JobRecord, error) {
	tx, err := q.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("queue: begin tx: %w", err)
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	var rec JobRecord
	var payloadBytes []byte
	err = tx.QueryRow(ctx,
		`SELECT id, tenant_id, job_type, payload_json, status, attempts, max_attempts,
		        last_error, scheduled_at, started_at, completed_at, created_at
		 FROM background_jobs
		 WHERE status IN ('pending','failed')
		   AND scheduled_at <= NOW()
		   AND attempts < max_attempts
		 ORDER BY scheduled_at ASC
		 LIMIT 1
		 FOR UPDATE SKIP LOCKED`,
	).Scan(
		&rec.ID, &rec.TenantID, &rec.JobType, &payloadBytes, &rec.Status,
		&rec.Attempts, &rec.MaxAttempts, &rec.LastError, &rec.ScheduledAt,
		&rec.StartedAt, &rec.CompletedAt, &rec.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("queue: dequeue scan: %w", err)
	}
	rec.PayloadJSON = payloadBytes

	_, err = tx.Exec(ctx,
		`UPDATE background_jobs SET status = 'processing', started_at = NOW() WHERE id = $1`,
		rec.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("queue: update processing: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("queue: commit dequeue: %w", err)
	}
	return &rec, nil
}

// ack marks a job as completed.
func (q *PostgresQueue) ack(ctx context.Context, jobID uuid.UUID) error {
	_, err := q.pool.Exec(ctx,
		`UPDATE background_jobs SET status = 'completed', completed_at = NOW() WHERE id = $1`,
		jobID,
	)
	return err
}

// recoverStaleJobs moves jobs stuck in 'processing' status (e.g. from a crash)
// back to 'failed' so they can be retried.
func (q *PostgresQueue) recoverStaleJobs(ctx context.Context) {
	const staleTimeout = 5 * time.Minute
	cutoff := time.Now().UTC().Add(-staleTimeout)

	tag, err := q.pool.Exec(ctx,
		`UPDATE background_jobs
		 SET status = 'failed', last_error = 'recovered: job was stuck in processing'
		 WHERE status = 'processing' AND started_at < $1`,
		cutoff,
	)
	if err != nil {
		slog.Error("queue: recover stale jobs failed", "error", err)
		return
	}
	if tag.RowsAffected() > 0 {
		slog.Warn("queue: recovered stale jobs", "count", tag.RowsAffected())
	}
}

// nack marks a job as failed with retry or dead.
func (q *PostgresQueue) nack(ctx context.Context, jobID uuid.UUID, attempts, maxAttempts int, jobErr error) error {
	newAttempts := attempts + 1
	errMsg := ""
	if jobErr != nil {
		errMsg = jobErr.Error()
	}

	if newAttempts >= maxAttempts {
		_, err := q.pool.Exec(ctx,
			`UPDATE background_jobs SET status = 'dead', attempts = $2, last_error = $3 WHERE id = $1`,
			jobID, newAttempts, errMsg,
		)
		return err
	}

	// Exponential backoff: 4^attempts seconds (4s, 16s, 64s, ...)
	backoff := time.Duration(math.Pow(4, float64(newAttempts))) * time.Second
	nextRun := time.Now().UTC().Add(backoff)

	_, err := q.pool.Exec(ctx,
		`UPDATE background_jobs SET status = 'failed', attempts = $2, last_error = $3, scheduled_at = $4 WHERE id = $1`,
		jobID, newAttempts, errMsg, nextRun,
	)
	return err
}
