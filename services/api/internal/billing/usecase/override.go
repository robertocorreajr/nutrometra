package usecase

import (
	"context"
	"errors"
	"fmt"

	"nutrometra/api/internal/billing/domain"
	"nutrometra/api/internal/billing/repository"
	"nutrometra/api/internal/platform/audit"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrOverrideNotFound = errors.New("billing: override not found")
	ErrOverrideExists   = errors.New("billing: override already exists for this feature")
	ErrReasonRequired   = errors.New("billing: reason is required for override mutations")
)

// OverrideUsecase manages tenant feature override CRUD.
type OverrideUsecase struct {
	repo     *repository.OverrideRepository
	pool     *pgxpool.Pool
	auditSvc *audit.Service
}

func NewOverrideUsecase(repo *repository.OverrideRepository, pool *pgxpool.Pool, auditSvc *audit.Service) *OverrideUsecase {
	return &OverrideUsecase{repo: repo, pool: pool, auditSvc: auditSvc}
}

// Create creates a new feature override.
func (uc *OverrideUsecase) Create(ctx context.Context, o domain.FeatureOverride, actorID uuid.UUID, actorScope audit.Scope) error {
	if o.Reason == "" {
		return ErrReasonRequired
	}

	existing, err := uc.repo.GetByTenantAndKey(ctx, o.TenantID, o.FeatureKey)
	if err != nil {
		return fmt.Errorf("override: check_existing: %w", err)
	}
	if existing != nil {
		return ErrOverrideExists
	}

	o.ID = uuid.New()
	o.CreatedByUserID = &actorID

	if err := uc.repo.Create(ctx, o); err != nil {
		return fmt.Errorf("override: create: %w", err)
	}

	entry := audit.NewEntry(
		audit.WithTenantID(o.TenantID),
		audit.WithActor(actorID, actorScope),
		audit.WithEntity("tenant_feature_override", o.ID),
		audit.WithAction("override_created"),
		audit.WithReason(o.Reason),
		audit.WithMetadata(map[string]string{"feature_key": o.FeatureKey}),
	)
	_ = uc.auditSvc.Write(ctx, uc.pool, entry)
	return nil
}

// Update modifies an existing feature override.
func (uc *OverrideUsecase) Update(ctx context.Context, o domain.FeatureOverride, actorID uuid.UUID, actorScope audit.Scope) error {
	if o.Reason == "" {
		return ErrReasonRequired
	}

	existing, err := uc.repo.GetByTenantAndKey(ctx, o.TenantID, o.FeatureKey)
	if err != nil {
		return fmt.Errorf("override: check_existing: %w", err)
	}
	if existing == nil {
		return ErrOverrideNotFound
	}

	if err := uc.repo.Update(ctx, o); err != nil {
		return fmt.Errorf("override: update: %w", err)
	}

	entry := audit.NewEntry(
		audit.WithTenantID(o.TenantID),
		audit.WithActor(actorID, actorScope),
		audit.WithEntity("tenant_feature_override", existing.ID),
		audit.WithAction("override_updated"),
		audit.WithReason(o.Reason),
		audit.WithMetadata(map[string]string{"feature_key": o.FeatureKey}),
	)
	_ = uc.auditSvc.Write(ctx, uc.pool, entry)
	return nil
}

// Delete removes a feature override.
func (uc *OverrideUsecase) Delete(ctx context.Context, tenantID uuid.UUID, featureKey string, actorID uuid.UUID, actorScope audit.Scope, reason string) error {
	if reason == "" {
		return ErrReasonRequired
	}

	existing, err := uc.repo.GetByTenantAndKey(ctx, tenantID, featureKey)
	if err != nil {
		return fmt.Errorf("override: check_existing: %w", err)
	}
	if existing == nil {
		return ErrOverrideNotFound
	}

	if err := uc.repo.Delete(ctx, tenantID, featureKey); err != nil {
		return fmt.Errorf("override: delete: %w", err)
	}

	entry := audit.NewEntry(
		audit.WithTenantID(tenantID),
		audit.WithActor(actorID, actorScope),
		audit.WithEntity("tenant_feature_override", existing.ID),
		audit.WithAction("override_deleted"),
		audit.WithReason(reason),
		audit.WithMetadata(map[string]string{"feature_key": featureKey}),
	)
	_ = uc.auditSvc.Write(ctx, uc.pool, entry)
	return nil
}

// ListByTenantID returns all overrides for a tenant.
func (uc *OverrideUsecase) ListByTenantID(ctx context.Context, tenantID uuid.UUID) ([]domain.FeatureOverride, error) {
	return uc.repo.ListByTenantID(ctx, tenantID)
}
