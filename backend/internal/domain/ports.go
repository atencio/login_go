package domain

import (
	"context"
	"time"
)

// Puertos: los casos de uso dependen de estas abstracciones (Inversión de dependencias).
// Las implementaciones concretas viven en repository/ y security/.

type UserRepository interface {
	Create(ctx context.Context, u *User) error
	GetByID(ctx context.Context, id int64) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	List(ctx context.Context) ([]User, error)
	Search(ctx context.Context, q string) ([]User, error)
	Update(ctx context.Context, u *User) error
	Delete(ctx context.Context, id int64) error
	Count(ctx context.Context) (int, error)
}

type RoleRepository interface {
	List(ctx context.Context) ([]Role, error)
	GetByID(ctx context.Context, id int64) (*Role, error)
	Create(ctx context.Context, r *Role) error
	Update(ctx context.Context, r *Role) error
	Delete(ctx context.Context, id int64) error
	ListPermissions(ctx context.Context) ([]Permission, error)
}

type ReportRepository interface {
	UsersSummary(ctx context.Context) (*ReportSummary, error)
	InventorySummary(ctx context.Context) (*InventorySummary, error)
}

type PasswordHasher interface {
	Hash(plain string) (string, error)
	Compare(hash, plain string) bool
}

type TokenManager interface {
	Generate(userID int64, remember bool) (token string, expiresAt time.Time, err error)
	Parse(token string) (userID int64, err error)
}

type PositionRepository interface {
	Create(ctx context.Context, p *Position) error
	GetByID(ctx context.Context, id int64) (*Position, error)
	List(ctx context.Context, q string) ([]Position, error)
	Update(ctx context.Context, p *Position) error
	Delete(ctx context.Context, id int64) error
}

type EmployeeRepository interface {
	Create(ctx context.Context, e *Employee) error
	GetByID(ctx context.Context, id int64) (*Employee, error)
	List(ctx context.Context, f EmployeeFilter) ([]Employee, error)
	Update(ctx context.Context, e *Employee) error
	Delete(ctx context.Context, id int64) error
}

type EquipmentRepository interface {
	// CreateMany registra una unidad por cada número de serie recibido (vacío = sin serie);
	// cada unidad recibe su propio código.
	CreateMany(ctx context.Context, base *Equipment, serials []string) ([]Equipment, error)
	GetByID(ctx context.Context, id int64) (*Equipment, error)
	List(ctx context.Context, f EquipmentFilter) ([]Equipment, error)
	Update(ctx context.Context, e *Equipment) error
	Delete(ctx context.Context, id int64) error
}

type AssignmentRepository interface {
	// Assign, Return y Transfer son atómicos: movimiento + estado del equipo en una transacción.
	Assign(ctx context.Context, a *Assignment) error
	GetByID(ctx context.Context, id int64) (*Assignment, error)
	List(ctx context.Context, f AssignmentFilter) ([]Assignment, error)
	UpdateNotes(ctx context.Context, id int64, notes string) error
	Return(ctx context.Context, id, userID int64, notes, equipmentStatus string) error
	Transfer(ctx context.Context, id, toEmployeeID, userID int64, notes string) (newID int64, err error)
}

// CatalogRepository sirve tanto para categorías como para marcas (misma forma, distinta tabla).
type CatalogRepository interface {
	Create(ctx context.Context, c *CatalogItem) error
	GetByID(ctx context.Context, id int64) (*CatalogItem, error)
	List(ctx context.Context, q string) ([]CatalogItem, error)
	Update(ctx context.Context, c *CatalogItem) error
	Delete(ctx context.Context, id int64) error
}
