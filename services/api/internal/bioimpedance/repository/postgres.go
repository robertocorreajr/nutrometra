package repository

import (
	"context"
	"errors"
	"fmt"

	"nutrometra/api/internal/bioimpedance/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository provides bioimpedance data access.
type Repository struct {
	pool *pgxpool.Pool
}

// New creates a bioimpedance Repository.
func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// Create inserts a new body measurement.
func (r *Repository) Create(ctx context.Context, m *domain.BodyMeasurement) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO body_measurements
			(id, tenant_id, patient_id, professional_id, measured_at,
			 weight_kg, height_cm, bmi,
			 body_fat_pct, lean_mass_kg, fat_mass_kg, muscle_mass_kg, bone_mass_kg,
			 water_pct, visceral_fat, basal_metabolic_rate,
			 waist_cm, hip_cm, chest_cm, right_arm_cm, left_arm_cm,
			 right_thigh_cm, left_thigh_cm, right_calf_cm, left_calf_cm,
			 neck_cm, abdomen_cm,
			 triceps_sf_mm, biceps_sf_mm, subscapular_sf_mm, suprailiac_sf_mm,
			 abdominal_sf_mm, thigh_sf_mm, calf_sf_mm,
			 source, device_model, notes, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26,$27,$28,$29,$30,$31,$32,$33,$34,$35,$36,$37,$38,$39)`,
		m.ID, m.TenantID, m.PatientID, m.ProfessionalID, m.MeasuredAt,
		m.WeightKg, m.HeightCm, m.BMI,
		m.BodyFatPct, m.LeanMassKg, m.FatMassKg, m.MuscleMassKg, m.BoneMassKg,
		m.WaterPct, m.VisceralFat, m.BasalMetabolicRate,
		m.WaistCm, m.HipCm, m.ChestCm, m.RightArmCm, m.LeftArmCm,
		m.RightThighCm, m.LeftThighCm, m.RightCalfCm, m.LeftCalfCm,
		m.NeckCm, m.AbdomenCm,
		m.TricepsSfMm, m.BicepsSfMm, m.SubscapularSfMm, m.SuprailiacSfMm,
		m.AbdominalSfMm, m.ThighSfMm, m.CalfSfMm,
		string(m.Source), nilIfEmpty(m.DeviceModel), nilIfEmpty(m.Notes),
		m.CreatedAt, m.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("bioimpedance: create: %w", err)
	}
	return nil
}

// GetByID returns a measurement by ID.
func (r *Repository) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*domain.BodyMeasurement, error) {
	m := &domain.BodyMeasurement{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, tenant_id, patient_id, professional_id, measured_at,
			weight_kg, height_cm, bmi,
			body_fat_pct, lean_mass_kg, fat_mass_kg, muscle_mass_kg, bone_mass_kg,
			water_pct, visceral_fat, basal_metabolic_rate,
			waist_cm, hip_cm, chest_cm, right_arm_cm, left_arm_cm,
			right_thigh_cm, left_thigh_cm, right_calf_cm, left_calf_cm,
			neck_cm, abdomen_cm,
			triceps_sf_mm, biceps_sf_mm, subscapular_sf_mm, suprailiac_sf_mm,
			abdominal_sf_mm, thigh_sf_mm, calf_sf_mm,
			source, COALESCE(device_model,''), COALESCE(notes,''),
			created_at, updated_at
		 FROM body_measurements WHERE id = $1 AND tenant_id = $2`,
		id, tenantID,
	).Scan(&m.ID, &m.TenantID, &m.PatientID, &m.ProfessionalID, &m.MeasuredAt,
		&m.WeightKg, &m.HeightCm, &m.BMI,
		&m.BodyFatPct, &m.LeanMassKg, &m.FatMassKg, &m.MuscleMassKg, &m.BoneMassKg,
		&m.WaterPct, &m.VisceralFat, &m.BasalMetabolicRate,
		&m.WaistCm, &m.HipCm, &m.ChestCm, &m.RightArmCm, &m.LeftArmCm,
		&m.RightThighCm, &m.LeftThighCm, &m.RightCalfCm, &m.LeftCalfCm,
		&m.NeckCm, &m.AbdomenCm,
		&m.TricepsSfMm, &m.BicepsSfMm, &m.SubscapularSfMm, &m.SuprailiacSfMm,
		&m.AbdominalSfMm, &m.ThighSfMm, &m.CalfSfMm,
		&m.Source, &m.DeviceModel, &m.Notes,
		&m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("bioimpedance: get_by_id: %w", err)
	}
	return m, nil
}

