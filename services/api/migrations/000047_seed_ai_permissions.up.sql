INSERT INTO permissions (id, code, application_scope, description) VALUES
    ('10000000-0000-0000-0000-000000000030', 'ai:suggest', 'tenant', 'Solicitar sugestao de IA'),
    ('10000000-0000-0000-0000-000000000031', 'ai:read',    'tenant', 'Ler sugestoes de IA')
ON CONFLICT (id) DO NOTHING;

-- Grant to nutritionist role (role_id from seed 000014)
INSERT INTO role_permissions (id, role_id, permission_id)
SELECT gen_random_uuid(), '00000000-0000-0000-0000-000000000002', id
FROM permissions WHERE code IN ('ai:suggest', 'ai:read')
ON CONFLICT (role_id, permission_id) DO NOTHING;
