-- Remove overrides:manage from backoffice_admin
DELETE FROM role_permissions
WHERE role_id = '00000000-0000-0000-0000-000000000005'
  AND permission_id = (SELECT id FROM permissions WHERE code = 'overrides:manage');

-- Remove role_permissions for new backoffice roles
DELETE FROM role_permissions
WHERE role_id IN (
    '00000000-0000-0000-0000-000000000006',
    '00000000-0000-0000-0000-000000000007',
    '00000000-0000-0000-0000-000000000008'
);

DELETE FROM permissions WHERE id = '10000000-0000-0000-0000-000000000035';

DELETE FROM roles WHERE id IN (
    '00000000-0000-0000-0000-000000000006',
    '00000000-0000-0000-0000-000000000007',
    '00000000-0000-0000-0000-000000000008'
);
