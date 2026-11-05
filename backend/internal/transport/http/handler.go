package http

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"nexo/internal/domain"
	"nexo/internal/usecase"
)

// Interfaces mínimas que necesitan los handlers (Segregación de interfaces).
type AuthUseCase interface {
	Login(ctx context.Context, email, password string, remember bool) (*usecase.LoginResult, error)
}

type UserUseCase interface {
	List(ctx context.Context) ([]domain.User, error)
	Search(ctx context.Context, q string) ([]domain.User, error)
	Get(ctx context.Context, id int64) (*domain.User, error)
	Create(ctx context.Context, in usecase.UserInput) (*domain.User, error)
	Update(ctx context.Context, id int64, in usecase.UserInput) (*domain.User, error)
	Delete(ctx context.Context, actorID, id int64) error
}

type RoleUseCase interface {
	List(ctx context.Context) ([]domain.Role, error)
	Get(ctx context.Context, id int64) (*domain.Role, error)
	Create(ctx context.Context, in usecase.RoleInput) (*domain.Role, error)
	Update(ctx context.Context, id int64, in usecase.RoleInput) (*domain.Role, error)
	Delete(ctx context.Context, id int64) error
	Permissions(ctx context.Context) ([]domain.Permission, error)
}

type ReportUseCase interface {
	UsersSummary(ctx context.Context) (*domain.ReportSummary, error)
}

type Handler struct {
	auth    AuthUseCase
	users   UserUseCase
	roles   RoleUseCase
	reports ReportUseCase
}

func NewHandler(a AuthUseCase, u UserUseCase, ro RoleUseCase, r ReportUseCase) *Handler {
	return &Handler{auth: a, users: u, roles: ro, reports: r}
}

// ---- DTOs ----

type userDTO struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	RoleID    int64     `json:"role_id"`
	Role      string    `json:"role"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
}

type profileDTO struct {
	User        userDTO  `json:"user"`
	Permissions []string `json:"permissions"`
}

type userRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	RoleID   int64  `json:"role_id"`
	Active   bool   `json:"active"`
}

func toUserDTO(u domain.User) userDTO {
	return userDTO{u.ID, u.Name, u.Email, u.RoleID, u.RoleName, u.Active, u.CreatedAt}
}

func toUserDTOs(users []domain.User) []userDTO {
	out := make([]userDTO, 0, len(users))
	for _, u := range users {
		out = append(out, toUserDTO(u))
	}
	return out
}

func toProfileDTO(p *domain.Profile) profileDTO {
	perms := p.Permissions
	if perms == nil {
		perms = []string{}
	}
	return profileDTO{User: toUserDTO(p.User), Permissions: perms}
}

func (r userRequest) toInput() usecase.UserInput {
	return usecase.UserInput{Name: r.Name, Email: r.Email, Password: r.Password, RoleID: r.RoleID, Active: r.Active}
}

func pathID(r *http.Request) (int64, error) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		return 0, domain.ValidationError{Msg: "identificador inválido"}
	}
	return id, nil
}

// ---- Auth ----

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Remember bool   `json:"remember"`
	}
	if err := decode(r, &in); err != nil {
		writeError(w, err)
		return
	}
	res, err := h.auth.Login(r.Context(), in.Email, in.Password, in.Remember)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"token": res.Token, "expires_at": res.ExpiresAt, "profile": toProfileDTO(res.Profile),
	})
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, toProfileDTO(profileFrom(r.Context())))
}

// ---- Usuarios ----

func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.users.List(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toUserDTOs(users))
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req userRequest
	if err := decode(r, &req); err != nil {
		writeError(w, err)
		return
	}
	u, err := h.users.Create(r.Context(), req.toInput())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, toUserDTO(*u))
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeError(w, err)
		return
	}
	var req userRequest
	if err := decode(r, &req); err != nil {
		writeError(w, err)
		return
	}
	u, err := h.users.Update(r.Context(), id, req.toInput())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toUserDTO(*u))
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeError(w, err)
		return
	}
	if err := h.users.Delete(r.Context(), profileFrom(r.Context()).User.ID, id); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) SearchUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.users.Search(r.Context(), r.URL.Query().Get("q"))
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toUserDTOs(users))
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeError(w, err)
		return
	}
	u, err := h.users.Get(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toUserDTO(*u))
}

// ---- Perfiles (roles) y permisos ----

type roleDTO struct {
	ID          int64    `json:"id"`
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
	UserCount   int      `json:"user_count"`
	Locked      bool     `json:"locked"` // el perfil Administrador no se edita ni se elimina
}

type roleRequest struct {
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
}

func toRoleDTO(r domain.Role) roleDTO {
	perms := r.Permissions
	if perms == nil {
		perms = []string{}
	}
	return roleDTO{r.ID, r.Name, perms, r.UserCount, r.Name == domain.AdminRoleName}
}

func (h *Handler) ListRoles(w http.ResponseWriter, r *http.Request) {
	roles, err := h.roles.List(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	out := make([]roleDTO, 0, len(roles))
	for _, ro := range roles {
		out = append(out, toRoleDTO(ro))
	}
	writeJSON(w, http.StatusOK, out)
}

// RoleOptions devuelve solo id y nombre, para los selectores de formularios.
func (h *Handler) RoleOptions(w http.ResponseWriter, r *http.Request) {
	roles, err := h.roles.List(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	type option struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	}
	out := make([]option, 0, len(roles))
	for _, ro := range roles {
		out = append(out, option{ro.ID, ro.Name})
	}
	writeJSON(w, http.StatusOK, out)
}

func (h *Handler) GetRole(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeError(w, err)
		return
	}
	ro, err := h.roles.Get(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toRoleDTO(*ro))
}

func (h *Handler) CreateRole(w http.ResponseWriter, r *http.Request) {
	var req roleRequest
	if err := decode(r, &req); err != nil {
		writeError(w, err)
		return
	}
	ro, err := h.roles.Create(r.Context(), usecase.RoleInput{Name: req.Name, Permissions: req.Permissions})
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, toRoleDTO(*ro))
}

func (h *Handler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeError(w, err)
		return
	}
	var req roleRequest
	if err := decode(r, &req); err != nil {
		writeError(w, err)
		return
	}
	ro, err := h.roles.Update(r.Context(), id, usecase.RoleInput{Name: req.Name, Permissions: req.Permissions})
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toRoleDTO(*ro))
}

func (h *Handler) DeleteRole(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeError(w, err)
		return
	}
	if err := h.roles.Delete(r.Context(), id); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ListPermissions devuelve el catálogo de módulos y acciones que forman la matriz.
func (h *Handler) ListPermissions(w http.ResponseWriter, r *http.Request) {
	perms, err := h.roles.Permissions(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	type permDTO struct {
		Code        string `json:"code"`
		Module      string `json:"module"`
		ModuleLabel string `json:"module_label"`
		Action      string `json:"action"`
	}
	out := make([]permDTO, 0, len(perms))
	for _, p := range perms {
		out = append(out, permDTO{p.Code, p.Module, p.ModuleLabel, p.Action})
	}
	writeJSON(w, http.StatusOK, out)
}

// ---- Reportes ----

func (h *Handler) UsersSummary(w http.ResponseWriter, r *http.Request) {
	s, err := h.reports.UsersSummary(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	type roleCount struct {
		Role  string `json:"role"`
		Total int    `json:"total"`
	}
	by := make([]roleCount, 0, len(s.ByRole))
	for _, rc := range s.ByRole {
		by = append(by, roleCount{rc.Role, rc.Total})
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"total_users": s.TotalUsers, "active_users": s.ActiveUsers, "by_role": by,
	})
}
