package queue

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

// Start begins the processing loop in a goroutine.
func (q *PostgresQueue) Start(ctx context.Context) {
	procCtx, cancel := context.WithCancel(ctx)
	q.cancel = cancel
	q.done = make(chan struct{})

	go q.processLoop(procCtx)
	slog.Info("queue: processor started", "poll_interval", q.pollInterval)
}

// Stop signals the processor to stop and waits for it to finish.
func (q *PostgresQueue) Stop() {
	if q.cancel != nil {
		q.cancel()
	}
	if q.done != nil {
		<-q.done
	}
	slog.Info("queue: processor stopped")
}

func (q *PostgresQueue) processLoop(ctx context.Context) {
	defer close(q.done)
	ticker := time.NewTicker(q.pollInterval)
	defer ticker.Stop()

	// Recover stale processing jobs on startup
	q.recoverStaleJobs(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			q.processOne(ctx)
		}
	}
}

func (q *PostgresQueue) processOne(ctx context.Context) {
	rec, err := q.dequeue(ctx)
	if err != nil {
		slog.Error("queue: dequeue error", "error", err)
		return
	}
	if rec == nil {
		return // no jobs available
	}

	handler, ok := q.handlers[rec.JobType]
	if !ok {
		slog.Error("queue: no handler registered", "job_type", rec.JobType, "job_id", rec.ID)
		if nackErr := q.nack(ctx, rec.ID, rec.Attempts, rec.MaxAttempts, fmt.Errorf("no handler for job type: %s", rec.JobType)); nackErr != nil {
			slog.Error("queue: nack failed", "job_id", rec.ID, "error", nackErr)
		}
		return
	}

	slog.Info("queue: processing job", "job_id", rec.ID, "type", rec.JobType, "attempt", rec.Attempts+1)

	if err := handler(ctx, rec); err != nil {
		slog.Error("queue: job failed", "job_id", rec.ID, "type", rec.JobType, "error", err)
		if nackErr := q.nack(ctx, rec.ID, rec.Attempts, rec.MaxAttempts, err); nackErr != nil {
			slog.Error("queue: nack failed", "job_id", rec.ID, "error", nackErr)
		}
		return
	}

	if err := q.ack(ctx, rec.ID); err != nil {
		slog.Error("queue: ack failed", "job_id", rec.ID, "error", err)
	} else {
		slog.Info("queue: job completed", "job_id", rec.ID, "type", rec.JobType)
	}
}
