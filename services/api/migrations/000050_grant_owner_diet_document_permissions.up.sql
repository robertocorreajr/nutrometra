-- Grant missing diet, document and export permissions to Proprietário role.
-- The Proprietário (owner) should have all tenant-scoped permissions.

INSERT INTO role_permissions (id, role_id, permission_id)
SELECT gen_random_uuid(), r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'Proprietário'
  AND p.code IN ('diet:write', 'diet:publish', 'document:write', 'document:publish', 'export:pdf')
ON CONFLICT DO NOTHING;
