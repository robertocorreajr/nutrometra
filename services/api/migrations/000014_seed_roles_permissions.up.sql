-- Inserir roles
INSERT INTO roles (id, code, application_scope, name, description) VALUES
  ('00000000-0000-0000-0000-000000000001', 'owner',            'tenant',    'Proprietário',       'Acesso total ao tenant'),
  ('00000000-0000-0000-0000-000000000002', 'nutritionist',     'tenant',    'Nutricionista',      'Atendimento clínico e agenda'),
  ('00000000-0000-0000-0000-000000000003', 'receptionist',     'tenant',    'Recepcionista',      'Agenda e leitura de pacientes'),
  ('00000000-0000-0000-0000-000000000004', 'patient',          'tenant',    'Paciente',           'Portal do paciente'),
  ('00000000-0000-0000-0000-000000000005', 'backoffice_admin', 'backoffice','Admin Backoffice',   'Gestão da plataforma');

-- Inserir permissions
INSERT INTO permissions (id, code, application_scope, description) VALUES
  -- Tenant management
  ('10000000-0000-0000-0000-000000000001', 'tenant:manage',        'tenant', 'Gerenciar dados do tenant'),
  ('10000000-0000-0000-0000-000000000002', 'members:invite',       'tenant', 'Convidar membros'),
  ('10000000-0000-0000-0000-000000000003', 'members:remove',       'tenant', 'Remover membros'),
  ('10000000-0000-0000-0000-000000000004', 'roles:assign',         'tenant', 'Atribuir roles'),
  ('10000000-0000-0000-0000-000000000005', 'subscription:manage',  'tenant', 'Gerenciar assinatura'),
  -- Clinical
  ('10000000-0000-0000-0000-000000000010', 'patients:read',        'tenant', 'Ler dados de pacientes'),
  ('10000000-0000-0000-0000-000000000011', 'patients:write',       'tenant', 'Escrever dados de pacientes'),
  ('10000000-0000-0000-0000-000000000012', 'clinical:write',       'tenant', 'Escrever prontuário'),
  ('10000000-0000-0000-0000-000000000013', 'schedule:manage',      'tenant', 'Gerenciar agenda'),
  ('10000000-0000-0000-0000-000000000014', 'schedule:read',        'tenant', 'Ler agenda'),
  ('10000000-0000-0000-0000-000000000015', 'diet:read',            'tenant', 'Ler dieta'),
  ('10000000-0000-0000-0000-000000000016', 'appointments:read',    'tenant', 'Ler consultas'),
  ('10000000-0000-0000-0000-000000000017', 'self:read',            'tenant', 'Ler próprio perfil'),
  -- Backoffice
  ('10000000-0000-0000-0000-000000000020', 'tenants:manage',       'backoffice', 'Gerenciar tenants'),
  ('10000000-0000-0000-0000-000000000021', 'billing:manage',       'backoffice', 'Gerenciar billing'),
  ('10000000-0000-0000-0000-000000000022', 'support:manage',       'backoffice', 'Gerenciar suporte');

-- Mapear role_permissions para owner (todas as permissões de tenant)
INSERT INTO role_permissions (id, role_id, permission_id)
SELECT gen_random_uuid(), '00000000-0000-0000-0000-000000000001', id
FROM permissions WHERE application_scope = 'tenant';

-- nutritionist
INSERT INTO role_permissions (id, role_id, permission_id)
SELECT gen_random_uuid(), '00000000-0000-0000-0000-000000000002', id
FROM permissions WHERE code IN ('patients:read','patients:write','clinical:write','schedule:manage','diet:read','appointments:read','self:read');

-- receptionist
INSERT INTO role_permissions (id, role_id, permission_id)
SELECT gen_random_uuid(), '00000000-0000-0000-0000-000000000003', id
FROM permissions WHERE code IN ('schedule:read','patients:read','appointments:read');

-- patient
INSERT INTO role_permissions (id, role_id, permission_id)
SELECT gen_random_uuid(), '00000000-0000-0000-0000-000000000004', id
FROM permissions WHERE code IN ('self:read','diet:read','appointments:read');

-- backoffice_admin
INSERT INTO role_permissions (id, role_id, permission_id)
SELECT gen_random_uuid(), '00000000-0000-0000-0000-000000000005', id
FROM permissions WHERE application_scope = 'backoffice';
