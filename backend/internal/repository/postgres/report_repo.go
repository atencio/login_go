package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"nexo/internal/domain"
)

type ReportRepo struct{ db *pgxpool.Pool }

func NewReportRepo(db *pgxpool.Pool) *ReportRepo { return &ReportRepo{db: db} }

func (r *ReportRepo) UsersSummary(ctx context.Context) (*domain.ReportSummary, error) {
	s := &domain.ReportSummary{ByRole: []domain.RoleCount{}}
	err := r.db.QueryRow(ctx, `SELECT count(*), count(*) FILTER (WHERE active) FROM users`).
		Scan(&s.TotalUsers, &s.ActiveUsers)
	if err != nil {
		return nil, err
	}
	rows, err := r.db.Query(ctx,
		`SELECT r.name, count(u.id) FROM roles r LEFT JOIN users u ON u.role_id = r.id GROUP BY r.id ORDER BY r.id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var rc domain.RoleCount
		if err := rows.Scan(&rc.Role, &rc.Total); err != nil {
			return nil, err
		}
		s.ByRole = append(s.ByRole, rc)
	}
	return s, rows.Err()
}
