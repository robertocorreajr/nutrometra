package usecase

import (
	"context"
	"fmt"
	"time"

	"nutrometra/api/internal/scheduling/domain"

	"github.com/google/uuid"
)

// SchedulingRepository defines the data access contract.
type SchedulingRepository interface {
	CreateAvailabilityRule(ctx context.Context, rule *domain.AvailabilityRule) error
	ListAvailabilityRules(ctx context.Context, tenantID, professionalID uuid.UUID) ([]domain.AvailabilityRule, error)
	GetAvailabilityRulesByDay(ctx context.Context, tenantID, professionalID uuid.UUID, dayOfWeek int) ([]domain.AvailabilityRule, error)
	UpdateAvailabilityRule(ctx context.Context, rule *domain.AvailabilityRule) error
	DeleteAvailabilityRule(ctx context.Context, tenantID, id uuid.UUID) error

	CreateBlock(ctx context.Context, b *domain.ScheduleBlock) error
	ListBlocks(ctx context.Context, tenantID, professionalID uuid.UUID, from, to time.Time) ([]domain.ScheduleBlock, error)
	DeleteBlock(ctx context.Context, tenantID, id uuid.UUID) error

	CreateAppointment(ctx context.Context, a *domain.Appointment) error
	GetAppointmentByID(ctx context.Context, tenantID, id uuid.UUID) (*domain.Appointment, error)
	ListAppointments(ctx context.Context, tenantID uuid.UUID, professionalID *uuid.UUID, from, to *time.Time, status *domain.AppointmentStatus) ([]domain.Appointment, error)
	UpdateAppointmentStatus(ctx context.Context, tenantID, id uuid.UUID, status domain.AppointmentStatus, cancellationReason string) error
	RescheduleAppointment(ctx context.Context, tenantID, id uuid.UUID, startAt, endAt time.Time) error
	FindConflictingAppointments(ctx context.Context, tenantID, professionalID uuid.UUID, startAt, endAt time.Time, excludeID *uuid.UUID) ([]domain.Appointment, error)
	CreateAuditEvent(ctx context.Context, e *domain.AppointmentAuditEvent) error
}

// Usecase contains business logic for scheduling.
type Usecase struct {
	repo SchedulingRepository
}

// New creates a scheduling Usecase.
func New(repo SchedulingRepository) *Usecase {
	return &Usecase{repo: repo}
}

// --- Availability Rules ---

// CreateAvailabilityRule creates a new availability rule.
func (uc *Usecase) CreateAvailabilityRule(ctx context.Context, rule *domain.AvailabilityRule) error {
	if rule.DayOfWeek < 0 || rule.DayOfWeek > 6 {
		return fmt.Errorf("scheduling: day_of_week must be 0-6")
	}
	if !rule.EndTime.After(rule.StartTime) {
		return fmt.Errorf("scheduling: end_time must be after start_time")
	}
	if rule.ServiceMode == domain.ModeOnsite && rule.AddressID == nil {
		return fmt.Errorf("scheduling: onsite mode requires address_id")
	}

	rule.ID = uuid.New()
	now := time.Now().UTC()
	rule.CreatedAt = now
	rule.UpdatedAt = now
	rule.Active = true

	return uc.repo.CreateAvailabilityRule(ctx, rule)
}

// ListAvailabilityRules returns all rules for a professional.
func (uc *Usecase) ListAvailabilityRules(ctx context.Context, tenantID, professionalID uuid.UUID) ([]domain.AvailabilityRule, error) {
	return uc.repo.ListAvailabilityRules(ctx, tenantID, professionalID)
}

// UpdateAvailabilityRule updates a rule.
func (uc *Usecase) UpdateAvailabilityRule(ctx context.Context, rule *domain.AvailabilityRule) error {
	return uc.repo.UpdateAvailabilityRule(ctx, rule)
}

// DeleteAvailabilityRule soft-deletes a rule.
func (uc *Usecase) DeleteAvailabilityRule(ctx context.Context, tenantID, id uuid.UUID) error {
	return uc.repo.DeleteAvailabilityRule(ctx, tenantID, id)
}

// --- Schedule Blocks ---

