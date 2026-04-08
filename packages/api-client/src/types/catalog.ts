import type { UUID } from "./common"

export type FoodSource = "system" | "taco" | "ibge" | "tenant"

export interface NutritionFacts {
  calories_kcal?: number
  protein_g?: number
  carbs_g?: number
  fiber_g?: number
  sugar_g?: number
  total_fat_g?: number
  saturated_fat_g?: number
  trans_fat_g?: number
  cholesterol_mg?: number
  sodium_mg?: number
  potassium_mg?: number
  calcium_mg?: number
  iron_mg?: number
}

export interface HouseholdMeasure {
  id: UUID
  label: string
  grams: number
}

export interface FoodItem {
  id: UUID
  tenant_id?: string
  name: string
  food_group: string
  brand?: string
  barcode?: string
  serving_size_g: number
  serving_label: string
  source: FoodSource
  active: boolean
  nutrition_facts?: NutritionFacts
  household_measures?: HouseholdMeasure[]
  created_at: string
  updated_at: string
}

export interface CreateFoodItemRequest {
  name: string
  food_group: string
  brand?: string
  barcode?: string
  serving_size_g: number
  serving_label?: string
  nutrition_facts?: NutritionFacts
}
