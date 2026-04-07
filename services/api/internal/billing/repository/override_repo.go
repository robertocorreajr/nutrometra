package repository

import (
	"context"
	"errors"
	"fmt"

	"nutrometra/api/internal/billing/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// OverrideRepository handles tenant feature override persistence.
type OverrideRepository struct {
	pool *pgxpool.Pool
}

func NewOverrideRepository(pool *pgxpool.Pool) *OverrideRepository {
	return &OverrideRepository{pool: pool}
}

// Create inserts a new feature override.
func (r *OverrideRepository) Create(ctx context.Context, o domain.FeatureOverride) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO tenant_feature_overrides
		 (id, tenant_id, feature_key, enabled, limit_value, starts_at, ends_at, reason, created_by_user_id)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
		o.ID, o.TenantID, o.FeatureKey, o.Enabled, o.LimitValue,
		o.StartsAt, o.EndsAt, o.Reason, o.CreatedByUserID,
	)
	if err != nil {
		return fmt.Errorf("billing: create_override: %w", err)
	}
	return nil
}

// Update modifies an existing feature override.
func (r *OverrideRepository) Update(ctx context.Context, o domain.FeatureOverride) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE tenant_feature_overrides
		 SET enabled = $1, limit_value = $2, starts_at = $3, ends_at = $4, reason = $5
		 WHERE tenant_id = $6 AND feature_key = $7`,
		o.Enabled, o.LimitValue, o.StartsAt, o.EndsAt, o.Reason,
		o.TenantID, o.FeatureKey,
	)
	if err != nil {
		return fmt.Errorf("billing: update_override: %w", err)
	}
	return nil
}

// Delete removes a feature override.
func (r *OverrideRepository) Delete(ctx context.Context, tenantID uuid.UUID, featureKey string) error {
	_, err := r.pool.Exec(ctx,
		`DELETE FROM tenant_feature_overrides WHERE tenant_id = $1 AND feature_key = $2`,
		tenantID, featureKey,
	)
	if err != nil {
		return fmt.Errorf("billing: delete_override: %w", err)
	}
	return nil
}

// ListByTenantID returns all overrides for a tenant.
func (r *OverrideRepository) ListByTenantID(ctx context.Context, tenantID uuid.UUID) ([]domain.FeatureOverride, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, tenant_id, feature_key, enabled, limit_value, starts_at, ends_at, reason, created_by_user_id
		 FROM tenant_feature_overrides
		 WHERE tenant_id = $1
		 ORDER BY feature_key`,
		tenantID,
	)
	if err != nil {
		return nil, fmt.Errorf("billing: list_overrides: %w", err)
	}
	defer rows.Close()

	var overrides []domain.FeatureOverride
	for rows.Next() {
		var o domain.FeatureOverride
		if err := rows.Scan(&o.ID, &o.TenantID, &o.FeatureKey, &o.Enabled, &o.LimitValue,
			&o.StartsAt, &o.EndsAt, &o.Reason, &o.CreatedByUserID); err != nil {
			return nil, fmt.Errorf("billing: scan_override: %w", err)
		}
		overrides = append(overrides, o)
	}
	return overrides, rows.Err()
}

// GetByTenantAndKey returns a specific override.
func (r *OverrideRepository) GetByTenantAndKey(ctx context.Context, tenantID uuid.UUID, featureKey string) (*domain.FeatureOverride, error) {
	var o domain.FeatureOverride
	err := r.pool.QueryRow(ctx,
		`SELECT id, tenant_id, feature_key, enabled, limit_value, starts_at, ends_at, reason, created_by_user_id
		 FROM tenant_feature_overrides
		 WHERE tenant_id = $1 AND feature_key = $2`,
		tenantID, featureKey,
	).Scan(&o.ID, &o.TenantID, &o.FeatureKey, &o.Enabled, &o.LimitValue,
		&o.StartsAt, &o.EndsAt, &o.Reason, &o.CreatedByUserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("billing: get_override: %w", err)
	}
	return &o, nil
}
