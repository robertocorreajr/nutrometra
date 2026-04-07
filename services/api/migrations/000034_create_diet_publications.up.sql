CREATE TABLE diet_publications (
    id                   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id            UUID NOT NULL REFERENCES tenants(id),
    diet_id              UUID NOT NULL REFERENCES diets(id),
    patient_id           UUID NOT NULL REFERENCES patients(id),
    published_by_user_id UUID NOT NULL REFERENCES users(id),
    published_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    status               TEXT NOT NULL DEFAULT 'active'
                         CHECK (status IN ('active', 'superseded'))
);
CREATE INDEX idx_diet_pub_patient ON diet_publications(tenant_id, patient_id);
CREATE UNIQUE INDEX idx_diet_pub_active ON diet_publications(diet_id) WHERE status = 'active';
