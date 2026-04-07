CREATE TABLE appointment_audit_events (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id      UUID NOT NULL REFERENCES tenants(id),
    appointment_id UUID NOT NULL REFERENCES appointments(id),
    actor_user_id  UUID REFERENCES users(id),
    event_type     TEXT NOT NULL,              -- scheduled, confirmed, cancelled, rescheduled, completed, no_show
    old_status     TEXT,
    new_status     TEXT,
    old_start_at   TIMESTAMPTZ,
    new_start_at   TIMESTAMPTZ,
    notes          TEXT,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_appt_audit_tenant      ON appointment_audit_events(tenant_id);
CREATE INDEX idx_appt_audit_appointment ON appointment_audit_events(appointment_id);
