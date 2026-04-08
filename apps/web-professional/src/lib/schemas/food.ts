import { z } from "zod"

export const foodSearchSchema = z.object({
  query: z.string().optional().default(""),
  group: z.string().optional().default(""),
})

export type FoodSearchValues = z.infer<typeof foodSearchSchema>

export const createFoodSchema = z.object({
  name: z.string().min(2, "Nome é obrigatório"),
  food_group: z.string().min(1, "Grupo alimentar é obrigatório"),
  serving_size_g: z.coerce.number().min(0.1, "Porção deve ser maior que 0"),
  serving_label: z.string().optional().default("100g"),
  brand: z.string().optional().default(""),
  barcode: z.string().optional().default(""),
})

export type CreateFoodValues = z.infer<typeof createFoodSchema>
