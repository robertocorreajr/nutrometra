package worker

import (
	"context"
	"log/slog"
	"sync"

	"github.com/google/uuid"
)

// Job represents a unit of work.
type Job struct {
	ID      uuid.UUID
	Type    string
	Payload any
}

// HandlerFunc processes a job.
type HandlerFunc func(ctx context.Context, job Job) error

// Worker processes background jobs via a buffered channel.
type Worker struct {
	handlers map[string]HandlerFunc
	queue    chan Job
	wg       sync.WaitGroup
	ctx      context.Context
	cancel   context.CancelFunc
}

// New creates a Worker with the given queue buffer size.
func New(bufferSize int) *Worker {
	return &Worker{
		handlers: make(map[string]HandlerFunc),
		queue:    make(chan Job, bufferSize),
	}
}

// Register adds a handler for a job type.
func (w *Worker) Register(jobType string, handler HandlerFunc) {
	w.handlers[jobType] = handler
}

// Enqueue adds a job to the queue. Non-blocking — drops if full (logs warning).
func (w *Worker) Enqueue(job Job) {
	select {
	case w.queue <- job:
		slog.Info("worker: job enqueued", "job_id", job.ID, "type", job.Type)
	default:
		slog.Warn("worker: queue full, job dropped", "job_id", job.ID, "type", job.Type)
	}
}

// Start begins processing jobs. Call this before the HTTP server starts.
func (w *Worker) Start(ctx context.Context) {
	w.ctx, w.cancel = context.WithCancel(ctx)
	w.wg.Add(1)
	go w.run()
	slog.Info("worker: started")
}

// Stop signals the worker to stop and waits for the current job to finish.
func (w *Worker) Stop() {
	w.cancel()
	close(w.queue)
	w.wg.Wait()
	slog.Info("worker: stopped")
}

// run processes jobs from the queue until stopped.
func (w *Worker) run() {
	defer w.wg.Done()
	for job := range w.queue {
		handler, ok := w.handlers[job.Type]
		if !ok {
			slog.Error("worker: no handler for job type", "type", job.Type)
			continue
		}
		if err := handler(w.ctx, job); err != nil {
			slog.Error("worker: job failed", "job_id", job.ID, "type", job.Type, "error", err)
		} else {
			slog.Info("worker: job completed", "job_id", job.ID, "type", job.Type)
		}
	}
}
