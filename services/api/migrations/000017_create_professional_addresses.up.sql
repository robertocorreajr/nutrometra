CREATE TABLE professional_addresses (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES tenants(id),
    professional_id UUID NOT NULL REFERENCES professionals(id),
    label           TEXT NOT NULL,                          -- ex: "Consultório Centro", "Clínica Zona Sul"
    street          TEXT NOT NULL,
    number          TEXT,
    complement      TEXT,
    neighborhood    TEXT,
    city            TEXT NOT NULL,
    state           TEXT NOT NULL,
    zip_code        TEXT NOT NULL,
    country         TEXT NOT NULL DEFAULT 'BR',
    latitude        NUMERIC(10,7),
    longitude       NUMERIC(10,7),
    phone           TEXT,
    notes           TEXT,
    active          BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_prof_addresses_tenant       ON professional_addresses(tenant_id);
CREATE INDEX idx_prof_addresses_professional ON professional_addresses(professional_id);
