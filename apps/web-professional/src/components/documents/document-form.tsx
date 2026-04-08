"use client"

import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { Button, Input, Card, CardContent } from "@nutrometra/ui"
import { useCreateDocument, useProfessionalMe } from "@nutrometra/api-client/hooks"
import {
  createDocumentSchema,
  type CreateDocumentFormValues,
  documentTypeOptions,
} from "@/lib/schemas/document"
import type { ClinicalDocument } from "@nutrometra/api-client"

interface DocumentFormProps {
  patientId: string
  onSuccess: (doc: ClinicalDocument) => void
  onCancel: () => void
}

export function DocumentForm({ patientId, onSuccess, onCancel }: DocumentFormProps) {
  const { data: professional } = useProfessionalMe()
  const createDocument = useCreateDocument(patientId)

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<CreateDocumentFormValues>({
    resolver: zodResolver(createDocumentSchema),
    defaultValues: {
      document_type: "exam_request",
      title: "",
    },
  })

  async function handleFormSubmit(data: CreateDocumentFormValues) {
    if (!professional?.id) return

    try {
      const doc = await createDocument.mutateAsync({
        professional_id: professional.id,
        document_type: data.document_type,
        title: data.title,
        content_json: { body: "" },
      })
      onSuccess(doc)
    } catch {
      // Error handled by mutation state
    }
  }

  return (
    <Card>
      <CardContent className="pt-6">
        <form onSubmit={handleSubmit(handleFormSubmit)} className="space-y-4">
          <div className="space-y-2">
            <label htmlFor="document_type" className="text-sm font-medium">
              Tipo de documento *
            </label>
            <select
              id="document_type"
              {...register("document_type")}
              className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
            >
              {documentTypeOptions.map((opt) => (
                <option key={opt.value} value={opt.value}>
                  {opt.label}
                </option>
              ))}
            </select>
            {errors.document_type && (
              <p className="text-sm text-destructive">{errors.document_type.message}</p>
            )}
          </div>

          <div className="space-y-2">
            <label htmlFor="title" className="text-sm font-medium">
              Título *
            </label>
            <Input
              id="title"
              {...register("title")}
              placeholder="Ex: Solicitação de hemograma completo"
            />
            {errors.title && (
              <p className="text-sm text-destructive">{errors.title.message}</p>
            )}
          </div>

          {createDocument.isError && (
            <p className="text-sm text-destructive">
              Erro ao criar documento. Tente novamente.
            </p>
          )}

          <div className="flex gap-2 justify-end">
            <Button type="button" variant="ghost" onClick={onCancel}>
              Cancelar
            </Button>
            <Button type="submit" disabled={createDocument.isPending}>
              {createDocument.isPending ? "Criando..." : "Criar Documento"}
            </Button>
          </div>
        </form>
      </CardContent>
    </Card>
  )
}
