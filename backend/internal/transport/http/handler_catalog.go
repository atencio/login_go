package http

import (
	"context"
	"net/http"
	"time"

	"nexo/internal/domain"
	"nexo/internal/usecase"
)

type CatalogUseCase interface {
	List(ctx context.Context, q string) ([]domain.CatalogItem, error)
	Get(ctx context.Context, id int64) (*domain.CatalogItem, error)
	Create(ctx context.Context, in usecase.CatalogInput) (*domain.CatalogItem, error)
	Update(ctx context.Context, id int64, in usecase.CatalogInput) (*domain.CatalogItem, error)
	Delete(ctx context.Context, id int64) error
}

// CatalogHandler atiende un catálogo (categorías o marcas); se crea una instancia por cada uno.
type CatalogHandler struct {
	uc       CatalogUseCase
	findPerm string // permiso "Encontrar" que exige la búsqueda
}

func NewCatalogHandler(uc CatalogUseCase, findPerm string) *CatalogHandler {
	return &CatalogHandler{uc: uc, findPerm: findPerm}
}

type catalogDTO struct {
	ID             int64     `json:"id"`
	Name           string    `json:"name"`
	Active         bool      `json:"active"`
	EquipmentCount int       `json:"equipment_count"`
	CreatedAt      time.Time `json:"created_at"`
}

type catalogRequest struct {
	Name   string `json:"name"`
	Active bool   `json:"active"`
}

func toCatalogDTO(c domain.CatalogItem) catalogDTO {
	return catalogDTO{c.ID, c.Name, c.Active, c.EquipmentCount, c.CreatedAt}
}

func (h *CatalogHandler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if err := requireFind(r, h.findPerm, q != ""); err != nil {
		writeError(w, err)
		return
	}
	list, err := h.uc.List(r.Context(), q)
	if err != nil {
		writeError(w, err)
		return
	}
	out := make([]catalogDTO, 0, len(list))
	for _, c := range list {
		out = append(out, toCatalogDTO(c))
	}
	writeJSON(w, http.StatusOK, out)
}

// Options alimenta los selectores del formulario de equipos.
func (h *CatalogHandler) Options(w http.ResponseWriter, r *http.Request) {
	list, err := h.uc.List(r.Context(), "")
	if err != nil {
		writeError(w, err)
		return
	}
	type option struct {
		ID     int64  `json:"id"`
		Name   string `json:"name"`
		Active bool   `json:"active"`
	}
	out := make([]option, 0, len(list))
	for _, c := range list {
		out = append(out, option{c.ID, c.Name, c.Active})
	}
	writeJSON(w, http.StatusOK, out)
}

func (h *CatalogHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeError(w, err)
		return
	}
	c, err := h.uc.Get(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toCatalogDTO(*c))
}

func (h *CatalogHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req catalogRequest
	if err := decode(r, &req); err != nil {
		writeError(w, err)
		return
	}
	c, err := h.uc.Create(r.Context(), usecase.CatalogInput{Name: req.Name, Active: req.Active})
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, toCatalogDTO(*c))
}

func (h *CatalogHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeError(w, err)
		return
	}
	var req catalogRequest
	if err := decode(r, &req); err != nil {
		writeError(w, err)
		return
	}
	c, err := h.uc.Update(r.Context(), id, usecase.CatalogInput{Name: req.Name, Active: req.Active})
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toCatalogDTO(*c))
}

func (h *CatalogHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeError(w, err)
		return
	}
	if err := h.uc.Delete(r.Context(), id); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
