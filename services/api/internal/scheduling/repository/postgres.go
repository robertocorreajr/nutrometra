package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"nutrometra/api/internal/scheduling/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository provides scheduling data access.
type Repository struct {
	pool *pgxpool.Pool
}

// New creates a scheduling Repository.
func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// --- Availability Rules ---

// CreateAvailabilityRule inserts a new availability rule.
func (r *Repository) CreateAvailabilityRule(ctx context.Context, rule *domain.AvailabilityRule) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO availability_rules
			(id, tenant_id, professional_id, day_of_week, start_time, end_time,
			 service_mode, address_id, active, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
		rule.ID, rule.TenantID, rule.ProfessionalID, rule.DayOfWeek,
		rule.StartTime.Format("15:04:05"), rule.EndTime.Format("15:04:05"),
		string(rule.ServiceMode), rule.AddressID, rule.Active,
		rule.CreatedAt, rule.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("scheduling: create_availability_rule: %w", err)
	}
	return nil
}

// ListAvailabilityRules returns rules for a professional.
func (r *Repository) ListAvailabilityRules(ctx context.Context, tenantID, professionalID uuid.UUID) ([]domain.AvailabilityRule, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, tenant_id, professional_id, day_of_week, start_time, end_time,
			service_mode, address_id, active, created_at, updated_at
		 FROM availability_rules
		 WHERE tenant_id = $1 AND professional_id = $2
		 ORDER BY day_of_week, start_time`,
		tenantID, professionalID,
	)
	if err != nil {
		return nil, fmt.Errorf("scheduling: list_availability_rules: %w", err)
	}
	defer rows.Close()

	var result []domain.AvailabilityRule
	for rows.Next() {
		var rule domain.AvailabilityRule
		if err := rows.Scan(&rule.ID, &rule.TenantID, &rule.ProfessionalID,
			&rule.DayOfWeek, &rule.StartTime, &rule.EndTime,
			&rule.ServiceMode, &rule.AddressID, &rule.Active,
			&rule.CreatedAt, &rule.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scheduling: scan_availability_rule: %w", err)
		}
		result = append(result, rule)
	}
	return result, rows.Err()
}

// GetAvailabilityRulesByDay returns active rules for a professional on a specific day of week.
func (r *Repository) GetAvailabilityRulesByDay(ctx context.Context, tenantID, professionalID uuid.UUID, dayOfWeek int) ([]domain.AvailabilityRule, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, tenant_id, professional_id, day_of_week, start_time, end_time,
			service_mode, address_id, active, created_at, updated_at
		 FROM availability_rules
		 WHERE tenant_id = $1 AND professional_id = $2 AND day_of_week = $3 AND active = TRUE
		 ORDER BY start_time`,
		tenantID, professionalID, dayOfWeek,
	)
	if err != nil {
		return nil, fmt.Errorf("scheduling: get_rules_by_day: %w", err)
	}
	defer rows.Close()

	var result []domain.AvailabilityRule
	for rows.Next() {
		var rule domain.AvailabilityRule
		if err := rows.Scan(&rule.ID, &rule.TenantID, &rule.ProfessionalID,
			&rule.DayOfWeek, &rule.StartTime, &rule.EndTime,
			&rule.ServiceMode, &rule.AddressID, &rule.Active,
			&rule.CreatedAt, &rule.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scheduling: scan_availability_rule: %w", err)
		}
		result = append(result, rule)
	}
	return result, rows.Err()
}

