import { z } from "zod"
import type { DietStatus } from "@nutrometra/api-client"

export const createDietSchema = z.object({
  title: z.string().min(2, "Título deve ter pelo menos 2 caracteres"),
  objective: z.string().optional().default(""),
  valid_from: z.string().optional().default(""),
  valid_until: z.string().optional().default(""),
})

export type CreateDietFormValues = z.infer<typeof createDietSchema>

export const mealSchema = z.object({
  meal_name: z.string().min(1, "Nome da refeição é obrigatório"),
  meal_order: z.coerce.number().min(1, "Ordem deve ser pelo menos 1"),
  notes: z.string().optional().default(""),
})

export type MealFormValues = z.infer<typeof mealSchema>

export const mealItemSchema = z.object({
  food_item_id: z.string().min(1, "Alimento é obrigatório"),
  quantity_value: z.coerce.number().positive("Quantidade deve ser maior que 0"),
  quantity_unit: z.string().min(1, "Unidade é obrigatória"),
  sort_order: z.coerce.number().min(0),
  preparation_notes: z.string().optional().default(""),
})

export type MealItemFormValues = z.infer<typeof mealItemSchema>

export const dietStatusLabels: Record<DietStatus, string> = {
  draft: "Rascunho",
  published: "Publicada",
  archived: "Arquivada",
}

export const dietStatusColors: Record<DietStatus, string> = {
  draft: "bg-blue-100 text-blue-700 border-blue-200",
  published: "bg-green-100 text-green-700 border-green-200",
  archived: "bg-gray-100 text-gray-600 border-gray-200",
}
