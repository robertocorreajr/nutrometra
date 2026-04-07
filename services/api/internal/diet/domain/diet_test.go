package domain

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDiet_Validate(t *testing.T) {
	tests := []struct {
		name    string
		diet    Diet
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid diet with title",
			diet:    Diet{Title: "Plano alimentar semanal"},
			wantErr: false,
		},
		{
			name:    "missing title fails",
			diet:    Diet{},
			wantErr: true,
			errMsg:  "title is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.diet.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDietMeal_Validate(t *testing.T) {
	tests := []struct {
		name    string
		meal    DietMeal
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid meal",
			meal:    DietMeal{MealName: "Cafe da manha", MealOrder: 1},
			wantErr: false,
		},
		{
			name:    "missing meal_name fails",
			meal:    DietMeal{MealOrder: 1},
			wantErr: true,
			errMsg:  "meal_name is required",
		},
		{
			name:    "zero meal_order fails",
			meal:    DietMeal{MealName: "Almoco", MealOrder: 0},
			wantErr: true,
			errMsg:  "meal_order must be greater than zero",
		},
		{
			name:    "negative meal_order fails",
			meal:    DietMeal{MealName: "Almoco", MealOrder: -1},
			wantErr: true,
			errMsg:  "meal_order must be greater than zero",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.meal.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDietMealItem_Validate(t *testing.T) {
	tests := []struct {
		name    string
		item    DietMealItem
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid item",
			item:    DietMealItem{QuantityValue: 100, QuantityUnit: "g"},
			wantErr: false,
		},
		{
			name:    "zero quantity fails",
			item:    DietMealItem{QuantityValue: 0, QuantityUnit: "g"},
			wantErr: true,
			errMsg:  "quantity_value must be greater than zero",
		},
		{
			name:    "negative quantity fails",
			item:    DietMealItem{QuantityValue: -5, QuantityUnit: "g"},
			wantErr: true,
			errMsg:  "quantity_value must be greater than zero",
		},
		{
			name:    "missing unit fails",
			item:    DietMealItem{QuantityValue: 100},
			wantErr: true,
			errMsg:  "quantity_unit is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.item.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDiet_HasContent(t *testing.T) {
	tests := []struct {
		name string
		diet Diet
		want bool
	}{
		{
			name: "diet with meal and item has content",
			diet: Diet{
				Meals: []DietMeal{
					{
						Items: []DietMealItem{
							{ID: uuid.New()},
						},
					},
				},
			},
			want: true,
		},
		{
			name: "diet with empty meal has no content",
			diet: Diet{
				Meals: []DietMeal{
					{Items: nil},
				},
			},
			want: false,
		},
		{
			name: "diet with no meals has no content",
			diet: Diet{Meals: nil},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.diet.HasContent())
		})
	}
}
