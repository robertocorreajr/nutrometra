package usecase

import (
	"context"
	"time"

	"nutrometra/api/internal/bioimpedance/domain"

	"github.com/google/uuid"
)

// BioimpedanceRepository defines the data access contract.
type BioimpedanceRepository interface {
	Create(ctx context.Context, m *domain.BodyMeasurement) error
	GetByID(ctx context.Context, tenantID, id uuid.UUID) (*domain.BodyMeasurement, error)
	ListByPatient(ctx context.Context, tenantID, patientID uuid.UUID) ([]domain.BodyMeasurement, error)
	PublishMeasurement(ctx context.Context, pub *domain.MeasurementPublication) error
}

// Usecase contains business logic for bioimpedance.
type Usecase struct {
	repo BioimpedanceRepository
}

// New creates a bioimpedance Usecase.
func New(repo BioimpedanceRepository) *Usecase {
	return &Usecase{repo: repo}
}

// Create creates a new measurement, validating ranges and auto-calculating BMI.
func (uc *Usecase) Create(ctx context.Context, m *domain.BodyMeasurement) error {
	if err := m.Validate(); err != nil {
		return err
	}

	m.CalculateBMI()
	m.ID = uuid.New()
	now := time.Now().UTC()
	m.CreatedAt = now
	m.UpdatedAt = now
	if m.MeasuredAt.IsZero() {
		m.MeasuredAt = now
	}
	if m.Source == "" {
		m.Source = domain.SourceManual
	}

	return uc.repo.Create(ctx, m)
}

// GetByID returns a measurement by ID.
func (uc *Usecase) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*domain.BodyMeasurement, error) {
	return uc.repo.GetByID(ctx, tenantID, id)
}

// ListByPatient returns the measurement history for a patient.
func (uc *Usecase) ListByPatient(ctx context.Context, tenantID, patientID uuid.UUID) ([]domain.BodyMeasurement, error) {
	return uc.repo.ListByPatient(ctx, tenantID, patientID)
}

// PublishToPatient publishes a measurement so the patient can see it.
func (uc *Usecase) PublishToPatient(ctx context.Context, tenantID, measurementID, publishedBy uuid.UUID, message string) error {
	m, err := uc.repo.GetByID(ctx, tenantID, measurementID)
	if err != nil {
		return err
	}

	pub := &domain.MeasurementPublication{
		ID:            uuid.New(),
		TenantID:      tenantID,
		MeasurementID: measurementID,
		PatientID:     m.PatientID,
		PublishedBy:   publishedBy,
		PublishedAt:   time.Now().UTC(),
		Message:       message,
	}

	return uc.repo.PublishMeasurement(ctx, pub)
}
