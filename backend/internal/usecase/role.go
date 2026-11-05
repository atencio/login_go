package usecase

import (
	"context"
	"strings"

	"nexo/internal/domain"
)

type RoleService struct{ roles domain.RoleRepository }

func NewRoleService(r domain.RoleRepository) *RoleService { return &RoleService{roles: r} }

type RoleInput struct {
	Name        string
	Permissions []string
}

func (in *RoleInput) normalize() error {
	in.Name = strings.TrimSpace(in.Name)
	if in.Name == "" {
		return domain.ValidationError{Msg: "el nombre del perfil es obligatorio"}
	}
	return nil
}

func (s *RoleService) List(ctx context.Context) ([]domain.Role, error)  { return s.roles.List(ctx) }
func (s *RoleService) Get(ctx context.Context, id int64) (*domain.Role, error) {
	return s.roles.GetByID(ctx, id)
}
func (s *RoleService) Permissions(ctx context.Context) ([]domain.Permission, error) {
	return s.roles.ListPermissions(ctx)
}

func (s *RoleService) Create(ctx context.Context, in RoleInput) (*domain.Role, error) {
	if err := in.normalize(); err != nil {
		return nil, err
	}
	r := &domain.Role{Name: in.Name, Permissions: in.Permissions}
	if err := s.roles.Create(ctx, r); err != nil {
		return nil, err
	}
	return s.roles.GetByID(ctx, r.ID)
}

func (s *RoleService) Update(ctx context.Context, id int64, in RoleInput) (*domain.Role, error) {
	if err := in.normalize(); err != nil {
		return nil, err
	}
	current, err := s.roles.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if current.Name == domain.AdminRoleName {
		return nil, domain.ValidationError{Msg: "el perfil Administrador no se puede modificar"}
	}
	if err := s.roles.Update(ctx, &domain.Role{ID: id, Name: in.Name, Permissions: in.Permissions}); err != nil {
		return nil, err
	}
	return s.roles.GetByID(ctx, id)
}

func (s *RoleService) Delete(ctx context.Context, id int64) error {
	current, err := s.roles.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if current.Name == domain.AdminRoleName {
		return domain.ValidationError{Msg: "el perfil Administrador no se puede eliminar"}
	}
	return s.roles.Delete(ctx, id)
}
