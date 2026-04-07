package repository

import (
	"context"
	"errors"
	"fmt"

	"nutrometra/api/internal/clinical/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository provides clinical data access.
type Repository struct {
	pool *pgxpool.Pool
}

// New creates a clinical Repository.
func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// --- Anamnesis ---

// CreateAnamnesis inserts a new anamnesis.
func (r *Repository) CreateAnamnesis(ctx context.Context, a *domain.Anamnesis) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO anamneses
			(id, tenant_id, patient_id, professional_id, status,
			 chief_complaint, history_present_illness, past_medical_history,
			 family_history, social_history, dietary_history,
			 physical_activity, sleep_pattern, bowel_habits,
			 water_intake, supplements, observations,
			 finalized_at, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20)`,
		a.ID, a.TenantID, a.PatientID, a.ProfessionalID, string(a.Status),
		nilIfEmpty(a.ChiefComplaint), nilIfEmpty(a.HistoryPresentIllness),
		nilIfEmpty(a.PastMedicalHistory), nilIfEmpty(a.FamilyHistory),
		nilIfEmpty(a.SocialHistory), nilIfEmpty(a.DietaryHistory),
		nilIfEmpty(a.PhysicalActivity), nilIfEmpty(a.SleepPattern),
		nilIfEmpty(a.BowelHabits), nilIfEmpty(a.WaterIntake),
		nilIfEmpty(a.Supplements), nilIfEmpty(a.Observations),
		a.FinalizedAt, a.CreatedAt, a.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("clinical: create_anamnesis: %w", err)
	}
	return nil
}

// GetAnamnesisByID returns an anamnesis by ID.
func (r *Repository) GetAnamnesisByID(ctx context.Context, tenantID, id uuid.UUID) (*domain.Anamnesis, error) {
	a := &domain.Anamnesis{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, tenant_id, patient_id, professional_id, status,
			COALESCE(chief_complaint,''), COALESCE(history_present_illness,''),
			COALESCE(past_medical_history,''), COALESCE(family_history,''),
			COALESCE(social_history,''), COALESCE(dietary_history,''),
			COALESCE(physical_activity,''), COALESCE(sleep_pattern,''),
			COALESCE(bowel_habits,''), COALESCE(water_intake,''),
			COALESCE(supplements,''), COALESCE(observations,''),
			finalized_at, created_at, updated_at
		 FROM anamneses WHERE id = $1 AND tenant_id = $2`,
		id, tenantID,
	).Scan(&a.ID, &a.TenantID, &a.PatientID, &a.ProfessionalID, &a.Status,
		&a.ChiefComplaint, &a.HistoryPresentIllness, &a.PastMedicalHistory,
		&a.FamilyHistory, &a.SocialHistory, &a.DietaryHistory,
		&a.PhysicalActivity, &a.SleepPattern, &a.BowelHabits,
		&a.WaterIntake, &a.Supplements, &a.Observations,
		&a.FinalizedAt, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("clinical: get_anamnesis: %w", err)
	}
	return a, nil
}

