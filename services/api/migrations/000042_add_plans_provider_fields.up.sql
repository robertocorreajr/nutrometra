ALTER TABLE subscription_plans
    ADD COLUMN provider_price_id VARCHAR(255),
    ADD COLUMN provider_product_id VARCHAR(255);
