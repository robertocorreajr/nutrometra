import { z } from "zod"

export const profileSchema = z.object({
  full_name: z.string().min(2, "Nome completo é obrigatório"),
  registration_type: z.string().min(1, "Tipo de registro é obrigatório"),
  registration_number: z.string().min(1, "Número de registro é obrigatório"),
  registration_state: z.string().optional().default(""),
  specialty: z.string().optional().default(""),
  bio: z.string().optional().default(""),
  phone: z.string().optional().default(""),
})

export type ProfileFormValues = z.infer<typeof profileSchema>

export const registrationTypes = [
  { value: "CRN", label: "CRN (Nutricionista)" },
  { value: "CRM", label: "CRM (Médico)" },
  { value: "other", label: "Outro" },
] as const

export const brazilianStates = [
  "AC","AL","AP","AM","BA","CE","DF","ES","GO","MA","MT","MS","MG",
  "PA","PB","PR","PE","PI","RJ","RN","RS","RO","RR","SC","SP","SE","TO",
] as const

export const addressSchema = z.object({
  label: z.string().min(1, "Nome do endereço é obrigatório"),
  street: z.string().min(1, "Rua é obrigatória"),
  number: z.string().optional().default(""),
  complement: z.string().optional().default(""),
  neighborhood: z.string().optional().default(""),
  city: z.string().min(1, "Cidade é obrigatória"),
  state: z.string().min(2, "Estado é obrigatório"),
  zip_code: z.string().min(8, "CEP é obrigatório"),
  country: z.string().optional().default("BR"),
  phone: z.string().optional().default(""),
  notes: z.string().optional().default(""),
})

export type AddressFormValues = z.infer<typeof addressSchema>

export const serviceModeSchema = z.object({
  mode: z.enum(["onsite", "online", "home_visit"]),
  address_id: z.string().optional(),
  duration_min: z.coerce.number().min(1, "Duração mínima é 1 minuto"),
})

export type ServiceModeFormValues = z.infer<typeof serviceModeSchema>

export const serviceModeLabels: Record<string, string> = {
  onsite: "Presencial",
  online: "Online",
  home_visit: "Domiciliar",
}
