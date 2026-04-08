"use client"

import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { Button, Input, Card, CardContent } from "@nutrometra/ui"
import { noteSchema, type NoteFormValues } from "@/lib/schemas/note"

interface NoteFormProps {
  onSubmit: (data: NoteFormValues) => void
  isSubmitting: boolean
  onCancel?: () => void
}

export function NoteForm({ onSubmit, isSubmitting, onCancel }: NoteFormProps) {
  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<NoteFormValues>({
    resolver: zodResolver(noteSchema),
    defaultValues: {
      title: "",
      content: "",
      visible_to_patient: false,
    },
  })

  function handleFormSubmit(data: NoteFormValues) {
    onSubmit(data)
    reset()
  }

  return (
    <Card>
      <CardContent className="pt-6">
        <form onSubmit={handleSubmit(handleFormSubmit)} className="space-y-4">
          <div className="space-y-2">
            <label htmlFor="title" className="text-sm font-medium">Título *</label>
            <Input id="title" {...register("title")} placeholder="Ex: Consulta de retorno, Reavaliação..." />
            {errors.title && <p className="text-sm text-destructive">{errors.title.message}</p>}
          </div>

          <div className="space-y-2">
            <label htmlFor="content" className="text-sm font-medium">Conteúdo *</label>
            <textarea
              id="content"
              {...register("content")}
              rows={6}
              placeholder="Descreva a evolução do paciente, condutas, orientações..."
              className="flex w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
            />
            {errors.content && <p className="text-sm text-destructive">{errors.content.message}</p>}
          </div>

          <div className="flex items-center gap-2">
            <input type="checkbox" id="visible_to_patient" {...register("visible_to_patient")} className="h-4 w-4 rounded border-input" />
            <label htmlFor="visible_to_patient" className="text-sm">Visível para o paciente</label>
          </div>

          <div className="flex gap-2 justify-end">
            {onCancel && <Button type="button" variant="ghost" onClick={onCancel}>Cancelar</Button>}
            <Button type="submit" disabled={isSubmitting}>{isSubmitting ? "Salvando..." : "Salvar Evolução"}</Button>
          </div>
        </form>
      </CardContent>
    </Card>
  )
}
