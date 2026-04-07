CREATE TABLE appointments (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES tenants(id),
    professional_id UUID NOT NULL REFERENCES professionals(id),
    patient_id      UUID,                                              -- FK added in migration 000023
    start_at        TIMESTAMPTZ NOT NULL,
    end_at          TIMESTAMPTZ NOT NULL,
    service_mode    TEXT NOT NULL CHECK (service_mode IN ('onsite','online','home_visit')),
    address_id      UUID REFERENCES professional_addresses(id),
    status          TEXT NOT NULL DEFAULT 'scheduled'
                    CHECK (status IN ('scheduled','confirmed','completed','cancelled','no_show')),
    source          TEXT NOT NULL DEFAULT 'professional'
                    CHECK (source IN ('professional','patient','system')),
    notes           TEXT,
    cancellation_reason TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_appointment_time CHECK (end_at > start_at)
);

CREATE INDEX idx_appointments_tenant        ON appointments(tenant_id);
CREATE INDEX idx_appointments_professional  ON appointments(professional_id, start_at, end_at);
CREATE INDEX idx_appointments_patient       ON appointments(patient_id);
CREATE INDEX idx_appointments_status        ON appointments(tenant_id, status);
CREATE INDEX idx_appointments_date          ON appointments(tenant_id, start_at);
