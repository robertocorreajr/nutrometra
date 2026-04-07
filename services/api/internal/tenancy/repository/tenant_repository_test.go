package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"nutrometra/api/internal/tenancy/domain"
	"nutrometra/api/internal/tenancy/repository"
	"nutrometra/api/internal/tenancy/usecase"
)

// fakeTenantRepo is an in-memory implementation of TenantRepository for tests.
type fakeTenantRepo struct {
	tenants map[uuid.UUID]*domain.Tenant
}

func (f *fakeTenantRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Tenant, error) {
	t, ok := f.tenants[id]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return t, nil
}

func (f *fakeTenantRepo) GetBySlug(ctx context.Context, slug string) (*domain.Tenant, error) {
	for _, t := range f.tenants {
		if t.Slug == slug {
			return t, nil
		}
	}
	return nil, repository.ErrNotFound
}

func newFakeRepo(tenants ...*domain.Tenant) *fakeTenantRepo {
	m := make(map[uuid.UUID]*domain.Tenant, len(tenants))
	for _, t := range tenants {
		m[t.ID] = t
	}
	return &fakeTenantRepo{tenants: m}
}

func sampleTenant() *domain.Tenant {
	planID := uuid.New()
	return &domain.Tenant{
		ID:        uuid.New(),
		Name:      "Acme Nutrition",
		Slug:      "acme-nutrition",
		PlanID:    &planID,
		Status:    domain.TenantStatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func TestTenantUsecase_GetByID_Found(t *testing.T) {
	tenant := sampleTenant()
	uc := usecase.NewTenantUsecase(newFakeRepo(tenant))

	got, err := uc.GetByID(context.Background(), tenant.ID)

	require.NoError(t, err)
	assert.Equal(t, tenant.ID, got.ID)
	assert.Equal(t, tenant.Name, got.Name)
	assert.Equal(t, tenant.Slug, got.Slug)
}

func TestTenantUsecase_GetByID_NotFound(t *testing.T) {
	uc := usecase.NewTenantUsecase(newFakeRepo())

	got, err := uc.GetByID(context.Background(), uuid.New())

	assert.Nil(t, got)
	assert.ErrorIs(t, err, repository.ErrNotFound)
}

func TestTenantUsecase_GetBySlug_Found(t *testing.T) {
	tenant := sampleTenant()
	uc := usecase.NewTenantUsecase(newFakeRepo(tenant))

	got, err := uc.GetBySlug(context.Background(), tenant.Slug)

	require.NoError(t, err)
	assert.Equal(t, tenant.ID, got.ID)
	assert.Equal(t, tenant.Slug, got.Slug)
}

func TestTenantUsecase_GetBySlug_NotFound(t *testing.T) {
	uc := usecase.NewTenantUsecase(newFakeRepo())

	got, err := uc.GetBySlug(context.Background(), "non-existent-slug")

	assert.Nil(t, got)
	assert.ErrorIs(t, err, repository.ErrNotFound)
}

func TestTenant_IsActive(t *testing.T) {
	active := &domain.Tenant{Status: domain.TenantStatusActive}
	suspended := &domain.Tenant{Status: domain.TenantStatusSuspended}
	cancelled := &domain.Tenant{Status: domain.TenantStatusCancelled}

	assert.True(t, active.IsActive())
	assert.False(t, suspended.IsActive())
	assert.False(t, cancelled.IsActive())
}
