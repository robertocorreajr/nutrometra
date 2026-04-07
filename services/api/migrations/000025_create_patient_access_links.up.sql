CREATE TABLE patient_access_links (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES tenants(id),
    patient_id      UUID NOT NULL REFERENCES patients(id),
    user_id         UUID NOT NULL REFERENCES users(id),
    invite_id       UUID NOT NULL REFERENCES patient_invites(id),
    active          BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE (patient_id, user_id)
);

CREATE INDEX idx_patient_access_links_tenant  ON patient_access_links(tenant_id);
CREATE INDEX idx_patient_access_links_patient ON patient_access_links(patient_id);
CREATE INDEX idx_patient_access_links_user    ON patient_access_links(user_id);
