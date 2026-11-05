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

## Perfiles y permisos

Usuario → Perfil (rol) → Permisos. Cada permiso es una celda de la matriz **módulo × acción**
(`users.create`, `users.list`, `users.find`, `users.update`, `users.delete`, `roles.*`, `reports.list`):
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
| GET | /api/reports/users-summary | reports.list |
