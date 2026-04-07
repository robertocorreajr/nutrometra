package usecase

import (
	"context"

	"github.com/google/uuid"

	"nutrometra/api/internal/tenancy/domain"
	"nutrometra/api/internal/tenancy/repository"
)

type TenantUsecase interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Tenant, error)
	GetBySlug(ctx context.Context, slug string) (*domain.Tenant, error)
	IsMember(ctx context.Context, tenantID, userID uuid.UUID) (bool, error)
}

type tenantUsecase struct {
	repo repository.TenantRepository
}

func NewTenantUsecase(repo repository.TenantRepository) TenantUsecase {
	return &tenantUsecase{repo: repo}
}

func (uc *tenantUsecase) GetByID(ctx context.Context, id uuid.UUID) (*domain.Tenant, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *tenantUsecase) GetBySlug(ctx context.Context, slug string) (*domain.Tenant, error) {
	return uc.repo.GetBySlug(ctx, slug)
}

func (uc *tenantUsecase) IsMember(ctx context.Context, tenantID, userID uuid.UUID) (bool, error) {
	return uc.repo.IsMember(ctx, tenantID, userID)
}
