CREATE TABLE tenant_users (
    id                 UUID PRIMARY KEY,
    tenant_id          UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    user_id            UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status             VARCHAR(32) NOT NULL DEFAULT 'active' CHECK (status IN ('active','suspended','removed')),
    joined_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    invited_by_user_id UUID REFERENCES users(id),
    created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (tenant_id, user_id)
);

CREATE INDEX idx_tenant_users_tenant_id     ON tenant_users(tenant_id);
CREATE INDEX idx_tenant_users_user_id       ON tenant_users(user_id);
CREATE INDEX idx_tenant_users_tenant_status ON tenant_users(tenant_id, status);
