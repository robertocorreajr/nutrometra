package repository

import (
	"context"
	"errors"
	"fmt"

	"nutrometra/api/internal/diet/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository provides diet data access.
type Repository struct {
	pool *pgxpool.Pool
}

// New creates a diet Repository.
func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// --- Diet ---

// CreateDiet inserts a new diet.
func (r *Repository) CreateDiet(ctx context.Context, d *domain.Diet) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO diets
			(id, tenant_id, patient_id, professional_id, title, objective,
			 status, version_number, previous_version_id,
			 published_at, valid_from, valid_until, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)`,
		d.ID, d.TenantID, d.PatientID, d.ProfessionalID,
		d.Title, nilIfEmpty(d.Objective),
		string(d.Status), d.VersionNumber, d.PreviousVersionID,
		d.PublishedAt, d.ValidFrom, d.ValidUntil,
		d.CreatedAt, d.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("diet: create_diet: %w", err)
	}
	return nil
}

// GetByID returns a diet with its full graph (meals, items, substitutions).
func (r *Repository) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*domain.Diet, error) {
	d, err := r.scanDiet(ctx, r.pool, tenantID, id)
	if err != nil {
		return nil, err
	}
	if err := r.loadGraph(ctx, r.pool, d); err != nil {
		return nil, err
	}
	return d, nil
}

// UpdateDiet updates diet-level fields only.
func (r *Repository) UpdateDiet(ctx context.Context, d *domain.Diet) error {
	tag, err := r.pool.Exec(ctx,
		`UPDATE diets SET
			title = $3, objective = $4,
			valid_from = $5, valid_until = $6, updated_at = NOW()
		 WHERE id = $1 AND tenant_id = $2`,
		d.ID, d.TenantID,
		d.Title, nilIfEmpty(d.Objective),
		d.ValidFrom, d.ValidUntil,
	)
	if err != nil {
		return fmt.Errorf("diet: update_diet: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// DeleteDiet deletes a diet. Caller must ensure status=draft.
func (r *Repository) DeleteDiet(ctx context.Context, tenantID, id uuid.UUID) error {
	tag, err := r.pool.Exec(ctx,
		`DELETE FROM diets WHERE id = $1 AND tenant_id = $2 AND status = 'draft'`,
		id, tenantID,
	)
	if err != nil {
		return fmt.Errorf("diet: delete_diet: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// ListByPatient returns diets for a patient (without nested graph).
func (r *Repository) ListByPatient(ctx context.Context, tenantID, patientID uuid.UUID) ([]domain.Diet, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, tenant_id, patient_id, professional_id, title,
			COALESCE(objective,''), status, version_number, previous_version_id,
			published_at, valid_from, valid_until, created_at, updated_at
		 FROM diets WHERE tenant_id = $1 AND patient_id = $2
		 ORDER BY created_at DESC`,
		tenantID, patientID,
	)
	if err != nil {
		return nil, fmt.Errorf("diet: list_by_patient: %w", err)
	}
	defer rows.Close()

	var result []domain.Diet
	for rows.Next() {
		var d domain.Diet
		if err := rows.Scan(
			&d.ID, &d.TenantID, &d.PatientID, &d.ProfessionalID,
			&d.Title, &d.Objective, &d.Status, &d.VersionNumber,
			&d.PreviousVersionID, &d.PublishedAt, &d.ValidFrom,
			&d.ValidUntil, &d.CreatedAt, &d.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("diet: scan_diet: %w", err)
		}
		result = append(result, d)
	}
	return result, rows.Err()
}

// --- Meals ---

