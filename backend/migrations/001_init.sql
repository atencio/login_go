CREATE TABLE IF NOT EXISTS roles (
    id   BIGSERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS permissions (
    id          BIGSERIAL PRIMARY KEY,
    code        TEXT UNIQUE NOT NULL,
    description TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS role_permissions (
    role_id       BIGINT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id BIGINT NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    PRIMARY KEY (role_id, permission_id)
);

CREATE TABLE IF NOT EXISTS users (
    id            BIGSERIAL PRIMARY KEY,
    name          TEXT NOT NULL,
    email         TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    role_id       BIGINT NOT NULL REFERENCES roles(id),
    active        BOOLEAN NOT NULL DEFAULT TRUE,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

INSERT INTO permissions (code, description) VALUES
    ('users.read',     'Ver usuarios'),
    ('users.create',   'Crear usuarios'),
    ('users.update',   'Editar usuarios'),
    ('users.delete',   'Eliminar usuarios'),
    ('reports.read',   'Ver reportes'),
    ('reports.export', 'Exportar reportes')
ON CONFLICT (code) DO NOTHING;

INSERT INTO roles (name) VALUES ('Administrador'), ('Supervisor'), ('Consulta')
ON CONFLICT (name) DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p WHERE r.name = 'Administrador'
ON CONFLICT DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p
WHERE r.name = 'Supervisor' AND p.code IN ('users.read', 'reports.read', 'reports.export')
ON CONFLICT DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p
WHERE r.name = 'Consulta' AND p.code IN ('reports.read')
ON CONFLICT DO NOTHING;
