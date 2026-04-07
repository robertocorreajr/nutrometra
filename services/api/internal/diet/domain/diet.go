package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// DietStatus represents the lifecycle state of a diet.
type DietStatus string

const (
	StatusDraft     DietStatus = "draft"
	StatusPublished DietStatus = "published"
	StatusArchived  DietStatus = "archived"
)

// Sentinel errors.
var (
	ErrNotFound         = errors.New("diet: not found")
	ErrDietNotDraft     = errors.New("diet: operation requires draft status")
	ErrDietEmpty        = errors.New("diet: publish requires at least one meal with one item")
	ErrAlreadyPublished = errors.New("diet: diet is already published")
	ErrNotPublished     = errors.New("diet: diet is not published")
)

// Diet represents a nutritional plan prescribed to a patient.
type Diet struct {
	ID                uuid.UUID
	TenantID          uuid.UUID
	PatientID         uuid.UUID
	ProfessionalID    uuid.UUID
	Title             string
	Objective         string
	Status            DietStatus
	VersionNumber     int
	PreviousVersionID *uuid.UUID
	PublishedAt       *time.Time
	ValidFrom         *time.Time
	ValidUntil        *time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
	Meals             []DietMeal
}

// Validate checks business rules for a Diet.
func (d *Diet) Validate() error {
	if d.Title == "" {
		return errors.New("diet: title is required")
	}
	return nil
}

// HasContent returns true if the diet has at least one meal with at least one item.
func (d *Diet) HasContent() bool {
	for i := range d.Meals {
		if len(d.Meals[i].Items) > 0 {
			return true
		}
	}
	return false
}

// DietMeal represents a named meal within a diet (e.g. "Breakfast", "Lunch").
type DietMeal struct {
	ID        uuid.UUID
	TenantID  uuid.UUID
	DietID    uuid.UUID
	MealName  string
	MealOrder int
	Notes     string
	Items     []DietMealItem
}

// Validate checks business rules for a DietMeal.
func (m *DietMeal) Validate() error {
	if m.MealName == "" {
		return errors.New("diet: meal_name is required")
	}
	if m.MealOrder <= 0 {
		return errors.New("diet: meal_order must be greater than zero")
	}
	return nil
}

// DietMealItem represents a food item within a meal.
type DietMealItem struct {
	ID                 uuid.UUID
	TenantID           uuid.UUID
	DietMealID         uuid.UUID
	FoodItemID         uuid.UUID
	QuantityValue      float64
	QuantityUnit       string
	HouseholdMeasureID *uuid.UUID
	AmountDescription  string
	PreparationNotes   string
	SortOrder          int
	Substitutions      []DietSubstitution
}

// Validate checks business rules for a DietMealItem.
func (item *DietMealItem) Validate() error {
	if item.QuantityValue <= 0 {
		return errors.New("diet: quantity_value must be greater than zero")
	}
	if item.QuantityUnit == "" {
		return errors.New("diet: quantity_unit is required")
	}
	return nil
}

// DietSubstitution represents a food substitution option for a meal item.
type DietSubstitution struct {
	ID                   uuid.UUID
	TenantID             uuid.UUID
	DietMealItemID       uuid.UUID
	SubstituteFoodItemID uuid.UUID
	QuantityValue        *float64
	QuantityUnit         *string
	Notes                string
	SortOrder            int
}

// DietPublication represents a publication record linking a published diet to a patient.
type DietPublication struct {
	ID                uuid.UUID
	TenantID          uuid.UUID
	DietID            uuid.UUID
	PatientID         uuid.UUID
	PublishedByUserID uuid.UUID
	PublishedAt       time.Time
	Status            string
}
