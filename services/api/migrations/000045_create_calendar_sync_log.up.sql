CREATE TABLE calendar_sync_log (
    id                     UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id              UUID NOT NULL REFERENCES tenants(id),
    calendar_connection_id UUID NOT NULL REFERENCES calendar_connections(id),
    entity_type            VARCHAR(64) NOT NULL,
    entity_id              UUID NOT NULL,
    provider_event_id      VARCHAR(255),
    operation              VARCHAR(32) NOT NULL CHECK (operation IN ('create','update','delete')),
    status                 VARCHAR(32) NOT NULL DEFAULT 'pending'
                           CHECK (status IN ('pending','synced','failed')),
    error_message          TEXT,
    created_at             TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_calendar_sync_entity ON calendar_sync_log(entity_type, entity_id);
