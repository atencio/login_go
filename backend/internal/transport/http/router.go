package http

import (
	"net/http"

	"nexo/internal/domain"
)

func NewRouter(h *Handler, m *Middleware, allowedOrigin string) http.Handler {
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

	// Reportes
	mux.Handle("GET /api/reports/users-summary", m.Require(domain.PermReportsList, h.UsersSummary))

	return CORS(allowedOrigin, mux)
}
