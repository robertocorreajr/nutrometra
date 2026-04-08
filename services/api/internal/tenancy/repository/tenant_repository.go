package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"nutrometra/api/internal/tenancy/domain"
)

type TenantRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Tenant, error)
	GetBySlug(ctx context.Context, slug string) (*domain.Tenant, error)
	IsMember(ctx context.Context, tenantID, userID uuid.UUID) (bool, error)
	ListByUser(ctx context.Context, userID uuid.UUID) ([]domain.Tenant, error)
}

type postgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) TenantRepository {
	return &postgresRepository{pool: pool}
}

const queryGetByID = `SELECT id, type, legal_name, display_name, slug, status, timezone, locale, trial_ends_at, created_at, updated_at FROM tenants WHERE id = $1`

const queryGetBySlug = `SELECT id, type, legal_name, display_name, slug, status, timezone, locale, trial_ends_at, created_at, updated_at FROM tenants WHERE slug = $1`

const queryIsMember = `SELECT EXISTS(SELECT 1 FROM tenant_users WHERE tenant_id = $1 AND user_id = $2 AND status = 'active')`

func scanTenant(row pgx.Row) (*domain.Tenant, error) {
	var t domain.Tenant
	err := row.Scan(
		&t.ID,
		&t.Type,
		&t.LegalName,
		&t.DisplayName,
		&t.Slug,
		&t.Status,
		&t.Timezone,
		&t.Locale,
		&t.TrialEndsAt,
		&t.CreatedAt,
		&t.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrTenantNotFound
		}
		return nil, err
	}
	return &t, nil
}

func (r *postgresRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Tenant, error) {
	row := r.pool.QueryRow(ctx, queryGetByID, id)
	return scanTenant(row)
}

func (r *postgresRepository) GetBySlug(ctx context.Context, slug string) (*domain.Tenant, error) {
	row := r.pool.QueryRow(ctx, queryGetBySlug, slug)
	return scanTenant(row)
}

func (r *postgresRepository) IsMember(ctx context.Context, tenantID, userID uuid.UUID) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx, queryIsMember, tenantID, userID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

const queryListByUser = `
SELECT t.id, t.type, t.legal_name, t.display_name, t.slug, t.status, t.timezone, t.locale, t.trial_ends_at, t.created_at, t.updated_at
FROM tenants t
JOIN tenant_users tu ON tu.tenant_id = t.id
WHERE tu.user_id = $1 AND tu.status = 'active' AND t.status = 'active'
ORDER BY t.display_name`

func (r *postgresRepository) ListByUser(ctx context.Context, userID uuid.UUID) ([]domain.Tenant, error) {
	rows, err := r.pool.Query(ctx, queryListByUser, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tenants []domain.Tenant
	for rows.Next() {
		var t domain.Tenant
		if err := rows.Scan(
			&t.ID, &t.Type, &t.LegalName, &t.DisplayName, &t.Slug,
			&t.Status, &t.Timezone, &t.Locale, &t.TrialEndsAt,
			&t.CreatedAt, &t.UpdatedAt,
		); err != nil {
			return nil, err
		}
		tenants = append(tenants, t)
	}
	return tenants, rows.Err()
}
