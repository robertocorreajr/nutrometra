package repository

import (
	"context"
	"errors"
	"fmt"

	"nutrometra/api/internal/document/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository provides clinical document data access.
type Repository struct {
	pool *pgxpool.Pool
}

// New creates a document Repository.
func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// querier is satisfied by both *pgxpool.Pool and pgx.Tx.
type querier interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

// Create inserts a new clinical document.
func (r *Repository) Create(ctx context.Context, d *domain.ClinicalDocument) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO clinical_documents
			(id, tenant_id, patient_id, professional_id, appointment_id,
			 document_type, title, status, content_json,
			 version_number, previous_version_id, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)`,
		d.ID, d.TenantID, d.PatientID, d.ProfessionalID, d.AppointmentID,
		string(d.DocumentType), d.Title, string(d.Status), d.ContentJSON,
		d.VersionNumber, d.PreviousVersionID,
		d.CreatedAt, d.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("document: create: %w", err)
	}
	return nil
}

// GetByID returns a clinical document by ID with tenant isolation.
func (r *Repository) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*domain.ClinicalDocument, error) {
	return r.scanDocument(ctx, r.pool, tenantID, id)
}

// Update updates a draft clinical document.
func (r *Repository) Update(ctx context.Context, d *domain.ClinicalDocument) error {
	tag, err := r.pool.Exec(ctx,
		`UPDATE clinical_documents SET
			title = $3, content_json = $4, appointment_id = $5, updated_at = NOW()
		 WHERE id = $1 AND tenant_id = $2 AND status = 'draft'`,
		d.ID, d.TenantID,
		d.Title, d.ContentJSON, d.AppointmentID,
	)
	if err != nil {
		return fmt.Errorf("document: update: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// ListByPatient returns clinical documents for a patient.
func (r *Repository) ListByPatient(ctx context.Context, tenantID, patientID uuid.UUID) ([]domain.ClinicalDocument, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, tenant_id, patient_id, professional_id, appointment_id,
			document_type, title, status, content_json,
			version_number, previous_version_id, created_at, updated_at
		 FROM clinical_documents WHERE tenant_id = $1 AND patient_id = $2
		 ORDER BY created_at DESC`,
		tenantID, patientID,
	)
	if err != nil {
		return nil, fmt.Errorf("document: list_by_patient: %w", err)
	}
	defer rows.Close()

	var result []domain.ClinicalDocument
	for rows.Next() {
		var d domain.ClinicalDocument
		if err := rows.Scan(
			&d.ID, &d.TenantID, &d.PatientID, &d.ProfessionalID, &d.AppointmentID,
			&d.DocumentType, &d.Title, &d.Status, &d.ContentJSON,
			&d.VersionNumber, &d.PreviousVersionID, &d.CreatedAt, &d.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("document: scan: %w", err)
		}
		result = append(result, d)
	}
	return result, rows.Err()
}

