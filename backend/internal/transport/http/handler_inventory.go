package http

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"nexo/internal/domain"
	"nexo/internal/usecase"
)

type PositionUseCase interface {
	List(ctx context.Context, q string) ([]domain.Position, error)
	Get(ctx context.Context, id int64) (*domain.Position, error)
	Create(ctx context.Context, in usecase.PositionInput) (*domain.Position, error)
	Update(ctx context.Context, id int64, in usecase.PositionInput) (*domain.Position, error)
	Delete(ctx context.Context, id int64) error
}

type EmployeeUseCase interface {
	List(ctx context.Context, f domain.EmployeeFilter) ([]domain.Employee, error)
	Get(ctx context.Context, id int64) (*domain.Employee, error)
	Create(ctx context.Context, in usecase.EmployeeInput) (*domain.Employee, error)
	Update(ctx context.Context, id int64, in usecase.EmployeeInput) (*domain.Employee, error)
	Delete(ctx context.Context, id int64) error
}

type EquipmentUseCase interface {
	List(ctx context.Context, f domain.EquipmentFilter) ([]domain.Equipment, error)
	Get(ctx context.Context, id int64) (*domain.Equipment, error)
	Create(ctx context.Context, in usecase.EquipmentInput, qty int, serials []string) ([]domain.Equipment, error)
	Update(ctx context.Context, id int64, in usecase.EquipmentInput) (*domain.Equipment, error)
	Delete(ctx context.Context, id int64) error
}

type AssignmentUseCase interface {
	List(ctx context.Context, f domain.AssignmentFilter) ([]domain.Assignment, error)
	Get(ctx context.Context, id int64) (*domain.Assignment, error)
	Assign(ctx context.Context, actorID, equipmentID, employeeID int64, notes string) (*domain.Assignment, error)
	Update(ctx context.Context, actorID, id, employeeID int64, notes string) (*domain.Assignment, error)
	Return(ctx context.Context, actorID, id int64, notes, status string) (*domain.Assignment, error)
}

// InventoryHandler agrupa los endpoints del inventario de activos.
type InventoryHandler struct {
	positions   PositionUseCase
	employees   EmployeeUseCase
	equipment   EquipmentUseCase
	assignments AssignmentUseCase
}

func NewInventoryHandler(p PositionUseCase, e EmployeeUseCase, q EquipmentUseCase, a AssignmentUseCase) *InventoryHandler {
	return &InventoryHandler{positions: p, employees: e, equipment: q, assignments: a}
}

// ---- utilidades ----

const dateLayout = "2006-01-02"

func fmtDate(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(dateLayout)
}

func parseDate(s string) (*time.Time, error) {
	if s == "" {
		return nil, nil
	}
	t, err := time.Parse(dateLayout, s)
	if err != nil {
		return nil, domain.ValidationError{Msg: "la fecha no es válida (usa AAAA-MM-DD)"}
	}
	return &t, nil
}

func queryInt(r *http.Request, key string) int64 {
	n, _ := strconv.ParseInt(r.URL.Query().Get(key), 10, 64)
	return n
}

// requireFind exige el permiso "Encontrar" cuando la petición usa búsqueda o filtros.
func requireFind(r *http.Request, code string, used bool) error {
	if used && !profileFrom(r.Context()).Has(code) {
		return domain.ErrForbidden
	}
	return nil
}

func actorID(r *http.Request) int64 { return profileFrom(r.Context()).User.ID }

// ---- Cargos ----

