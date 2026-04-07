CREATE TABLE patients (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES tenants(id),
    professional_id UUID NOT NULL REFERENCES professionals(id),
    full_name       TEXT NOT NULL,
    email           TEXT,
    phone           TEXT,
    cpf             TEXT,
    date_of_birth   DATE,
    gender          TEXT CHECK (gender IS NULL OR gender IN ('male','female','other','prefer_not_to_say')),
    notes           TEXT,
    active          BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE (tenant_id, cpf)
);

CREATE INDEX idx_patients_tenant       ON patients(tenant_id);
CREATE INDEX idx_patients_professional ON patients(professional_id);
CREATE INDEX idx_patients_email        ON patients(tenant_id, email);
CREATE INDEX idx_patients_name         ON patients(tenant_id, full_name);

-- Add FK from appointments to patients
ALTER TABLE appointments ADD CONSTRAINT fk_appointments_patient
    FOREIGN KEY (patient_id) REFERENCES patients(id);
