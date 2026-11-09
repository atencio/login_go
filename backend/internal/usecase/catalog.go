package usecase

import (
	"context"
	"strings"

	"nexo/internal/domain"
)

// CatalogService gestiona un catálogo simple (categorías o marcas); noun es el nombre en singular para los mensajes.
type CatalogService struct {
	repo domain.CatalogRepository
	noun string
}

func NewCatalogService(r domain.CatalogRepository, noun string) *CatalogService {
	return &CatalogService{repo: r, noun: noun}
}

type CatalogInput struct {
	Name   string
	Active bool
}

func (s *CatalogService) normalize(in *CatalogInput) error {
	in.Name = strings.TrimSpace(in.Name)
	if in.Name == "" {
		return domain.ValidationError{Msg: "el nombre de la " + s.noun + " es obligatorio"}
	}
	return nil
}

func (s *CatalogService) List(ctx context.Context, q string) ([]domain.CatalogItem, error) {
	return s.repo.List(ctx, strings.TrimSpace(q))
}

func (s *CatalogService) Get(ctx context.Context, id int64) (*domain.CatalogItem, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *CatalogService) Create(ctx context.Context, in CatalogInput) (*domain.CatalogItem, error) {
	if err := s.normalize(&in); err != nil {
		return nil, err
	}
	c := &domain.CatalogItem{Name: in.Name, Active: in.Active}
	if err := s.repo.Create(ctx, c); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, c.ID)
}

func (s *CatalogService) Update(ctx context.Context, id int64, in CatalogInput) (*domain.CatalogItem, error) {
	if err := s.normalize(&in); err != nil {
		return nil, err
	}
	if err := s.repo.Update(ctx, &domain.CatalogItem{ID: id, Name: in.Name, Active: in.Active}); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, id)
}

func (s *CatalogService) Delete(ctx context.Context, id int64) error { return s.repo.Delete(ctx, id) }
