CREATE TABLE tenant_subscriptions (
    id                       UUID PRIMARY KEY,
    tenant_id                UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    plan_id                  UUID NOT NULL REFERENCES subscription_plans(id),
    coupon_id                UUID REFERENCES coupons(id),
    status                   VARCHAR(32) NOT NULL DEFAULT 'active'
                                 CHECK (status IN ('trialing','active','past_due','cancelled','expired')),
    started_at               TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    trial_ends_at            TIMESTAMPTZ,
    renews_at                TIMESTAMPTZ,
    canceled_at              TIMESTAMPTZ,
    provider_customer_id     VARCHAR(255),
    provider_subscription_id VARCHAR(255),
    created_at               TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at               TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_tenant_subscriptions_tenant_status ON tenant_subscriptions(tenant_id, status);
CREATE INDEX idx_tenant_subscriptions_tenant_id     ON tenant_subscriptions(tenant_id);
