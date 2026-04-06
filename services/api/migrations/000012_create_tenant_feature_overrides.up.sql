CREATE TABLE tenant_feature_overrides (
    id                 UUID PRIMARY KEY,
    tenant_id          UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    feature_key        VARCHAR(128) NOT NULL,
    enabled            BOOLEAN,
    limit_value        BIGINT,
    starts_at          TIMESTAMPTZ,
    ends_at            TIMESTAMPTZ,
    reason             TEXT NOT NULL,
    created_by_user_id UUID REFERENCES users(id),
    created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (tenant_id, feature_key),
    CHECK (enabled IS NOT NULL OR limit_value IS NOT NULL)
);

CREATE INDEX idx_tenant_feature_overrides_tenant_key ON tenant_feature_overrides(tenant_id, feature_key);
