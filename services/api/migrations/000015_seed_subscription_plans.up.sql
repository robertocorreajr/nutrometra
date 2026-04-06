-- Planos
INSERT INTO subscription_plans (id, code, name, active, billing_cycle, currency, price_cents) VALUES
  ('20000000-0000-0000-0000-000000000001', 'free',       'Free',       TRUE, 'monthly', 'BRL',      0),
  ('20000000-0000-0000-0000-000000000002', 'starter',    'Starter',    TRUE, 'monthly', 'BRL',   9900),
  ('20000000-0000-0000-0000-000000000003', 'pro',        'Pro',        TRUE, 'monthly', 'BRL',  19900),
  ('20000000-0000-0000-0000-000000000004', 'enterprise', 'Enterprise', TRUE, 'monthly', 'BRL',     -1);

-- Features do plano Free
INSERT INTO plan_features (id, plan_id, feature_key, enabled, limit_value, trial_enabled, trial_days) VALUES
  (gen_random_uuid(), '20000000-0000-0000-0000-000000000001', 'patients:create',     TRUE,   5, FALSE, NULL),
  (gen_random_uuid(), '20000000-0000-0000-0000-000000000001', 'appointments:create', TRUE,  20, FALSE, NULL),
  (gen_random_uuid(), '20000000-0000-0000-0000-000000000001', 'pdf:export',          FALSE, NULL, TRUE,   7),
  (gen_random_uuid(), '20000000-0000-0000-0000-000000000001', 'ai:assist',           FALSE, NULL, TRUE,   7);

-- Features do plano Starter
INSERT INTO plan_features (id, plan_id, feature_key, enabled, limit_value, trial_enabled, trial_days) VALUES
  (gen_random_uuid(), '20000000-0000-0000-0000-000000000002', 'patients:create',     TRUE,  50, FALSE, NULL),
  (gen_random_uuid(), '20000000-0000-0000-0000-000000000002', 'appointments:create', TRUE, NULL, FALSE, NULL),
  (gen_random_uuid(), '20000000-0000-0000-0000-000000000002', 'pdf:export',          TRUE, NULL, FALSE, NULL),
  (gen_random_uuid(), '20000000-0000-0000-0000-000000000002', 'ai:assist',           FALSE, NULL, TRUE,  14);

-- Features do plano Pro
INSERT INTO plan_features (id, plan_id, feature_key, enabled, limit_value, trial_enabled, trial_days) VALUES
  (gen_random_uuid(), '20000000-0000-0000-0000-000000000003', 'patients:create',     TRUE, NULL, FALSE, NULL),
  (gen_random_uuid(), '20000000-0000-0000-0000-000000000003', 'appointments:create', TRUE, NULL, FALSE, NULL),
  (gen_random_uuid(), '20000000-0000-0000-0000-000000000003', 'pdf:export',          TRUE, NULL, FALSE, NULL),
  (gen_random_uuid(), '20000000-0000-0000-0000-000000000003', 'ai:assist',           TRUE, NULL, FALSE, NULL),
  (gen_random_uuid(), '20000000-0000-0000-0000-000000000003', 'google_calendar:sync',TRUE, NULL, FALSE, NULL);

-- Features do plano Enterprise (sem limites, todas habilitadas)
INSERT INTO plan_features (id, plan_id, feature_key, enabled, limit_value, trial_enabled, trial_days) VALUES
  (gen_random_uuid(), '20000000-0000-0000-0000-000000000004', 'patients:create',     TRUE, NULL, FALSE, NULL),
  (gen_random_uuid(), '20000000-0000-0000-0000-000000000004', 'appointments:create', TRUE, NULL, FALSE, NULL),
  (gen_random_uuid(), '20000000-0000-0000-0000-000000000004', 'pdf:export',          TRUE, NULL, FALSE, NULL),
  (gen_random_uuid(), '20000000-0000-0000-0000-000000000004', 'ai:assist',           TRUE, NULL, FALSE, NULL),
  (gen_random_uuid(), '20000000-0000-0000-0000-000000000004', 'google_calendar:sync',TRUE, NULL, FALSE, NULL),
  (gen_random_uuid(), '20000000-0000-0000-0000-000000000004', 'custom_branding',     TRUE, NULL, FALSE, NULL),
  (gen_random_uuid(), '20000000-0000-0000-0000-000000000004', 'sso',                 TRUE, NULL, FALSE, NULL);
