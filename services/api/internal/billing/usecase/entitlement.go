package usecase

import (
	"context"
	"fmt"

	"nutrometra/api/internal/billing/domain"
	"nutrometra/api/internal/billing/repository"

	"github.com/google/uuid"
)

type EntitlementService struct {
	repo *repository.Repository
}

func NewEntitlementService(repo *repository.Repository) *EntitlementService {
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

	override, err := s.repo.GetActiveOverride(ctx, tenantID, featureKey)
	if err != nil {
		return nil, fmt.Errorf("entitlement: get_override: %w", err)
	}

	planFeature, err := s.repo.GetPlanFeature(ctx, sub.PlanID, featureKey)
	if err != nil {
		return nil, fmt.Errorf("entitlement: get_plan_feature: %w", err)
	}

	return ResolveEntitlement(featureKey, planFeature, override), nil
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