// ListAnamnesesByPatient returns anamneses for a patient.
func (r *Repository) ListAnamnesesByPatient(ctx context.Context, tenantID, patientID uuid.UUID) ([]domain.Anamnesis, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, tenant_id, patient_id, professional_id, status,
			COALESCE(chief_complaint,''), COALESCE(history_present_illness,''),
			COALESCE(past_medical_history,''), COALESCE(family_history,''),
			COALESCE(social_history,''), COALESCE(dietary_history,''),
			COALESCE(physical_activity,''), COALESCE(sleep_pattern,''),
			COALESCE(bowel_habits,''), COALESCE(water_intake,''),
			COALESCE(supplements,''), COALESCE(observations,''),
			finalized_at, created_at, updated_at
		 FROM anamneses WHERE tenant_id = $1 AND patient_id = $2
		 ORDER BY created_at DESC`,
		tenantID, patientID,
	)
	if err != nil {
		return nil, fmt.Errorf("clinical: list_anamneses: %w", err)
	}
	defer rows.Close()

	var result []domain.Anamnesis
	for rows.Next() {
		var a domain.Anamnesis
		if err := rows.Scan(&a.ID, &a.TenantID, &a.PatientID, &a.ProfessionalID, &a.Status,
			&a.ChiefComplaint, &a.HistoryPresentIllness, &a.PastMedicalHistory,
			&a.FamilyHistory, &a.SocialHistory, &a.DietaryHistory,
			&a.PhysicalActivity, &a.SleepPattern, &a.BowelHabits,
			&a.WaterIntake, &a.Supplements, &a.Observations,
			&a.FinalizedAt, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, fmt.Errorf("clinical: scan_anamnesis: %w", err)
		}
		result = append(result, a)
	}
	return result, rows.Err()
}

// UpdateAnamnesis updates a draft anamnesis.
func (r *Repository) UpdateAnamnesis(ctx context.Context, a *domain.Anamnesis) error {
	tag, err := r.pool.Exec(ctx,
		`UPDATE anamneses SET
			chief_complaint = $3, history_present_illness = $4,
			past_medical_history = $5, family_history = $6,
			social_history = $7, dietary_history = $8,
			physical_activity = $9, sleep_pattern = $10,
			bowel_habits = $11, water_intake = $12,
			supplements = $13, observations = $14,
			status = $15, finalized_at = $16, updated_at = NOW()
		 WHERE id = $1 AND tenant_id = $2`,
		a.ID, a.TenantID,
		nilIfEmpty(a.ChiefComplaint), nilIfEmpty(a.HistoryPresentIllness),
		nilIfEmpty(a.PastMedicalHistory), nilIfEmpty(a.FamilyHistory),
		nilIfEmpty(a.SocialHistory), nilIfEmpty(a.DietaryHistory),
		nilIfEmpty(a.PhysicalActivity), nilIfEmpty(a.SleepPattern),
		nilIfEmpty(a.BowelHabits), nilIfEmpty(a.WaterIntake),
		nilIfEmpty(a.Supplements), nilIfEmpty(a.Observations),
		string(a.Status), a.FinalizedAt,
	)
	if err != nil {
		return fmt.Errorf("clinical: update_anamnesis: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// --- Progress Notes ---

// CreateProgressNote inserts a progress note.
func (r *Repository) CreateProgressNote(ctx context.Context, n *domain.ProgressNote) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO progress_notes
			(id, tenant_id, patient_id, professional_id, appointment_id,
			 title, content, visible_to_patient, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
		n.ID, n.TenantID, n.PatientID, n.ProfessionalID, n.AppointmentID,
		n.Title, n.Content, n.VisibleToPatient, n.CreatedAt, n.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("clinical: create_progress_note: %w", err)
	}
	return nil
}

// ListProgressNotesByPatient returns progress notes for a patient.
func (r *Repository) ListProgressNotesByPatient(ctx context.Context, tenantID, patientID uuid.UUID) ([]domain.ProgressNote, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, tenant_id, patient_id, professional_id, appointment_id,
			title, content, visible_to_patient, created_at, updated_at
		 FROM progress_notes WHERE tenant_id = $1 AND patient_id = $2
		 ORDER BY created_at DESC`,
		tenantID, patientID,
	)
	if err != nil {
		return nil, fmt.Errorf("clinical: list_progress_notes: %w", err)
	}
	defer rows.Close()

	var result []domain.ProgressNote
	for rows.Next() {
		var n domain.ProgressNote
		if err := rows.Scan(&n.ID, &n.TenantID, &n.PatientID, &n.ProfessionalID,
			&n.AppointmentID, &n.Title, &n.Content, &n.VisibleToPatient,
			&n.CreatedAt, &n.UpdatedAt); err != nil {
			return nil, fmt.Errorf("clinical: scan_progress_note: %w", err)
		}
		result = append(result, n)
	}
	return result, rows.Err()
}

// UpdateProgressNote updates a progress note.
func (r *Repository) UpdateProgressNote(ctx context.Context, n *domain.ProgressNote) error {
	tag, err := r.pool.Exec(ctx,
		`UPDATE progress_notes SET
			title = $3, content = $4, visible_to_patient = $5, updated_at = NOW()
		 WHERE id = $1 AND tenant_id = $2`,
		n.ID, n.TenantID, n.Title, n.Content, n.VisibleToPatient,
	)
	if err != nil {
		return fmt.Errorf("clinical: update_progress_note: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// --- Clinical Attachments ---

// CreateAttachment inserts a clinical attachment.
func (r *Repository) CreateAttachment(ctx context.Context, a *domain.ClinicalAttachment) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO clinical_attachments
			(id, tenant_id, patient_id, professional_id, file_name, file_type,
			 file_size_bytes, storage_key, category, description, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
		a.ID, a.TenantID, a.PatientID, a.ProfessionalID,
		a.FileName, a.FileType, a.FileSizeBytes, a.StorageKey,
		string(a.Category), nilIfEmpty(a.Description), a.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("clinical: create_attachment: %w", err)
	}
	return nil
}

// ListAttachmentsByPatient returns attachments for a patient.
func (r *Repository) ListAttachmentsByPatient(ctx context.Context, tenantID, patientID uuid.UUID) ([]domain.ClinicalAttachment, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, tenant_id, patient_id, professional_id, file_name, file_type,
			file_size_bytes, storage_key, category, COALESCE(description,''), created_at
		 FROM clinical_attachments WHERE tenant_id = $1 AND patient_id = $2
		 ORDER BY created_at DESC`,
		tenantID, patientID,
	)
	if err != nil {
		return nil, fmt.Errorf("clinical: list_attachments: %w", err)
	}
	defer rows.Close()

	var result []domain.ClinicalAttachment
	for rows.Next() {
		var a domain.ClinicalAttachment
		if err := rows.Scan(&a.ID, &a.TenantID, &a.PatientID, &a.ProfessionalID,
			&a.FileName, &a.FileType, &a.FileSizeBytes, &a.StorageKey,
			&a.Category, &a.Description, &a.CreatedAt); err != nil {
			return nil, fmt.Errorf("clinical: scan_attachment: %w", err)
		}
		result = append(result, a)
	}
	return result, rows.Err()
}

func nilIfEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
