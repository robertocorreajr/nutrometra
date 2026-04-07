package usecase

import (
	"context"
	"fmt"
	"time"

	"nutrometra/api/internal/diet/domain"
	"nutrometra/api/internal/platform/db"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DietRepository defines the data access contract for diets.
type DietRepository interface {
	CreateDiet(ctx context.Context, d *domain.Diet) error
	GetByID(ctx context.Context, tenantID, id uuid.UUID) (*domain.Diet, error)
	UpdateDiet(ctx context.Context, d *domain.Diet) error
	DeleteDiet(ctx context.Context, tenantID, id uuid.UUID) error
	ListByPatient(ctx context.Context, tenantID, patientID uuid.UUID) ([]domain.Diet, error)

	CreateMeal(ctx context.Context, m *domain.DietMeal) error
	UpdateMeal(ctx context.Context, m *domain.DietMeal) error
	DeleteMeal(ctx context.Context, tenantID, id uuid.UUID) error

	CreateMealItem(ctx context.Context, item *domain.DietMealItem) error
	UpdateMealItem(ctx context.Context, item *domain.DietMealItem) error
	DeleteMealItem(ctx context.Context, tenantID, id uuid.UUID) error

	CreateSubstitution(ctx context.Context, sub *domain.DietSubstitution) error
	DeleteSubstitution(ctx context.Context, tenantID, id uuid.UUID) error

	// Transactional methods
	PublishDiet(ctx context.Context, tx pgx.Tx, d *domain.Diet) error
	SupersedePublications(ctx context.Context, tx pgx.Tx, tenantID, patientID uuid.UUID) error
	CreatePublication(ctx context.Context, tx pgx.Tx, pub *domain.DietPublication) error
	GetDietForVersion(ctx context.Context, tx pgx.Tx, tenantID, id uuid.UUID) (*domain.Diet, error)
	ArchiveDiet(ctx context.Context, dietID, tenantID uuid.UUID) error
	ArchiveDietTx(ctx context.Context, tx pgx.Tx, dietID, tenantID uuid.UUID) error
	CreateDietTx(ctx context.Context, tx pgx.Tx, d *domain.Diet) error
	CreateMealTx(ctx context.Context, tx pgx.Tx, m *domain.DietMeal) error
	CreateMealItemTx(ctx context.Context, tx pgx.Tx, item *domain.DietMealItem) error
	CreateSubstitutionTx(ctx context.Context, tx pgx.Tx, sub *domain.DietSubstitution) error
}

// Usecase contains business logic for diets.
type Usecase struct {
	repo DietRepository
	pool *pgxpool.Pool
}

// New creates a diet Usecase.
func New(repo DietRepository, pool *pgxpool.Pool) *Usecase {
	return &Usecase{repo: repo, pool: pool}
}

// --- Diet CRUD ---

// CreateDiet creates a new diet in draft status.
func (uc *Usecase) CreateDiet(ctx context.Context, d *domain.Diet) error {
	if err := d.Validate(); err != nil {
		return err
	}
	d.ID = uuid.New()
	now := time.Now().UTC()
	d.CreatedAt = now
	d.UpdatedAt = now
	d.Status = domain.StatusDraft
	d.VersionNumber = 1
	return uc.repo.CreateDiet(ctx, d)
}

// GetByID returns a diet with its full graph.
func (uc *Usecase) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*domain.Diet, error) {
	return uc.repo.GetByID(ctx, tenantID, id)
}

// UpdateDiet updates a draft diet's top-level fields.
func (uc *Usecase) UpdateDiet(ctx context.Context, d *domain.Diet) error {
	if err := d.Validate(); err != nil {
		return err
	}
	existing, err := uc.repo.GetByID(ctx, d.TenantID, d.ID)
	if err != nil {
		return err
	}
	if existing.Status != domain.StatusDraft {
		return domain.ErrDietNotDraft
	}
	return uc.repo.UpdateDiet(ctx, d)
}

