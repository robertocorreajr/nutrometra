package usecase

import (
	"context"
	"encoding/json"
	"time"

	"nutrometra/api/internal/document/domain"
	"nutrometra/api/internal/platform/db"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DocumentRepository defines the data access contract for clinical documents.
type DocumentRepository interface {
	Create(ctx context.Context, d *domain.ClinicalDocument) error
	GetByID(ctx context.Context, tenantID, id uuid.UUID) (*domain.ClinicalDocument, error)
	Update(ctx context.Context, d *domain.ClinicalDocument) error
	ListByPatient(ctx context.Context, tenantID, patientID uuid.UUID) ([]domain.ClinicalDocument, error)
	Finalize(ctx context.Context, tenantID, id uuid.UUID) error
	ListVersions(ctx context.Context, tenantID, id uuid.UUID) ([]domain.ClinicalDocument, error)

	// Transactional methods
	Publish(ctx context.Context, tx pgx.Tx, tenantID, id uuid.UUID) error
	CreatePublication(ctx context.Context, tx pgx.Tx, pub *domain.DocumentPublication) error
	GetForVersion(ctx context.Context, tx pgx.Tx, tenantID, id uuid.UUID) (*domain.ClinicalDocument, error)
	CreateTx(ctx context.Context, tx pgx.Tx, d *domain.ClinicalDocument) error
}

// Usecase contains business logic for clinical documents.
type Usecase struct {
	repo DocumentRepository
	pool *pgxpool.Pool
}

// New creates a document Usecase.
func New(repo DocumentRepository, pool *pgxpool.Pool) *Usecase {
	return &Usecase{repo: repo, pool: pool}
}

// Create creates a new clinical document in draft status.
func (uc *Usecase) Create(ctx context.Context, d *domain.ClinicalDocument) error {
	if err := d.Validate(); err != nil {
		return err
	}
	d.ID = uuid.New()
	now := time.Now().UTC()
	d.CreatedAt = now
	d.UpdatedAt = now
	d.Status = domain.StatusDraft
	d.VersionNumber = 1
	if d.ContentJSON == nil {
		d.ContentJSON = json.RawMessage(`{}`)
	}
	return uc.repo.Create(ctx, d)
}

// GetByID returns a clinical document.
func (uc *Usecase) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*domain.ClinicalDocument, error) {
	return uc.repo.GetByID(ctx, tenantID, id)
}

// Update updates a draft clinical document.
func (uc *Usecase) Update(ctx context.Context, d *domain.ClinicalDocument) error {
	if err := d.Validate(); err != nil {
		return err
	}
	existing, err := uc.repo.GetByID(ctx, d.TenantID, d.ID)
	if err != nil {
		return err
	}
	if existing.Status != domain.StatusDraft {
		return domain.ErrNotDraft
	}
	return uc.repo.Update(ctx, d)
}

// ListByPatient returns clinical documents for a patient.
func (uc *Usecase) ListByPatient(ctx context.Context, tenantID, patientID uuid.UUID) ([]domain.ClinicalDocument, error) {
	return uc.repo.ListByPatient(ctx, tenantID, patientID)
}

// Finalize transitions a document from draft to finalized (irreversible).
func (uc *Usecase) Finalize(ctx context.Context, tenantID, id uuid.UUID) error {
	existing, err := uc.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return err
	}
	if existing.Status == domain.StatusFinalized || existing.Status == domain.StatusPublished {
		return domain.ErrAlreadyFinalized
	}
	if existing.Status != domain.StatusDraft {
		return domain.ErrNotDraft
	}
	return uc.repo.Finalize(ctx, tenantID, id)
}

// Publish publishes a finalized document and creates a publication record.
// This is a transactional operation that:
// 1. Verifies the document is finalized
// 2. Updates status to published
// 3. Creates a publication record
func (uc *Usecase) Publish(ctx context.Context, tenantID, docID, userID uuid.UUID) error {
	return db.RunInTx(ctx, uc.pool, func(ctx context.Context, tx pgx.Tx) error {
		doc, err := uc.repo.GetForVersion(ctx, tx, tenantID, docID)
		if err != nil {
			return err
		}

		if doc.Status != domain.StatusFinalized {
			return domain.ErrNotFinalized
		}

		if err := uc.repo.Publish(ctx, tx, tenantID, docID); err != nil {
			return err
		}

		now := time.Now().UTC()
		pub := &domain.DocumentPublication{
			ID:                 uuid.New(),
			TenantID:           tenantID,
			ClinicalDocumentID: docID,
			PatientID:          doc.PatientID,
			PublishedByUserID:  userID,
			PublishedAt:        now,
		}
		return uc.repo.CreatePublication(ctx, tx, pub)
	})
}

// CreateNewVersion creates a deep copy of a document with incremented version number.
// This is a transactional operation that:
// 1. Loads the document within the transaction
// 2. Creates a new document with new ID, version_number+1, previous_version_id set, status=draft
func (uc *Usecase) CreateNewVersion(ctx context.Context, tenantID, docID uuid.UUID) (*domain.ClinicalDocument, error) {
	var newDoc *domain.ClinicalDocument

	err := db.RunInTx(ctx, uc.pool, func(ctx context.Context, tx pgx.Tx) error {
		original, err := uc.repo.GetForVersion(ctx, tx, tenantID, docID)
		if err != nil {
			return err
		}

		now := time.Now().UTC()
		newID := uuid.New()
		originalID := original.ID

		// Deep copy content_json
		var contentCopy json.RawMessage
		if original.ContentJSON != nil {
			contentCopy = make(json.RawMessage, len(original.ContentJSON))
			copy(contentCopy, original.ContentJSON)
		}

		newDoc = &domain.ClinicalDocument{
			ID:                newID,
			TenantID:          original.TenantID,
			PatientID:         original.PatientID,
			ProfessionalID:    original.ProfessionalID,
			AppointmentID:     original.AppointmentID,
			DocumentType:      original.DocumentType,
			Title:             original.Title,
			Status:            domain.StatusDraft,
			ContentJSON:       contentCopy,
			VersionNumber:     original.VersionNumber + 1,
			PreviousVersionID: &originalID,
			CreatedAt:         now,
			UpdatedAt:         now,
		}

		return uc.repo.CreateTx(ctx, tx, newDoc)
	})
	if err != nil {
		return nil, err
	}

	return newDoc, nil
}

// ListVersions returns the version chain for a document.
func (uc *Usecase) ListVersions(ctx context.Context, tenantID, id uuid.UUID) ([]domain.ClinicalDocument, error) {
	return uc.repo.ListVersions(ctx, tenantID, id)
}
