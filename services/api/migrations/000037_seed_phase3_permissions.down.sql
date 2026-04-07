DELETE FROM role_permissions WHERE permission_id IN (
    SELECT id FROM permissions WHERE code IN ('diet:write','diet:publish','document:write','document:publish','export:pdf')
);
DELETE FROM role_permissions WHERE role_id = '00000000-0000-0000-0000-000000000004' AND permission_id = (
    SELECT id FROM permissions WHERE code = 'clinical:read'
);
DELETE FROM permissions WHERE code IN ('diet:write','diet:publish','document:write','document:publish','export:pdf');
