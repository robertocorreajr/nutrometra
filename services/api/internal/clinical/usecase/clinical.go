package usecase

import (
	"context"
	"time"

	"nutrometra/api/internal/clinical/domain"

	"github.com/google/uuid"
)

// ClinicalRepository defines the data access contract.
type ClinicalRepository interface {
	CreateAnamnesis(ctx context.Context, a *domain.Anamnesis) error
	GetAnamnesisByID(ctx context.Context, tenantID, id uuid.UUID) (*domain.Anamnesis, error)
	ListAnamnesesByPatient(ctx context.Context, tenantID, patientID uuid.UUID) ([]domain.Anamnesis, error)
	UpdateAnamnesis(ctx context.Context, a *domain.Anamnesis) error

	CreateProgressNote(ctx context.Context, n *domain.ProgressNote) error
	ListProgressNotesByPatient(ctx context.Context, tenantID, patientID uuid.UUID) ([]domain.ProgressNote, error)
	UpdateProgressNote(ctx context.Context, n *domain.ProgressNote) error

	CreateAttachment(ctx context.Context, a *domain.ClinicalAttachment) error
	ListAttachmentsByPatient(ctx context.Context, tenantID, patientID uuid.UUID) ([]domain.ClinicalAttachment, error)
}

// Usecase contains business logic for clinical records.
type Usecase struct {
	repo ClinicalRepository
}

// New creates a clinical Usecase.
func New(repo ClinicalRepository) *Usecase {
	return &Usecase{repo: repo}
}

// --- Anamnesis ---

// CreateAnamnesis creates a new anamnesis in draft status.
func (uc *Usecase) CreateAnamnesis(ctx context.Context, a *domain.Anamnesis) error {
	a.ID = uuid.New()
	now := time.Now().UTC()
	a.CreatedAt = now
	a.UpdatedAt = now
	a.Status = domain.AnamnesisStatusDraft
	return uc.repo.CreateAnamnesis(ctx, a)
}

// GetAnamnesisByID returns an anamnesis by ID.
func (uc *Usecase) GetAnamnesisByID(ctx context.Context, tenantID, id uuid.UUID) (*domain.Anamnesis, error) {
	return uc.repo.GetAnamnesisByID(ctx, tenantID, id)
}

// ListAnamnesesByPatient returns anamneses for a patient.
func (uc *Usecase) ListAnamnesesByPatient(ctx context.Context, tenantID, patientID uuid.UUID) ([]domain.Anamnesis, error) {
	return uc.repo.ListAnamnesesByPatient(ctx, tenantID, patientID)
}

// UpdateAnamnesis updates a draft anamnesis. Returns error if already finalized.
func (uc *Usecase) UpdateAnamnesis(ctx context.Context, a *domain.Anamnesis) error {
	existing, err := uc.repo.GetAnamnesisByID(ctx, a.TenantID, a.ID)
	if err != nil {
		return err
	}
	if existing.Status == domain.AnamnesisStatusFinalized {
		return domain.ErrAlreadyFinalized
	}
	return uc.repo.UpdateAnamnesis(ctx, a)
}

// FinalizeAnamnesis marks an anamnesis as finalized (irreversible).
func (uc *Usecase) FinalizeAnamnesis(ctx context.Context, tenantID, id uuid.UUID) error {
	existing, err := uc.repo.GetAnamnesisByID(ctx, tenantID, id)
	if err != nil {
		return err
	}
	if existing.Status == domain.AnamnesisStatusFinalized {
		return domain.ErrAlreadyFinalized
	}

	now := time.Now().UTC()
	existing.Status = domain.AnamnesisStatusFinalized
	existing.FinalizedAt = &now

	return uc.repo.UpdateAnamnesis(ctx, existing)
}

// --- Progress Notes ---

// CreateProgressNote creates a new progress note.
func (uc *Usecase) CreateProgressNote(ctx context.Context, n *domain.ProgressNote) error {
	if err := n.Validate(); err != nil {
		return err
	}
	n.ID = uuid.New()
	now := time.Now().UTC()
	n.CreatedAt = now
	n.UpdatedAt = now
	return uc.repo.CreateProgressNote(ctx, n)
}

// ListProgressNotesByPatient returns progress notes for a patient.
func (uc *Usecase) ListProgressNotesByPatient(ctx context.Context, tenantID, patientID uuid.UUID) ([]domain.ProgressNote, error) {
	return uc.repo.ListProgressNotesByPatient(ctx, tenantID, patientID)
}

// UpdateProgressNote updates a progress note.
func (uc *Usecase) UpdateProgressNote(ctx context.Context, n *domain.ProgressNote) error {
	if err := n.Validate(); err != nil {
		return err
	}
	return uc.repo.UpdateProgressNote(ctx, n)
}

// --- Attachments ---

// CreateAttachment creates a clinical attachment record.
func (uc *Usecase) CreateAttachment(ctx context.Context, a *domain.ClinicalAttachment) error {
	a.ID = uuid.New()
	a.CreatedAt = time.Now().UTC()
	if a.Category == "" {
		a.Category = domain.CategoryGeneral
	}
	return uc.repo.CreateAttachment(ctx, a)
}

// ListAttachmentsByPatient returns attachments for a patient.
func (uc *Usecase) ListAttachmentsByPatient(ctx context.Context, tenantID, patientID uuid.UUID) ([]domain.ClinicalAttachment, error) {
	return uc.repo.ListAttachmentsByPatient(ctx, tenantID, patientID)
}
