CREATE TABLE tenants (
    id           UUID PRIMARY KEY,
    type         VARCHAR(32) NOT NULL CHECK (type IN ('solo_professional','clinic','company')),
    legal_name   VARCHAR(255) NOT NULL,
    display_name VARCHAR(255) NOT NULL,
    slug         VARCHAR(100) NOT NULL UNIQUE,
    status       VARCHAR(32) NOT NULL DEFAULT 'active' CHECK (status IN ('active','suspended','cancelled')),
    timezone     VARCHAR(64) NOT NULL DEFAULT 'America/Sao_Paulo',
    locale       VARCHAR(10) NOT NULL DEFAULT 'pt-BR',
    trial_ends_at TIMESTAMPTZ,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_tenants_slug   ON tenants(slug);
CREATE INDEX idx_tenants_status ON tenants(status);
