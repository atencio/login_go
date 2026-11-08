package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"nexo/internal/domain"
)

type EmployeeRepo struct{ db *pgxpool.Pool }

func NewEmployeeRepo(db *pgxpool.Pool) *EmployeeRepo { return &EmployeeRepo{db: db} }

var employeeUnique = map[string]string{
	"employees_email_key":   "ya existe un funcionario con ese correo",
	"employees_document_uq": "ya existe un funcionario con ese documento",
}

const (
	employeeSaveFK   = "el cargo seleccionado no existe"
	employeeDeleteFK = "el funcionario tiene historial de asignaciones; desactívalo en lugar de eliminarlo"
)

const employeeSelect = `SELECT e.id, e.first_name, e.last_name, e.email, e.document, e.phone, e.department,
  e.position_id, p.name, e.hire_date, e.active, e.created_at,
  (SELECT count(*) FROM assignments a WHERE a.employee_id = e.id AND a.returned_at IS NULL)
FROM employees e JOIN positions p ON p.id = e.position_id`

func scanEmployee(row pgx.Row) (*domain.Employee, error) {
	var e domain.Employee
	err := row.Scan(&e.ID, &e.FirstName, &e.LastName, &e.Email, &e.Document, &e.Phone, &e.Department,
		&e.PositionID, &e.PositionName, &e.HireDate, &e.Active, &e.CreatedAt, &e.ActiveAssignments)
	if err != nil {
		return nil, scanErr(err)
	}
	return &e, nil
}

func (r *EmployeeRepo) Create(ctx context.Context, e *domain.Employee) error {
	err := r.db.QueryRow(ctx,
		`INSERT INTO employees (first_name, last_name, email, document, phone, department, position_id, hire_date, active)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING id, created_at`,
		e.FirstName, e.LastName, e.Email, e.Document, e.Phone, e.Department, e.PositionID, e.HireDate, e.Active,
	).Scan(&e.ID, &e.CreatedAt)
	return pgErr(err, employeeUnique, employeeSaveFK)
}

func (r *EmployeeRepo) GetByID(ctx context.Context, id int64) (*domain.Employee, error) {
	return scanEmployee(r.db.QueryRow(ctx, employeeSelect+` WHERE e.id = $1`, id))
}

func (r *EmployeeRepo) List(ctx context.Context, f domain.EmployeeFilter) ([]domain.Employee, error) {
	rows, err := r.db.Query(ctx, employeeSelect+`
 WHERE ($1::text = '' OR (e.first_name || ' ' || e.last_name) ILIKE '%'||$1||'%' OR e.email ILIKE '%'||$1||'%'
        OR e.document ILIKE '%'||$1||'%' OR e.department ILIKE '%'||$1||'%')
   AND ($2::bigint = 0 OR e.position_id = $2)
 ORDER BY e.first_name, e.last_name`, f.Q, f.PositionID)
	return collect(rows, err, scanEmployee)
}

func (r *EmployeeRepo) Update(ctx context.Context, e *domain.Employee) error {
	tag, err := r.db.Exec(ctx,
		`UPDATE employees SET first_name=$1, last_name=$2, email=$3, document=$4, phone=$5, department=$6,
		   position_id=$7, hire_date=$8, active=$9 WHERE id=$10`,
		e.FirstName, e.LastName, e.Email, e.Document, e.Phone, e.Department, e.PositionID, e.HireDate, e.Active, e.ID)
	return affected(tag, err, employeeUnique, employeeSaveFK)
}

func (r *EmployeeRepo) Delete(ctx context.Context, id int64) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM employees WHERE id = $1`, id)
	return affected(tag, err, nil, employeeDeleteFK)
}
