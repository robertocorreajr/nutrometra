"use client"

import { useParams } from "next/navigation"
import { useState } from "react"
import { LoadingState, ErrorState, EmptyState, Button } from "@nutrometra/ui"
import { useAttachments, useCreateAttachment, useProfessionalMe } from "@nutrometra/api-client/hooks"
import { AttachmentForm } from "@/components/attachments/attachment-form"
import { AttachmentGrid } from "@/components/attachments/attachment-grid"
import type { AttachmentFormValues } from "@/lib/schemas/attachment"
import { Plus, Paperclip } from "lucide-react"

export default function AnexosPage() {
  const params = useParams()
  const patientId = params.id as string
  const { data: professional } = useProfessionalMe()

  const { data: attachments, isLoading, isError, refetch } = useAttachments(patientId)
  const createAttachment = useCreateAttachment(patientId)

  const [showForm, setShowForm] = useState(false)

  async function handleSubmit(data: AttachmentFormValues) {
    const storageKey = `uploads/${patientId}/${Date.now()}-${data.file_name}`

    try {
      await createAttachment.mutateAsync({
        professional_id: professional?.id ?? "",
        file_name: data.file_name,
        file_type: data.file_type,
        file_size_bytes: data.file_size_bytes,
        storage_key: storageKey,
        category: data.category,
        description: data.description ?? "",
      })
      setShowForm(false)
    } catch {
      // Error handled by mutation
    }
  }

  if (isLoading) return <LoadingState lines={6} />
  if (isError) return <ErrorState message="Erro ao carregar anexos." onRetry={refetch} />

  const attachmentList = attachments ?? []

  return (
    <div className="space-y-4">
      {showForm ? (
        <AttachmentForm onSubmit={handleSubmit} isSubmitting={createAttachment.isPending} onCancel={() => setShowForm(false)} />
      ) : (
        <div className="flex justify-end">
          <Button onClick={() => setShowForm(true)}>
            <Plus className="h-4 w-4 mr-2" />
            Novo Anexo
          </Button>
        </div>
      )}

      {createAttachment.isError && (
        <p className="text-sm text-destructive">Erro ao registrar anexo. Tente novamente.</p>
      )}

      {attachmentList.length === 0 && !showForm ? (
        <EmptyState
          icon={<Paperclip className="h-12 w-12" />}
          title="Nenhum anexo"
          description="Adicione exames, fotos ou outros documentos do paciente."
          action={
            <Button onClick={() => setShowForm(true)}>
              <Plus className="h-4 w-4 mr-2" />
              Novo Anexo
            </Button>
          }
        />
      ) : (
        <AttachmentGrid attachments={attachmentList} />
      )}
    </div>
  )
}
