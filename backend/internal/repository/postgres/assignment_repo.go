package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"nexo/internal/domain"
)

type AssignmentRepo struct{ db *pgxpool.Pool }

func NewAssignmentRepo(db *pgxpool.Pool) *AssignmentRepo { return &AssignmentRepo{db: db} }

var assignmentUnique = map[string]string{"assignments_one_active": "el equipo ya está asignado a otro funcionario"}

const assignmentFK = "el equipo o el funcionario seleccionado no existe"

const assignmentSelect = `SELECT a.id, a.equipment_id, q.code, q.name, a.employee_id,
  (e.first_name || ' ' || e.last_name), a.assigned_at, COALESCE(u1.name, ''), a.delivery_notes,
  a.returned_at, COALESCE(u2.name, ''), a.return_notes
FROM assignments a
JOIN equipment q ON q.id = a.equipment_id
JOIN employees e ON e.id = a.employee_id
LEFT JOIN users u1 ON u1.id = a.assigned_by
LEFT JOIN users u2 ON u2.id = a.received_by`

func scanAssignment(row pgx.Row) (*domain.Assignment, error) {
	var a domain.Assignment
	err := row.Scan(&a.ID, &a.EquipmentID, &a.EquipmentCode, &a.EquipmentName, &a.EmployeeID, &a.EmployeeName,
		&a.AssignedAt, &a.AssignedBy, &a.DeliveryNotes, &a.ReturnedAt, &a.ReceivedBy, &a.ReturnNotes)
	if err != nil {
		return nil, scanErr(err)
	}
	return &a, nil
}

func (r *AssignmentRepo) Assign(ctx context.Context, a *domain.Assignment) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) //nolint:errcheck
	// El UPDATE condicionado evita asignar dos veces el mismo equipo en peticiones simultáneas.
	tag, err := tx.Exec(ctx, `UPDATE equipment SET status = 'assigned' WHERE id = $1 AND status = 'available'`, a.EquipmentID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ValidationError{Msg: "el equipo no está disponible para asignar"}
	}
	err = tx.QueryRow(ctx,
		`INSERT INTO assignments (equipment_id, employee_id, assigned_by, delivery_notes)
		 VALUES ($1, $2, NULLIF($3::bigint, 0), $4) RETURNING id, assigned_at`,
		a.EquipmentID, a.EmployeeID, a.AssignedByID, a.DeliveryNotes).Scan(&a.ID, &a.AssignedAt)
	if err != nil {
		return pgErr(err, assignmentUnique, assignmentFK)
	}
	return tx.Commit(ctx)
}

func (r *AssignmentRepo) GetByID(ctx context.Context, id int64) (*domain.Assignment, error) {
	return scanAssignment(r.db.QueryRow(ctx, assignmentSelect+` WHERE a.id = $1`, id))
}

func (r *AssignmentRepo) List(ctx context.Context, f domain.AssignmentFilter) ([]domain.Assignment, error) {
	rows, err := r.db.Query(ctx, assignmentSelect+`
 WHERE ($1::text = '' OR q.code ILIKE '%'||$1||'%' OR q.name ILIKE '%'||$1||'%' OR q.serial_number ILIKE '%'||$1||'%'
        OR (e.first_name || ' ' || e.last_name) ILIKE '%'||$1||'%')
   AND ($2::text = '' OR ($2 = 'active' AND a.returned_at IS NULL) OR ($2 = 'returned' AND a.returned_at IS NOT NULL))
   AND ($3::bigint = 0 OR a.employee_id = $3)
   AND ($4::bigint = 0 OR a.equipment_id = $4)
 ORDER BY a.assigned_at DESC, a.id DESC`, f.Q, f.State, f.EmployeeID, f.EquipmentID)
	return collect(rows, err, scanAssignment)
}

func (r *AssignmentRepo) UpdateNotes(ctx context.Context, id int64, notes string) error {
	tag, err := r.db.Exec(ctx, `UPDATE assignments SET delivery_notes = $1 WHERE id = $2 AND returned_at IS NULL`, notes, id)
	return affected(tag, err, nil, "")
}

func (r *AssignmentRepo) Return(ctx context.Context, id, userID int64, notes, status string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) //nolint:errcheck
	var equipmentID int64
	err = tx.QueryRow(ctx,
		`UPDATE assignments SET returned_at = now(), received_by = NULLIF($2::bigint, 0), return_notes = $3
		 WHERE id = $1 AND returned_at IS NULL RETURNING equipment_id`, id, userID, notes).Scan(&equipmentID)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ValidationError{Msg: "la asignación ya fue devuelta"}
	}
	if err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `UPDATE equipment SET status = $1 WHERE id = $2`, status, equipmentID); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// Transfer cierra la asignación actual y abre otra para el nuevo funcionario; el equipo sigue "asignado".
func (r *AssignmentRepo) Transfer(ctx context.Context, id, toEmployeeID, userID int64, notes string) (int64, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx) //nolint:errcheck
	var equipmentID int64
	err = tx.QueryRow(ctx,
		`UPDATE assignments SET returned_at = now(), received_by = NULLIF($2::bigint, 0),
		   return_notes = 'Transferido a otro funcionario'
		 WHERE id = $1 AND returned_at IS NULL RETURNING equipment_id`, id, userID).Scan(&equipmentID)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, domain.ValidationError{Msg: "la asignación ya fue devuelta"}
	}
	if err != nil {
		return 0, err
	}
	var newID int64
	err = tx.QueryRow(ctx,
		`INSERT INTO assignments (equipment_id, employee_id, assigned_by, delivery_notes)
		 VALUES ($1, $2, NULLIF($3::bigint, 0), $4) RETURNING id`, equipmentID, toEmployeeID, userID, notes).Scan(&newID)
	if err != nil {
		return 0, pgErr(err, assignmentUnique, assignmentFK)
	}
	return newID, tx.Commit(ctx)
}
