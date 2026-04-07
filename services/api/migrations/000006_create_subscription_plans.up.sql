CREATE TABLE subscription_plans (
    id            UUID PRIMARY KEY,
    code          VARCHAR(64) NOT NULL UNIQUE,
    name          VARCHAR(128) NOT NULL,
    active        BOOLEAN NOT NULL DEFAULT TRUE,
    billing_cycle VARCHAR(32) NOT NULL CHECK (billing_cycle IN ('monthly','yearly','lifetime')),
    currency      VARCHAR(3) NOT NULL DEFAULT 'BRL',
    price_cents   BIGINT NOT NULL DEFAULT 0,
    metadata_json JSONB,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_subscription_plans_active ON subscription_plans(active) WHERE active = TRUE;
