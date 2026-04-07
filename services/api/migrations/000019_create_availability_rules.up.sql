CREATE TABLE availability_rules (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES tenants(id),
    professional_id UUID NOT NULL REFERENCES professionals(id),
    day_of_week     INT NOT NULL CHECK (day_of_week BETWEEN 0 AND 6),  -- 0=Sunday
    start_time      TIME NOT NULL,
    end_time        TIME NOT NULL,
    service_mode    TEXT NOT NULL CHECK (service_mode IN ('onsite','online','home_visit')),
    address_id      UUID REFERENCES professional_addresses(id),
    active          BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_availability_time CHECK (end_time > start_time),
    CONSTRAINT chk_availability_onsite CHECK (
        (service_mode = 'onsite' AND address_id IS NOT NULL) OR service_mode <> 'onsite'
    )
);

CREATE INDEX idx_availability_rules_tenant       ON availability_rules(tenant_id);
CREATE INDEX idx_availability_rules_professional ON availability_rules(professional_id);
CREATE INDEX idx_availability_rules_day          ON availability_rules(professional_id, day_of_week);
