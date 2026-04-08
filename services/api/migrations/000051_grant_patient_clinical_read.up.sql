-- Grant clinical:read and patients:read to the patient role.
-- Needed for patient portal: documents, measurements, and profile access.

INSERT INTO role_permissions (id, role_id, permission_id)
SELECT gen_random_uuid(), r.id, p.id
FROM roles r, permissions p
WHERE r.code = 'patient'
  AND p.code IN ('clinical:read', 'patients:read')
ON CONFLICT (role_id, permission_id) DO NOTHING;