// CreateMeal inserts a diet meal.
func (r *Repository) CreateMeal(ctx context.Context, m *domain.DietMeal) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO diet_meals (id, tenant_id, diet_id, meal_name, meal_order, notes)
		 VALUES ($1,$2,$3,$4,$5,$6)`,
		m.ID, m.TenantID, m.DietID, m.MealName, m.MealOrder, nilIfEmpty(m.Notes),
	)
	if err != nil {
		return fmt.Errorf("diet: create_meal: %w", err)
	}
	return nil
}

// UpdateMeal updates a diet meal.
func (r *Repository) UpdateMeal(ctx context.Context, m *domain.DietMeal) error {
	tag, err := r.pool.Exec(ctx,
		`UPDATE diet_meals SET meal_name = $3, meal_order = $4, notes = $5
		 WHERE id = $1 AND tenant_id = $2`,
		m.ID, m.TenantID, m.MealName, m.MealOrder, nilIfEmpty(m.Notes),
	)
	if err != nil {
		return fmt.Errorf("diet: update_meal: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// DeleteMeal deletes a diet meal.
func (r *Repository) DeleteMeal(ctx context.Context, tenantID, id uuid.UUID) error {
	tag, err := r.pool.Exec(ctx,
		`DELETE FROM diet_meals WHERE id = $1 AND tenant_id = $2`,
		id, tenantID,
	)
	if err != nil {
		return fmt.Errorf("diet: delete_meal: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// --- Meal Items ---

// CreateMealItem inserts a diet meal item.
func (r *Repository) CreateMealItem(ctx context.Context, item *domain.DietMealItem) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO diet_meal_items
			(id, tenant_id, diet_meal_id, food_item_id,
			 quantity_value, quantity_unit, household_measure_id,
			 amount_description, preparation_notes, sort_order)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
		item.ID, item.TenantID, item.DietMealID, item.FoodItemID,
		item.QuantityValue, item.QuantityUnit, item.HouseholdMeasureID,
		nilIfEmpty(item.AmountDescription), nilIfEmpty(item.PreparationNotes),
		item.SortOrder,
	)
	if err != nil {
		return fmt.Errorf("diet: create_meal_item: %w", err)
	}
	return nil
}

// UpdateMealItem updates a diet meal item.
func (r *Repository) UpdateMealItem(ctx context.Context, item *domain.DietMealItem) error {
	tag, err := r.pool.Exec(ctx,
		`UPDATE diet_meal_items SET
			food_item_id = $3, quantity_value = $4, quantity_unit = $5,
			household_measure_id = $6, amount_description = $7,
			preparation_notes = $8, sort_order = $9
		 WHERE id = $1 AND tenant_id = $2`,
		item.ID, item.TenantID,
		item.FoodItemID, item.QuantityValue, item.QuantityUnit,
		item.HouseholdMeasureID, nilIfEmpty(item.AmountDescription),
		nilIfEmpty(item.PreparationNotes), item.SortOrder,
	)
	if err != nil {
		return fmt.Errorf("diet: update_meal_item: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// DeleteMealItem deletes a diet meal item.
func (r *Repository) DeleteMealItem(ctx context.Context, tenantID, id uuid.UUID) error {
	tag, err := r.pool.Exec(ctx,
		`DELETE FROM diet_meal_items WHERE id = $1 AND tenant_id = $2`,
		id, tenantID,
	)
	if err != nil {
		return fmt.Errorf("diet: delete_meal_item: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// --- Substitutions ---

// CreateSubstitution inserts a diet substitution.
func (r *Repository) CreateSubstitution(ctx context.Context, sub *domain.DietSubstitution) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO diet_substitutions
			(id, tenant_id, diet_meal_item_id, substitute_food_item_id,
			 quantity_value, quantity_unit, notes, sort_order)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		sub.ID, sub.TenantID, sub.DietMealItemID, sub.SubstituteFoodItemID,
		sub.QuantityValue, sub.QuantityUnit,
		nilIfEmpty(sub.Notes), sub.SortOrder,
	)
	if err != nil {
		return fmt.Errorf("diet: create_substitution: %w", err)
	}
	return nil
}

// DeleteSubstitution deletes a diet substitution.
func (r *Repository) DeleteSubstitution(ctx context.Context, tenantID, id uuid.UUID) error {
	tag, err := r.pool.Exec(ctx,
		`DELETE FROM diet_substitutions WHERE id = $1 AND tenant_id = $2`,
		id, tenantID,
	)
	if err != nil {
		return fmt.Errorf("diet: delete_substitution: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// --- Transactional (publish) ---

// PublishDiet updates a diet's status to published within a transaction.
func (r *Repository) PublishDiet(ctx context.Context, tx pgx.Tx, d *domain.Diet) error {
	tag, err := tx.Exec(ctx,
		`UPDATE diets SET status = $3, published_at = $4, updated_at = NOW()
		 WHERE id = $1 AND tenant_id = $2`,
		d.ID, d.TenantID, string(domain.StatusPublished), d.PublishedAt,
	)
	if err != nil {
		return fmt.Errorf("diet: publish_diet: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// SupersedePublications marks all active publications for a patient as superseded.
func (r *Repository) SupersedePublications(ctx context.Context, tx pgx.Tx, tenantID, patientID uuid.UUID) error {
	_, err := tx.Exec(ctx,
		`UPDATE diet_publications SET status = 'superseded'
		 WHERE tenant_id = $1 AND patient_id = $2 AND status = 'active'`,
		tenantID, patientID,
	)
	if err != nil {
		return fmt.Errorf("diet: supersede_publications: %w", err)
	}
	return nil
}

// CreatePublication inserts a new diet publication record within a transaction.
func (r *Repository) CreatePublication(ctx context.Context, tx pgx.Tx, pub *domain.DietPublication) error {
	_, err := tx.Exec(ctx,
		`INSERT INTO diet_publications
			(id, tenant_id, diet_id, patient_id, published_by_user_id, published_at, status)
		 VALUES ($1,$2,$3,$4,$5,$6,$7)`,
		pub.ID, pub.TenantID, pub.DietID, pub.PatientID,
		pub.PublishedByUserID, pub.PublishedAt, pub.Status,
	)
	if err != nil {
		return fmt.Errorf("diet: create_publication: %w", err)
	}
	return nil
}

// GetDietForVersion loads a full diet graph within a transaction (for CreateNewVersion).
func (r *Repository) GetDietForVersion(ctx context.Context, tx pgx.Tx, tenantID, id uuid.UUID) (*domain.Diet, error) {
	d, err := r.scanDiet(ctx, tx, tenantID, id)
	if err != nil {
		return nil, err
	}
	if err := r.loadGraph(ctx, tx, d); err != nil {
		return nil, err
	}
	return d, nil
}

// ArchiveDiet updates a diet's status to archived.
func (r *Repository) ArchiveDiet(ctx context.Context, dietID, tenantID uuid.UUID) error {
	tag, err := r.pool.Exec(ctx,
		`UPDATE diets SET status = 'archived', updated_at = NOW()
		 WHERE id = $1 AND tenant_id = $2`,
		dietID, tenantID,
	)
	if err != nil {
		return fmt.Errorf("diet: archive_diet: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// ArchiveDietTx updates a diet's status to archived within a transaction.
func (r *Repository) ArchiveDietTx(ctx context.Context, tx pgx.Tx, dietID, tenantID uuid.UUID) error {
	tag, err := tx.Exec(ctx,
		`UPDATE diets SET status = 'archived', updated_at = NOW()
		 WHERE id = $1 AND tenant_id = $2`,
		dietID, tenantID,
	)
	if err != nil {
		return fmt.Errorf("diet: archive_diet_tx: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// CreateDietTx inserts a new diet within a transaction (for CreateNewVersion).
func (r *Repository) CreateDietTx(ctx context.Context, tx pgx.Tx, d *domain.Diet) error {
	_, err := tx.Exec(ctx,
		`INSERT INTO diets
			(id, tenant_id, patient_id, professional_id, title, objective,
			 status, version_number, previous_version_id,
			 published_at, valid_from, valid_until, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)`,
		d.ID, d.TenantID, d.PatientID, d.ProfessionalID,
		d.Title, nilIfEmpty(d.Objective),
		string(d.Status), d.VersionNumber, d.PreviousVersionID,
		d.PublishedAt, d.ValidFrom, d.ValidUntil,
		d.CreatedAt, d.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("diet: create_diet_tx: %w", err)
	}
	return nil
}

// CreateMealTx inserts a diet meal within a transaction.
func (r *Repository) CreateMealTx(ctx context.Context, tx pgx.Tx, m *domain.DietMeal) error {
	_, err := tx.Exec(ctx,
		`INSERT INTO diet_meals (id, tenant_id, diet_id, meal_name, meal_order, notes)
		 VALUES ($1,$2,$3,$4,$5,$6)`,
		m.ID, m.TenantID, m.DietID, m.MealName, m.MealOrder, nilIfEmpty(m.Notes),
	)
	if err != nil {
		return fmt.Errorf("diet: create_meal_tx: %w", err)
	}
	return nil
}

// CreateMealItemTx inserts a diet meal item within a transaction.
func (r *Repository) CreateMealItemTx(ctx context.Context, tx pgx.Tx, item *domain.DietMealItem) error {
	_, err := tx.Exec(ctx,
		`INSERT INTO diet_meal_items
			(id, tenant_id, diet_meal_id, food_item_id,
			 quantity_value, quantity_unit, household_measure_id,
			 amount_description, preparation_notes, sort_order)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
		item.ID, item.TenantID, item.DietMealID, item.FoodItemID,
		item.QuantityValue, item.QuantityUnit, item.HouseholdMeasureID,
		nilIfEmpty(item.AmountDescription), nilIfEmpty(item.PreparationNotes),
		item.SortOrder,
	)
	if err != nil {
		return fmt.Errorf("diet: create_meal_item_tx: %w", err)
	}
	return nil
}

