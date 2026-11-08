package usecase

import (
	"context"

	"nexo/internal/domain"
)

// ReportService agrupa los reportes. Cada reporte nuevo es un método más
// (y un permiso nuevo en la tabla permissions).
type ReportService struct{ repo domain.ReportRepository }

func NewReportService(r domain.ReportRepository) *ReportService { return &ReportService{repo: r} }

func (s *ReportService) UsersSummary(ctx context.Context) (*domain.ReportSummary, error) {
	return s.repo.UsersSummary(ctx)
}

func (s *ReportService) InventorySummary(ctx context.Context) (*domain.InventorySummary, error) {
	return s.repo.InventorySummary(ctx)
}