// DeleteDiet deletes a draft diet.
func (uc *Usecase) DeleteDiet(ctx context.Context, tenantID, id uuid.UUID) error {
	existing, err := uc.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return err
	}
	if existing.Status != domain.StatusDraft {
		return domain.ErrDietNotDraft
	}
	return uc.repo.DeleteDiet(ctx, tenantID, id)
}

// ListByPatient returns diets for a patient (without nested graph).
func (uc *Usecase) ListByPatient(ctx context.Context, tenantID, patientID uuid.UUID) ([]domain.Diet, error) {
	return uc.repo.ListByPatient(ctx, tenantID, patientID)
}

// --- Meals ---

// AddMeal adds a meal to a diet.
func (uc *Usecase) AddMeal(ctx context.Context, m *domain.DietMeal) error {
	if err := m.Validate(); err != nil {
		return err
	}
	m.ID = uuid.New()
	return uc.repo.CreateMeal(ctx, m)
}

// UpdateMeal updates a diet meal.
func (uc *Usecase) UpdateMeal(ctx context.Context, m *domain.DietMeal) error {
	if err := m.Validate(); err != nil {
		return err
	}
	return uc.repo.UpdateMeal(ctx, m)
}

// RemoveMeal removes a meal from a diet.
func (uc *Usecase) RemoveMeal(ctx context.Context, tenantID, id uuid.UUID) error {
	return uc.repo.DeleteMeal(ctx, tenantID, id)
}

// --- Meal Items ---

// AddMealItem adds an item to a meal.
func (uc *Usecase) AddMealItem(ctx context.Context, item *domain.DietMealItem) error {
	if err := item.Validate(); err != nil {
		return err
	}
	item.ID = uuid.New()
	return uc.repo.CreateMealItem(ctx, item)
}

// UpdateMealItem updates a meal item.
func (uc *Usecase) UpdateMealItem(ctx context.Context, item *domain.DietMealItem) error {
	if err := item.Validate(); err != nil {
		return err
	}
	return uc.repo.UpdateMealItem(ctx, item)
}

// RemoveMealItem removes an item from a meal.
func (uc *Usecase) RemoveMealItem(ctx context.Context, tenantID, id uuid.UUID) error {
	return uc.repo.DeleteMealItem(ctx, tenantID, id)
}

// --- Substitutions ---

// AddSubstitution adds a substitution to a meal item.
func (uc *Usecase) AddSubstitution(ctx context.Context, sub *domain.DietSubstitution) error {
	sub.ID = uuid.New()
	return uc.repo.CreateSubstitution(ctx, sub)
}

// RemoveSubstitution removes a substitution.
func (uc *Usecase) RemoveSubstitution(ctx context.Context, tenantID, id uuid.UUID) error {
	return uc.repo.DeleteSubstitution(ctx, tenantID, id)
}

// --- Publishing ---

// PublishDiet publishes a diet and creates a publication record.
// This is a transactional operation that:
// 1. Loads the diet and validates it has content
// 2. Updates status to published
// 3. Supersedes existing active publications for this patient
// 4. Creates a new active publication record
func (uc *Usecase) PublishDiet(ctx context.Context, tenantID, dietID, userID uuid.UUID) error {
	return db.RunInTx(ctx, uc.pool, func(ctx context.Context, tx pgx.Tx) error {
		// Load full diet within transaction
		diet, err := uc.repo.GetDietForVersion(ctx, tx, tenantID, dietID)
		if err != nil {
			return err
		}

		if diet.Status != domain.StatusDraft {
			return domain.ErrAlreadyPublished
		}

		if !diet.HasContent() {
			return domain.ErrDietEmpty
		}

		now := time.Now().UTC()
		diet.PublishedAt = &now

		if err := uc.repo.PublishDiet(ctx, tx, diet); err != nil {
			return err
		}

		if err := uc.repo.SupersedePublications(ctx, tx, tenantID, diet.PatientID); err != nil {
			return err
		}

		pub := &domain.DietPublication{
			ID:                uuid.New(),
			TenantID:          tenantID,
			DietID:            dietID,
			PatientID:         diet.PatientID,
			PublishedByUserID: userID,
			PublishedAt:       now,
			Status:            "active",
		}
		return uc.repo.CreatePublication(ctx, tx, pub)
	})
}

