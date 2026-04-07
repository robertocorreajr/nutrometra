CREATE TABLE clinical_attachments (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES tenants(id),
    patient_id      UUID NOT NULL REFERENCES patients(id),
    professional_id UUID NOT NULL REFERENCES professionals(id),
    file_name       TEXT NOT NULL,
    file_type       TEXT NOT NULL,          -- MIME type
    file_size_bytes BIGINT NOT NULL,
    storage_key     TEXT NOT NULL,           -- object storage key
    category        TEXT NOT NULL DEFAULT 'general'
                    CHECK (category IN ('general','exam','lab_result','prescription','photo','other')),
    description     TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_clinical_attachments_tenant  ON clinical_attachments(tenant_id);
CREATE INDEX idx_clinical_attachments_patient ON clinical_attachments(patient_id);
