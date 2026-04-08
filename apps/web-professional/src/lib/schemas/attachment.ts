import { z } from "zod"

export const attachmentCategories = [
  { value: "general", label: "Geral" },
  { value: "exam", label: "Exame" },
  { value: "lab_result", label: "Resultado de Laboratório" },
  { value: "prescription", label: "Prescrição" },
  { value: "photo", label: "Foto" },
  { value: "other", label: "Outro" },
] as const

export const attachmentSchema = z.object({
  file_name: z.string().min(1, "Nome do arquivo é obrigatório"),
  file_type: z.string().min(1, "Tipo do arquivo é obrigatório"),
  file_size_bytes: z.number().min(1, "Tamanho deve ser maior que zero"),
  category: z.enum(["general", "exam", "lab_result", "prescription", "photo", "other"]),
  description: z.string().optional().default(""),
})

export type AttachmentFormValues = z.infer<typeof attachmentSchema>