// Finalize updates a clinical document status to finalized.
func (r *Repository) Finalize(ctx context.Context, tenantID, id uuid.UUID) error {
	tag, err := r.pool.Exec(ctx,
		`UPDATE clinical_documents SET status = 'finalized', updated_at = NOW()
		 WHERE id = $1 AND tenant_id = $2 AND status = 'draft'`,
		id, tenantID,
	)
	if err != nil {
		return fmt.Errorf("document: finalize: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// Publish updates a clinical document status to published within a transaction.
func (r *Repository) Publish(ctx context.Context, tx pgx.Tx, tenantID, id uuid.UUID) error {
	tag, err := tx.Exec(ctx,
		`UPDATE clinical_documents SET status = 'published', updated_at = NOW()
		 WHERE id = $1 AND tenant_id = $2 AND status = 'finalized'`,
		id, tenantID,
	)
	if err != nil {
		return fmt.Errorf("document: publish: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// CreatePublication inserts a new document publication record within a transaction.
func (r *Repository) CreatePublication(ctx context.Context, tx pgx.Tx, pub *domain.DocumentPublication) error {
	_, err := tx.Exec(ctx,
		`INSERT INTO document_publications
			(id, tenant_id, clinical_document_id, patient_id, published_by_user_id, published_at)
		 VALUES ($1,$2,$3,$4,$5,$6)`,
		pub.ID, pub.TenantID, pub.ClinicalDocumentID, pub.PatientID,
		pub.PublishedByUserID, pub.PublishedAt,
	)
	if err != nil {
		return fmt.Errorf("document: create_publication: %w", err)
	}
	return nil
}

// GetForVersion loads a clinical document within a transaction (for CreateNewVersion).
func (r *Repository) GetForVersion(ctx context.Context, tx pgx.Tx, tenantID, id uuid.UUID) (*domain.ClinicalDocument, error) {
	return r.scanDocument(ctx, tx, tenantID, id)
}

// CreateTx inserts a new clinical document within a transaction.
func (r *Repository) CreateTx(ctx context.Context, tx pgx.Tx, d *domain.ClinicalDocument) error {
	_, err := tx.Exec(ctx,
		`INSERT INTO clinical_documents
			(id, tenant_id, patient_id, professional_id, appointment_id,
			 document_type, title, status, content_json,
			 version_number, previous_version_id, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)`,
		d.ID, d.TenantID, d.PatientID, d.ProfessionalID, d.AppointmentID,
		string(d.DocumentType), d.Title, string(d.Status), d.ContentJSON,
		d.VersionNumber, d.PreviousVersionID,
		d.CreatedAt, d.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("document: create_tx: %w", err)
	}
	return nil
}

// ListVersions returns the version chain for a document by following previous_version_id.
func (r *Repository) ListVersions(ctx context.Context, tenantID, id uuid.UUID) ([]domain.ClinicalDocument, error) {
	// Use a recursive CTE to follow the previous_version_id chain in both directions.
	rows, err := r.pool.Query(ctx,
		`WITH RECURSIVE chain AS (
			-- Find the root of the version chain (walk backwards)
			SELECT id, previous_version_id
			FROM clinical_documents
			WHERE id = $1 AND tenant_id = $2
			UNION ALL
			SELECT cd.id, cd.previous_version_id
			FROM clinical_documents cd
			JOIN chain c ON cd.id = c.previous_version_id
			WHERE cd.tenant_id = $2
		), root AS (
			SELECT id FROM chain WHERE previous_version_id IS NULL
		), forward AS (
			-- Walk forward from root
			SELECT cd.id, cd.tenant_id, cd.patient_id, cd.professional_id,
				cd.appointment_id, cd.document_type, cd.title, cd.status,
				cd.content_json, cd.version_number, cd.previous_version_id,
				cd.created_at, cd.updated_at
			FROM clinical_documents cd
			JOIN root r ON cd.id = r.id
			WHERE cd.tenant_id = $2
			UNION ALL
			SELECT cd.id, cd.tenant_id, cd.patient_id, cd.professional_id,
				cd.appointment_id, cd.document_type, cd.title, cd.status,
				cd.content_json, cd.version_number, cd.previous_version_id,
				cd.created_at, cd.updated_at
			FROM clinical_documents cd
			JOIN forward f ON cd.previous_version_id = f.id
			WHERE cd.tenant_id = $2
		)
		SELECT id, tenant_id, patient_id, professional_id, appointment_id,
			document_type, title, status, content_json,
			version_number, previous_version_id, created_at, updated_at
		FROM forward
		ORDER BY version_number ASC`,
		id, tenantID,
	)
	if err != nil {
		return nil, fmt.Errorf("document: list_versions: %w", err)
	}
	defer rows.Close()

	var result []domain.ClinicalDocument
	for rows.Next() {
		var d domain.ClinicalDocument
		if err := rows.Scan(
			&d.ID, &d.TenantID, &d.PatientID, &d.ProfessionalID, &d.AppointmentID,
			&d.DocumentType, &d.Title, &d.Status, &d.ContentJSON,
			&d.VersionNumber, &d.PreviousVersionID, &d.CreatedAt, &d.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("document: scan_version: %w", err)
		}
		result = append(result, d)
	}
	return result, rows.Err()
}

// --- internal helpers ---

// scanDocument loads a single clinical document row.
func (r *Repository) scanDocument(ctx context.Context, q querier, tenantID, id uuid.UUID) (*domain.ClinicalDocument, error) {
	d := &domain.ClinicalDocument{}
	err := q.QueryRow(ctx,
		`SELECT id, tenant_id, patient_id, professional_id, appointment_id,
			document_type, title, status, content_json,
			version_number, previous_version_id, created_at, updated_at
		 FROM clinical_documents WHERE id = $1 AND tenant_id = $2`,
		id, tenantID,
	).Scan(
		&d.ID, &d.TenantID, &d.PatientID, &d.ProfessionalID, &d.AppointmentID,
		&d.DocumentType, &d.Title, &d.Status, &d.ContentJSON,
		&d.VersionNumber, &d.PreviousVersionID, &d.CreatedAt, &d.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("document: get_by_id: %w", err)
	}
	return d, nil
}
