package usecase

import (
	"context"
	"fmt"
	"time"

	"nutrometra/api/internal/billing/domain"
	"nutrometra/api/internal/platform/cache"

	"github.com/google/uuid"
)

const entitlementCacheTTL = 5 * time.Minute

// CachedEntitlementService wraps EntitlementService with a cache layer.
type CachedEntitlementService struct {
	inner *EntitlementService
	cache cache.Cache
}

// NewCachedEntitlementService creates a cached wrapper around EntitlementService.
func NewCachedEntitlementService(inner *EntitlementService, c cache.Cache) *CachedEntitlementService {
	return &CachedEntitlementService{inner: inner, cache: c}
}

func entitlementCacheKey(tenantID uuid.UUID) string {
	return fmt.Sprintf("ent:%s", tenantID.String())
}

// GetAllEntitlements returns the full entitlement map for a tenant, using cache when available.
func (s *CachedEntitlementService) GetAllEntitlements(ctx context.Context, tenantID uuid.UUID) (map[string]*domain.Entitlement, error) {
	key := entitlementCacheKey(tenantID)

	var cached map[string]*domain.Entitlement
	if err := s.cache.Get(ctx, key, &cached); err == nil {
		return cached, nil
	}

	result, err := s.inner.GetAllEntitlements(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	_ = s.cache.Set(ctx, key, result, entitlementCacheTTL)
	return result, nil
}

// CheckEntitlement resolves whether a feature is enabled for a tenant, using the cached entitlement map.
func (s *CachedEntitlementService) CheckEntitlement(ctx context.Context, tenantID uuid.UUID, featureKey string) (*domain.Entitlement, error) {
	all, err := s.GetAllEntitlements(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	if ent, ok := all[featureKey]; ok {
		return ent, nil
	}
	// Feature not in plan — deny by default
	return &domain.Entitlement{
		FeatureKey: featureKey,
		Enabled:    false,
		Source:     "default",
	}, nil
}

// InvalidateEntitlements removes the cached entitlements for a tenant.
func (s *CachedEntitlementService) InvalidateEntitlements(ctx context.Context, tenantID uuid.UUID) {
	_ = s.cache.Delete(ctx, entitlementCacheKey(tenantID))
}
