import { z } from "zod"

export const anamnesisSchema = z.object({
  chief_complaint: z.string().optional().default(""),
  history_present_illness: z.string().optional().default(""),
  past_medical_history: z.string().optional().default(""),
  family_history: z.string().optional().default(""),
  social_history: z.string().optional().default(""),
  dietary_history: z.string().optional().default(""),
  physical_activity: z.string().optional().default(""),
  sleep_pattern: z.string().optional().default(""),
  bowel_habits: z.string().optional().default(""),
  water_intake: z.string().optional().default(""),
  supplements: z.string().optional().default(""),
  observations: z.string().optional().default(""),
})

export type AnamnesisFormValues = z.infer<typeof anamnesisSchema>

export const stepFields: Record<number, (keyof AnamnesisFormValues)[]> = {
  0: ["chief_complaint"],
  1: ["history_present_illness", "past_medical_history", "family_history"],
  2: ["social_history", "dietary_history", "physical_activity", "sleep_pattern", "bowel_habits", "water_intake"],
  3: ["supplements"],
  4: ["observations"],
}