// CreateSubstitutionTx inserts a diet substitution within a transaction.
func (r *Repository) CreateSubstitutionTx(ctx context.Context, tx pgx.Tx, sub *domain.DietSubstitution) error {
	_, err := tx.Exec(ctx,
		`INSERT INTO diet_substitutions
			(id, tenant_id, diet_meal_item_id, substitute_food_item_id,
			 quantity_value, quantity_unit, notes, sort_order)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		sub.ID, sub.TenantID, sub.DietMealItemID, sub.SubstituteFoodItemID,
		sub.QuantityValue, sub.QuantityUnit,
		nilIfEmpty(sub.Notes), sub.SortOrder,
	)
	if err != nil {
		return fmt.Errorf("diet: create_substitution_tx: %w", err)
	}
	return nil
}

// --- internal helpers ---

// querier is satisfied by both *pgxpool.Pool and pgx.Tx.
type querier interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

// scanDiet loads a single diet row.
func (r *Repository) scanDiet(ctx context.Context, q querier, tenantID, id uuid.UUID) (*domain.Diet, error) {
	d := &domain.Diet{}
	err := q.QueryRow(ctx,
		`SELECT id, tenant_id, patient_id, professional_id, title,
			COALESCE(objective,''), status, version_number, previous_version_id,
			published_at, valid_from, valid_until, created_at, updated_at
		 FROM diets WHERE id = $1 AND tenant_id = $2`,
		id, tenantID,
	).Scan(
		&d.ID, &d.TenantID, &d.PatientID, &d.ProfessionalID,
		&d.Title, &d.Objective, &d.Status, &d.VersionNumber,
		&d.PreviousVersionID, &d.PublishedAt, &d.ValidFrom,
		&d.ValidUntil, &d.CreatedAt, &d.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("diet: get_by_id: %w", err)
	}
	return d, nil
}

// loadGraph loads meals, items, and substitutions for a diet.
func (r *Repository) loadGraph(ctx context.Context, q querier, d *domain.Diet) error {
	// 1. Load meals
	mealRows, err := q.Query(ctx,
		`SELECT id, tenant_id, diet_id, meal_name, meal_order, COALESCE(notes,'')
		 FROM diet_meals WHERE diet_id = $1 AND tenant_id = $2
		 ORDER BY meal_order`,
		d.ID, d.TenantID,
	)
	if err != nil {
		return fmt.Errorf("diet: load_meals: %w", err)
	}
	defer mealRows.Close()

	var meals []domain.DietMeal
	mealIndex := map[uuid.UUID]int{}
	for mealRows.Next() {
		var m domain.DietMeal
		if err := mealRows.Scan(&m.ID, &m.TenantID, &m.DietID, &m.MealName, &m.MealOrder, &m.Notes); err != nil {
			return fmt.Errorf("diet: scan_meal: %w", err)
		}
		mealIndex[m.ID] = len(meals)
		meals = append(meals, m)
	}
	if err := mealRows.Err(); err != nil {
		return fmt.Errorf("diet: meals_rows: %w", err)
	}

	if len(meals) == 0 {
		d.Meals = meals
		return nil
	}

	// 2. Collect meal IDs for items query
	mealIDs := make([]uuid.UUID, len(meals))
	for i := range meals {
		mealIDs[i] = meals[i].ID
	}

	// 3. Load items
	itemRows, err := q.Query(ctx,
		`SELECT id, tenant_id, diet_meal_id, food_item_id,
			quantity_value, quantity_unit, household_measure_id,
			COALESCE(amount_description,''), COALESCE(preparation_notes,''), sort_order
		 FROM diet_meal_items WHERE diet_meal_id = ANY($1) AND tenant_id = $2
		 ORDER BY sort_order`,
		mealIDs, d.TenantID,
	)
	if err != nil {
		return fmt.Errorf("diet: load_items: %w", err)
	}
	defer itemRows.Close()

	var allItems []domain.DietMealItem
	itemIndex := map[uuid.UUID]int{}
	itemToMeal := map[uuid.UUID]uuid.UUID{}
	for itemRows.Next() {
		var item domain.DietMealItem
		if err := itemRows.Scan(
			&item.ID, &item.TenantID, &item.DietMealID, &item.FoodItemID,
			&item.QuantityValue, &item.QuantityUnit, &item.HouseholdMeasureID,
			&item.AmountDescription, &item.PreparationNotes, &item.SortOrder,
		); err != nil {
			return fmt.Errorf("diet: scan_item: %w", err)
		}
		itemIndex[item.ID] = len(allItems)
		itemToMeal[item.ID] = item.DietMealID
		allItems = append(allItems, item)
	}
	if err := itemRows.Err(); err != nil {
		return fmt.Errorf("diet: items_rows: %w", err)
	}

	// 4. Load substitutions
	if len(allItems) > 0 {
		itemIDs := make([]uuid.UUID, len(allItems))
		for i := range allItems {
			itemIDs[i] = allItems[i].ID
		}

		subRows, err := q.Query(ctx,
			`SELECT id, tenant_id, diet_meal_item_id, substitute_food_item_id,
				quantity_value, quantity_unit, COALESCE(notes,''), sort_order
			 FROM diet_substitutions WHERE diet_meal_item_id = ANY($1) AND tenant_id = $2
			 ORDER BY sort_order`,
			itemIDs, d.TenantID,
		)
		if err != nil {
			return fmt.Errorf("diet: load_substitutions: %w", err)
		}
		defer subRows.Close()

		for subRows.Next() {
			var sub domain.DietSubstitution
			if err := subRows.Scan(
				&sub.ID, &sub.TenantID, &sub.DietMealItemID, &sub.SubstituteFoodItemID,
				&sub.QuantityValue, &sub.QuantityUnit, &sub.Notes, &sub.SortOrder,
			); err != nil {
				return fmt.Errorf("diet: scan_substitution: %w", err)
			}
			if idx, ok := itemIndex[sub.DietMealItemID]; ok {
				allItems[idx].Substitutions = append(allItems[idx].Substitutions, sub)
			}
		}
		if err := subRows.Err(); err != nil {
			return fmt.Errorf("diet: subs_rows: %w", err)
		}
	}

	// 5. Assemble items into meals
	for i := range allItems {
		mealID := itemToMeal[allItems[i].ID]
		if idx, ok := mealIndex[mealID]; ok {
			meals[idx].Items = append(meals[idx].Items, allItems[i])
		}
	}

	d.Meals = meals
	return nil
}

func nilIfEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
