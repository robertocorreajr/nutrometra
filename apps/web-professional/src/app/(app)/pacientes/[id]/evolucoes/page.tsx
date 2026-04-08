"use client"

import { useParams } from "next/navigation"
import { useState } from "react"
import { LoadingState, ErrorState, EmptyState, Button } from "@nutrometra/ui"
import { useProgressNotes, useCreateProgressNote, useProfessionalMe } from "@nutrometra/api-client/hooks"
import { NoteForm } from "@/components/notes/note-form"
import { NoteList } from "@/components/notes/note-list"
import type { NoteFormValues } from "@/lib/schemas/note"
import { Plus, FileText } from "lucide-react"

export default function EvolucoesPage() {
  const params = useParams()
  const patientId = params.id as string
  const { data: professional } = useProfessionalMe()

  const { data: notes, isLoading, isError, refetch } = useProgressNotes(patientId)
  const createNote = useCreateProgressNote(patientId)

  const [showForm, setShowForm] = useState(false)

  async function handleSubmit(data: NoteFormValues) {
    try {
      await createNote.mutateAsync({
        professional_id: professional?.id ?? "",
        title: data.title,
        content: data.content,
        visible_to_patient: data.visible_to_patient,
      })
      setShowForm(false)
    } catch {
      // Error handled by mutation
    }
  }

  if (isLoading) return <LoadingState lines={6} />
  if (isError) return <ErrorState message="Erro ao carregar evoluções." onRetry={refetch} />

  const noteList = notes ?? []

  return (
    <div className="space-y-4">
      {showForm ? (
        <NoteForm onSubmit={handleSubmit} isSubmitting={createNote.isPending} onCancel={() => setShowForm(false)} />
      ) : (
        <div className="flex justify-end">
          <Button onClick={() => setShowForm(true)}>
            <Plus className="h-4 w-4 mr-2" />
            Nova Evolução
          </Button>
        </div>
      )}

      {createNote.isError && (
        <p className="text-sm text-destructive">Erro ao salvar evolução. Tente novamente.</p>
      )}

      {noteList.length === 0 && !showForm ? (
        <EmptyState
          icon={<FileText className="h-12 w-12" />}
          title="Nenhuma evolução"
          description="Registre a primeira evolução para este paciente."
          action={
            <Button onClick={() => setShowForm(true)}>
              <Plus className="h-4 w-4 mr-2" />
              Nova Evolução
            </Button>
          }
        />
      ) : (
        <NoteList notes={noteList} />
      )}
    </div>
  )
}
