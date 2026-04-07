package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"nutrometra/api/internal/tenancy/domain"
)

var ErrNotFound = errors.New("tenant not found")

// TenantRepository defines read operations for tenants.
type TenantRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Tenant, error)
	GetBySlug(ctx context.Context, slug string) (*domain.Tenant, error)
}

type postgresRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresRepository returns a TenantRepository backed by PostgreSQL.
func NewPostgresRepository(pool *pgxpool.Pool) TenantRepository {
	return &postgresRepository{pool: pool}
}

const selectColumns = `id, name, slug, plan_id, status, trial_ends_at, created_at, updated_at`

func scanTenant(row pgx.Row) (*domain.Tenant, error) {
	var t domain.Tenant
	err := row.Scan(
		&t.ID,
		&t.Name,
		&t.Slug,
		&t.PlanID,
		&t.Status,
		&t.TrialEndsAt,
		&t.CreatedAt,
		&t.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &t, nil
}

func (r *postgresRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Tenant, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT `+selectColumns+` FROM tenants WHERE id = $1`,
		id,
	)
	return scanTenant(row)
}

func (r *postgresRepository) GetBySlug(ctx context.Context, slug string) (*domain.Tenant, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT `+selectColumns+` FROM tenants WHERE slug = $1`,
		slug,
	)
	return scanTenant(row)
}
