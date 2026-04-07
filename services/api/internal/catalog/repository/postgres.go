package repository

import (
	"context"
	"errors"
	"fmt"

	"nutrometra/api/internal/catalog/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository provides food catalog data access.
type Repository struct {
	pool *pgxpool.Pool
}

// New creates a catalog Repository.
func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// Search returns food items matching a query, combining global + tenant-scoped items.
func (r *Repository) Search(ctx context.Context, tenantID *uuid.UUID, query string, limit int) ([]domain.FoodItem, error) {
	if limit <= 0 {
		limit = 50
	}
	sql := `SELECT id, tenant_id, name, food_group, COALESCE(brand,''), COALESCE(barcode,''),
			serving_size_g, serving_label, source, active, created_at, updated_at
		 FROM food_items
		 WHERE active = TRUE
		   AND (tenant_id IS NULL OR tenant_id = $1)
		   AND to_tsvector('portuguese', name) @@ plainto_tsquery('portuguese', $2)
		 ORDER BY name
		 LIMIT $3`

	rows, err := r.pool.Query(ctx, sql, tenantID, query, limit)
	if err != nil {
		return nil, fmt.Errorf("catalog: search: %w", err)
	}
	defer rows.Close()

	return scanFoodItems(rows)
}

// List returns food items optionally filtered by food group.
func (r *Repository) List(ctx context.Context, tenantID *uuid.UUID, foodGroup string, limit int) ([]domain.FoodItem, error) {
	if limit <= 0 {
		limit = 100
	}

	sql := `SELECT id, tenant_id, name, food_group, COALESCE(brand,''), COALESCE(barcode,''),
			serving_size_g, serving_label, source, active, created_at, updated_at
		 FROM food_items
		 WHERE active = TRUE AND (tenant_id IS NULL OR tenant_id = $1)`
	args := []any{tenantID}
	argIdx := 2

	if foodGroup != "" {
		sql += fmt.Sprintf(" AND food_group = $%d", argIdx)
		args = append(args, foodGroup)
		argIdx++
	}
	sql += fmt.Sprintf(" ORDER BY name LIMIT $%d", argIdx)
	args = append(args, limit)

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("catalog: list: %w", err)
	}
	defer rows.Close()

	return scanFoodItems(rows)
}

// GetByID returns a food item by ID.
func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*domain.FoodItem, error) {
	f := &domain.FoodItem{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, tenant_id, name, food_group, COALESCE(brand,''), COALESCE(barcode,''),
			serving_size_g, serving_label, source, active, created_at, updated_at
		 FROM food_items WHERE id = $1`,
		id,
	).Scan(&f.ID, &f.TenantID, &f.Name, &f.FoodGroup, &f.Brand, &f.Barcode,
		&f.ServingSizeG, &f.ServingLabel, &f.Source, &f.Active, &f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("catalog: get_by_id: %w", err)
	}

	// Load nutrition facts
	nf, err := r.GetNutritionFacts(ctx, id)
	if err == nil {
		f.NutritionFacts = nf
	}

	// Load household measures
	measures, err := r.ListHouseholdMeasures(ctx, id)
	if err == nil {
		f.HouseholdMeasures = measures
	}

	return f, nil
}

// Create inserts a new food item (tenant-scoped only).
func (r *Repository) Create(ctx context.Context, f *domain.FoodItem) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO food_items
			(id, tenant_id, name, food_group, brand, barcode,
			 serving_size_g, serving_label, source, active, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`,
		f.ID, f.TenantID, f.Name, f.FoodGroup,
		nilIfEmpty(f.Brand), nilIfEmpty(f.Barcode),
		f.ServingSizeG, f.ServingLabel, string(f.Source),
		f.Active, f.CreatedAt, f.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("catalog: create: %w", err)
	}
	return nil
}

