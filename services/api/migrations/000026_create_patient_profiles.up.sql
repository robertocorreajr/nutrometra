CREATE TABLE patient_profiles (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES tenants(id),
    patient_id      UUID NOT NULL REFERENCES patients(id),
    occupation      TEXT,
    marital_status  TEXT,
    ethnicity       TEXT,
    blood_type      TEXT,
    allergies       TEXT[],
    chronic_conditions TEXT[],
    medications     TEXT[],
    emergency_contact_name  TEXT,
    emergency_contact_phone TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE (tenant_id, patient_id)
);

CREATE TABLE patient_consents (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES tenants(id),
    patient_id      UUID NOT NULL REFERENCES patients(id),
    consent_type    TEXT NOT NULL,        -- 'data_processing', 'portal_terms', 'communication'
    accepted        BOOLEAN NOT NULL DEFAULT FALSE,
    accepted_at     TIMESTAMPTZ,
    ip_address      INET,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_patient_profiles_tenant  ON patient_profiles(tenant_id);
CREATE INDEX idx_patient_consents_tenant  ON patient_consents(tenant_id);
CREATE INDEX idx_patient_consents_patient ON patient_consents(patient_id);
