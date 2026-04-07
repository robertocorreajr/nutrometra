ALTER TABLE tenant_subscriptions
    ADD COLUMN previous_plan_id UUID REFERENCES subscription_plans(id),
    ADD COLUMN plan_changed_at TIMESTAMPTZ;
