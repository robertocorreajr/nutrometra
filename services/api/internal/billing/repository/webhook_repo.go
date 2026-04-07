package repository

import (
	"context"
	"fmt"
	"time"

	"nutrometra/api/internal/billing/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// WebhookRepository handles webhook event persistence.
type WebhookRepository struct {
	pool *pgxpool.Pool
}

func NewWebhookRepository(pool *pgxpool.Pool) *WebhookRepository {
	return &WebhookRepository{pool: pool}
}

// InsertIdempotent attempts to insert a webhook event. Returns true if inserted
// (new event), false if already exists (duplicate). Uses ON CONFLICT DO NOTHING
// on the unique provider_event_id constraint.
func (r *WebhookRepository) InsertIdempotent(ctx context.Context, evt domain.WebhookEvent) (bool, error) {
	tag, err := r.pool.Exec(ctx,
		`INSERT INTO billing_webhook_events
		 (id, provider, provider_event_id, event_type, payload_json, status, attempts, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, 1, NOW())
		 ON CONFLICT (provider_event_id) DO NOTHING`,
		evt.ID, evt.Provider, evt.ProviderEventID, evt.EventType, evt.PayloadJSON, string(domain.WebhookStatusPending),
	)
	if err != nil {
		return false, fmt.Errorf("billing: insert_webhook_event: %w", err)
	}
	return tag.RowsAffected() > 0, nil
}

// MarkProcessed sets the event status to processed.
func (r *WebhookRepository) MarkProcessed(ctx context.Context, eventID uuid.UUID) error {
	now := time.Now().UTC()
	_, err := r.pool.Exec(ctx,
		`UPDATE billing_webhook_events SET status = 'processed', processed_at = $1 WHERE id = $2`,
		now, eventID,
	)
	if err != nil {
		return fmt.Errorf("billing: mark_webhook_processed: %w", err)
	}
	return nil
}

// MarkFailed sets the event status to failed with an error message.
func (r *WebhookRepository) MarkFailed(ctx context.Context, eventID uuid.UUID, errMsg string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE billing_webhook_events SET status = 'failed', error_message = $1 WHERE id = $2`,
		errMsg, eventID,
	)
	if err != nil {
		return fmt.Errorf("billing: mark_webhook_failed: %w", err)
	}
	return nil
}
