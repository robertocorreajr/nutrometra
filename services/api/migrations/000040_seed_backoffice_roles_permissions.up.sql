INSERT INTO roles (id, code, application_scope, name, description) VALUES
  ('00000000-0000-0000-0000-000000000006', 'backoffice_support', 'backoffice', 'Suporte Backoffice',    'Suporte operacional'),
  ('00000000-0000-0000-0000-000000000007', 'backoffice_finance', 'backoffice', 'Financeiro Backoffice', 'Gestão financeira'),
  ('00000000-0000-0000-0000-000000000008', 'backoffice_sales',   'backoffice', 'Vendas Backoffice',     'Gestão comercial')
ON CONFLICT (id) DO NOTHING;

INSERT INTO permissions (id, code, application_scope, description) VALUES
  ('10000000-0000-0000-0000-000000000035', 'overrides:manage', 'backoffice', 'Gerenciar feature overrides')
ON CONFLICT (id) DO NOTHING;

-- backoffice_support: tenants:manage, support:manage
INSERT INTO role_permissions (id, role_id, permission_id)
SELECT gen_random_uuid(), '00000000-0000-0000-0000-000000000006', id
FROM permissions WHERE code IN ('tenants:manage','support:manage')
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- backoffice_finance: billing:manage
INSERT INTO role_permissions (id, role_id, permission_id)
SELECT gen_random_uuid(), '00000000-0000-0000-0000-000000000007', id
FROM permissions WHERE code IN ('billing:manage')
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- backoffice_sales: billing:manage, overrides:manage
INSERT INTO role_permissions (id, role_id, permission_id)
SELECT gen_random_uuid(), '00000000-0000-0000-0000-000000000008', id
FROM permissions WHERE code IN ('billing:manage','overrides:manage')
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- backoffice_admin: add overrides:manage
INSERT INTO role_permissions (id, role_id, permission_id)
SELECT gen_random_uuid(), '00000000-0000-0000-0000-000000000005', id
FROM permissions WHERE code = 'overrides:manage'
ON CONFLICT (role_id, permission_id) DO NOTHING;
