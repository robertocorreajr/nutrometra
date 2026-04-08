import { z } from "zod"

export const createDocumentSchema = z.object({
  document_type: z.enum(["exam_request", "prescription", "letter", "other"]),
  title: z.string().min(2, "Título é obrigatório"),
})

export type CreateDocumentFormValues = z.infer<typeof createDocumentSchema>

export const documentTypeLabels: Record<string, string> = {
  exam_request: "Solicitação de Exame",
  prescription: "Prescrição",
  letter: "Carta",
  other: "Outro",
}

export const documentTypeOptions = [
  { value: "exam_request", label: "Solicitação de Exame" },
  { value: "prescription", label: "Prescrição" },
  { value: "letter", label: "Carta" },
  { value: "other", label: "Outro" },
] as const

export const documentStatusLabels: Record<string, string> = {
  draft: "Rascunho",
  finalized: "Finalizado",
  published: "Publicado",
}

export const documentStatusColors: Record<string, string> = {
  draft: "bg-blue-100 text-blue-700 border-blue-200",
  finalized: "bg-amber-100 text-amber-700 border-amber-200",
  published: "bg-green-100 text-green-700 border-green-200",
}
