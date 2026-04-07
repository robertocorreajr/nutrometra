ALTER TABLE tenant_subscriptions
    DROP COLUMN IF EXISTS plan_changed_at,
    DROP COLUMN IF EXISTS previous_plan_id;
