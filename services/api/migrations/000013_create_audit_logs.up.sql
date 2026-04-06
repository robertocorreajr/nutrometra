CREATE TABLE audit_logs (
    id            UUID PRIMARY KEY,
    tenant_id     UUID,
    actor_user_id UUID,
    actor_scope   VARCHAR(32) NOT NULL CHECK (actor_scope IN ('tenant','backoffice','system')),
    entity_type   VARCHAR(128) NOT NULL,
    entity_id     UUID,
    action        VARCHAR(128) NOT NULL,
    reason        TEXT,
    metadata_json JSONB,
    ip_address    INET,
    user_agent    TEXT,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_audit_logs_tenant     ON audit_logs(tenant_id);
CREATE INDEX idx_audit_logs_entity     ON audit_logs(entity_type, entity_id);
CREATE INDEX idx_audit_logs_actor      ON audit_logs(actor_user_id);
CREATE INDEX idx_audit_logs_created_at        ON audit_logs(created_at DESC);
CREATE INDEX idx_audit_logs_tenant_created_at ON audit_logs(tenant_id, created_at DESC);
