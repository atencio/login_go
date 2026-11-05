package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"nexo/internal/domain"
)

type RoleRepo struct{ db *pgxpool.Pool }

func NewRoleRepo(db *pgxpool.Pool) *RoleRepo { return &RoleRepo{db: db} }

const roleSelect = `SELECT r.id, r.name,
  COALESCE(array_agg(p.code ORDER BY p.sort, p.id) FILTER (WHERE p.code IS NOT NULL), '{}'),
  (SELECT count(*) FROM users u WHERE u.role_id = r.id)
FROM roles r
LEFT JOIN role_permissions rp ON rp.role_id = r.id
LEFT JOIN permissions p ON p.id = rp.permission_id`

func (r *RoleRepo) List(ctx context.Context) ([]domain.Role, error) {
	rows, err := r.db.Query(ctx, roleSelect+` GROUP BY r.id ORDER BY r.id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	roles := []domain.Role{}
	for rows.Next() {
		var ro domain.Role
		if err := rows.Scan(&ro.ID, &ro.Name, &ro.Permissions, &ro.UserCount); err != nil {
			return nil, err
		}
		roles = append(roles, ro)
	}
	return roles, rows.Err()
}

func (r *RoleRepo) GetByID(ctx context.Context, id int64) (*domain.Role, error) {
	var ro domain.Role
	err := r.db.QueryRow(ctx, roleSelect+` WHERE r.id = $1 GROUP BY r.id`, id).
		Scan(&ro.ID, &ro.Name, &ro.Permissions, &ro.UserCount)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &ro, nil
}

func (r *RoleRepo) Create(ctx context.Context, ro *domain.Role) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) //nolint:errcheck
	if err := tx.QueryRow(ctx, `INSERT INTO roles (name) VALUES ($1) RETURNING id`, ro.Name).Scan(&ro.ID); err != nil {
		return mapRoleErr(err)
	}
	if err := setPermissions(ctx, tx, ro.ID, ro.Permissions); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *RoleRepo) Update(ctx context.Context, ro *domain.Role) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) //nolint:errcheck
	tag, err := tx.Exec(ctx, `UPDATE roles SET name = $1 WHERE id = $2`, ro.Name, ro.ID)
	if err != nil {
		return mapRoleErr(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	if err := setPermissions(ctx, tx, ro.ID, ro.Permissions); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *RoleRepo) Delete(ctx context.Context, id int64) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM roles WHERE id = $1`, id)
	if err != nil {
		return mapRoleErr(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *RoleRepo) ListPermissions(ctx context.Context) ([]domain.Permission, error) {
	rows, err := r.db.Query(ctx, `SELECT code, module, module_label, action FROM permissions ORDER BY sort, id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []domain.Permission{}
	for rows.Next() {
		var p domain.Permission
		if err := rows.Scan(&p.Code, &p.Module, &p.ModuleLabel, &p.Action); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

// setPermissions reemplaza los permisos del perfil por los códigos indicados.
func setPermissions(ctx context.Context, tx pgx.Tx, roleID int64, codes []string) error {
	if _, err := tx.Exec(ctx, `DELETE FROM role_permissions WHERE role_id = $1`, roleID); err != nil {
		return err
	}
	_, err := tx.Exec(ctx,
		`INSERT INTO role_permissions (role_id, permission_id) SELECT $1, id FROM permissions WHERE code = ANY($2)`,
		roleID, codes)
	return err
}

func mapRoleErr(err error) error {
	var pg *pgconn.PgError
	if errors.As(err, &pg) {
		switch pg.Code {
		case "23505":
			return domain.ValidationError{Msg: "ya existe un perfil con ese nombre"}
		case "23503":
			return domain.ValidationError{Msg: "no puedes eliminar un perfil que tiene usuarios asignados"}
		}
	}
	return err
}
