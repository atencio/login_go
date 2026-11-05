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
}

type PasswordHasher interface {
	Hash(plain string) (string, error)
	Compare(hash, plain string) bool
}

type TokenManager interface {
	Generate(userID int64, remember bool) (token string, expiresAt time.Time, err error)
	Parse(token string) (userID int64, err error)
}
