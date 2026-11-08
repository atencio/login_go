package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"nexo/internal/domain"
)

type PositionRepo struct{ db *pgxpool.Pool }

func NewPositionRepo(db *pgxpool.Pool) *PositionRepo { return &PositionRepo{db: db} }

var positionUnique = map[string]string{"positions_name_key": "ya existe un cargo con ese nombre"}

const positionDeleteFK = "el cargo tiene funcionarios asignados; desactívalo en lugar de eliminarlo"

const positionSelect = `SELECT p.id, p.name, p.description, p.active, p.created_at,
  (SELECT count(*) FROM employees e WHERE e.position_id = p.id)
FROM positions p`

func scanPosition(row pgx.Row) (*domain.Position, error) {
	var p domain.Position
	if err := row.Scan(&p.ID, &p.Name, &p.Description, &p.Active, &p.CreatedAt, &p.EmployeeCount); err != nil {
		return nil, scanErr(err)
	}
	return &p, nil
}

func (r *PositionRepo) Create(ctx context.Context, p *domain.Position) error {
	err := r.db.QueryRow(ctx,
		`INSERT INTO positions (name, description, active) VALUES ($1,$2,$3) RETURNING id, created_at`,
		p.Name, p.Description, p.Active).Scan(&p.ID, &p.CreatedAt)
	return pgErr(err, positionUnique, "")
}

func (r *PositionRepo) GetByID(ctx context.Context, id int64) (*domain.Position, error) {
	return scanPosition(r.db.QueryRow(ctx, positionSelect+` WHERE p.id = $1`, id))
}

func (r *PositionRepo) List(ctx context.Context, q string) ([]domain.Position, error) {
	rows, err := r.db.Query(ctx, positionSelect+
		` WHERE ($1::text = '' OR p.name ILIKE '%'||$1||'%' OR p.description ILIKE '%'||$1||'%') ORDER BY p.name`, q)
	return collect(rows, err, scanPosition)
}

func (r *PositionRepo) Update(ctx context.Context, p *domain.Position) error {
	tag, err := r.db.Exec(ctx, `UPDATE positions SET name=$1, description=$2, active=$3 WHERE id=$4`,
		p.Name, p.Description, p.Active, p.ID)
	return affected(tag, err, positionUnique, "")
}

func (r *PositionRepo) Delete(ctx context.Context, id int64) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM positions WHERE id = $1`, id)
	return affected(tag, err, nil, positionDeleteFK)
}
