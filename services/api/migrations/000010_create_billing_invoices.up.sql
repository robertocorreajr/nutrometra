CREATE TABLE billing_invoices (
    id                     UUID PRIMARY KEY,
    tenant_id              UUID NOT NULL REFERENCES tenants(id),
    tenant_subscription_id UUID NOT NULL REFERENCES tenant_subscriptions(id),
    provider_invoice_id    VARCHAR(255),
    amount_cents           INTEGER NOT NULL,
    currency               VARCHAR(3) NOT NULL DEFAULT 'BRL',
    status                 VARCHAR(32) NOT NULL CHECK (status IN ('draft','open','paid','void','uncollectible')),
    due_at                 TIMESTAMPTZ,
    paid_at                TIMESTAMPTZ,
    hosted_url             TEXT,
    created_at             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at             TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_billing_invoices_tenant_status ON billing_invoices(tenant_id, status);
