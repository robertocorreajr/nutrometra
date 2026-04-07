CREATE TABLE ai_suggestions (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id         UUID NOT NULL REFERENCES tenants(id),
    user_id           UUID NOT NULL REFERENCES users(id),
    suggestion_type   VARCHAR(64) NOT NULL
                      CHECK (suggestion_type IN (
                          'diet_draft','meal_structure','substitutions',
                          'clinical_summary','review_checklist')),
    status            VARCHAR(32) NOT NULL DEFAULT 'pending'
                      CHECK (status IN (
                          'pending','generating','completed','failed',
                          'accepted','rejected')),
    input_context_json JSONB NOT NULL,
    prompt_used       TEXT,
    response_text     TEXT,
    model_id          VARCHAR(64),
    input_tokens      INTEGER,
    output_tokens     INTEGER,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at      TIMESTAMPTZ,
    reviewed_at       TIMESTAMPTZ,
    reviewed_by       UUID REFERENCES users(id)
);
CREATE INDEX idx_ai_suggestions_tenant ON ai_suggestions(tenant_id, user_id);
