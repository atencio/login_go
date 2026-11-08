-- Inventario de asignación de activos: cargos, funcionarios, equipos y asignaciones

CREATE TABLE IF NOT EXISTS positions (
    id          BIGSERIAL PRIMARY KEY,
    name        TEXT UNIQUE NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    active      BOOLEAN NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS employees (
    id          BIGSERIAL PRIMARY KEY,
    first_name  TEXT NOT NULL,
    last_name   TEXT NOT NULL,
    email       TEXT UNIQUE NOT NULL,
    document    TEXT NOT NULL DEFAULT '',
    phone       TEXT NOT NULL DEFAULT '',
    department  TEXT NOT NULL DEFAULT '',
    position_id BIGINT NOT NULL REFERENCES positions(id),
    hire_date   DATE,
    active      BOOLEAN NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE UNIQUE INDEX IF NOT EXISTS employees_document_uq ON employees (document) WHERE document <> '';
CREATE INDEX IF NOT EXISTS employees_position_idx ON employees (position_id);

CREATE SEQUENCE IF NOT EXISTS equipment_code_seq;

CREATE TABLE IF NOT EXISTS equipment (
    id            BIGSERIAL PRIMARY KEY,
    code          TEXT UNIQUE NOT NULL,
    name          TEXT NOT NULL,
    category      TEXT NOT NULL DEFAULT '',
    brand         TEXT NOT NULL DEFAULT '',
    model         TEXT NOT NULL DEFAULT '',
    serial_number TEXT NOT NULL DEFAULT '',
    status        TEXT NOT NULL DEFAULT 'available'
                  CHECK (status IN ('available', 'assigned', 'maintenance', 'retired')),
    purchase_date DATE,
    cost          NUMERIC(12,2) CHECK (cost IS NULL OR cost >= 0),
    notes         TEXT NOT NULL DEFAULT '',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE UNIQUE INDEX IF NOT EXISTS equipment_serial_uq ON equipment (serial_number) WHERE serial_number <> '';

-- Cada fila es un movimiento: entrega (assigned_at) y, cuando termina, devolución (returned_at).
CREATE TABLE IF NOT EXISTS assignments (
    id             BIGSERIAL PRIMARY KEY,
    equipment_id   BIGINT NOT NULL REFERENCES equipment(id),
    employee_id    BIGINT NOT NULL REFERENCES employees(id),
    assigned_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    assigned_by    BIGINT REFERENCES users(id) ON DELETE SET NULL,
    delivery_notes TEXT NOT NULL DEFAULT '',
    returned_at    TIMESTAMPTZ,
    received_by    BIGINT REFERENCES users(id) ON DELETE SET NULL,
    return_notes   TEXT NOT NULL DEFAULT ''
);
-- Un equipo solo puede tener una asignación activa a la vez
CREATE UNIQUE INDEX IF NOT EXISTS assignments_one_active ON assignments (equipment_id) WHERE returned_at IS NULL;
CREATE INDEX IF NOT EXISTS assignments_employee_idx ON assignments (employee_id);

INSERT INTO permissions (code, description, module, action, module_label, sort) VALUES
    ('positions.create',   'Crear cargos',            'positions',   'create', 'Cargos', 40),
    ('positions.list',     'Listar cargos',           'positions',   'list',   'Cargos', 41),
    ('positions.find',     'Buscar cargos',           'positions',   'find',   'Cargos', 42),
    ('positions.update',   'Editar cargos',           'positions',   'update', 'Cargos', 43),
    ('positions.delete',   'Eliminar cargos',         'positions',   'delete', 'Cargos', 44),
    ('employees.create',   'Crear funcionarios',      'employees',   'create', 'Funcionarios', 50),
    ('employees.list',     'Listar funcionarios',     'employees',   'list',   'Funcionarios', 51),
    ('employees.find',     'Buscar funcionarios',     'employees',   'find',   'Funcionarios', 52),
    ('employees.update',   'Editar funcionarios',     'employees',   'update', 'Funcionarios', 53),
    ('employees.delete',   'Eliminar funcionarios',   'employees',   'delete', 'Funcionarios', 54),
    ('equipment.create',   'Registrar equipos',       'equipment',   'create', 'Equipos de oficina', 60),
    ('equipment.list',     'Listar equipos',          'equipment',   'list',   'Equipos de oficina', 61),
    ('equipment.find',     'Buscar equipos',          'equipment',   'find',   'Equipos de oficina', 62),
    ('equipment.update',   'Editar equipos',          'equipment',   'update', 'Equipos de oficina', 63),
    ('equipment.delete',   'Eliminar equipos',        'equipment',   'delete', 'Equipos de oficina', 64),
    ('assignments.create', 'Asignar equipos',         'assignments', 'create', 'Asignaciones', 70),
    ('assignments.list',   'Ver asignaciones e historial', 'assignments', 'list',   'Asignaciones', 71),
    ('assignments.find',   'Buscar asignaciones',     'assignments', 'find',   'Asignaciones', 72),
    ('assignments.update', 'Actualizar y devolver asignaciones', 'assignments', 'update', 'Asignaciones', 73)
ON CONFLICT (code) DO UPDATE SET
    description = EXCLUDED.description, module = EXCLUDED.module, action = EXCLUDED.action,
    module_label = EXCLUDED.module_label, sort = EXCLUDED.sort;

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p WHERE r.name = 'Administrador'
ON CONFLICT DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p
WHERE r.name = 'Supervisor' AND p.code IN (
    'positions.list', 'positions.find', 'employees.list', 'employees.find',
    'equipment.list', 'equipment.find',
    'assignments.create', 'assignments.list', 'assignments.find', 'assignments.update')
ON CONFLICT DO NOTHING;
