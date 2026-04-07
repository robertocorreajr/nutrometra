CREATE TABLE schedule_blocks (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES tenants(id),
    professional_id UUID NOT NULL REFERENCES professionals(id),
    start_at        TIMESTAMPTZ NOT NULL,
    end_at          TIMESTAMPTZ NOT NULL,
    reason          TEXT,
    all_day         BOOLEAN NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_block_time CHECK (end_at > start_at)
);

CREATE INDEX idx_schedule_blocks_tenant       ON schedule_blocks(tenant_id);
CREATE INDEX idx_schedule_blocks_professional ON schedule_blocks(professional_id);
CREATE INDEX idx_schedule_blocks_range        ON schedule_blocks(professional_id, start_at, end_at);
