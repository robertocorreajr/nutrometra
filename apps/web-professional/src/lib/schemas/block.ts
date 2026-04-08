import { z } from "zod"

export const blockSchema = z.object({
  start_at: z.string().min(1, "Data de início é obrigatória"),
  end_at: z.string().min(1, "Data de término é obrigatória"),
  reason: z.string().optional().default(""),
  all_day: z.boolean().default(false),
})

export type BlockFormValues = z.infer<typeof blockSchema>
