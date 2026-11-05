// Package domain contiene las entidades y reglas del negocio. No depende de nada externo.
package domain

import (
	"errors"
	"time"
)

// Códigos de permisos. Cada módulo nuevo (p. ej. un reporte) agrega los suyos aquí
// y los inserta en la tabla permissions mediante una migración.
const (
	PermUsersCreate = "users.create"
	PermUsersList   = "users.list"
	PermUsersFind   = "users.find"
	PermUsersUpdate = "users.update"
	PermUsersDelete = "users.delete"

	PermRolesCreate = "roles.create"
	PermRolesList   = "roles.list"
	PermRolesFind   = "roles.find"
	PermRolesUpdate = "roles.update"
	PermRolesDelete = "roles.delete"

	PermReportsList = "reports.list"
)

const AdminRoleName = "Administrador"

type User struct {
	ID           int64
	Name         string
	Email        string
	PasswordHash string
	RoleID       int64
	RoleName     string
	Active       bool
	CreatedAt    time.Time
}

type Role struct {
	ID          int64
	Name        string
	Permissions []string
	UserCount   int
}

// Permission es una celda de la matriz: un módulo y una acción.
type Permission struct {
	Code        string
	Module      string
	ModuleLabel string
	Action      string
}

// Profile es el usuario autenticado junto con los permisos de su rol.
type Profile struct {
	User        User
	Permissions []string
}

func (p Profile) Has(code string) bool {
	for _, c := range p.Permissions {
		if c == code {
			return true
		}
	}
	return false
}

type RoleCount struct {
	Role  string
	Total int
}

type ReportSummary struct {
	TotalUsers  int
	ActiveUsers int
	ByRole      []RoleCount
}

var (
	ErrNotFound           = errors.New("recurso no encontrado")
	ErrInvalidCredentials = errors.New("correo o contraseña incorrectos")
	ErrEmailTaken         = errors.New("ya existe un usuario con ese correo")
	ErrForbidden          = errors.New("no tienes permiso para realizar esta acción")
	ErrUnauthorized       = errors.New("la sesión no es válida o expiró")
	ErrInactive           = errors.New("tu usuario está inactivo, contacta al administrador")
)

type ValidationError struct{ Msg string }

func (e ValidationError) Error() string { return e.Msg }
