import { z } from "zod"

export const patientSchema = z.object({
  full_name: z.string().min(2, "Nome deve ter pelo menos 2 caracteres"),
  email: z.string().email("E-mail inválido").or(z.literal("")),
  phone: z.string().optional().default(""),
  cpf: z
    .string()
    .regex(/^\d{11}$/, "CPF deve ter 11 dígitos")
    .or(z.literal(""))
    .optional()
    .default(""),
  date_of_birth: z.string().optional().default(""),
  gender: z
    .enum(["male", "female", "other", "prefer_not_to_say"])
    .optional()
    .default("prefer_not_to_say"),
  notes: z.string().optional().default(""),
})

export type PatientFormValues = z.infer<typeof patientSchema>

export const genderOptions = [
  { value: "male", label: "Masculino" },
  { value: "female", label: "Feminino" },
  { value: "other", label: "Outro" },
  { value: "prefer_not_to_say", label: "Prefiro não dizer" },
] as const
