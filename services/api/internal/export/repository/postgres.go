package repository

import (
	"context"
	"errors"
	"fmt"

	"nutrometra/api/internal/export/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository provides export data access.
type Repository struct {
	pool *pgxpool.Pool
}

// New creates an export Repository.
func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// Create inserts a new exported file record.
func (r *Repository) Create(ctx context.Context, ef *domain.ExportedFile) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO exported_files
			(id, tenant_id, related_entity_type, related_entity_id,
			 export_type, file_key, status, requested_by_user_id,
			 created_at, completed_at, failure_reason)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
		ef.ID, ef.TenantID, ef.RelatedEntityType, ef.RelatedEntityID,
		string(ef.ExportType), nilIfEmpty(ef.FileKey), string(ef.Status),
		ef.RequestedByUserID,
		ef.CreatedAt, ef.CompletedAt, nilIfEmpty(ef.FailureReason),
	)
	if err != nil {
		return fmt.Errorf("export: create: %w", err)
	}
	return nil
}

// GetByID returns an exported file by ID with tenant isolation.
func (r *Repository) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*domain.ExportedFile, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id, tenant_id, related_entity_type, related_entity_id,
				export_type, COALESCE(file_key,''), status, requested_by_user_id,
				created_at, completed_at, COALESCE(failure_reason,'')
		 FROM exported_files
		 WHERE id = $1 AND tenant_id = $2`,
		id, tenantID,
	)

	var ef domain.ExportedFile
	var exportType, status string
	err := row.Scan(
		&ef.ID, &ef.TenantID, &ef.RelatedEntityType, &ef.RelatedEntityID,
		&exportType, &ef.FileKey, &status, &ef.RequestedByUserID,
		&ef.CreatedAt, &ef.CompletedAt, &ef.FailureReason,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("export: get_by_id: %w", err)
	}
	ef.ExportType = domain.ExportType(exportType)
	ef.Status = domain.ExportStatus(status)
	return &ef, nil
}

// ListByTenant returns all exports for a tenant ordered by creation date descending.
func (r *Repository) ListByTenant(ctx context.Context, tenantID uuid.UUID) ([]domain.ExportedFile, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, tenant_id, related_entity_type, related_entity_id,
				export_type, COALESCE(file_key,''), status, requested_by_user_id,
				created_at, completed_at, COALESCE(failure_reason,'')
		 FROM exported_files
		 WHERE tenant_id = $1
		 ORDER BY created_at DESC`,
		tenantID,
	)
	if err != nil {
		return nil, fmt.Errorf("export: list_by_tenant: %w", err)
	}
	defer rows.Close()

	var result []domain.ExportedFile
	for rows.Next() {
		var ef domain.ExportedFile
		var exportType, status string
		if err := rows.Scan(
			&ef.ID, &ef.TenantID, &ef.RelatedEntityType, &ef.RelatedEntityID,
			&exportType, &ef.FileKey, &status, &ef.RequestedByUserID,
			&ef.CreatedAt, &ef.CompletedAt, &ef.FailureReason,
		); err != nil {
			return nil, fmt.Errorf("export: list_by_tenant scan: %w", err)
		}
		ef.ExportType = domain.ExportType(exportType)
		ef.Status = domain.ExportStatus(status)
		result = append(result, ef)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("export: list_by_tenant rows: %w", err)
	}
	return result, nil
}

// UpdateStatus updates the status, file_key, completed_at, and failure_reason of an export.
// completed_at is set to NOW() when status is completed or failed.
func (r *Repository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.ExportStatus, fileKey string, failureReason string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE exported_files SET
			status = $2,
			file_key = COALESCE(NULLIF($3,''), file_key),
			completed_at = CASE WHEN $2 IN ('completed','failed') THEN NOW() ELSE completed_at END,
			failure_reason = CASE WHEN $4 = '' THEN failure_reason ELSE $4 END
		 WHERE id = $1`,
		id, string(status), fileKey, failureReason,
	)
	if err != nil {
		return fmt.Errorf("export: update_status: %w", err)
	}
	return nil
}

func nilIfEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
