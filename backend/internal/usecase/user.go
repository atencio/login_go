package usecase

import (
	"context"
	"strings"

	"nexo/internal/domain"
)

type UserService struct {
	users  domain.UserRepository
	roles  domain.RoleRepository
	hasher domain.PasswordHasher
}

func NewUserService(u domain.UserRepository, r domain.RoleRepository, h domain.PasswordHasher) *UserService {
	return &UserService{users: u, roles: r, hasher: h}
}

type UserInput struct {
	Name     string
	Email    string
	Password string // obligatoria al crear, opcional al editar
	RoleID   int64
	Active   bool
}

func (in *UserInput) normalize(passwordRequired bool) error {
	in.Name = strings.TrimSpace(in.Name)
	in.Email = strings.ToLower(strings.TrimSpace(in.Email))
	switch {
	case in.Name == "":
		return domain.ValidationError{Msg: "el nombre es obligatorio"}
	case !strings.Contains(in.Email, "@"):
		return domain.ValidationError{Msg: "el correo no es válido"}
	case in.RoleID == 0:
		return domain.ValidationError{Msg: "selecciona un rol"}
	case (passwordRequired || in.Password != "") && len(in.Password) < 8:
		return domain.ValidationError{Msg: "la contraseña debe tener al menos 8 caracteres"}
	}
	return nil
}

func (s *UserService) List(ctx context.Context) ([]domain.User, error) { return s.users.List(ctx) }

func (s *UserService) Search(ctx context.Context, q string) ([]domain.User, error) {
	return s.users.Search(ctx, strings.TrimSpace(q))
}

func (s *UserService) Get(ctx context.Context, id int64) (*domain.User, error) {
	return s.users.GetByID(ctx, id)
}

func (s *UserService) Roles(ctx context.Context) ([]domain.Role, error) { return s.roles.List(ctx) }

func (s *UserService) Create(ctx context.Context, in UserInput) (*domain.User, error) {
	if err := in.normalize(true); err != nil {
		return nil, err
	}
	hash, err := s.hasher.Hash(in.Password)
	if err != nil {
		return nil, err
	}
	u := &domain.User{Name: in.Name, Email: in.Email, PasswordHash: hash, RoleID: in.RoleID, Active: in.Active}
	if err := s.users.Create(ctx, u); err != nil {
		return nil, err
	}
	return s.users.GetByID(ctx, u.ID)
}

func (s *UserService) Update(ctx context.Context, id int64, in UserInput) (*domain.User, error) {
	if err := in.normalize(false); err != nil {
		return nil, err
	}
	u, err := s.users.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	u.Name, u.Email, u.RoleID, u.Active = in.Name, in.Email, in.RoleID, in.Active
	if in.Password != "" {
		if u.PasswordHash, err = s.hasher.Hash(in.Password); err != nil {
			return nil, err
		}
	}
	if err := s.users.Update(ctx, u); err != nil {
		return nil, err
	}
	return s.users.GetByID(ctx, id)
}

func (s *UserService) Delete(ctx context.Context, actorID, id int64) error {
	if actorID == id {
		return domain.ValidationError{Msg: "no puedes eliminar tu propio usuario"}
	}
	return s.users.Delete(ctx, id)
}

// EnsureAdmin crea el primer administrador si la base de datos está vacía.
func (s *UserService) EnsureAdmin(ctx context.Context, name, email, password string) (bool, error) {
	n, err := s.users.Count(ctx)
	if err != nil || n > 0 {
		return false, err
	}
	roles, err := s.roles.List(ctx)
	if err != nil {
		return false, err
	}
	for _, r := range roles {
		if r.Name == domain.AdminRoleName {
			_, err = s.Create(ctx, UserInput{Name: name, Email: email, Password: password, RoleID: r.ID, Active: true})
			return err == nil, err
		}
	}
	return false, domain.ErrNotFound
}
