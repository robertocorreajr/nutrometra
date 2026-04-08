import { z } from "zod"

export const serviceModeOptions = [
  { value: "onsite", label: "Presencial" },
  { value: "online", label: "Online" },
  { value: "home_visit", label: "Domiciliar" },
] as const

export const appointmentStatusLabels: Record<string, string> = {
  scheduled: "Agendada",
  confirmed: "Confirmada",
  completed: "Concluída",
  cancelled: "Cancelada",
  no_show: "Não Compareceu",
}

export const appointmentStatusColors: Record<string, string> = {
  scheduled: "bg-blue-100 text-blue-800 border-blue-300",
  confirmed: "bg-green-100 text-green-800 border-green-300",
  completed: "bg-gray-100 text-gray-600 border-gray-300",
  cancelled: "bg-red-100 text-red-800 border-red-300",
  no_show: "bg-orange-100 text-orange-800 border-orange-300",
}

export const appointmentSchema = z.object({
  patient_id: z.string().min(1, "Paciente é obrigatório"),
  start_at: z.string().min(1, "Data/hora de início é obrigatória"),
  end_at: z.string().min(1, "Data/hora de término é obrigatória"),
  service_mode: z.enum(["onsite", "online", "home_visit"]),
  address_id: z.string().optional(),
  notes: z.string().optional().default(""),
})

export type AppointmentFormValues = z.infer<typeof appointmentSchema>
