-- Categorías y marcas parametrizables (antes eran texto libre en el equipo)

CREATE TABLE IF NOT EXISTS equipment_categories (
    id         BIGSERIAL PRIMARY KEY,
    name       TEXT UNIQUE NOT NULL,
    active     BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS equipment_brands (
    id         BIGSERIAL PRIMARY KEY,
    name       TEXT UNIQUE NOT NULL,
    active     BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

INSERT INTO equipment_categories (name) VALUES
    ('Computador portátil'), ('Computador de escritorio'), ('Monitor'), ('Teclado'), ('Mouse'),
    ('Teléfono'), ('Impresora'), ('Escáner'), ('Proyector'), ('Tablet'), ('Celular'),
    ('Silla'), ('Escritorio'), ('Otro')
ON CONFLICT (name) DO NOTHING;

INSERT INTO equipment_brands (name) VALUES
    ('Dell'), ('HP'), ('Lenovo'), ('Apple'), ('Asus'), ('Acer'), ('Samsung'), ('LG'), ('Microsoft'),
    ('Logitech'), ('Epson'), ('Canon'), ('Brother'), ('Xerox'), ('Ricoh'), ('Kyocera'), ('Cisco'),
    ('Huawei'), ('Xiaomi'), ('Sony'), ('Toshiba'), ('Philips'), ('BenQ'), ('ViewSonic'), ('AOC'),
    ('Genius'), ('Kingston'), ('Jabra'), ('Plantronics'), ('Polycom'), ('Panasonic'), ('Steelcase'),
    ('Herman Miller')
ON CONFLICT (name) DO NOTHING;

-- Conserva lo que ya estuviera escrito en los equipos existentes
INSERT INTO equipment_categories (name)
SELECT DISTINCT btrim(category) FROM equipment WHERE btrim(category) <> '' ON CONFLICT (name) DO NOTHING;
INSERT INTO equipment_brands (name)
SELECT DISTINCT btrim(brand) FROM equipment WHERE btrim(brand) <> '' ON CONFLICT (name) DO NOTHING;

ALTER TABLE equipment ADD COLUMN IF NOT EXISTS category_id BIGINT REFERENCES equipment_categories(id);
ALTER TABLE equipment ADD COLUMN IF NOT EXISTS brand_id    BIGINT REFERENCES equipment_brands(id);

UPDATE equipment e SET category_id = c.id FROM equipment_categories c WHERE btrim(e.category) = c.name;
UPDATE equipment e SET brand_id    = b.id FROM equipment_brands b     WHERE btrim(e.brand)    = b.name;

ALTER TABLE equipment DROP COLUMN category;
ALTER TABLE equipment DROP COLUMN brand;
CREATE INDEX IF NOT EXISTS equipment_category_idx ON equipment (category_id);
CREATE INDEX IF NOT EXISTS equipment_brand_idx ON equipment (brand_id);

INSERT INTO permissions (code, description, module, action, module_label, sort) VALUES
    ('categories.create', 'Crear categorías',    'categories', 'create', 'Categorías de equipos', 80),
    ('categories.list',   'Listar categorías',   'categories', 'list',   'Categorías de equipos', 81),
    ('categories.find',   'Buscar categorías',   'categories', 'find',   'Categorías de equipos', 82),
    ('categories.update', 'Editar categorías',   'categories', 'update', 'Categorías de equipos', 83),
    ('categories.delete', 'Eliminar categorías', 'categories', 'delete', 'Categorías de equipos', 84),
    ('brands.create',     'Crear marcas',        'brands',     'create', 'Marcas de equipos', 90),
    ('brands.list',       'Listar marcas',       'brands',     'list',   'Marcas de equipos', 91),
    ('brands.find',       'Buscar marcas',       'brands',     'find',   'Marcas de equipos', 92),
    ('brands.update',     'Editar marcas',       'brands',     'update', 'Marcas de equipos', 93),
    ('brands.delete',     'Eliminar marcas',     'brands',     'delete', 'Marcas de equipos', 94)
ON CONFLICT (code) DO UPDATE SET
    description = EXCLUDED.description, module = EXCLUDED.module, action = EXCLUDED.action,
    module_label = EXCLUDED.module_label, sort = EXCLUDED.sort;

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p WHERE r.name = 'Administrador'
ON CONFLICT DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p
WHERE r.name = 'Supervisor' AND p.code IN ('categories.list', 'categories.find', 'brands.list', 'brands.find')
ON CONFLICT DO NOTHING;
