INSERT INTO permissions (id, code, application_scope, description) VALUES
  ('10000000-0000-0000-0000-000000000030', 'diet:write',      'tenant', 'Criar e editar dietas'),
  ('10000000-0000-0000-0000-000000000031', 'diet:publish',    'tenant', 'Publicar dietas ao paciente'),
  ('10000000-0000-0000-0000-000000000032', 'document:write',  'tenant', 'Criar e editar documentos clínicos'),
  ('10000000-0000-0000-0000-000000000033', 'document:publish','tenant', 'Publicar documentos ao paciente'),
  ('10000000-0000-0000-0000-000000000034', 'export:pdf',      'tenant', 'Exportar PDF')
ON CONFLICT (id) DO NOTHING;

-- nutritionist role: grant all phase 3 permissions
INSERT INTO role_permissions (id, role_id, permission_id)
SELECT gen_random_uuid(), '00000000-0000-0000-0000-000000000002', id
FROM permissions WHERE code IN ('diet:write','diet:publish','document:write','document:publish','export:pdf')
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- patient role: grant clinical:read so patients can view published documents
INSERT INTO role_permissions (id, role_id, permission_id)
SELECT gen_random_uuid(), '00000000-0000-0000-0000-000000000004', id
FROM permissions WHERE code = 'clinical:read'
ON CONFLICT (role_id, permission_id) DO NOTHING;