type positionDTO struct {
	ID            int64     `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	Active        bool      `json:"active"`
	EmployeeCount int       `json:"employee_count"`
	CreatedAt     time.Time `json:"created_at"`
}

type positionRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Active      bool   `json:"active"`
}

func toPositionDTO(p domain.Position) positionDTO {
	return positionDTO{p.ID, p.Name, p.Description, p.Active, p.EmployeeCount, p.CreatedAt}
}

func (req positionRequest) input() usecase.PositionInput {
	return usecase.PositionInput{Name: req.Name, Description: req.Description, Active: req.Active}
}

func (h *InventoryHandler) ListPositions(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if err := requireFind(r, domain.PermPositionsFind, q != ""); err != nil {
		writeError(w, err)
		return
	}
	list, err := h.positions.List(r.Context(), q)
	if err != nil {
		writeError(w, err)
		return
	}
	out := make([]positionDTO, 0, len(list))
	for _, p := range list {
		out = append(out, toPositionDTO(p))
	}
	writeJSON(w, http.StatusOK, out)
}

func (h *InventoryHandler) PositionOptions(w http.ResponseWriter, r *http.Request) {
	list, err := h.positions.List(r.Context(), "")
	if err != nil {
		writeError(w, err)
		return
	}
	type option struct {
		ID     int64  `json:"id"`
		Name   string `json:"name"`
		Active bool   `json:"active"`
	}
	out := make([]option, 0, len(list))
	for _, p := range list {
		out = append(out, option{p.ID, p.Name, p.Active})
	}
	writeJSON(w, http.StatusOK, out)
}

func (h *InventoryHandler) GetPosition(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeError(w, err)
		return
	}
	p, err := h.positions.Get(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toPositionDTO(*p))
}

func (h *InventoryHandler) CreatePosition(w http.ResponseWriter, r *http.Request) {
	var req positionRequest
	if err := decode(r, &req); err != nil {
		writeError(w, err)
		return
	}
	p, err := h.positions.Create(r.Context(), req.input())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, toPositionDTO(*p))
}

func (h *InventoryHandler) UpdatePosition(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeError(w, err)
		return
	}
	var req positionRequest
	if err := decode(r, &req); err != nil {
		writeError(w, err)
		return
	}
	p, err := h.positions.Update(r.Context(), id, req.input())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toPositionDTO(*p))
}

func (h *InventoryHandler) DeletePosition(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeError(w, err)
		return
	}
	if err := h.positions.Delete(r.Context(), id); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ---- Funcionarios ----

type employeeDTO struct {
	ID                int64     `json:"id"`
	FirstName         string    `json:"first_name"`
	LastName          string    `json:"last_name"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	Document          string    `json:"document"`
	Phone             string    `json:"phone"`
	Department        string    `json:"department"`
	PositionID        int64     `json:"position_id"`
	Position          string    `json:"position"`
	HireDate          string    `json:"hire_date"`
	Active            bool      `json:"active"`
	ActiveAssignments int       `json:"active_assignments"`
	CreatedAt         time.Time `json:"created_at"`
}

type employeeRequest struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Email      string `json:"email"`
	Document   string `json:"document"`
	Phone      string `json:"phone"`
	Department string `json:"department"`
	PositionID int64  `json:"position_id"`
	HireDate   string `json:"hire_date"`
	Active     bool   `json:"active"`
}

func toEmployeeDTO(e domain.Employee) employeeDTO {
	return employeeDTO{e.ID, e.FirstName, e.LastName, e.FullName(), e.Email, e.Document, e.Phone, e.Department,
		e.PositionID, e.PositionName, fmtDate(e.HireDate), e.Active, e.ActiveAssignments, e.CreatedAt}
}

func (req employeeRequest) input() (usecase.EmployeeInput, error) {
	hire, err := parseDate(req.HireDate)
	if err != nil {
		return usecase.EmployeeInput{}, err
	}
	return usecase.EmployeeInput{
		FirstName: req.FirstName, LastName: req.LastName, Email: req.Email, Document: req.Document,
		Phone: req.Phone, Department: req.Department, PositionID: req.PositionID, HireDate: hire, Active: req.Active,
	}, nil
}

func (h *InventoryHandler) ListEmployees(w http.ResponseWriter, r *http.Request) {
	f := domain.EmployeeFilter{Q: r.URL.Query().Get("q"), PositionID: queryInt(r, "position_id")}
	if err := requireFind(r, domain.PermEmployeesFind, f.Q != "" || f.PositionID != 0); err != nil {
		writeError(w, err)
		return
	}
	list, err := h.employees.List(r.Context(), f)
	if err != nil {
		writeError(w, err)
		return
	}
	out := make([]employeeDTO, 0, len(list))
	for _, e := range list {
		out = append(out, toEmployeeDTO(e))
	}
	writeJSON(w, http.StatusOK, out)
}

// EmployeeOptions lista los funcionarios activos para los selectores de asignación.
func (h *InventoryHandler) EmployeeOptions(w http.ResponseWriter, r *http.Request) {
	list, err := h.employees.List(r.Context(), domain.EmployeeFilter{})
	if err != nil {
		writeError(w, err)
		return
	}
	type option struct {
		ID       int64  `json:"id"`
		Name     string `json:"name"`
		Position string `json:"position"`
	}
	out := []option{}
	for _, e := range list {
		if e.Active {
			out = append(out, option{e.ID, e.FullName(), e.PositionName})
		}
	}
	writeJSON(w, http.StatusOK, out)
}

func (h *InventoryHandler) GetEmployee(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeError(w, err)
		return
	}
	e, err := h.employees.Get(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toEmployeeDTO(*e))
}

func (h *InventoryHandler) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	var req employeeRequest
	if err := decode(r, &req); err != nil {
		writeError(w, err)
		return
	}
	in, err := req.input()
	if err != nil {
		writeError(w, err)
		return
	}
	e, err := h.employees.Create(r.Context(), in)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, toEmployeeDTO(*e))
}

