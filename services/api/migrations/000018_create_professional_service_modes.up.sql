CREATE TABLE professional_service_modes (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES tenants(id),
    professional_id UUID NOT NULL REFERENCES professionals(id),
    mode            TEXT NOT NULL CHECK (mode IN ('onsite','online','home_visit')),
    address_id      UUID REFERENCES professional_addresses(id),  -- obrigatório se mode = 'onsite'
    duration_min    INT NOT NULL DEFAULT 50 CHECK (duration_min > 0),
    active          BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE (professional_id, mode, address_id)
);

CREATE INDEX idx_prof_svc_modes_tenant       ON professional_service_modes(tenant_id);
CREATE INDEX idx_prof_svc_modes_professional ON professional_service_modes(professional_id);

-- Enforce address_id when mode = 'onsite'
ALTER TABLE professional_service_modes ADD CONSTRAINT chk_onsite_address
    CHECK (
        (mode = 'onsite' AND address_id IS NOT NULL)
        OR mode <> 'onsite'
    );
