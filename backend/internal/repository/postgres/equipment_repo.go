package postgres

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"nexo/internal/domain"
)

type EquipmentRepo struct{ db *pgxpool.Pool }

func NewEquipmentRepo(db *pgxpool.Pool) *EquipmentRepo { return &EquipmentRepo{db: db} }

var (
	equipmentUnique = map[string]string{"equipment_serial_uq": "ya existe un equipo con ese número de serie"}
	serialDetail    = regexp.MustCompile(`\(serial_number\)=\((.*)\)`)
)

const (
	equipmentSaveFK   = "la categoría o la marca seleccionada no existe"
	equipmentDeleteFK = "el equipo tiene historial de asignaciones; márcalo como dado de baja en lugar de eliminarlo"
)

const equipmentSelect = `SELECT q.id, q.code, q.name, COALESCE(q.category_id, 0), COALESCE(c.name, ''),
  COALESCE(q.brand_id, 0), COALESCE(b.name, ''), q.model, q.serial_number, q.status,
  q.purchase_date, q.cost, q.notes, q.created_at, a.employee_id, (e.first_name || ' ' || e.last_name)
FROM equipment q
LEFT JOIN equipment_categories c ON c.id = q.category_id
LEFT JOIN equipment_brands b ON b.id = q.brand_id
LEFT JOIN assignments a ON a.equipment_id = q.id AND a.returned_at IS NULL
LEFT JOIN employees e ON e.id = a.employee_id`

func scanEquipment(row pgx.Row) (*domain.Equipment, error) {
	var q domain.Equipment
	err := row.Scan(&q.ID, &q.Code, &q.Name, &q.CategoryID, &q.CategoryName, &q.BrandID, &q.BrandName, &q.Model,
		&q.Serial, &q.Status, &q.PurchaseDate, &q.Cost, &q.Notes, &q.CreatedAt, &q.HolderID, &q.HolderName)
	if err != nil {
		return nil, scanErr(err)
	}
	return &q, nil
}

// saveErr indica cuál número de serie está repetido cuando la base lo rechaza.
func saveErr(err error) error {
	var pg *pgconn.PgError
	if errors.As(err, &pg) && pg.Code == "23505" && pg.ConstraintName == "equipment_serial_uq" {
		if m := serialDetail.FindStringSubmatch(pg.Detail); m != nil {
			return domain.ValidationError{Msg: fmt.Sprintf("el número de serie %q ya está registrado en otro equipo", m[1])}
		}
	}
	return pgErr(err, equipmentUnique, equipmentSaveFK)
}

func (r *EquipmentRepo) CreateMany(ctx context.Context, b *domain.Equipment, serials []string) ([]domain.Equipment, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx) //nolint:errcheck
	ids := make([]int64, 0, len(serials))
	for _, serial := range serials {
		var id int64
		err := tx.QueryRow(ctx,
			`INSERT INTO equipment (code, name, category_id, brand_id, model, serial_number, status, purchase_date, cost, notes)
			 VALUES ('EQ-' || lpad(nextval('equipment_code_seq')::text, 5, '0'), $1, NULLIF($2::bigint, 0), NULLIF($3::bigint, 0),
			         $4, $5, $6, $7, $8, $9) RETURNING id`,
			b.Name, b.CategoryID, b.BrandID, b.Model, serial, b.Status, b.PurchaseDate, b.Cost, b.Notes).Scan(&id)
		if err != nil {
			return nil, saveErr(err)
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
 WHERE ($1::text = '' OR q.code ILIKE '%'||$1||'%' OR q.name ILIKE '%'||$1||'%' OR c.name ILIKE '%'||$1||'%'
        OR b.name ILIKE '%'||$1||'%' OR q.model ILIKE '%'||$1||'%' OR q.serial_number ILIKE '%'||$1||'%'
        OR (e.first_name || ' ' || e.last_name) ILIKE '%'||$1||'%')
   AND ($2::text = '' OR q.status = $2)
 ORDER BY q.id DESC`, f.Q, f.Status)
	return collect(rows, err, scanEquipment)
}

func (r *EquipmentRepo) Update(ctx context.Context, q *domain.Equipment) error {
	tag, err := r.db.Exec(ctx,
		`UPDATE equipment SET name=$1, category_id=NULLIF($2::bigint, 0), brand_id=NULLIF($3::bigint, 0), model=$4,
		   serial_number=$5, status=$6, purchase_date=$7, cost=$8, notes=$9 WHERE id=$10`,
		q.Name, q.CategoryID, q.BrandID, q.Model, q.Serial, q.Status, q.PurchaseDate, q.Cost, q.Notes, q.ID)
	if err != nil {
		return saveErr(err)
	}
	return affected(tag, nil, nil, "")
}

func (r *EquipmentRepo) Delete(ctx context.Context, id int64) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM equipment WHERE id = $1`, id)
	return affected(tag, err, nil, equipmentDeleteFK)
}
