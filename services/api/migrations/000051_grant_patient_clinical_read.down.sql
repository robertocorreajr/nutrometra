-- Remove clinical:read and patients:read from the patient role.

DELETE FROM role_permissions
WHERE role_id = (SELECT id FROM roles WHERE code = 'patient')
  AND permission_id IN (
    SELECT id FROM permissions WHERE code IN ('clinical:read', 'patients:read')
  );
