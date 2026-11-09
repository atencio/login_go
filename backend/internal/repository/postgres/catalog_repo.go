package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"nexo/internal/domain"
)

// CatalogSpec describe una tabla de catálogo (categorías o marcas). Son constantes del código, nunca entrada del usuario.
type CatalogSpec struct {
	Table     string // tabla del catálogo
	FKColumn  string // columna de equipment que la referencia
	UniqueMsg string
	DeleteMsg string
}

var (
	CategorySpec = CatalogSpec{
		Table: "equipment_categories", FKColumn: "category_id",
		UniqueMsg: "ya existe una categoría con ese nombre",
		DeleteMsg: "la categoría tiene equipos asociados; desactívala en lugar de eliminarla",
	}
	BrandSpec = CatalogSpec{
		Table: "equipment_brands", FKColumn: "brand_id",
		UniqueMsg: "ya existe una marca con ese nombre",
		DeleteMsg: "la marca tiene equipos asociados; desactívala en lugar de eliminarla",
	}
)

type CatalogRepo struct {
	db     *pgxpool.Pool
	spec   CatalogSpec
	unique map[string]string
	sel    string
}

func NewCatalogRepo(db *pgxpool.Pool, spec CatalogSpec) *CatalogRepo {
	return &CatalogRepo{
		db: db, spec: spec,
		unique: map[string]string{spec.Table + "_name_key": spec.UniqueMsg},
		sel: fmt.Sprintf(`SELECT t.id, t.name, t.active, t.created_at,
  (SELECT count(*) FROM equipment e WHERE e.%s = t.id) FROM %s t`, spec.FKColumn, spec.Table),
	}
}

func scanCatalog(row pgx.Row) (*domain.CatalogItem, error) {
	var c domain.CatalogItem
	if err := row.Scan(&c.ID, &c.Name, &c.Active, &c.CreatedAt, &c.EquipmentCount); err != nil {
		return nil, scanErr(err)
	}
	return &c, nil
}

func (r *CatalogRepo) Create(ctx context.Context, c *domain.CatalogItem) error {
	err := r.db.QueryRow(ctx,
		fmt.Sprintf(`INSERT INTO %s (name, active) VALUES ($1, $2) RETURNING id, created_at`, r.spec.Table),
		c.Name, c.Active).Scan(&c.ID, &c.CreatedAt)
	return pgErr(err, r.unique, "")
}

func (r *CatalogRepo) GetByID(ctx context.Context, id int64) (*domain.CatalogItem, error) {
	return scanCatalog(r.db.QueryRow(ctx, r.sel+` WHERE t.id = $1`, id))
}

func (r *CatalogRepo) List(ctx context.Context, q string) ([]domain.CatalogItem, error) {
	rows, err := r.db.Query(ctx, r.sel+` WHERE ($1::text = '' OR t.name ILIKE '%'||$1||'%') ORDER BY t.name`, q)
	return collect(rows, err, scanCatalog)
}

func (r *CatalogRepo) Update(ctx context.Context, c *domain.CatalogItem) error {
	tag, err := r.db.Exec(ctx, fmt.Sprintf(`UPDATE %s SET name = $1, active = $2 WHERE id = $3`, r.spec.Table),
		c.Name, c.Active, c.ID)
	return affected(tag, err, r.unique, "")
}

func (r *CatalogRepo) Delete(ctx context.Context, id int64) error {
	tag, err := r.db.Exec(ctx, fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, r.spec.Table), id)
	return affected(tag, err, nil, r.spec.DeleteMsg)
}
