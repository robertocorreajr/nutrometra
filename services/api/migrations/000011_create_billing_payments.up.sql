CREATE TABLE billing_payments (
    id                  UUID PRIMARY KEY,
    tenant_id           UUID NOT NULL REFERENCES tenants(id) ON DELETE RESTRICT,
    billing_invoice_id  UUID NOT NULL REFERENCES billing_invoices(id) ON DELETE RESTRICT,
    provider_payment_id VARCHAR(255),
    amount_cents        BIGINT NOT NULL,
    currency            VARCHAR(3) NOT NULL DEFAULT 'BRL',
    status              VARCHAR(32) NOT NULL CHECK (status IN ('pending','succeeded','failed','refunded')),
    paid_at             TIMESTAMPTZ,
    failure_reason      TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_billing_payments_tenant_status ON billing_payments(tenant_id, status);
CREATE INDEX idx_billing_payments_invoice        ON billing_payments(billing_invoice_id);
