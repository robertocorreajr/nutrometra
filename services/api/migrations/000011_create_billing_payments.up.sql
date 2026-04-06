CREATE TABLE billing_payments (
    id                  UUID PRIMARY KEY,
    tenant_id           UUID NOT NULL REFERENCES tenants(id),
    billing_invoice_id  UUID NOT NULL REFERENCES billing_invoices(id),
    provider_payment_id VARCHAR(255),
    amount_cents        INTEGER NOT NULL,
    currency            VARCHAR(3) NOT NULL DEFAULT 'BRL',
    status              VARCHAR(32) NOT NULL CHECK (status IN ('pending','succeeded','failed','refunded')),
    paid_at             TIMESTAMPTZ,
    failure_reason      TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
