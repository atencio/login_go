-- Modelo de permisos por módulo y acción (Crear, Listar, Encontrar, Editar, Eliminar)
ALTER TABLE permissions ADD COLUMN IF NOT EXISTS module       TEXT NOT NULL DEFAULT '';
ALTER TABLE permissions ADD COLUMN IF NOT EXISTS action       TEXT NOT NULL DEFAULT '';
ALTER TABLE permissions ADD COLUMN IF NOT EXISTS module_label TEXT NOT NULL DEFAULT '';
ALTER TABLE permissions ADD COLUMN IF NOT EXISTS sort         INT  NOT NULL DEFAULT 0;

UPDATE permissions SET code = 'users.list'   WHERE code = 'users.read';
UPDATE permissions SET code = 'reports.list' WHERE code = 'reports.read';
DELETE FROM permissions WHERE code = 'reports.export';

INSERT INTO permissions (code, description, module, action, module_label, sort) VALUES
    ('users.create',   'Crear usuarios',     'users',   'create', 'Usuarios', 10),
    ('users.list',     'Listar usuarios',    'users',   'list',   'Usuarios', 11),
    ('users.find',     'Buscar usuarios',    'users',   'find',   'Usuarios', 12),
    ('users.update',   'Editar usuarios',    'users',   'update', 'Usuarios', 13),
    ('users.delete',   'Eliminar usuarios',  'users',   'delete', 'Usuarios', 14),
    ('roles.create',   'Crear perfiles',     'roles',   'create', 'Perfiles', 20),
    ('roles.list',     'Listar perfiles',    'roles',   'list',   'Perfiles', 21),
    ('roles.find',     'Ver perfiles',       'roles',   'find',   'Perfiles', 22),
    ('roles.update',   'Editar perfiles',    'roles',   'update', 'Perfiles', 23),
    ('roles.delete',   'Eliminar perfiles',  'roles',   'delete', 'Perfiles', 24),
    ('reports.list',   'Ver reportes',       'reports', 'list',   'Reportes', 31)
ON CONFLICT (code) DO UPDATE SET
    description = EXCLUDED.description, module = EXCLUDED.module, action = EXCLUDED.action,
    module_label = EXCLUDED.module_label, sort = EXCLUDED.sort;

-- El administrador siempre tiene todo
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p WHERE r.name = 'Administrador'
ON CONFLICT DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p
WHERE r.name = 'Supervisor' AND p.code IN ('users.list', 'users.find', 'reports.list')
ON CONFLICT DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p
WHERE r.name = 'Consulta' AND p.code IN ('reports.list')
ON CONFLICT DO NOTHING;
