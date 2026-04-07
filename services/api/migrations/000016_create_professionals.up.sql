CREATE TABLE professionals (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES tenants(id),
    user_id         UUID NOT NULL REFERENCES users(id),
    full_name       TEXT NOT NULL,
    registration_type TEXT NOT NULL CHECK (registration_type IN ('CRN','CRM','other')),
    registration_number TEXT NOT NULL,
    registration_state  TEXT,          -- UF obrigatório se CRN/CRM
    specialty       TEXT,
    bio             TEXT,
    phone           TEXT,
    avatar_url      TEXT,
    active          BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE (tenant_id, user_id),
    UNIQUE (tenant_id, registration_type, registration_number)
);

CREATE INDEX idx_professionals_tenant ON professionals(tenant_id);
CREATE INDEX idx_professionals_user   ON professionals(user_id);

-- Enforce registration_state when registration_type requires it
ALTER TABLE professionals ADD CONSTRAINT chk_registration_state
    CHECK (
        (registration_type IN ('CRN','CRM') AND registration_state IS NOT NULL AND registration_state <> '')
        OR registration_type = 'other'
    );
