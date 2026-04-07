package usecase

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"nutrometra/api/internal/export/domain"
	"nutrometra/api/internal/platform/pdfgen"
	"nutrometra/api/internal/platform/worker"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// EntitlementChecker checks feature entitlements.
type EntitlementChecker interface {
	CheckEntitlement(ctx context.Context, tenantID uuid.UUID, featureKey string) (bool, *int64, error)
}

// ExportRepository defines the data access contract.
type ExportRepository interface {
	Create(ctx context.Context, ef *domain.ExportedFile) error
	GetByID(ctx context.Context, tenantID, id uuid.UUID) (*domain.ExportedFile, error)
	ListByTenant(ctx context.Context, tenantID uuid.UUID) ([]domain.ExportedFile, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status domain.ExportStatus, fileKey string, failureReason string) error
}

// Usecase orchestrates export operations.
type Usecase struct {
	repo         ExportRepository
	pool         *pgxpool.Pool
	entitlements EntitlementChecker
	pdfGen       pdfgen.Generator
	worker       *worker.Worker
}

// New creates an export Usecase and registers the PDF job handler.
func New(repo ExportRepository, pool *pgxpool.Pool, entitlements EntitlementChecker, pdfGen pdfgen.Generator, w *worker.Worker) *Usecase {
	uc := &Usecase{
		repo:         repo,
		pool:         pool,
		entitlements: entitlements,
		pdfGen:       pdfGen,
		worker:       w,
	}
	w.Register("pdf_export", uc.handlePDFJob)
	return uc
}

// RequestExport creates a pending export and enqueues the PDF generation job.
func (uc *Usecase) RequestExport(ctx context.Context, tenantID, userID uuid.UUID, entityType string, entityID uuid.UUID) (*domain.ExportedFile, error) {
	allowed, _, err := uc.entitlements.CheckEntitlement(ctx, tenantID, "pdf:export")
	if err != nil {
		return nil, fmt.Errorf("export: entitlement check: %w", err)
	}
	if !allowed {
		return nil, fmt.Errorf("export: feature pdf:export not allowed for tenant")
	}

	ef := &domain.ExportedFile{
		ID:                uuid.New(),
		TenantID:          tenantID,
		RelatedEntityType: entityType,
		RelatedEntityID:   entityID,
		ExportType:        domain.TypePDF,
		Status:            domain.StatusPending,
		RequestedByUserID: userID,
		CreatedAt:         time.Now().UTC(),
	}

	if err := uc.repo.Create(ctx, ef); err != nil {
		return nil, fmt.Errorf("export: create: %w", err)
	}

	uc.worker.Enqueue(worker.Job{
		ID:      ef.ID,
		Type:    "pdf_export",
		Payload: ef,
	})

	return ef, nil
}

// GetByID returns an export by ID with tenant isolation.
func (uc *Usecase) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*domain.ExportedFile, error) {
	return uc.repo.GetByID(ctx, tenantID, id)
}

// ListByTenant returns all exports for a tenant.
func (uc *Usecase) ListByTenant(ctx context.Context, tenantID uuid.UUID) ([]domain.ExportedFile, error) {
	return uc.repo.ListByTenant(ctx, tenantID)
}

// handlePDFJob is the worker handler for pdf_export jobs.
func (uc *Usecase) handlePDFJob(ctx context.Context, job worker.Job) error {
	ef, ok := job.Payload.(*domain.ExportedFile)
	if !ok {
		return fmt.Errorf("export: invalid job payload")
	}

	_ = uc.repo.UpdateStatus(ctx, ef.ID, domain.StatusProcessing, "", "")

	// Generate placeholder PDF (actual data loading will be wired in a future step).
	data := pdfgen.DocumentPDFData{
		Title:            fmt.Sprintf("Export %s", ef.RelatedEntityType),
		DocumentType:     ef.RelatedEntityType,
		PatientName:      "\u2014",
		ProfessionalName: "\u2014",
		Content:          map[string]any{"entity_id": ef.RelatedEntityID.String()},
		CreatedAt:        ef.CreatedAt,
	}
	pdfBytes, err := uc.pdfGen.GenerateDocumentPDF(ctx, data)
	if err != nil {
		_ = uc.repo.UpdateStatus(ctx, ef.ID, domain.StatusFailed, "", err.Error())
		return fmt.Errorf("export: pdf generation failed: %w", err)
	}

	// Save to filesystem.
	dir := filepath.Join("exports", ef.TenantID.String())
	if err := os.MkdirAll(dir, 0o750); err != nil {
		_ = uc.repo.UpdateStatus(ctx, ef.ID, domain.StatusFailed, "", err.Error())
		return fmt.Errorf("export: mkdir failed: %w", err)
	}
	fileKey := filepath.Join(dir, ef.ID.String()+".pdf")
	if err := os.WriteFile(fileKey, pdfBytes, 0o640); err != nil {
		_ = uc.repo.UpdateStatus(ctx, ef.ID, domain.StatusFailed, "", err.Error())
		return fmt.Errorf("export: write failed: %w", err)
	}

	_ = uc.repo.UpdateStatus(ctx, ef.ID, domain.StatusCompleted, fileKey, "")
	return nil
}