// CreateBlock creates a schedule block.
func (uc *Usecase) CreateBlock(ctx context.Context, b *domain.ScheduleBlock) error {
	if !b.EndAt.After(b.StartAt) {
		return fmt.Errorf("scheduling: end_at must be after start_at")
	}

	b.ID = uuid.New()
	b.CreatedAt = time.Now().UTC()

	return uc.repo.CreateBlock(ctx, b)
}

// ListBlocks returns blocks for a professional in a date range.
func (uc *Usecase) ListBlocks(ctx context.Context, tenantID, professionalID uuid.UUID, from, to time.Time) ([]domain.ScheduleBlock, error) {
	return uc.repo.ListBlocks(ctx, tenantID, professionalID, from, to)
}

// DeleteBlock removes a block.
func (uc *Usecase) DeleteBlock(ctx context.Context, tenantID, id uuid.UUID) error {
	return uc.repo.DeleteBlock(ctx, tenantID, id)
}

// --- Appointments ---

// CreateAppointment creates a new appointment after checking for conflicts.
func (uc *Usecase) CreateAppointment(ctx context.Context, a *domain.Appointment, actorID *uuid.UUID) error {
	if !a.EndAt.After(a.StartAt) {
		return fmt.Errorf("scheduling: end_at must be after start_at")
	}

	// Check for conflicts
	conflicts, err := uc.repo.FindConflictingAppointments(ctx, a.TenantID, a.ProfessionalID, a.StartAt, a.EndAt, nil)
	if err != nil {
		return err
	}
	if len(conflicts) > 0 {
		return domain.ErrConflict
	}

	a.ID = uuid.New()
	now := time.Now().UTC()
	a.CreatedAt = now
	a.UpdatedAt = now
	a.Status = domain.StatusScheduled

	if err := uc.repo.CreateAppointment(ctx, a); err != nil {
		return err
	}

	// Audit event
	event := &domain.AppointmentAuditEvent{
		ID:            uuid.New(),
		TenantID:      a.TenantID,
		AppointmentID: a.ID,
		ActorUserID:   actorID,
		EventType:     "scheduled",
		NewStatus:     string(domain.StatusScheduled),
		CreatedAt:     now,
	}
	_ = uc.repo.CreateAuditEvent(ctx, event)

	return nil
}

// GetAppointmentByID returns an appointment by ID.
func (uc *Usecase) GetAppointmentByID(ctx context.Context, tenantID, id uuid.UUID) (*domain.Appointment, error) {
	return uc.repo.GetAppointmentByID(ctx, tenantID, id)
}

// ListAppointments returns appointments with optional filters.
func (uc *Usecase) ListAppointments(ctx context.Context, tenantID uuid.UUID, professionalID *uuid.UUID, from, to *time.Time, status *domain.AppointmentStatus) ([]domain.Appointment, error) {
	return uc.repo.ListAppointments(ctx, tenantID, professionalID, from, to, status)
}

// UpdateAppointmentStatus changes an appointment's status, enforcing valid transitions.
func (uc *Usecase) UpdateAppointmentStatus(ctx context.Context, tenantID, id uuid.UUID, newStatus domain.AppointmentStatus, cancellationReason string, actorID *uuid.UUID) error {
	appt, err := uc.repo.GetAppointmentByID(ctx, tenantID, id)
	if err != nil {
		return err
	}

	if !domain.IsValidTransition(appt.Status, newStatus) {
		return domain.ErrInvalidTransition
	}

	if err := uc.repo.UpdateAppointmentStatus(ctx, tenantID, id, newStatus, cancellationReason); err != nil {
		return err
	}

	// Audit event
	event := &domain.AppointmentAuditEvent{
		ID:            uuid.New(),
		TenantID:      tenantID,
		AppointmentID: id,
		ActorUserID:   actorID,
		EventType:     string(newStatus),
		OldStatus:     string(appt.Status),
		NewStatus:     string(newStatus),
		CreatedAt:     time.Now().UTC(),
	}
	_ = uc.repo.CreateAuditEvent(ctx, event)

	return nil
}

