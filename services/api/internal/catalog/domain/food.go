package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// FoodSource defines the origin of a food item.
type FoodSource string

const (
	SourceSystem FoodSource = "system"
	SourceTACO   FoodSource = "taco"
	SourceIBGE   FoodSource = "ibge"
	SourceTenant FoodSource = "tenant"
)

// Sentinel errors.
var (
	ErrNotFound      = errors.New("catalog: food item not found")
	ErrNotEditable   = errors.New("catalog: cannot edit system food items")
)

// FoodItem represents a food in the catalog.
type FoodItem struct {
	ID            uuid.UUID
	TenantID      *uuid.UUID // nil = global/system item
	Name          string
	FoodGroup     string
	Brand         string
	Barcode       string
	ServingSizeG  float64
	ServingLabel  string
	Source        FoodSource
	Active        bool
	CreatedAt     time.Time
	UpdatedAt     time.Time

	// Nested (loaded when needed)
	NutritionFacts    *NutritionFacts
	HouseholdMeasures []HouseholdMeasure
}

// Validate checks business rules for a FoodItem.
func (f *FoodItem) Validate() error {
	if f.Name == "" {
		return errors.New("catalog: name is required")
	}
	if f.FoodGroup == "" {
		return errors.New("catalog: food_group is required")
	}
	if f.ServingSizeG <= 0 {
		return errors.New("catalog: serving_size_g must be positive")
	}
	return nil
}

// NutritionFacts stores nutritional data per serving for a food item.
type NutritionFacts struct {
	ID             uuid.UUID
	FoodItemID     uuid.UUID
	CaloriesKcal   *float64
	ProteinG       *float64
	CarbsG         *float64
	FiberG         *float64
	SugarG         *float64
	TotalFatG      *float64
	SaturatedFatG  *float64
	TransFatG      *float64
	CholesterolMg  *float64
	SodiumMg       *float64
	PotassiumMg    *float64
	CalciumMg      *float64
	IronMg         *float64
	VitaminAMcg    *float64
	VitaminCMg     *float64
	VitaminDMcg    *float64
	VitaminB12Mcg  *float64
	ZincMg         *float64
	MagnesiumMg    *float64
}

// HouseholdMeasure represents a common serving measure for a food item.
type HouseholdMeasure struct {
	ID         uuid.UUID
	FoodItemID uuid.UUID
	Label      string  // e.g., "colher de sopa", "xícara"
	Grams      float64
}
