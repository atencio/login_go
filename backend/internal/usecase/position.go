package usecase

import (
	"context"
	"strings"

	"nexo/internal/domain"
)

type PositionService struct{ repo domain.PositionRepository }

func NewPositionService(r domain.PositionRepository) *PositionService { return &PositionService{repo: r} }

type PositionInput struct {
	Name        string
	Description string
	Active      bool
}

func (in *PositionInput) normalize() error {
	in.Name = strings.TrimSpace(in.Name)
	in.Description = strings.TrimSpace(in.Description)
	if in.Name == "" {
		return domain.ValidationError{Msg: "el nombre del cargo es obligatorio"}
	}
	return nil
}

func (s *PositionService) List(ctx context.Context, q string) ([]domain.Position, error) {
	return s.repo.List(ctx, strings.TrimSpace(q))
}

func (s *PositionService) Get(ctx context.Context, id int64) (*domain.Position, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *PositionService) Create(ctx context.Context, in PositionInput) (*domain.Position, error) {
	if err := in.normalize(); err != nil {
		return nil, err
	}
	p := &domain.Position{Name: in.Name, Description: in.Description, Active: in.Active}
	if err := s.repo.Create(ctx, p); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, p.ID)
}

func (s *PositionService) Update(ctx context.Context, id int64, in PositionInput) (*domain.Position, error) {
	if err := in.normalize(); err != nil {
		return nil, err
	}
	p := &domain.Position{ID: id, Name: in.Name, Description: in.Description, Active: in.Active}
	if err := s.repo.Update(ctx, p); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, id)
}

func (s *PositionService) Delete(ctx context.Context, id int64) error { return s.repo.Delete(ctx, id) }
