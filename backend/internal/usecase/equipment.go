package usecase

import (
	"context"
	"strings"
	"time"

	"nexo/internal/domain"
)

const maxBatch = 200

type EquipmentService struct{ repo domain.EquipmentRepository }

func NewEquipmentService(r domain.EquipmentRepository) *EquipmentService { return &EquipmentService{repo: r} }

type EquipmentInput struct {
	Name         string
	CategoryID   int64
	BrandID      int64
	Model        string
	Serial       string // solo al editar; al crear se usan las series por unidad
	Status       string
	PurchaseDate *time.Time
	Cost         *float64
	Notes        string
}

func (in *EquipmentInput) normalize() error {
	in.Name = strings.TrimSpace(in.Name)
	in.Model = strings.TrimSpace(in.Model)
	in.Serial = strings.TrimSpace(in.Serial)
	in.Notes = strings.TrimSpace(in.Notes)
	switch {
	case in.Name == "":
		return domain.ValidationError{Msg: "el nombre del equipo es obligatorio"}
	case in.CategoryID == 0:
		return domain.ValidationError{Msg: "selecciona una categoría"}
	case in.Cost != nil && *in.Cost < 0:
		return domain.ValidationError{Msg: "el costo no puede ser negativo"}
	}
	return nil
}

func (in EquipmentInput) toEntity(id int64) *domain.Equipment {
	return &domain.Equipment{
		ID: id, Name: in.Name, CategoryID: in.CategoryID, BrandID: in.BrandID, Model: in.Model, Serial: in.Serial,
		Status: in.Status, PurchaseDate: in.PurchaseDate, Cost: in.Cost, Notes: in.Notes,
	}
}

func (s *EquipmentService) List(ctx context.Context, f domain.EquipmentFilter) ([]domain.Equipment, error) {
	f.Q = strings.TrimSpace(f.Q)
	return s.repo.List(ctx, f)
}

func (s *EquipmentService) Get(ctx context.Context, id int64) (*domain.Equipment, error) {
	return s.repo.GetByID(ctx, id)
}

// Create registra qty unidades: una por cada número de serie recibido (vacío = aún sin serie).
// La serie es única por producto, por eso no se comparte entre unidades.
// La cantidad solo existe al crear: después cada unidad se edita por separado.
func (s *EquipmentService) Create(ctx context.Context, in EquipmentInput, qty int, serials []string) ([]domain.Equipment, error) {
	if err := in.normalize(); err != nil {
		return nil, err
	}
	if qty == 0 {
		qty = 1
	}
	if qty < 1 || qty > maxBatch {
		return nil, domain.ValidationError{Msg: "la cantidad debe estar entre 1 y 200"}
	}
	if len(serials) == 0 && qty == 1 && in.Serial != "" {
		serials = []string{in.Serial}
	}
	if len(serials) == 0 {
		serials = make([]string, qty)
	}
	if len(serials) != qty {
		return nil, domain.ValidationError{Msg: "indica un número de serie por cada unidad (puedes dejar vacías las que aún no tengas)"}
	}
	seen := map[string]bool{}
	for i, v := range serials {
		v = strings.TrimSpace(v)
		serials[i] = v
		if v == "" {
			continue
		}
		key := strings.ToLower(v)
		if seen[key] {
			return nil, domain.ValidationError{Msg: "el número de serie " + `"` + v + `"` + " está repetido en la lista"}
		}
		seen[key] = true
	}
	if in.Status == "" {
		in.Status = domain.EquipAvailable
	}
	if in.Status != domain.EquipAvailable && in.Status != domain.EquipMaintenance {
		return nil, domain.ValidationError{Msg: "un equipo nuevo solo puede crearse como disponible o en mantenimiento"}
	}
	return s.repo.CreateMany(ctx, in.toEntity(0), serials)
}

func (s *EquipmentService) Update(ctx context.Context, id int64, in EquipmentInput) (*domain.Equipment, error) {
	if err := in.normalize(); err != nil {
		return nil, err
	}
	current, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if in.Status == "" {
		in.Status = current.Status
	}
	if !domain.ValidEquipmentStatus(in.Status) {
		return nil, domain.ValidationError{Msg: "estado no válido"}
	}
	// El estado "asignado" solo lo gestionan las asignaciones, para que inventario e historial no se desfasen.
	if current.Status == domain.EquipAssigned && in.Status != domain.EquipAssigned {
		return nil, domain.ValidationError{Msg: "el equipo está asignado: devuélvelo desde Asignaciones antes de cambiar su estado"}
	}
	if current.Status != domain.EquipAssigned && in.Status == domain.EquipAssigned {
		return nil, domain.ValidationError{Msg: "para asignar un equipo usa la vista de Asignaciones"}
	}
	if err := s.repo.Update(ctx, in.toEntity(id)); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, id)
}

func (s *EquipmentService) Delete(ctx context.Context, id int64) error {
	current, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if current.Status == domain.EquipAssigned {
		return domain.ValidationError{Msg: "el equipo está asignado: devuélvelo antes de eliminarlo"}
	}
	return s.repo.Delete(ctx, id)
}
