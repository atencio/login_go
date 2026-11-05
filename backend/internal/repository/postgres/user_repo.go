package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"nexo/internal/domain"
)

type UserRepo struct{ db *pgxpool.Pool }

func NewUserRepo(db *pgxpool.Pool) *UserRepo { return &UserRepo{db: db} }

const userSelect = `SELECT u.id, u.name, u.email, u.password_hash, u.role_id, r.name, u.active, u.created_at
FROM users u JOIN roles r ON r.id = u.role_id`

func scanUser(row pgx.Row) (*domain.User, error) {
	var u domain.User
	err := row.Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.RoleID, &u.RoleName, &u.Active, &u.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) Create(ctx context.Context, u *domain.User) error {
	return mapErr(r.db.QueryRow(ctx,
		`INSERT INTO users (name, email, password_hash, role_id, active) VALUES ($1,$2,$3,$4,$5) RETURNING id, created_at`,
		u.Name, u.Email, u.PasswordHash, u.RoleID, u.Active).Scan(&u.ID, &u.CreatedAt))
}

func (r *UserRepo) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	return scanUser(r.db.QueryRow(ctx, userSelect+` WHERE u.id = $1`, id))
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return scanUser(r.db.QueryRow(ctx, userSelect+` WHERE u.email = $1`, email))
}

func (r *UserRepo) List(ctx context.Context) ([]domain.User, error) {
	rows, err := r.db.Query(ctx, userSelect+` ORDER BY u.id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	users := []domain.User{}
	for rows.Next() {
		u, err := scanUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, *u)
	}
	return users, rows.Err()
}

func (r *UserRepo) Search(ctx context.Context, q string) ([]domain.User, error) {
	rows, err := r.db.Query(ctx, userSelect+` WHERE u.name ILIKE $1 OR u.email ILIKE $1 OR r.name ILIKE $1 ORDER BY u.id`, "%"+q+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	users := []domain.User{}
	for rows.Next() {
		u, err := scanUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, *u)
	}
	return users, rows.Err()
}

func (r *UserRepo) Update(ctx context.Context, u *domain.User) error {
	tag, err := r.db.Exec(ctx,
		`UPDATE users SET name=$1, email=$2, password_hash=$3, role_id=$4, active=$5 WHERE id=$6`,
		u.Name, u.Email, u.PasswordHash, u.RoleID, u.Active, u.ID)
	if err != nil {
		return mapErr(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *UserRepo) Delete(ctx context.Context, id int64) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *UserRepo) Count(ctx context.Context) (int, error) {
	var n int
	err := r.db.QueryRow(ctx, `SELECT count(*) FROM users`).Scan(&n)
	return n, err
}
