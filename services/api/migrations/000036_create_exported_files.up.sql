CREATE TABLE exported_files (
    id                   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id            UUID NOT NULL REFERENCES tenants(id),
    related_entity_type  TEXT NOT NULL,
    related_entity_id    UUID NOT NULL,
    export_type          TEXT NOT NULL DEFAULT 'pdf'
                         CHECK (export_type IN ('pdf', 'print_job')),
    file_key             TEXT,
    status               TEXT NOT NULL DEFAULT 'pending'
                         CHECK (status IN ('pending', 'processing', 'completed', 'failed')),
    requested_by_user_id UUID NOT NULL REFERENCES users(id),
    created_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at         TIMESTAMPTZ,
    failure_reason       TEXT
);
CREATE INDEX idx_exported_files_entity ON exported_files(related_entity_type, related_entity_id);
CREATE INDEX idx_exported_files_status ON exported_files(tenant_id, status);
