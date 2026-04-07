CREATE TABLE tenant_user_roles (
    id             UUID PRIMARY KEY,
    tenant_id      UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    tenant_user_id UUID NOT NULL REFERENCES tenant_users(id) ON DELETE CASCADE,
    role_id        UUID NOT NULL REFERENCES roles(id),
    UNIQUE (tenant_user_id, role_id)
);

CREATE INDEX idx_tenant_user_roles_tenant_user ON tenant_user_roles(tenant_user_id);
CREATE INDEX idx_tenant_user_roles_tenant      ON tenant_user_roles(tenant_id);