// --- Archive ---

// ArchiveDiet archives a published diet.
func (uc *Usecase) ArchiveDiet(ctx context.Context, tenantID, dietID uuid.UUID) error {
	existing, err := uc.repo.GetByID(ctx, tenantID, dietID)
	if err != nil {
		return err
	}
	if existing.Status != domain.StatusPublished {
		return domain.ErrNotPublished
	}
	return uc.repo.ArchiveDiet(ctx, dietID, tenantID)
}

// --- Versioning ---

// CreateNewVersion creates a deep copy of a diet with incremented version number.
// This is a transactional operation that:
// 1. Loads the full diet graph within the transaction
// 2. Creates a new diet with new IDs, version_number+1, previous_version_id set
// 3. Archives the old diet if it was published
func (uc *Usecase) CreateNewVersion(ctx context.Context, tenantID, dietID uuid.UUID) (*domain.Diet, error) {
	var newDiet *domain.Diet

	err := db.RunInTx(ctx, uc.pool, func(ctx context.Context, tx pgx.Tx) error {
		// Load full graph within transaction
		original, err := uc.repo.GetDietForVersion(ctx, tx, tenantID, dietID)
		if err != nil {
			return err
		}

		// Archive old diet if published
		if original.Status == domain.StatusPublished {
			if err := uc.repo.ArchiveDietTx(ctx, tx, original.ID, tenantID); err != nil {
				return fmt.Errorf("diet: archive_original: %w", err)
			}
		}

		// Create deep copy
		now := time.Now().UTC()
		newID := uuid.New()
		originalID := original.ID

		newDiet = &domain.Diet{
			ID:                newID,
			TenantID:          original.TenantID,
			PatientID:         original.PatientID,
			ProfessionalID:    original.ProfessionalID,
			Title:             original.Title,
			Objective:         original.Objective,
			Status:            domain.StatusDraft,
			VersionNumber:     original.VersionNumber + 1,
			PreviousVersionID: &originalID,
			ValidFrom:         original.ValidFrom,
			ValidUntil:        original.ValidUntil,
			CreatedAt:         now,
			UpdatedAt:         now,
		}

		if err := uc.repo.CreateDietTx(ctx, tx, newDiet); err != nil {
			return err
		}

		// Deep copy meals, items, and substitutions
		for _, meal := range original.Meals {
			newMealID := uuid.New()
			newMeal := &domain.DietMeal{
				ID:        newMealID,
				TenantID:  meal.TenantID,
				DietID:    newID,
				MealName:  meal.MealName,
				MealOrder: meal.MealOrder,
				Notes:     meal.Notes,
			}
			if err := uc.repo.CreateMealTx(ctx, tx, newMeal); err != nil {
				return err
			}

			for _, item := range meal.Items {
				newItemID := uuid.New()
				newItem := &domain.DietMealItem{
					ID:                 newItemID,
					TenantID:           item.TenantID,
					DietMealID:         newMealID,
					FoodItemID:         item.FoodItemID,
					QuantityValue:      item.QuantityValue,
					QuantityUnit:       item.QuantityUnit,
					HouseholdMeasureID: item.HouseholdMeasureID,
					AmountDescription:  item.AmountDescription,
					PreparationNotes:   item.PreparationNotes,
					SortOrder:          item.SortOrder,
				}
				if err := uc.repo.CreateMealItemTx(ctx, tx, newItem); err != nil {
					return err
				}

				for _, sub := range item.Substitutions {
					newSub := &domain.DietSubstitution{
						ID:                   uuid.New(),
						TenantID:             sub.TenantID,
						DietMealItemID:       newItemID,
						SubstituteFoodItemID: sub.SubstituteFoodItemID,
						QuantityValue:        sub.QuantityValue,
						QuantityUnit:         sub.QuantityUnit,
						Notes:                sub.Notes,
						SortOrder:            sub.SortOrder,
					}
					if err := uc.repo.CreateSubstitutionTx(ctx, tx, newSub); err != nil {
						return err
					}
				}
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return newDiet, nil
}
