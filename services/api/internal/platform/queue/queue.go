package queue

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Job is the input for enqueueing a new background job.
type Job struct {
	TenantID    *uuid.UUID
	JobType     string
	Payload     json.RawMessage
	MaxAttempts int
	ScheduleAt  time.Time
}

// JobRecord represents a job row from the database.
type JobRecord struct {
	ID          uuid.UUID
	TenantID    *uuid.UUID
	JobType     string
	PayloadJSON json.RawMessage
	Status      string
	Attempts    int
	MaxAttempts int
	LastError   *string
	ScheduledAt time.Time
	StartedAt   *time.Time
	CompletedAt *time.Time
	CreatedAt   time.Time
}

// HandlerFunc processes a dequeued job.
type HandlerFunc func(ctx context.Context, job *JobRecord) error

// Queue defines the background job queue interface.
type Queue interface {
	Enqueue(ctx context.Context, job Job) error
	RegisterHandler(jobType string, handler HandlerFunc)
	Start(ctx context.Context)
	Stop()
}
