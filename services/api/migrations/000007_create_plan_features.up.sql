CREATE TABLE plan_features (
    id            UUID PRIMARY KEY,
    plan_id       UUID NOT NULL REFERENCES subscription_plans(id) ON DELETE CASCADE,
    feature_key   VARCHAR(128) NOT NULL,
    enabled       BOOLEAN NOT NULL DEFAULT TRUE,
    limit_value   BIGINT,
    trial_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    trial_days    INTEGER,
    metadata_json JSONB,
    UNIQUE (plan_id, feature_key)
);

CREATE INDEX idx_plan_features_plan_key ON plan_features(plan_id, feature_key);