// UpdateAvailabilityRule updates a rule.
func (r *Repository) UpdateAvailabilityRule(ctx context.Context, rule *domain.AvailabilityRule) error {
	tag, err := r.pool.Exec(ctx,
		`UPDATE availability_rules SET
			day_of_week = $3, start_time = $4, end_time = $5,
			service_mode = $6, address_id = $7, active = $8, updated_at = NOW()
		 WHERE id = $1 AND tenant_id = $2`,
		rule.ID, rule.TenantID, rule.DayOfWeek,
		rule.StartTime.Format("15:04:05"), rule.EndTime.Format("15:04:05"),
		string(rule.ServiceMode), rule.AddressID, rule.Active,
	)
	if err != nil {
		return fmt.Errorf("scheduling: update_availability_rule: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// DeleteAvailabilityRule soft-deletes a rule.
func (r *Repository) DeleteAvailabilityRule(ctx context.Context, tenantID, id uuid.UUID) error {
	tag, err := r.pool.Exec(ctx,
		`UPDATE availability_rules SET active = FALSE, updated_at = NOW() WHERE id = $1 AND tenant_id = $2`,
		id, tenantID,
	)
	if err != nil {
		return fmt.Errorf("scheduling: delete_availability_rule: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// --- Schedule Blocks ---

// CreateBlock inserts a schedule block.
func (r *Repository) CreateBlock(ctx context.Context, b *domain.ScheduleBlock) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO schedule_blocks
			(id, tenant_id, professional_id, start_at, end_at, reason, all_day, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		b.ID, b.TenantID, b.ProfessionalID, b.StartAt, b.EndAt,
		nilIfEmpty(b.Reason), b.AllDay, b.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("scheduling: create_block: %w", err)
	}
	return nil
}

// ListBlocks returns schedule blocks for a professional in a date range.
func (r *Repository) ListBlocks(ctx context.Context, tenantID, professionalID uuid.UUID, from, to time.Time) ([]domain.ScheduleBlock, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, tenant_id, professional_id, start_at, end_at,
			COALESCE(reason,''), all_day, created_at
		 FROM schedule_blocks
		 WHERE tenant_id = $1 AND professional_id = $2
		   AND start_at < $4 AND end_at > $3
		 ORDER BY start_at`,
		tenantID, professionalID, from, to,
	)
	if err != nil {
		return nil, fmt.Errorf("scheduling: list_blocks: %w", err)
	}
	defer rows.Close()

	var result []domain.ScheduleBlock
	for rows.Next() {
		var b domain.ScheduleBlock
		if err := rows.Scan(&b.ID, &b.TenantID, &b.ProfessionalID,
			&b.StartAt, &b.EndAt, &b.Reason, &b.AllDay, &b.CreatedAt); err != nil {
			return nil, fmt.Errorf("scheduling: scan_block: %w", err)
		}
		result = append(result, b)
	}
	return result, rows.Err()
}

// DeleteBlock removes a schedule block.
func (r *Repository) DeleteBlock(ctx context.Context, tenantID, id uuid.UUID) error {
	tag, err := r.pool.Exec(ctx,
		`DELETE FROM schedule_blocks WHERE id = $1 AND tenant_id = $2`,
		id, tenantID,
	)
	if err != nil {
		return fmt.Errorf("scheduling: delete_block: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// --- Appointments ---

// CreateAppointment inserts a new appointment.
func (r *Repository) CreateAppointment(ctx context.Context, a *domain.Appointment) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO appointments
			(id, tenant_id, professional_id, patient_id, start_at, end_at,
			 service_mode, address_id, status, source, notes, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)`,
		a.ID, a.TenantID, a.ProfessionalID, a.PatientID,
		a.StartAt, a.EndAt, string(a.ServiceMode), a.AddressID,
		string(a.Status), string(a.Source), nilIfEmpty(a.Notes),
		a.CreatedAt, a.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("scheduling: create_appointment: %w", err)
	}
	return nil
}

// GetAppointmentByID returns an appointment by ID.
func (r *Repository) GetAppointmentByID(ctx context.Context, tenantID, id uuid.UUID) (*domain.Appointment, error) {
	a := &domain.Appointment{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, tenant_id, professional_id, patient_id, start_at, end_at,
			service_mode, address_id, status, source, COALESCE(notes,''),
			COALESCE(cancellation_reason,''), created_at, updated_at
		 FROM appointments WHERE id = $1 AND tenant_id = $2`,
		id, tenantID,
	).Scan(&a.ID, &a.TenantID, &a.ProfessionalID, &a.PatientID,
		&a.StartAt, &a.EndAt, &a.ServiceMode, &a.AddressID,
		&a.Status, &a.Source, &a.Notes, &a.CancellationReason,
		&a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("scheduling: get_appointment: %w", err)
	}
	return a, nil
}

// ListAppointments returns appointments for a tenant with optional filters.
func (r *Repository) ListAppointments(ctx context.Context, tenantID uuid.UUID, professionalID *uuid.UUID, from, to *time.Time, status *domain.AppointmentStatus) ([]domain.Appointment, error) {
	query := `SELECT id, tenant_id, professional_id, patient_id, start_at, end_at,
		service_mode, address_id, status, source, COALESCE(notes,''),
		COALESCE(cancellation_reason,''), created_at, updated_at
	 FROM appointments WHERE tenant_id = $1`
	args := []any{tenantID}
	argIdx := 2

	if professionalID != nil {
		query += fmt.Sprintf(" AND professional_id = $%d", argIdx)
		args = append(args, *professionalID)
		argIdx++
	}
	if from != nil {
		query += fmt.Sprintf(" AND start_at >= $%d", argIdx)
		args = append(args, *from)
		argIdx++
	}
	if to != nil {
		query += fmt.Sprintf(" AND start_at <= $%d", argIdx)
		args = append(args, *to)
		argIdx++
	}
	if status != nil {
		query += fmt.Sprintf(" AND status = $%d", argIdx)
		args = append(args, string(*status))
		argIdx++
	}
	query += " ORDER BY start_at"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("scheduling: list_appointments: %w", err)
	}
	defer rows.Close()

	var result []domain.Appointment
	for rows.Next() {
		var a domain.Appointment
		if err := rows.Scan(&a.ID, &a.TenantID, &a.ProfessionalID, &a.PatientID,
			&a.StartAt, &a.EndAt, &a.ServiceMode, &a.AddressID,
			&a.Status, &a.Source, &a.Notes, &a.CancellationReason,
			&a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scheduling: scan_appointment: %w", err)
		}
		result = append(result, a)
	}
	return result, rows.Err()
}

// UpdateAppointmentStatus updates an appointment's status.
func (r *Repository) UpdateAppointmentStatus(ctx context.Context, tenantID, id uuid.UUID, status domain.AppointmentStatus, cancellationReason string) error {
	tag, err := r.pool.Exec(ctx,
		`UPDATE appointments SET status = $3, cancellation_reason = $4, updated_at = NOW()
		 WHERE id = $1 AND tenant_id = $2`,
		id, tenantID, string(status), nilIfEmpty(cancellationReason),
	)
	if err != nil {
		return fmt.Errorf("scheduling: update_appointment_status: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// RescheduleAppointment updates an appointment's time.
func (r *Repository) RescheduleAppointment(ctx context.Context, tenantID, id uuid.UUID, startAt, endAt time.Time) error {
	tag, err := r.pool.Exec(ctx,
		`UPDATE appointments SET start_at = $3, end_at = $4, updated_at = NOW()
		 WHERE id = $1 AND tenant_id = $2`,
		id, tenantID, startAt, endAt,
	)
	if err != nil {
		return fmt.Errorf("scheduling: reschedule_appointment: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// FindConflictingAppointments returns appointments that overlap with the given time range.
func (r *Repository) FindConflictingAppointments(ctx context.Context, tenantID, professionalID uuid.UUID, startAt, endAt time.Time, excludeID *uuid.UUID) ([]domain.Appointment, error) {
	query := `SELECT id, tenant_id, professional_id, patient_id, start_at, end_at,
		service_mode, address_id, status, source, COALESCE(notes,''),
		COALESCE(cancellation_reason,''), created_at, updated_at
	 FROM appointments
	 WHERE tenant_id = $1 AND professional_id = $2
	   AND status NOT IN ('cancelled','no_show')
	   AND start_at < $4 AND end_at > $3`
	args := []any{tenantID, professionalID, startAt, endAt}

	if excludeID != nil {
		query += " AND id <> $5"
		args = append(args, *excludeID)
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("scheduling: find_conflicts: %w", err)
	}
	defer rows.Close()

	var result []domain.Appointment
	for rows.Next() {
		var a domain.Appointment
		if err := rows.Scan(&a.ID, &a.TenantID, &a.ProfessionalID, &a.PatientID,
			&a.StartAt, &a.EndAt, &a.ServiceMode, &a.AddressID,
			&a.Status, &a.Source, &a.Notes, &a.CancellationReason,
			&a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scheduling: scan_conflict: %w", err)
		}
		result = append(result, a)
	}
	return result, rows.Err()
}

// CreateAuditEvent inserts an appointment audit event.
func (r *Repository) CreateAuditEvent(ctx context.Context, e *domain.AppointmentAuditEvent) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO appointment_audit_events
			(id, tenant_id, appointment_id, actor_user_id, event_type,
			 old_status, new_status, old_start_at, new_start_at, notes, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
		e.ID, e.TenantID, e.AppointmentID, e.ActorUserID, e.EventType,
		nilIfEmpty(e.OldStatus), nilIfEmpty(e.NewStatus),
		e.OldStartAt, e.NewStartAt, nilIfEmpty(e.Notes), e.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("scheduling: create_audit_event: %w", err)
	}
	return nil
}

func nilIfEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
