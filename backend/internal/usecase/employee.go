package usecase

import (
	"context"
	"strings"
	"time"

	"nexo/internal/domain"
)

type EmployeeService struct{ repo domain.EmployeeRepository }

func NewEmployeeService(r domain.EmployeeRepository) *EmployeeService { return &EmployeeService{repo: r} }

type EmployeeInput struct {
	FirstName  string
	LastName   string
	Email      string
	Document   string
	Phone      string
	Department string
	PositionID int64
	HireDate   *time.Time
	Active     bool
}

func (in *EmployeeInput) normalize() error {
	in.FirstName = strings.TrimSpace(in.FirstName)
	in.LastName = strings.TrimSpace(in.LastName)
	in.Email = strings.ToLower(strings.TrimSpace(in.Email))
	in.Document = strings.TrimSpace(in.Document)
	in.Phone = strings.TrimSpace(in.Phone)
	in.Department = strings.TrimSpace(in.Department)
	switch {
	case in.FirstName == "":
		return domain.ValidationError{Msg: "el nombre es obligatorio"}
	case in.LastName == "":
		return domain.ValidationError{Msg: "el apellido es obligatorio"}
	case !strings.Contains(in.Email, "@"):
		return domain.ValidationError{Msg: "el correo no es válido"}
	case in.PositionID == 0:
		return domain.ValidationError{Msg: "selecciona un cargo"}
	}
	return nil
}

func (in EmployeeInput) toEntity(id int64) *domain.Employee {
	return &domain.Employee{
		ID: id, FirstName: in.FirstName, LastName: in.LastName, Email: in.Email, Document: in.Document,
		Phone: in.Phone, Department: in.Department, PositionID: in.PositionID, HireDate: in.HireDate, Active: in.Active,
	}
}

func (s *EmployeeService) List(ctx context.Context, f domain.EmployeeFilter) ([]domain.Employee, error) {
	f.Q = strings.TrimSpace(f.Q)
	return s.repo.List(ctx, f)
}

func (s *EmployeeService) Get(ctx context.Context, id int64) (*domain.Employee, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *EmployeeService) Create(ctx context.Context, in EmployeeInput) (*domain.Employee, error) {
	if err := in.normalize(); err != nil {
		return nil, err
	}
	e := in.toEntity(0)
	if err := s.repo.Create(ctx, e); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, e.ID)
}

func (s *EmployeeService) Update(ctx context.Context, id int64, in EmployeeInput) (*domain.Employee, error) {
	if err := in.normalize(); err != nil {
		return nil, err
	}
	if err := s.repo.Update(ctx, in.toEntity(id)); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, id)
}

func (s *EmployeeService) Delete(ctx context.Context, id int64) error { return s.repo.Delete(ctx, id) }
