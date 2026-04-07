package usecase

import (
	"context"

	"github.com/google/uuid"

	"nutrometra/api/internal/tenancy/domain"
	"nutrometra/api/internal/tenancy/repository"
)

// TenantUsecase defines application-level operations for tenants.
type TenantUsecase interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Tenant, error)
	GetBySlug(ctx context.Context, slug string) (*domain.Tenant, error)
}

type tenantUsecase struct {
	repo repository.TenantRepository
}

// NewTenantUsecase returns a TenantUsecase wrapping the provided repository.
func NewTenantUsecase(repo repository.TenantRepository) TenantUsecase {
	return &tenantUsecase{repo: repo}
}

func (uc *tenantUsecase) GetByID(ctx context.Context, id uuid.UUID) (*domain.Tenant, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *tenantUsecase) GetBySlug(ctx context.Context, slug string) (*domain.Tenant, error) {
	return uc.repo.GetBySlug(ctx, slug)
}
