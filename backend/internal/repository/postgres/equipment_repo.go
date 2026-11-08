package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"nexo/internal/domain"
)

type EquipmentRepo struct{ db *pgxpool.Pool }

func NewEquipmentRepo(db *pgxpool.Pool) *EquipmentRepo { return &EquipmentRepo{db: db} }

var equipmentUnique = map[string]string{"equipment_serial_uq": "ya existe un equipo con ese número de serie"}

const equipmentDeleteFK = "el equipo tiene historial de asignaciones; márcalo como dado de baja en lugar de eliminarlo"

const equipmentSelect = `SELECT q.id, q.code, q.name, q.category, q.brand, q.model, q.serial_number, q.status,
  q.purchase_date, q.cost, q.notes, q.created_at, a.employee_id, (e.first_name || ' ' || e.last_name)
FROM equipment q
LEFT JOIN assignments a ON a.equipment_id = q.id AND a.returned_at IS NULL
LEFT JOIN employees e ON e.id = a.employee_id`

func scanEquipment(row pgx.Row) (*domain.Equipment, error) {
	var q domain.Equipment
	err := row.Scan(&q.ID, &q.Code, &q.Name, &q.Category, &q.Brand, &q.Model, &q.Serial, &q.Status,
		&q.PurchaseDate, &q.Cost, &q.Notes, &q.CreatedAt, &q.HolderID, &q.HolderName)
	if err != nil {
		return nil, scanErr(err)
	}
	return &q, nil
}

func (r *EquipmentRepo) CreateMany(ctx context.Context, b *domain.Equipment, qty int) ([]domain.Equipment, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx) //nolint:errcheck
	ids := make([]int64, 0, qty)
	for i := 0; i < qty; i++ {
		var id int64
		err := tx.QueryRow(ctx,
			`INSERT INTO equipment (code, name, category, brand, model, serial_number, status, purchase_date, cost, notes)
			 VALUES ('EQ-' || lpad(nextval('equipment_code_seq')::text, 5, '0'), $1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING id`,
			b.Name, b.Category, b.Brand, b.Model, b.Serial, b.Status, b.PurchaseDate, b.Cost, b.Notes).Scan(&id)
		if err != nil {
			return nil, pgErr(err, equipmentUnique, "")
		}
		ids = append(ids, id)
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	rows, err := r.db.Query(ctx, equipmentSelect+` WHERE q.id = ANY($1) ORDER BY q.id`, ids)
	return collect(rows, err, scanEquipment)
}

func (r *EquipmentRepo) GetByID(ctx context.Context, id int64) (*domain.Equipment, error) {
	return scanEquipment(r.db.QueryRow(ctx, equipmentSelect+` WHERE q.id = $1`, id))
}

func (r *EquipmentRepo) List(ctx context.Context, f domain.EquipmentFilter) ([]domain.Equipment, error) {
	rows, err := r.db.Query(ctx, equipmentSelect+`
 WHERE ($1::text = '' OR q.code ILIKE '%'||$1||'%' OR q.name ILIKE '%'||$1||'%' OR q.category ILIKE '%'||$1||'%'
        OR q.brand ILIKE '%'||$1||'%' OR q.model ILIKE '%'||$1||'%' OR q.serial_number ILIKE '%'||$1||'%'
        OR (e.first_name || ' ' || e.last_name) ILIKE '%'||$1||'%')
   AND ($2::text = '' OR q.status = $2)
 ORDER BY q.id DESC`, f.Q, f.Status)
	return collect(rows, err, scanEquipment)
}

func (r *EquipmentRepo) Update(ctx context.Context, q *domain.Equipment) error {
	tag, err := r.db.Exec(ctx,
		`UPDATE equipment SET name=$1, category=$2, brand=$3, model=$4, serial_number=$5, status=$6,
		   purchase_date=$7, cost=$8, notes=$9 WHERE id=$10`,
		q.Name, q.Category, q.Brand, q.Model, q.Serial, q.Status, q.PurchaseDate, q.Cost, q.Notes, q.ID)
	return affected(tag, err, equipmentUnique, "")
}

func (r *EquipmentRepo) Delete(ctx context.Context, id int64) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM equipment WHERE id = $1`, id)
	return affected(tag, err, nil, equipmentDeleteFK)
}
