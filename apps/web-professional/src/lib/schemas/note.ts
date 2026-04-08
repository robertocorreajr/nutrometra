import { z } from "zod"

export const noteSchema = z.object({
  title: z.string().min(1, "Título é obrigatório"),
  content: z.string().min(1, "Conteúdo é obrigatório"),
  visible_to_patient: z.boolean().default(false),
})

export type NoteFormValues = z.infer<typeof noteSchema>