// ListByPatient returns measurements for a patient ordered by date desc.
func (r *Repository) ListByPatient(ctx context.Context, tenantID, patientID uuid.UUID) ([]domain.BodyMeasurement, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, tenant_id, patient_id, professional_id, measured_at,
			weight_kg, height_cm, bmi,
			body_fat_pct, lean_mass_kg, fat_mass_kg, muscle_mass_kg, bone_mass_kg,
			water_pct, visceral_fat, basal_metabolic_rate,
			waist_cm, hip_cm, chest_cm, right_arm_cm, left_arm_cm,
			right_thigh_cm, left_thigh_cm, right_calf_cm, left_calf_cm,
			neck_cm, abdomen_cm,
			triceps_sf_mm, biceps_sf_mm, subscapular_sf_mm, suprailiac_sf_mm,
			abdominal_sf_mm, thigh_sf_mm, calf_sf_mm,
			source, COALESCE(device_model,''), COALESCE(notes,''),
			created_at, updated_at
		 FROM body_measurements WHERE tenant_id = $1 AND patient_id = $2
		 ORDER BY measured_at DESC`,
		tenantID, patientID,
	)
	if err != nil {
		return nil, fmt.Errorf("bioimpedance: list: %w", err)
	}
	defer rows.Close()

	var result []domain.BodyMeasurement
	for rows.Next() {
		var m domain.BodyMeasurement
		if err := rows.Scan(&m.ID, &m.TenantID, &m.PatientID, &m.ProfessionalID, &m.MeasuredAt,
			&m.WeightKg, &m.HeightCm, &m.BMI,
			&m.BodyFatPct, &m.LeanMassKg, &m.FatMassKg, &m.MuscleMassKg, &m.BoneMassKg,
			&m.WaterPct, &m.VisceralFat, &m.BasalMetabolicRate,
			&m.WaistCm, &m.HipCm, &m.ChestCm, &m.RightArmCm, &m.LeftArmCm,
			&m.RightThighCm, &m.LeftThighCm, &m.RightCalfCm, &m.LeftCalfCm,
			&m.NeckCm, &m.AbdomenCm,
			&m.TricepsSfMm, &m.BicepsSfMm, &m.SubscapularSfMm, &m.SuprailiacSfMm,
			&m.AbdominalSfMm, &m.ThighSfMm, &m.CalfSfMm,
			&m.Source, &m.DeviceModel, &m.Notes,
			&m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, fmt.Errorf("bioimpedance: scan: %w", err)
		}
		result = append(result, m)
	}
	return result, rows.Err()
}

// PublishMeasurement creates a publication record for a measurement.
func (r *Repository) PublishMeasurement(ctx context.Context, pub *domain.MeasurementPublication) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO measurement_publications
			(id, tenant_id, measurement_id, patient_id, published_by, published_at, message)
		 VALUES ($1,$2,$3,$4,$5,$6,$7)`,
		pub.ID, pub.TenantID, pub.MeasurementID, pub.PatientID,
		pub.PublishedBy, pub.PublishedAt, nilIfEmpty(pub.Message),
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.ErrAlreadyPublished
		}
		return fmt.Errorf("bioimpedance: publish: %w", err)
	}
	return nil
}

func nilIfEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
