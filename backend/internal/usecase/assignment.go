package usecase

import (
	"context"
	"errors"
	"strings"

	"nexo/internal/domain"
)

type AssignmentService struct {
	assignments domain.AssignmentRepository
	equipment   domain.EquipmentRepository
	employees   domain.EmployeeRepository
}

func NewAssignmentService(a domain.AssignmentRepository, q domain.EquipmentRepository, e domain.EmployeeRepository) *AssignmentService {
	return &AssignmentService{assignments: a, equipment: q, employees: e}
}

func (s *AssignmentService) requireActiveEmployee(ctx context.Context, id int64) error {
	e, err := s.employees.GetByID(ctx, id)
	if errors.Is(err, domain.ErrNotFound) {
		return domain.ValidationError{Msg: "el funcionario no existe"}
	}
	if err != nil {
		return err
	}
	if !e.Active {
		return domain.ValidationError{Msg: "el funcionario está inactivo y no puede recibir equipos"}
	}
	return nil
}

func (s *AssignmentService) List(ctx context.Context, f domain.AssignmentFilter) ([]domain.Assignment, error) {
	f.Q = strings.TrimSpace(f.Q)
	return s.assignments.List(ctx, f)
}

func (s *AssignmentService) Get(ctx context.Context, id int64) (*domain.Assignment, error) {
	return s.assignments.GetByID(ctx, id)
}

// Assign entrega un equipo disponible a un funcionario activo.
func (s *AssignmentService) Assign(ctx context.Context, actorID, equipmentID, employeeID int64, notes string) (*domain.Assignment, error) {
	eq, err := s.equipment.GetByID(ctx, equipmentID)
	if errors.Is(err, domain.ErrNotFound) {
		return nil, domain.ValidationError{Msg: "el equipo no existe"}
	}
	if err != nil {
		return nil, err
	}
	if eq.Status != domain.EquipAvailable {
		return nil, domain.ValidationError{Msg: "el equipo no está disponible para asignar"}
	}
	if err := s.requireActiveEmployee(ctx, employeeID); err != nil {
		return nil, err
	}
	a := &domain.Assignment{
		EquipmentID: equipmentID, EmployeeID: employeeID, AssignedByID: actorID, DeliveryNotes: strings.TrimSpace(notes),
	}
	if err := s.assignments.Assign(ctx, a); err != nil {
		return nil, err
	}
	return s.assignments.GetByID(ctx, a.ID)
}

// Update cambia las notas de una asignación activa o, si se indica otro funcionario, la transfiere
// (queda la devolución al anterior y una nueva entrega en el historial).
func (s *AssignmentService) Update(ctx context.Context, actorID, id, employeeID int64, notes string) (*domain.Assignment, error) {
	current, err := s.assignments.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if current.ReturnedAt != nil {
		return nil, domain.ValidationError{Msg: "la asignación ya fue devuelta y no se puede modificar"}
	}
	notes = strings.TrimSpace(notes)
	if employeeID != 0 && employeeID != current.EmployeeID {
		if err := s.requireActiveEmployee(ctx, employeeID); err != nil {
			return nil, err
		}
		newID, err := s.assignments.Transfer(ctx, id, employeeID, actorID, notes)
		if err != nil {
			return nil, err
		}
		return s.assignments.GetByID(ctx, newID)
	}
	if err := s.assignments.UpdateNotes(ctx, id, notes); err != nil {
		return nil, err
	}
	return s.assignments.GetByID(ctx, id)
}

// Return registra la devolución y deja el equipo disponible, en mantenimiento o dado de baja.
func (s *AssignmentService) Return(ctx context.Context, actorID, id int64, notes, status string) (*domain.Assignment, error) {
	if status == "" {
		status = domain.EquipAvailable
	}
	switch status {
	case domain.EquipAvailable, domain.EquipMaintenance, domain.EquipRetired:
	default:
		return nil, domain.ValidationError{Msg: "al devolver, el equipo debe quedar disponible, en mantenimiento o dado de baja"}
	}
	if err := s.assignments.Return(ctx, id, actorID, strings.TrimSpace(notes), status); err != nil {
		return nil, err
	}
	return s.assignments.GetByID(ctx, id)
}
