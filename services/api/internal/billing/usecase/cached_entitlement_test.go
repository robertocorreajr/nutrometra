package usecase_test

import (
	"context"
	"encoding/json"
	"sync"
	"testing"
	"time"

	"nutrometra/api/internal/billing/domain"
	"nutrometra/api/internal/billing/usecase"
	"nutrometra/api/internal/platform/cache"

	"github.com/google/uuid"
)

// memoryCache is a simple in-memory implementation of cache.Cache for testing.
type memoryCache struct {
	mu    sync.RWMutex
	store map[string][]byte
}

func newMemoryCache() *memoryCache {
	return &memoryCache{store: make(map[string][]byte)}
}

func (m *memoryCache) Get(_ context.Context, key string, dest any) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	data, ok := m.store[key]
	if !ok {
		return cache.ErrCacheMiss
	}
	return json.Unmarshal(data, dest)
}

func (m *memoryCache) Set(_ context.Context, key string, value any, _ time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	m.store[key] = data
	return nil
}

func (m *memoryCache) Delete(_ context.Context, key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.store, key)
	return nil
}

// stubBillingRepo implements usecase.BillingRepository for testing.
type stubBillingRepo struct {
	sub      *domain.Subscription
	features []domain.PlanFeature
	calls    int
}

func (r *stubBillingRepo) GetActiveSubscription(_ context.Context, _ uuid.UUID) (*domain.Subscription, error) {
	r.calls++
	return r.sub, nil
}

func (r *stubBillingRepo) GetPlanFeature(_ context.Context, _ uuid.UUID, featureKey string) (*domain.PlanFeature, error) {
	for _, f := range r.features {
		if f.FeatureKey == featureKey {
			return &f, nil
		}
	}
	return nil, nil
}

func (r *stubBillingRepo) GetAllPlanFeatures(_ context.Context, _ uuid.UUID) ([]domain.PlanFeature, error) {
	r.calls++
	return r.features, nil
}

func (r *stubBillingRepo) GetActiveOverride(_ context.Context, _ uuid.UUID, _ string, _ time.Time) (*domain.FeatureOverride, error) {
	return nil, nil
}

func TestCachedEntitlementService_CachesGetAllEntitlements(t *testing.T) {
	planID := uuid.New()
	tenantID := uuid.New()

	repo := &stubBillingRepo{
		sub: &domain.Subscription{
			ID:       uuid.New(),
			TenantID: tenantID,
			PlanID:   planID,
			Status:   domain.SubscriptionActive,
		},
		features: []domain.PlanFeature{
			{ID: uuid.New(), PlanID: planID, FeatureKey: "patients", Enabled: true},
			{ID: uuid.New(), PlanID: planID, FeatureKey: "diets", Enabled: true},
		},
	}

	inner := usecase.NewEntitlementService(repo)
	mc := newMemoryCache()
	cached := usecase.NewCachedEntitlementService(inner, mc)

	ctx := context.Background()

	// First call: hits the repo
	result1, err := cached.GetAllEntitlements(ctx, tenantID)
	if err != nil {
		t.Fatalf("first call failed: %v", err)
	}
	if len(result1) != 2 {
		t.Fatalf("expected 2 entitlements, got %d", len(result1))
	}

	callsAfterFirst := repo.calls

	// Second call: should hit the cache, not the repo
	result2, err := cached.GetAllEntitlements(ctx, tenantID)
	if err != nil {
		t.Fatalf("second call failed: %v", err)
	}
	if len(result2) != 2 {
		t.Fatalf("expected 2 entitlements from cache, got %d", len(result2))
	}

	if repo.calls != callsAfterFirst {
		t.Fatalf("expected repo not to be called again; calls before=%d, after=%d", callsAfterFirst, repo.calls)
	}
}

func TestCachedEntitlementService_InvalidateClears(t *testing.T) {
	planID := uuid.New()
	tenantID := uuid.New()

	repo := &stubBillingRepo{
		sub: &domain.Subscription{
			ID:       uuid.New(),
			TenantID: tenantID,
			PlanID:   planID,
			Status:   domain.SubscriptionActive,
		},
		features: []domain.PlanFeature{
			{ID: uuid.New(), PlanID: planID, FeatureKey: "patients", Enabled: true},
		},
	}

	inner := usecase.NewEntitlementService(repo)
	mc := newMemoryCache()
	cached := usecase.NewCachedEntitlementService(inner, mc)

	ctx := context.Background()

	// Populate cache
	_, err := cached.GetAllEntitlements(ctx, tenantID)
	if err != nil {
		t.Fatalf("populate cache failed: %v", err)
	}
	callsAfterPopulate := repo.calls

	// Invalidate
	cached.InvalidateEntitlements(ctx, tenantID)

	// Next call should hit repo again
	_, err = cached.GetAllEntitlements(ctx, tenantID)
	if err != nil {
		t.Fatalf("post-invalidation call failed: %v", err)
	}
	if repo.calls == callsAfterPopulate {
		t.Fatal("expected repo to be called again after invalidation")
	}
}

func TestCachedEntitlementService_CheckEntitlementDenyByDefault(t *testing.T) {
	planID := uuid.New()
	tenantID := uuid.New()

	repo := &stubBillingRepo{
		sub: &domain.Subscription{
			ID:       uuid.New(),
			TenantID: tenantID,
			PlanID:   planID,
			Status:   domain.SubscriptionActive,
		},
		features: []domain.PlanFeature{
			{ID: uuid.New(), PlanID: planID, FeatureKey: "patients", Enabled: true},
		},
	}

	inner := usecase.NewEntitlementService(repo)
	mc := newMemoryCache()
	cached := usecase.NewCachedEntitlementService(inner, mc)

	ctx := context.Background()

	// Check a feature that exists
	ent, err := cached.CheckEntitlement(ctx, tenantID, "patients")
	if err != nil {
		t.Fatalf("CheckEntitlement failed: %v", err)
	}
	if !ent.Enabled {
		t.Fatal("expected 'patients' to be enabled")
	}

	// Check a feature that does not exist — should deny by default
	ent, err = cached.CheckEntitlement(ctx, tenantID, "nonexistent_feature")
	if err != nil {
		t.Fatalf("CheckEntitlement for missing feature failed: %v", err)
	}
	if ent.Enabled {
		t.Fatal("expected unknown feature to be denied by default")
	}
	if ent.Source != "default" {
		t.Fatalf("expected source 'default', got %q", ent.Source)
	}
}
