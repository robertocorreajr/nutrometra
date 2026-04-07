CREATE TABLE calendar_connections (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id               UUID NOT NULL REFERENCES tenants(id),
    user_id                 UUID NOT NULL REFERENCES users(id),
    provider                VARCHAR(32) NOT NULL DEFAULT 'google',
    status                  VARCHAR(32) NOT NULL DEFAULT 'active'
                            CHECK (status IN ('active','revoked','expired')),
    external_account_id     VARCHAR(255),
    encrypted_access_token  TEXT NOT NULL,
    encrypted_refresh_token TEXT NOT NULL,
    token_expires_at        TIMESTAMPTZ,
    scopes                  TEXT NOT NULL,
    synced_at               TIMESTAMPTZ,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (tenant_id, user_id, provider)
);