// Update updates a food item.
func (r *Repository) Update(ctx context.Context, f *domain.FoodItem) error {
	tag, err := r.pool.Exec(ctx,
		`UPDATE food_items SET
			name = $2, food_group = $3, brand = $4, barcode = $5,
			serving_size_g = $6, serving_label = $7, active = $8, updated_at = NOW()
		 WHERE id = $1 AND tenant_id IS NOT NULL`,
		f.ID, f.Name, f.FoodGroup,
		nilIfEmpty(f.Brand), nilIfEmpty(f.Barcode),
		f.ServingSizeG, f.ServingLabel, f.Active,
	)
	if err != nil {
		return fmt.Errorf("catalog: update: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotEditable
	}
	return nil
}

// SoftDelete marks a tenant-scoped food item as inactive.
func (r *Repository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	tag, err := r.pool.Exec(ctx,
		`UPDATE food_items SET active = FALSE, updated_at = NOW() WHERE id = $1 AND tenant_id IS NOT NULL`,
		id,
	)
	if err != nil {
		return fmt.Errorf("catalog: soft_delete: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotEditable
	}
	return nil
}

// --- Nutrition Facts ---

// GetNutritionFacts returns nutrition facts for a food item.
func (r *Repository) GetNutritionFacts(ctx context.Context, foodItemID uuid.UUID) (*domain.NutritionFacts, error) {
	nf := &domain.NutritionFacts{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, sugar_g,
			total_fat_g, saturated_fat_g, trans_fat_g, cholesterol_mg,
			sodium_mg, potassium_mg, calcium_mg, iron_mg,
			vitamin_a_mcg, vitamin_c_mg, vitamin_d_mcg, vitamin_b12_mcg,
			zinc_mg, magnesium_mg
		 FROM nutrition_facts WHERE food_item_id = $1`,
		foodItemID,
	).Scan(&nf.ID, &nf.FoodItemID,
		&nf.CaloriesKcal, &nf.ProteinG, &nf.CarbsG, &nf.FiberG, &nf.SugarG,
		&nf.TotalFatG, &nf.SaturatedFatG, &nf.TransFatG, &nf.CholesterolMg,
		&nf.SodiumMg, &nf.PotassiumMg, &nf.CalciumMg, &nf.IronMg,
		&nf.VitaminAMcg, &nf.VitaminCMg, &nf.VitaminDMcg, &nf.VitaminB12Mcg,
		&nf.ZincMg, &nf.MagnesiumMg)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("catalog: get_nutrition: %w", err)
	}
	return nf, nil
}

// UpsertNutritionFacts creates or updates nutrition facts for a food item.
func (r *Repository) UpsertNutritionFacts(ctx context.Context, nf *domain.NutritionFacts) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO nutrition_facts
			(id, food_item_id, calories_kcal, protein_g, carbs_g, fiber_g, sugar_g,
			 total_fat_g, saturated_fat_g, trans_fat_g, cholesterol_mg,
			 sodium_mg, potassium_mg, calcium_mg, iron_mg,
			 vitamin_a_mcg, vitamin_c_mg, vitamin_d_mcg, vitamin_b12_mcg,
			 zinc_mg, magnesium_mg)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21)
		 ON CONFLICT (food_item_id) DO UPDATE SET
			calories_kcal = EXCLUDED.calories_kcal,
			protein_g = EXCLUDED.protein_g,
			carbs_g = EXCLUDED.carbs_g,
			fiber_g = EXCLUDED.fiber_g,
			sugar_g = EXCLUDED.sugar_g,
			total_fat_g = EXCLUDED.total_fat_g,
			saturated_fat_g = EXCLUDED.saturated_fat_g,
			trans_fat_g = EXCLUDED.trans_fat_g,
			cholesterol_mg = EXCLUDED.cholesterol_mg,
			sodium_mg = EXCLUDED.sodium_mg,
			potassium_mg = EXCLUDED.potassium_mg,
			calcium_mg = EXCLUDED.calcium_mg,
			iron_mg = EXCLUDED.iron_mg,
			vitamin_a_mcg = EXCLUDED.vitamin_a_mcg,
			vitamin_c_mg = EXCLUDED.vitamin_c_mg,
			vitamin_d_mcg = EXCLUDED.vitamin_d_mcg,
			vitamin_b12_mcg = EXCLUDED.vitamin_b12_mcg,
			zinc_mg = EXCLUDED.zinc_mg,
			magnesium_mg = EXCLUDED.magnesium_mg`,
		nf.ID, nf.FoodItemID,
		nf.CaloriesKcal, nf.ProteinG, nf.CarbsG, nf.FiberG, nf.SugarG,
		nf.TotalFatG, nf.SaturatedFatG, nf.TransFatG, nf.CholesterolMg,
		nf.SodiumMg, nf.PotassiumMg, nf.CalciumMg, nf.IronMg,
		nf.VitaminAMcg, nf.VitaminCMg, nf.VitaminDMcg, nf.VitaminB12Mcg,
		nf.ZincMg, nf.MagnesiumMg,
	)
	if err != nil {
		return fmt.Errorf("catalog: upsert_nutrition: %w", err)
	}
	return nil
}

// --- Household Measures ---

// ListHouseholdMeasures returns household measures for a food item.
func (r *Repository) ListHouseholdMeasures(ctx context.Context, foodItemID uuid.UUID) ([]domain.HouseholdMeasure, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, food_item_id, label, grams FROM household_measures WHERE food_item_id = $1 ORDER BY label`,
		foodItemID,
	)
	if err != nil {
		return nil, fmt.Errorf("catalog: list_measures: %w", err)
	}
	defer rows.Close()

	var result []domain.HouseholdMeasure
	for rows.Next() {
		var m domain.HouseholdMeasure
		if err := rows.Scan(&m.ID, &m.FoodItemID, &m.Label, &m.Grams); err != nil {
			return nil, fmt.Errorf("catalog: scan_measure: %w", err)
		}
		result = append(result, m)
	}
	return result, rows.Err()
}

// ListFoodGroups returns distinct food groups.
func (r *Repository) ListFoodGroups(ctx context.Context, tenantID *uuid.UUID) ([]string, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT DISTINCT food_group FROM food_items
		 WHERE active = TRUE AND (tenant_id IS NULL OR tenant_id = $1)
		 ORDER BY food_group`,
		tenantID,
	)
	if err != nil {
		return nil, fmt.Errorf("catalog: list_groups: %w", err)
	}
	defer rows.Close()

	var result []string
	for rows.Next() {
		var g string
		if err := rows.Scan(&g); err != nil {
			return nil, fmt.Errorf("catalog: scan_group: %w", err)
		}
		result = append(result, g)
	}
	return result, rows.Err()
}

func scanFoodItems(rows pgx.Rows) ([]domain.FoodItem, error) {
	var result []domain.FoodItem
	for rows.Next() {
		var f domain.FoodItem
		if err := rows.Scan(&f.ID, &f.TenantID, &f.Name, &f.FoodGroup, &f.Brand, &f.Barcode,
			&f.ServingSizeG, &f.ServingLabel, &f.Source, &f.Active,
			&f.CreatedAt, &f.UpdatedAt); err != nil {
			return nil, fmt.Errorf("catalog: scan: %w", err)
		}
		result = append(result, f)
	}
	return result, rows.Err()
}

func nilIfEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
