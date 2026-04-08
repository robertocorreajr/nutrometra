-- Adicionar features professionals:create e addresses:create para todos os planos

-- Free plan
INSERT INTO plan_features (id, plan_id, feature_key, enabled, "limit", trial_days)
VALUES
  (gen_random_uuid(), '20000000-0000-0000-0000-000000000001', 'professionals:create', TRUE, 1, NULL),
  (gen_random_uuid(), '20000000-0000-0000-0000-000000000001', 'addresses:create', TRUE, 1, NULL);

-- Starter plan
INSERT INTO plan_features (id, plan_id, feature_key, enabled, "limit", trial_days)
VALUES
  (gen_random_uuid(), '20000000-0000-0000-0000-000000000002', 'professionals:create', TRUE, 5, NULL),
  (gen_random_uuid(), '20000000-0000-0000-0000-000000000002', 'addresses:create', TRUE, 5, NULL);

-- Pro plan
INSERT INTO plan_features (id, plan_id, feature_key, enabled, "limit", trial_days)
VALUES
  (gen_random_uuid(), '20000000-0000-0000-0000-000000000003', 'professionals:create', TRUE, NULL, NULL),
  (gen_random_uuid(), '20000000-0000-0000-0000-000000000003', 'addresses:create', TRUE, NULL, NULL);

-- Enterprise plan
INSERT INTO plan_features (id, plan_id, feature_key, enabled, "limit", trial_days)
VALUES
  (gen_random_uuid(), '20000000-0000-0000-0000-000000000004', 'professionals:create', TRUE, NULL, NULL),
  (gen_random_uuid(), '20000000-0000-0000-0000-000000000004', 'addresses:create', TRUE, NULL, NULL);
