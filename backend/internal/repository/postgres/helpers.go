package postgres

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"nexo/internal/domain"
)

// pgErr traduce violaciones de restricciones de PostgreSQL a errores de validación legibles.
// unique mapea el nombre de la restricción única a su mensaje; fk es el mensaje para claves foráneas.
func pgErr(err error, unique map[string]string, fk string) error {
	var pg *pgconn.PgError
	if errors.As(err, &pg) {
		switch pg.Code {
		case "23505":
			if m, ok := unique[pg.ConstraintName]; ok {
				return domain.ValidationError{Msg: m}
			}
			return domain.ValidationError{Msg: "ya existe un registro con esos datos"}
		case "23503":
			if fk != "" {
				return domain.ValidationError{Msg: fk}
			}
		}
	}
	return err
}

// affected resuelve el resultado de un UPDATE/DELETE: error mapeado o ErrNotFound si no tocó filas.
func affected(tag pgconn.CommandTag, err error, unique map[string]string, fk string) error {
	if err != nil {
		return pgErr(err, unique, fk)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func scanErr(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ErrNotFound
	}
	return err
}

// collect recorre las filas con la función de escaneo de cada entidad.
func collect[T any](rows pgx.Rows, err error, scan func(pgx.Row) (*T, error)) ([]T, error) {
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []T{}
	for rows.Next() {
		v, err := scan(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *v)
	}
	return out, rows.Err()
}
