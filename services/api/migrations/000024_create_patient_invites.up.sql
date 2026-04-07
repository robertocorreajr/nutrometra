CREATE TABLE patient_invites (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES tenants(id),
    patient_id      UUID NOT NULL REFERENCES patients(id),
    invited_by      UUID NOT NULL REFERENCES users(id),
    code_hash       TEXT NOT NULL,                                -- SHA-256 of plaintext code
    max_uses        INT NOT NULL DEFAULT 1,
    used_count      INT NOT NULL DEFAULT 0,
    expires_at      TIMESTAMPTZ NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_patient_invites_tenant  ON patient_invites(tenant_id);
CREATE INDEX idx_patient_invites_patient ON patient_invites(patient_id);
CREATE INDEX idx_patient_invites_hash    ON patient_invites(code_hash);
