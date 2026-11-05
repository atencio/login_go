// Package usecase contiene la lógica de aplicación. Solo conoce interfaces del dominio.
package usecase

import (
	"context"
	"strings"
	"time"

	"nexo/internal/domain"
)

type AuthService struct {
	users  domain.UserRepository
	roles  domain.RoleRepository
	hasher domain.PasswordHasher
	tokens domain.TokenManager
}

func NewAuthService(u domain.UserRepository, r domain.RoleRepository, h domain.PasswordHasher, t domain.TokenManager) *AuthService {
	return &AuthService{users: u, roles: r, hasher: h, tokens: t}
}

type LoginResult struct {
	Token     string
	ExpiresAt time.Time
	Profile   *domain.Profile
}

func (s *AuthService) Login(ctx context.Context, email, password string, remember bool) (*LoginResult, error) {
	u, err := s.users.GetByEmail(ctx, strings.ToLower(strings.TrimSpace(email)))
	if err != nil || !s.hasher.Compare(u.PasswordHash, password) {
		return nil, domain.ErrInvalidCredentials
	}
	if !u.Active {
		return nil, domain.ErrInactive
	}
	profile, err := s.Profile(ctx, u.ID)
	if err != nil {
		return nil, err
	}
	token, exp, err := s.tokens.Generate(u.ID, remember)
	if err != nil {
		return nil, err
	}
	return &LoginResult{Token: token, ExpiresAt: exp, Profile: profile}, nil
}

// Profile carga el usuario y los permisos vigentes de su rol (siempre frescos desde la BD).
func (s *AuthService) Profile(ctx context.Context, userID int64) (*domain.Profile, error) {
	u, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return nil, domain.ErrUnauthorized
	}
	if !u.Active {
		return nil, domain.ErrInactive
	}
	role, err := s.roles.GetByID(ctx, u.RoleID)
	if err != nil {
		return nil, err
	}
	return &domain.Profile{User: *u, Permissions: role.Permissions}, nil
}
