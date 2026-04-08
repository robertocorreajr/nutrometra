package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"nutrometra/api/internal/tenancy/domain"
	"nutrometra/api/internal/tenancy/usecase"
)

type fakeTenantRepo struct {
	tenants map[uuid.UUID]*domain.Tenant
	members map[string]bool // key: "tenantID:userID"
}

func (f *fakeTenantRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Tenant, error) {
	t, ok := f.tenants[id]
	if !ok {
		return nil, domain.ErrTenantNotFound
	}
	return t, nil
}

func (f *fakeTenantRepo) GetBySlug(ctx context.Context, slug string) (*domain.Tenant, error) {
	for _, t := range f.tenants {
		if t.Slug == slug {
			return t, nil
		}
	}
	return nil, domain.ErrTenantNotFound
}

func (f *fakeTenantRepo) IsMember(ctx context.Context, tenantID, userID uuid.UUID) (bool, error) {
	key := tenantID.String() + ":" + userID.String()
	return f.members[key], nil
}

func (f *fakeTenantRepo) ListByUser(ctx context.Context, userID uuid.UUID) ([]domain.Tenant, error) {
	return nil, nil
}

func newFakeRepo(tenants ...*domain.Tenant) *fakeTenantRepo {
	m := make(map[uuid.UUID]*domain.Tenant, len(tenants))
	for _, t := range tenants {
		m[t.ID] = t
	}
	return &fakeTenantRepo{tenants: m, members: make(map[string]bool)}
}

func sampleTenant() *domain.Tenant {
	return &domain.Tenant{
		ID:          uuid.New(),
		Type:        domain.TenantTypeClinic,
		LegalName:   "Acme Nutrition Ltda",
		DisplayName: "Acme Nutrition",
		Slug:        "acme-nutrition",
		Status:      domain.TenantStatusActive,
		Timezone:    "America/Sao_Paulo",
		Locale:      "pt-BR",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func TestTenantUsecase_GetByID_Found(t *testing.T) {
	tenant := sampleTenant()
	uc := usecase.NewTenantUsecase(newFakeRepo(tenant))

	got, err := uc.GetByID(context.Background(), tenant.ID)

	require.NoError(t, err)
	assert.Equal(t, tenant.ID, got.ID)
	assert.Equal(t, tenant.DisplayName, got.DisplayName)
	assert.Equal(t, tenant.Slug, got.Slug)
}

func TestTenantUsecase_GetByID_NotFound(t *testing.T) {
	uc := usecase.NewTenantUsecase(newFakeRepo())

	got, err := uc.GetByID(context.Background(), uuid.New())

	assert.Nil(t, got)
	assert.ErrorIs(t, err, domain.ErrTenantNotFound)
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
	assert.ErrorIs(t, err, domain.ErrTenantNotFound)
}

func TestTenantUsecase_IsMember(t *testing.T) {
	tenant := sampleTenant()
	userID := uuid.New()
	repo := newFakeRepo(tenant)
	repo.members[tenant.ID.String()+":"+userID.String()] = true
	uc := usecase.NewTenantUsecase(repo)

	ok, err := uc.IsMember(context.Background(), tenant.ID, userID)
	require.NoError(t, err)
	assert.True(t, ok)

	ok, err = uc.IsMember(context.Background(), tenant.ID, uuid.New())
	require.NoError(t, err)
	assert.False(t, ok)
}
