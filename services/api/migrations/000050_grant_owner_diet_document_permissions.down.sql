-- Remove diet, document and export permissions from Proprietário role.

DELETE FROM role_permissions
WHERE role_id = (SELECT id FROM roles WHERE name = 'Proprietário')
  AND permission_id IN (
    SELECT id FROM permissions WHERE code IN ('diet:write', 'diet:publish', 'document:write', 'document:publish', 'export:pdf')
  );