func (h *InventoryHandler) UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeError(w, err)
		return
	}
	var req employeeRequest
	if err := decode(r, &req); err != nil {
		writeError(w, err)
		return
	}
	in, err := req.input()
	if err != nil {
		writeError(w, err)
		return
	}
	e, err := h.employees.Update(r.Context(), id, in)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toEmployeeDTO(*e))
}

func (h *InventoryHandler) DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeError(w, err)
		return
	}
	if err := h.employees.Delete(r.Context(), id); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ---- Equipos de oficina ----

type equipmentDTO struct {
	ID           int64     `json:"id"`
	Code         string    `json:"code"`
	Name         string    `json:"name"`
	CategoryID   int64     `json:"category_id"`
	Category     string    `json:"category"`
	BrandID      int64     `json:"brand_id"`
	Brand        string    `json:"brand"`
	Model        string    `json:"model"`
	Serial       string    `json:"serial_number"`
	Status       string    `json:"status"`
	PurchaseDate string    `json:"purchase_date"`
	Cost         *float64  `json:"cost"`
	Notes        string    `json:"notes"`
	HolderID     *int64    `json:"holder_id"`
	Holder       string    `json:"holder"`
	CreatedAt    time.Time `json:"created_at"`
}

type equipmentRequest struct {
	Name         string   `json:"name"`
	CategoryID   int64    `json:"category_id"`
	BrandID      int64    `json:"brand_id"`
	Model        string   `json:"model"`
	Serial       string   `json:"serial_number"` // al editar una unidad
	Serials      []string `json:"serials"`       // al crear: una serie por unidad
	Status       string   `json:"status"`
	PurchaseDate string   `json:"purchase_date"`
	Cost         *float64 `json:"cost"`
	Notes        string   `json:"notes"`
	Quantity     int      `json:"quantity"` // solo se usa al crear
}

func toEquipmentDTO(q domain.Equipment) equipmentDTO {
	holder := ""
	if q.HolderName != nil {
		holder = *q.HolderName
	}
	return equipmentDTO{
		ID: q.ID, Code: q.Code, Name: q.Name, CategoryID: q.CategoryID, Category: q.CategoryName,
		BrandID: q.BrandID, Brand: q.BrandName, Model: q.Model, Serial: q.Serial, Status: q.Status,
		PurchaseDate: fmtDate(q.PurchaseDate), Cost: q.Cost, Notes: q.Notes, HolderID: q.HolderID,
		Holder: holder, CreatedAt: q.CreatedAt,
	}
}

func (req equipmentRequest) input() (usecase.EquipmentInput, error) {
	pd, err := parseDate(req.PurchaseDate)
	if err != nil {
		return usecase.EquipmentInput{}, err
	}
	return usecase.EquipmentInput{
		Name: req.Name, CategoryID: req.CategoryID, BrandID: req.BrandID, Model: req.Model, Serial: req.Serial,
		Status: req.Status, PurchaseDate: pd, Cost: req.Cost, Notes: req.Notes,
	}, nil
}

func toEquipmentDTOs(list []domain.Equipment) []equipmentDTO {
	out := make([]equipmentDTO, 0, len(list))
	for _, q := range list {
		out = append(out, toEquipmentDTO(q))
	}
	return out
}

func (h *InventoryHandler) ListEquipment(w http.ResponseWriter, r *http.Request) {
	f := domain.EquipmentFilter{Q: r.URL.Query().Get("q"), Status: r.URL.Query().Get("status")}
	if err := requireFind(r, domain.PermEquipmentFind, f.Q != "" || f.Status != ""); err != nil {
		writeError(w, err)
		return
	}
	list, err := h.equipment.List(r.Context(), f)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toEquipmentDTOs(list))
}

// EquipmentOptions alimenta el selector de la asignación (normalmente ?status=available).
func (h *InventoryHandler) EquipmentOptions(w http.ResponseWriter, r *http.Request) {
	list, err := h.equipment.List(r.Context(), domain.EquipmentFilter{Status: r.URL.Query().Get("status")})
	if err != nil {
		writeError(w, err)
		return
	}
	type option struct {
		ID    int64  `json:"id"`
		Code  string `json:"code"`
		Name  string `json:"name"`
		Brand string `json:"brand"`
		Model string `json:"model"`
	}
	out := make([]option, 0, len(list))
	for _, q := range list {
		out = append(out, option{q.ID, q.Code, q.Name, q.BrandName, q.Model})
	}
	writeJSON(w, http.StatusOK, out)
}

