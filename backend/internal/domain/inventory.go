package domain

import "time"

const (
	PermPositionsCreate = "positions.create"
	PermPositionsList   = "positions.list"
	PermPositionsFind   = "positions.find"
	PermPositionsUpdate = "positions.update"
	PermPositionsDelete = "positions.delete"

	PermEmployeesCreate = "employees.create"
	PermEmployeesList   = "employees.list"
	PermEmployeesFind   = "employees.find"
	PermEmployeesUpdate = "employees.update"
	PermEmployeesDelete = "employees.delete"

	PermEquipmentCreate = "equipment.create"
	PermEquipmentList   = "equipment.list"
	PermEquipmentFind   = "equipment.find"
	PermEquipmentUpdate = "equipment.update"
	PermEquipmentDelete = "equipment.delete"

	PermAssignmentsCreate = "assignments.create"
	PermAssignmentsList   = "assignments.list"
	PermAssignmentsFind   = "assignments.find"
	PermAssignmentsUpdate = "assignments.update"

	PermCategoriesCreate = "categories.create"
	PermCategoriesList   = "categories.list"
	PermCategoriesFind   = "categories.find"
	PermCategoriesUpdate = "categories.update"
	PermCategoriesDelete = "categories.delete"

	PermBrandsCreate = "brands.create"
	PermBrandsList   = "brands.list"
	PermBrandsFind   = "brands.find"
	PermBrandsUpdate = "brands.update"
	PermBrandsDelete = "brands.delete"
)

const (
	EquipAvailable   = "available"
	EquipAssigned    = "assigned"
	EquipMaintenance = "maintenance"
	EquipRetired     = "retired"
)

func ValidEquipmentStatus(s string) bool {
	switch s {
	case EquipAvailable, EquipAssigned, EquipMaintenance, EquipRetired:
		return true
	}
	return false
}

type Position struct {
	ID            int64
	Name          string
	Description   string
	Active        bool
	EmployeeCount int
	CreatedAt     time.Time
}

type Employee struct {
	ID                int64
	FirstName         string
	LastName          string
	Email             string
	Document          string
	Phone             string
	Department        string
	PositionID        int64
	PositionName      string
	HireDate          *time.Time
	Active            bool
	CreatedAt         time.Time
	ActiveAssignments int
}

func (e Employee) FullName() string { return e.FirstName + " " + e.LastName }

type Equipment struct {
	ID           int64
	Code         string
	Name         string
	CategoryID   int64 // 0 = sin categoría
	CategoryName string
	BrandID      int64 // 0 = sin marca
	BrandName    string
	Model        string
	Serial       string
	Status       string
	PurchaseDate *time.Time
	Cost         *float64
	Notes        string
	CreatedAt    time.Time
	HolderID     *int64  // funcionario que lo tiene hoy, si está asignado
	HolderName   *string
}

// Assignment es un movimiento del historial: entrega y, cuando termina, devolución.
type Assignment struct {
	ID            int64
	EquipmentID   int64
	EquipmentCode string
	EquipmentName string
	EmployeeID    int64
	EmployeeName  string
	AssignedAt    time.Time
	AssignedByID  int64
	AssignedBy    string
	DeliveryNotes string
	ReturnedAt    *time.Time
	ReceivedBy    string
	ReturnNotes   string
}

// CatalogItem es un valor parametrizable de equipos (categoría o marca).
type CatalogItem struct {
	ID             int64
	Name           string
	Active         bool
	EquipmentCount int
	CreatedAt      time.Time
}

type EmployeeFilter struct {
	Q          string
	PositionID int64
}

type EquipmentFilter struct {
	Q      string
	Status string
}

type AssignmentFilter struct {
	Q           string
	State       string // "active", "returned" o vacío (todas)
	EmployeeID  int64
	EquipmentID int64
}

type NameCount struct {
	Name  string
	Total int
}

type InventorySummary struct {
	TotalEquipment int
	Available      int
	Assigned       int
	Maintenance    int
	Retired        int
	Employees      int
	ByCategory     []NameCount
}
