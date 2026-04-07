package usecase

import (
	"context"
	"time"

	"nutrometra/api/internal/catalog/domain"

	"github.com/google/uuid"
)

// CatalogRepository defines the data access contract.
type CatalogRepository interface {
	Search(ctx context.Context, tenantID *uuid.UUID, query string, limit int) ([]domain.FoodItem, error)
	List(ctx context.Context, tenantID *uuid.UUID, foodGroup string, limit int) ([]domain.FoodItem, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.FoodItem, error)
	Create(ctx context.Context, f *domain.FoodItem) error
	Update(ctx context.Context, f *domain.FoodItem) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
	UpsertNutritionFacts(ctx context.Context, nf *domain.NutritionFacts) error
	ListFoodGroups(ctx context.Context, tenantID *uuid.UUID) ([]string, error)
}

// Usecase contains business logic for the food catalog.
type Usecase struct {
	repo CatalogRepository
}

// New creates a catalog Usecase.
func New(repo CatalogRepository) *Usecase {
	return &Usecase{repo: repo}
}

// Search searches food items by name.
func (uc *Usecase) Search(ctx context.Context, tenantID *uuid.UUID, query string, limit int) ([]domain.FoodItem, error) {
	return uc.repo.Search(ctx, tenantID, query, limit)
}

// List returns food items, optionally filtered by food group.
func (uc *Usecase) List(ctx context.Context, tenantID *uuid.UUID, foodGroup string, limit int) ([]domain.FoodItem, error) {
	return uc.repo.List(ctx, tenantID, foodGroup, limit)
}

// GetByID returns a food item with nutrition facts and household measures.
func (uc *Usecase) GetByID(ctx context.Context, id uuid.UUID) (*domain.FoodItem, error) {
	return uc.repo.GetByID(ctx, id)
}

// Create creates a new tenant-scoped food item.
func (uc *Usecase) Create(ctx context.Context, f *domain.FoodItem, nf *domain.NutritionFacts) error {
	if err := f.Validate(); err != nil {
		return err
	}

	f.ID = uuid.New()
	now := time.Now().UTC()
	f.CreatedAt = now
	f.UpdatedAt = now
	f.Active = true
	f.Source = domain.SourceTenant

	if err := uc.repo.Create(ctx, f); err != nil {
		return err
	}

	if nf != nil {
		nf.ID = uuid.New()
		nf.FoodItemID = f.ID
		if err := uc.repo.UpsertNutritionFacts(ctx, nf); err != nil {
			return err
		}
	}

	return nil
}

// Update updates a tenant-scoped food item.
func (uc *Usecase) Update(ctx context.Context, f *domain.FoodItem, nf *domain.NutritionFacts) error {
	if err := f.Validate(); err != nil {
		return err
	}

	if err := uc.repo.Update(ctx, f); err != nil {
		return err
	}

	if nf != nil {
		nf.FoodItemID = f.ID
		if nf.ID == uuid.Nil {
			nf.ID = uuid.New()
		}
		if err := uc.repo.UpsertNutritionFacts(ctx, nf); err != nil {
			return err
		}
	}

	return nil
}

// Delete soft-deletes a tenant-scoped food item.
func (uc *Usecase) Delete(ctx context.Context, id uuid.UUID) error {
	return uc.repo.SoftDelete(ctx, id)
}

// ListFoodGroups returns distinct food groups.
func (uc *Usecase) ListFoodGroups(ctx context.Context, tenantID *uuid.UUID) ([]string, error) {
	return uc.repo.ListFoodGroups(ctx, tenantID)
}
