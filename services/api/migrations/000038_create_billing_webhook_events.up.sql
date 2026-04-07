CREATE TABLE billing_webhook_events (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    provider          VARCHAR(32) NOT NULL DEFAULT 'stripe',
    provider_event_id VARCHAR(255) NOT NULL UNIQUE,
    event_type        VARCHAR(128) NOT NULL,
    payload_json      JSONB NOT NULL,
    status            VARCHAR(32) NOT NULL DEFAULT 'pending'
                      CHECK (status IN ('pending','processed','failed','skipped')),
    error_message     TEXT,
    attempts          INTEGER NOT NULL DEFAULT 0,
    processed_at      TIMESTAMPTZ,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_webhook_events_status ON billing_webhook_events(status)
    WHERE status IN ('pending','failed');