func (h *InventoryHandler) GetEquipment(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeError(w, err)
		return
	}
	q, err := h.equipment.Get(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toEquipmentDTO(*q))
}

// CreateEquipment responde con la lista de unidades creadas (una por cada unidad de la cantidad).
func (h *InventoryHandler) CreateEquipment(w http.ResponseWriter, r *http.Request) {
	var req equipmentRequest
	if err := decode(r, &req); err != nil {
		writeError(w, err)
		return
	}
	in, err := req.input()
	if err != nil {
		writeError(w, err)
		return
	}
	list, err := h.equipment.Create(r.Context(), in, req.Quantity, req.Serials)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, toEquipmentDTOs(list))
}

func (h *InventoryHandler) UpdateEquipment(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeError(w, err)
		return
	}
	var req equipmentRequest
	if err := decode(r, &req); err != nil {
		writeError(w, err)
		return
	}
	in, err := req.input()
	if err != nil {
		writeError(w, err)
		return
	}
	q, err := h.equipment.Update(r.Context(), id, in)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toEquipmentDTO(*q))
}

func (h *InventoryHandler) DeleteEquipment(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeError(w, err)
		return
	}
	if err := h.equipment.Delete(r.Context(), id); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ---- Asignaciones e historial ----

type assignmentDTO struct {
	ID            int64      `json:"id"`
	EquipmentID   int64      `json:"equipment_id"`
	EquipmentCode string     `json:"equipment_code"`
	EquipmentName string     `json:"equipment_name"`
	EmployeeID    int64      `json:"employee_id"`
	Employee      string     `json:"employee"`
	AssignedAt    time.Time  `json:"assigned_at"`
	AssignedBy    string     `json:"assigned_by"`
	DeliveryNotes string     `json:"delivery_notes"`
	ReturnedAt    *time.Time `json:"returned_at"`
	ReceivedBy    string     `json:"received_by"`
	ReturnNotes   string     `json:"return_notes"`
	Active        bool       `json:"active"`
}

func toAssignmentDTO(a domain.Assignment) assignmentDTO {
	return assignmentDTO{a.ID, a.EquipmentID, a.EquipmentCode, a.EquipmentName, a.EmployeeID, a.EmployeeName,
		a.AssignedAt, a.AssignedBy, a.DeliveryNotes, a.ReturnedAt, a.ReceivedBy, a.ReturnNotes, a.ReturnedAt == nil}
}

func (h *InventoryHandler) ListAssignments(w http.ResponseWriter, r *http.Request) {
	f := domain.AssignmentFilter{
		Q: r.URL.Query().Get("q"), State: r.URL.Query().Get("state"),
		EmployeeID: queryInt(r, "employee_id"), EquipmentID: queryInt(r, "equipment_id"),
	}
	if err := requireFind(r, domain.PermAssignmentsFind, f.Q != ""); err != nil {
		writeError(w, err)
		return
	}
	list, err := h.assignments.List(r.Context(), f)
	if err != nil {
		writeError(w, err)
		return
	}
	out := make([]assignmentDTO, 0, len(list))
	for _, a := range list {
		out = append(out, toAssignmentDTO(a))
	}
	writeJSON(w, http.StatusOK, out)
}

func (h *InventoryHandler) GetAssignment(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeError(w, err)
		return
	}
	a, err := h.assignments.Get(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toAssignmentDTO(*a))
}

func (h *InventoryHandler) CreateAssignment(w http.ResponseWriter, r *http.Request) {
	var req struct {
		EquipmentID int64  `json:"equipment_id"`
		EmployeeID  int64  `json:"employee_id"`
		Notes       string `json:"notes"`
	}
	if err := decode(r, &req); err != nil {
		writeError(w, err)
		return
	}
	a, err := h.assignments.Assign(r.Context(), actorID(r), req.EquipmentID, req.EmployeeID, req.Notes)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, toAssignmentDTO(*a))
}

func (h *InventoryHandler) UpdateAssignment(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeError(w, err)
		return
	}
	var req struct {
		EmployeeID int64  `json:"employee_id"`
		Notes      string `json:"notes"`
	}
	if err := decode(r, &req); err != nil {
		writeError(w, err)
		return
	}
	a, err := h.assignments.Update(r.Context(), actorID(r), id, req.EmployeeID, req.Notes)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toAssignmentDTO(*a))
}

func (h *InventoryHandler) ReturnAssignment(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeError(w, err)
		return
	}
	var req struct {
		Notes  string `json:"notes"`
		Status string `json:"status"`
	}
	if err := decode(r, &req); err != nil {
		writeError(w, err)
		return
	}
	a, err := h.assignments.Return(r.Context(), actorID(r), id, req.Notes, req.Status)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toAssignmentDTO(*a))
}
