package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type WebhookEventStatus string

const (
	WebhookStatusPending   WebhookEventStatus = "pending"
	WebhookStatusProcessed WebhookEventStatus = "processed"
	WebhookStatusFailed    WebhookEventStatus = "failed"
	WebhookStatusSkipped   WebhookEventStatus = "skipped"
)

type WebhookEvent struct {
	ID              uuid.UUID
	Provider        string
	ProviderEventID string
	EventType       string
	PayloadJSON     json.RawMessage
	Status          WebhookEventStatus
	ErrorMessage    *string
	Attempts        int
	ProcessedAt     *time.Time
	CreatedAt       time.Time
}
