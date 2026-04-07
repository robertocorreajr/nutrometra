CREATE TABLE diets (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id           UUID NOT NULL REFERENCES tenants(id),
    patient_id          UUID NOT NULL REFERENCES patients(id),
    professional_id     UUID NOT NULL REFERENCES professionals(id),
    title               TEXT NOT NULL,
    objective           TEXT,
    status              TEXT NOT NULL DEFAULT 'draft'
                        CHECK (status IN ('draft', 'published', 'archived')),
    version_number      INT NOT NULL DEFAULT 1,
    previous_version_id UUID REFERENCES diets(id),
    published_at        TIMESTAMPTZ,
    valid_from          DATE,
    valid_until         DATE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_diets_tenant ON diets(tenant_id);
CREATE INDEX idx_diets_patient ON diets(tenant_id, patient_id);
CREATE INDEX idx_diets_professional ON diets(tenant_id, professional_id);
