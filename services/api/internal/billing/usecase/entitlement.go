package usecase

import (
	"context"
	"fmt"
	"time"

	"nutrometra/api/internal/billing/domain"

	"github.com/google/uuid"
)

// BillingRepository defines the data access needed by the entitlement service.
type BillingRepository interface {
	GetActiveSubscription(ctx context.Context, tenantID uuid.UUID) (*domain.Subscription, error)
	GetPlanFeature(ctx context.Context, planID uuid.UUID, featureKey string) (*domain.PlanFeature, error)
	GetAllPlanFeatures(ctx context.Context, planID uuid.UUID) ([]domain.PlanFeature, error)
	GetActiveOverride(ctx context.Context, tenantID uuid.UUID, featureKey string, now time.Time) (*domain.FeatureOverride, error)
}

type EntitlementService struct {
	repo BillingRepository
}

func NewEntitlementService(repo BillingRepository) *EntitlementService {
	return &EntitlementService{repo: repo}
}

// CheckEntitlement resolves whether a feature is enabled for a tenant.
// Precedence: active override > plan_features > default (deny).
func (s *EntitlementService) CheckEntitlement(ctx context.Context, tenantID uuid.UUID, featureKey string) (*domain.Entitlement, error) {
	sub, err := s.repo.GetActiveSubscription(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("entitlement: get_subscription: %w", err)
	}

	if sub == nil {
		return &domain.Entitlement{
			FeatureKey: featureKey,
			Enabled:    false,
			Source:     "default",
		}, nil
	}

	now := time.Now().UTC()
	override, err := s.repo.GetActiveOverride(ctx, tenantID, featureKey, now)
	if err != nil {
		return nil, fmt.Errorf("entitlement: get_override: %w", err)
	}

	planFeature, err := s.repo.GetPlanFeature(ctx, sub.PlanID, featureKey)
	if err != nil {
		return nil, fmt.Errorf("entitlement: get_plan_feature: %w", err)
	}

	return ResolveEntitlement(featureKey, planFeature, override), nil
}

// GetAllEntitlements returns the full entitlement map for a tenant.
func (s *EntitlementService) GetAllEntitlements(ctx context.Context, tenantID uuid.UUID) (map[string]*domain.Entitlement, error) {
	sub, err := s.repo.GetActiveSubscription(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("entitlement: get_subscription: %w", err)
	}
	if sub == nil {
		return map[string]*domain.Entitlement{}, nil
	}

	features, err := s.repo.GetAllPlanFeatures(ctx, sub.PlanID)
	if err != nil {
		return nil, fmt.Errorf("entitlement: get_all_features: %w", err)
	}

	now := time.Now().UTC()
	result := make(map[string]*domain.Entitlement, len(features))
	for _, pf := range features {
		override, _ := s.repo.GetActiveOverride(ctx, tenantID, pf.FeatureKey, now)
		result[pf.FeatureKey] = ResolveEntitlement(pf.FeatureKey, &pf, override)
	}
	return result, nil
}

// ResolveEntitlement applies the precedence logic. Exported for unit testing.
func ResolveEntitlement(featureKey string, planFeature *domain.PlanFeature, override *domain.FeatureOverride) *domain.Entitlement {
	// Override has absolute precedence
	if override != nil {
		enabled := true
		if override.Enabled != nil {
			enabled = *override.Enabled
		}
		return &domain.Entitlement{
			FeatureKey: featureKey,
			Enabled:    enabled,
			Limit:      override.LimitValue,
			Source:     "override",
		}
	}

	// Fall back to plan feature
	if planFeature != nil {
		return &domain.Entitlement{
			FeatureKey: featureKey,
			Enabled:    planFeature.Enabled,
			Limit:      planFeature.LimitValue,
			Source:     "plan",
		}
	}

	// No plan feature defined: deny-by-default
	return &domain.Entitlement{
		FeatureKey: featureKey,
		Enabled:    false,
		Source:     "default",
	}
}