// RescheduleAppointment reschedules an appointment to a new time, checking for conflicts.
func (uc *Usecase) RescheduleAppointment(ctx context.Context, tenantID, id uuid.UUID, newStart, newEnd time.Time, actorID *uuid.UUID) error {
	appt, err := uc.repo.GetAppointmentByID(ctx, tenantID, id)
	if err != nil {
		return err
	}

	// Only scheduled or confirmed can be rescheduled
	if appt.Status != domain.StatusScheduled && appt.Status != domain.StatusConfirmed {
		return domain.ErrInvalidTransition
	}

	// Check for conflicts (exclude current appointment)
	conflicts, err := uc.repo.FindConflictingAppointments(ctx, tenantID, appt.ProfessionalID, newStart, newEnd, &id)
	if err != nil {
		return err
	}
	if len(conflicts) > 0 {
		return domain.ErrConflict
	}

	oldStart := appt.StartAt
	if err := uc.repo.RescheduleAppointment(ctx, tenantID, id, newStart, newEnd); err != nil {
		return err
	}

	// Audit event
	event := &domain.AppointmentAuditEvent{
		ID:            uuid.New(),
		TenantID:      tenantID,
		AppointmentID: id,
		ActorUserID:   actorID,
		EventType:     "rescheduled",
		OldStartAt:    &oldStart,
		NewStartAt:    &newStart,
		CreatedAt:     time.Now().UTC(),
	}
	_ = uc.repo.CreateAuditEvent(ctx, event)

	return nil
}

// GetAvailableSlots computes free time slots for a professional on a given date.
func (uc *Usecase) GetAvailableSlots(ctx context.Context, tenantID, professionalID uuid.UUID, date time.Time, slotDurationMin int) ([]domain.TimeSlot, error) {
	dayOfWeek := int(date.Weekday())

	rules, err := uc.repo.GetAvailabilityRulesByDay(ctx, tenantID, professionalID, dayOfWeek)
	if err != nil {
		return nil, err
	}
	if len(rules) == 0 {
		return []domain.TimeSlot{}, nil
	}

	// Get blocks for this day
	dayStart := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	dayEnd := dayStart.Add(24 * time.Hour)
	blocks, err := uc.repo.ListBlocks(ctx, tenantID, professionalID, dayStart, dayEnd)
	if err != nil {
		return nil, err
	}

	// Get existing appointments for this day
	appts, err := uc.repo.FindConflictingAppointments(ctx, tenantID, professionalID, dayStart, dayEnd, nil)
	if err != nil {
		return nil, err
	}

	if slotDurationMin <= 0 {
		slotDurationMin = 50
	}
	slotDuration := time.Duration(slotDurationMin) * time.Minute

	var slots []domain.TimeSlot
	for _, rule := range rules {
		// Convert rule times to absolute times on this date
		ruleStart := time.Date(date.Year(), date.Month(), date.Day(),
			rule.StartTime.Hour(), rule.StartTime.Minute(), 0, 0, date.Location())
		ruleEnd := time.Date(date.Year(), date.Month(), date.Day(),
			rule.EndTime.Hour(), rule.EndTime.Minute(), 0, 0, date.Location())

		// Generate slots within this rule
		for slotStart := ruleStart; slotStart.Add(slotDuration).Before(ruleEnd) || slotStart.Add(slotDuration).Equal(ruleEnd); slotStart = slotStart.Add(slotDuration) {
			slotEnd := slotStart.Add(slotDuration)
			if isSlotFree(slotStart, slotEnd, blocks, appts) {
				slots = append(slots, domain.TimeSlot{
					Start:       slotStart,
					End:         slotEnd,
					ServiceMode: rule.ServiceMode,
					AddressID:   rule.AddressID,
				})
			}
		}
	}

	return slots, nil
}

// isSlotFree checks if a slot doesn't overlap with any blocks or appointments.
func isSlotFree(start, end time.Time, blocks []domain.ScheduleBlock, appts []domain.Appointment) bool {
	for _, b := range blocks {
		if start.Before(b.EndAt) && end.After(b.StartAt) {
			return false
		}
	}
	for _, a := range appts {
		if start.Before(a.EndAt) && end.After(a.StartAt) {
			return false
		}
	}
	return true
}
