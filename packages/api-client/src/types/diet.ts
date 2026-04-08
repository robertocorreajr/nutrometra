import type { UUID } from "./common"

export type DietStatus = "draft" | "published" | "archived"

export interface DietSubstitution {
  id: UUID
  diet_meal_item_id: UUID
  substitute_food_item_id: UUID
  quantity_value?: number
  quantity_unit?: string
  notes?: string
  sort_order: number
}

export interface DietMealItem {
  id: UUID
  diet_meal_id: UUID
  food_item_id: UUID
  quantity_value: number
  quantity_unit: string
  household_measure_id?: string
  amount_description?: string
  preparation_notes?: string
  sort_order: number
  substitutions?: DietSubstitution[]
}

export interface DietMeal {
  id: UUID
  diet_id: UUID
  meal_name: string
  meal_order: number
  notes?: string
  items?: DietMealItem[]
}

export interface Diet {
  id: UUID
  patient_id: UUID
  professional_id: UUID
  title: string
  objective?: string
  status: DietStatus
  version_number: number
  previous_version_id?: string
  published_at?: string
  valid_from?: string
  valid_until?: string
  created_at: string
  updated_at: string
  meals?: DietMeal[]
}

export interface CreateDietRequest {
  professional_id: string
  title: string
  objective?: string
  valid_from?: string
  valid_until?: string
}

export interface UpdateDietRequest {
  title: string
  objective?: string
  valid_from?: string
  valid_until?: string
}

export interface CreateMealRequest {
  meal_name: string
  meal_order: number
  notes?: string
}

export interface CreateMealItemRequest {
  food_item_id: string
  quantity_value: number
  quantity_unit: string
  household_measure_id?: string
  amount_description?: string
  preparation_notes?: string
  sort_order: number
}

export interface CreateSubstitutionRequest {
  substitute_food_item_id: string
  quantity_value?: number
  quantity_unit?: string
  notes?: string
  sort_order: number
}
