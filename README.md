# Nexo — Go + PostgreSQL + React

Panel con login (JWT), CRUD de usuarios, roles con permisos, perfil en modal y reportes.
Las alertas son modales tipo SweetAlert.

## Arrancar

```bash
# 1) Base de datos
docker compose up -d db

# 2) Backend (genera go.sum la primera vez)
cd backend
go mod tidy
go run ./cmd/api          # http://localhost:8080

# 3) Frontend
cd ../frontend
npm install
npm run dev               # http://localhost:5173
```

Usuario inicial (se crea solo si la tabla está vacía): `admin@nexo.local` / `Admin12345`.
Cámbialo con las variables `ADMIN_EMAIL` y `ADMIN_PASSWORD` (ver `backend/.env.example`).
También puedes levantar todo el backend con `docker compose up --build`.

## Arquitectura del backend (capas limpias / hexagonal ligera)

```
cmd/api                  → composición: único lugar donde se conectan implementaciones
internal/domain          → entidades, errores y puertos (interfaces). Sin dependencias externas
internal/usecase         → lógica de aplicación (Auth, Users, Reports). Dependen solo de puertos
internal/repository      → adaptadores de salida: PostgreSQL (pgx)
internal/security        → adaptadores de salida: bcrypt y JWT
internal/transport/http  → adaptador de entrada: handlers, middleware, router, DTOs
migrations               → SQL embebido, se aplica al arrancar
```

SOLID: cada capa tiene una sola responsabilidad (S), se extiende agregando casos de uso o
permisos sin tocar los existentes (O), los repositorios son intercambiables por sus interfaces (L),
los handlers piden interfaces pequeñas (I) y todo depende de abstracciones del dominio (D).

## Interfaz

- **Tablas fijas:** en todas las tablas solo se desplazan las filas; el encabezado, el buscador y el paginador no se mueven.
- **Paginación** en la parte inferior ("Mostrando 1–10 de 33 registros") y **"Mostrar N registros"** (10, 25, 50 o 100) junto al buscador; la elección vale para todas las tablas.
- **Tamaño de fuente** (100 %, 90 % o 80 %) en la barra superior: escala texto, espacios e iconos. Las dos preferencias se recuerdan en el navegador.
- La paginación se hace en el navegador sobre la lista que devuelve la API; si algún día una tabla supera varios miles de filas conviene paginar también en el servidor.

## Inventario de asignación de activos

Menú **Inventario de activos**:

- **Cargos**: CRUD de puestos (no se elimina un cargo con funcionarios).
- **Funcionarios**: CRUD con nombre, apellido, correo, documento, teléfono, cargo, departamento, fecha de ingreso y estado. Botón de historial por funcionario.
- **Categorías** y **Marcas**: catálogos parametrizables (CRUD) que alimentan las listas del formulario de equipos. Vienen cargados con categorías de oficina y marcas habituales; no se elimina una con equipos asociados (se desactiva y deja de ofrecerse).
- **Equipos de oficina**: CRUD de activos con categoría y marca elegidas de esos catálogos. Al **crear** se digita la **cantidad** N: se generan N unidades con su código (`EQ-00001`…) y el formulario habilita **N campos de serie**, uno por unidad (la serie es única por equipo; pueden quedar vacías). Después cada unidad se edita por separado (no hay cantidad al editar).
  Estados: disponible, asignado, en mantenimiento y dado de baja. Muestra quién tiene cada equipo.
- **Asignaciones**: pestaña *Asignaciones activas* (nueva asignación, **actualizar** = cambiar notas o transferir a otro funcionario, **devolver** dejando el equipo disponible, en mantenimiento o dado de baja) y pestaña *Historial* con todos los movimientos.

Reglas que garantiza el backend:
- Un equipo tiene una sola asignación activa (índice único parcial en la base de datos).
- Solo se asigna un equipo *disponible* a un funcionario *activo*; asignar, transferir y devolver son transacciones.
- El estado "asignado" solo lo cambian las asignaciones; el historial no se borra.
- Un funcionario o equipo con historial no se elimina: se desactiva o se da de baja.

## Perfiles y permisos

Usuario → Perfil (rol) → Permisos. Cada permiso es una celda de la matriz **módulo × acción**
(`users.*`, `roles.*`, `positions.*`, `employees.*`, `equipment.*`, `assignments.*`, `categories.*`, `brands.*`, `reports.list`):
Crear, Listar, Encontrar (buscar / ver detalle), Editar y Eliminar.

- La pantalla **Perfiles** crea, edita, consulta y elimina perfiles con la matriz de interruptores.
- El perfil **Administrador** está protegido: siempre tiene todo y no se edita ni se elimina.
- El middleware `Require("código", handler)` protege cada ruta y los permisos se leen de la BD en
  cada petición, así que un cambio de perfil se aplica de inmediato.
- El frontend usa `can("código")` para ocultar menús y botones; **Mi perfil** (menú del usuario)
  muestra la matriz de su perfil en modo lectura.
- Las migraciones se registran en `schema_migrations` y cada una corre una sola vez.

### Agregar un módulo nuevo (por ejemplo Candidatos)
1. Backend: constantes en `domain/entities.go`, migración `003_*.sql` que inserte sus permisos en
   `permissions` (`code`, `module`, `action`, `module_label`, `sort`), caso de uso, repositorio,
   handler y rutas con `m.Require(...)`.
2. Frontend: la página y el ítem en `GROUPS` del layout. La fila aparece sola en la matriz de
   perfiles porque el catálogo sale de `GET /api/permissions`.

## API

| Método | Ruta | Permiso |
|---|---|---|
| POST | /api/auth/login | pública (`remember: true` → sesión de 30 días) |
| GET | /api/auth/me, /api/permissions | sesión |
| GET | /api/users | users.list |
| GET | /api/users/search?q=, /api/users/{id} | users.find |
| POST / PUT / DELETE | /api/users, /api/users/{id} | users.create / update / delete |
| GET | /api/roles | roles.list |
| GET | /api/roles/options | users.create, users.update o roles.list |
| GET | /api/roles/{id} | roles.find |
| POST / PUT / DELETE | /api/roles, /api/roles/{id} | roles.create / update / delete |
| GET | /api/positions, /api/employees, /api/equipment | positions / employees / equipment `.list` (con `q` o filtros: `.find`) |
| GET / POST / PUT / DELETE | /api/categories, /api/brands (+ `/{id}`) | categories / brands `.list` / `.create` / `.update` / `.delete` (`q` y `/{id}`: `.find`) |
| GET | /api/positions/options, /api/employees/options, /api/equipment/options, /api/categories/options, /api/brands/options | para selectores de formularios |
| POST | /api/equipment | body con `quantity` y `serials` (una serie por unidad) |
| GET / POST / PUT / DELETE | /api/positions, /api/employees, /api/equipment (+ `/{id}`) | `.find` / `.create` / `.update` / `.delete` |
| GET | /api/assignments?state=&q=&employee_id=&equipment_id= | assignments.list (`q` requiere `.find`) |
| POST | /api/assignments | assignments.create |
| PUT | /api/assignments/{id} (notas o transferencia) | assignments.update |
| POST | /api/assignments/{id}/return | assignments.update |
| GET | /api/reports/users-summary, /api/reports/inventory-summary | reports.list |
