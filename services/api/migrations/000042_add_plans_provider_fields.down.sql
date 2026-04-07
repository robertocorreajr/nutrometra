ALTER TABLE subscription_plans
    DROP COLUMN IF EXISTS provider_product_id,
    DROP COLUMN IF EXISTS provider_price_id;
