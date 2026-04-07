CREATE TABLE anamneses (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES tenants(id),
    patient_id      UUID NOT NULL REFERENCES patients(id),
    professional_id UUID NOT NULL REFERENCES professionals(id),
    status          TEXT NOT NULL DEFAULT 'draft' CHECK (status IN ('draft','finalized')),
    chief_complaint TEXT,
    history_present_illness TEXT,
    past_medical_history    TEXT,
    family_history          TEXT,
    social_history          TEXT,
    dietary_history         TEXT,
    physical_activity       TEXT,
    sleep_pattern           TEXT,
    bowel_habits            TEXT,
    water_intake            TEXT,
    supplements             TEXT,
    observations            TEXT,
    finalized_at   TIMESTAMPTZ,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_anamneses_tenant   ON anamneses(tenant_id);
CREATE INDEX idx_anamneses_patient  ON anamneses(patient_id);
CREATE INDEX idx_anamneses_prof     ON anamneses(professional_id);
