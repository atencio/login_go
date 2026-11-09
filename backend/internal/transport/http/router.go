package http

import (
	"net/http"

	"nexo/internal/domain"
)

func NewRouter(h *Handler, ih *InventoryHandler, categories, brands *CatalogHandler, m *Middleware, allowedOrigin string) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	// Públicas
	mux.HandleFunc("POST /api/auth/login", h.Login)

	// Cualquier usuario con sesión
	mux.Handle("GET /api/auth/me", m.Authenticate(http.HandlerFunc(h.Me)))
	mux.Handle("GET /api/permissions", m.Authenticate(http.HandlerFunc(h.ListPermissions)))

	// Usuarios
	mux.Handle("GET /api/users", m.Require(domain.PermUsersList, h.ListUsers))
	mux.Handle("GET /api/users/search", m.Require(domain.PermUsersFind, h.SearchUsers))
	mux.Handle("GET /api/users/{id}", m.Require(domain.PermUsersFind, h.GetUser))
	mux.Handle("POST /api/users", m.Require(domain.PermUsersCreate, h.CreateUser))
	mux.Handle("PUT /api/users/{id}", m.Require(domain.PermUsersUpdate, h.UpdateUser))
	mux.Handle("DELETE /api/users/{id}", m.Require(domain.PermUsersDelete, h.DeleteUser))

	// Perfiles (roles con su matriz de permisos)
	mux.Handle("GET /api/roles", m.Require(domain.PermRolesList, h.ListRoles))
	mux.Handle("GET /api/roles/options", m.RequireAny(
		[]string{domain.PermUsersCreate, domain.PermUsersUpdate, domain.PermRolesList}, h.RoleOptions))
	mux.Handle("GET /api/roles/{id}", m.Require(domain.PermRolesFind, h.GetRole))
	mux.Handle("POST /api/roles", m.Require(domain.PermRolesCreate, h.CreateRole))
	mux.Handle("PUT /api/roles/{id}", m.Require(domain.PermRolesUpdate, h.UpdateRole))
	mux.Handle("DELETE /api/roles/{id}", m.Require(domain.PermRolesDelete, h.DeleteRole))

	// Cargos
	mux.Handle("GET /api/positions", m.Require(domain.PermPositionsList, ih.ListPositions))
	mux.Handle("GET /api/positions/options", m.RequireAny(
		[]string{domain.PermEmployeesCreate, domain.PermEmployeesUpdate, domain.PermPositionsList}, ih.PositionOptions))
	mux.Handle("GET /api/positions/{id}", m.Require(domain.PermPositionsFind, ih.GetPosition))
	mux.Handle("POST /api/positions", m.Require(domain.PermPositionsCreate, ih.CreatePosition))
	mux.Handle("PUT /api/positions/{id}", m.Require(domain.PermPositionsUpdate, ih.UpdatePosition))
	mux.Handle("DELETE /api/positions/{id}", m.Require(domain.PermPositionsDelete, ih.DeletePosition))

	// Funcionarios
	mux.Handle("GET /api/employees", m.Require(domain.PermEmployeesList, ih.ListEmployees))
	mux.Handle("GET /api/employees/options", m.RequireAny(
		[]string{domain.PermAssignmentsCreate, domain.PermAssignmentsUpdate, domain.PermEmployeesList}, ih.EmployeeOptions))
	mux.Handle("GET /api/employees/{id}", m.Require(domain.PermEmployeesFind, ih.GetEmployee))
	mux.Handle("POST /api/employees", m.Require(domain.PermEmployeesCreate, ih.CreateEmployee))
	mux.Handle("PUT /api/employees/{id}", m.Require(domain.PermEmployeesUpdate, ih.UpdateEmployee))
	mux.Handle("DELETE /api/employees/{id}", m.Require(domain.PermEmployeesDelete, ih.DeleteEmployee))

	// Equipos de oficina
	mux.Handle("GET /api/equipment", m.Require(domain.PermEquipmentList, ih.ListEquipment))
	mux.Handle("GET /api/equipment/options", m.RequireAny(
		[]string{domain.PermAssignmentsCreate, domain.PermEquipmentList}, ih.EquipmentOptions))
	mux.Handle("GET /api/equipment/{id}", m.Require(domain.PermEquipmentFind, ih.GetEquipment))
	mux.Handle("POST /api/equipment", m.Require(domain.PermEquipmentCreate, ih.CreateEquipment))
	mux.Handle("PUT /api/equipment/{id}", m.Require(domain.PermEquipmentUpdate, ih.UpdateEquipment))
	mux.Handle("DELETE /api/equipment/{id}", m.Require(domain.PermEquipmentDelete, ih.DeleteEquipment))

	// Categorías y marcas de equipos (catálogos parametrizables)
	catalogRoutes(mux, m, "/api/categories", categories, catalogPerms{
		domain.PermCategoriesCreate, domain.PermCategoriesList, domain.PermCategoriesFind,
		domain.PermCategoriesUpdate, domain.PermCategoriesDelete})
	catalogRoutes(mux, m, "/api/brands", brands, catalogPerms{
		domain.PermBrandsCreate, domain.PermBrandsList, domain.PermBrandsFind,
		domain.PermBrandsUpdate, domain.PermBrandsDelete})

	// Asignaciones e historial (no se borran: son el registro de auditoría)
	mux.Handle("GET /api/assignments", m.Require(domain.PermAssignmentsList, ih.ListAssignments))
	mux.Handle("GET /api/assignments/{id}", m.Require(domain.PermAssignmentsFind, ih.GetAssignment))
	mux.Handle("POST /api/assignments", m.Require(domain.PermAssignmentsCreate, ih.CreateAssignment))
	mux.Handle("PUT /api/assignments/{id}", m.Require(domain.PermAssignmentsUpdate, ih.UpdateAssignment))
	mux.Handle("POST /api/assignments/{id}/return", m.Require(domain.PermAssignmentsUpdate, ih.ReturnAssignment))

	// Reportes
	mux.Handle("GET /api/reports/users-summary", m.Require(domain.PermReportsList, h.UsersSummary))
	mux.Handle("GET /api/reports/inventory-summary", m.Require(domain.PermReportsList, h.InventorySummary))

	return CORS(allowedOrigin, mux)
}

type catalogPerms struct{ create, list, find, update, delete string }

// catalogRoutes registra el CRUD completo de un catálogo bajo la ruta base indicada.
func catalogRoutes(mux *http.ServeMux, m *Middleware, base string, h *CatalogHandler, p catalogPerms) {
	mux.Handle("GET "+base, m.Require(p.list, h.List))
	mux.Handle("GET "+base+"/options", m.RequireAny(
		[]string{domain.PermEquipmentCreate, domain.PermEquipmentUpdate, p.list}, h.Options))
	mux.Handle("GET "+base+"/{id}", m.Require(p.find, h.Get))
	mux.Handle("POST "+base, m.Require(p.create, h.Create))
	mux.Handle("PUT "+base+"/{id}", m.Require(p.update, h.Update))
	mux.Handle("DELETE "+base+"/{id}", m.Require(p.delete, h.Delete))
}
