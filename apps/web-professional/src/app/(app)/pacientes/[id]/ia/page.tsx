"use client"

import { useParams } from "next/navigation"
import { useState } from "react"
import { LoadingState, ErrorState, EmptyState, Button } from "@nutrometra/ui"
import { useAISuggestions, useCreateAISuggestion, useProfessionalMe } from "@nutrometra/api-client/hooks"
import { SuggestionForm } from "@/components/ai/suggestion-form"
import { SuggestionList } from "@/components/ai/suggestion-list"
import { SuggestionDetail } from "@/components/ai/suggestion-detail"
import type { AISuggestion, SuggestionType } from "@nutrometra/api-client"
import { Plus, Sparkles } from "lucide-react"

export default function IAPage() {
  const params = useParams()
  const patientId = params.id as string
  const { data: professional } = useProfessionalMe()

  const { data: suggestionsData, isLoading, isError, refetch } = useAISuggestions()
  const createSuggestion = useCreateAISuggestion()

  const [showForm, setShowForm] = useState(false)
  const [selectedSuggestion, setSelectedSuggestion] = useState<AISuggestion | null>(null)

  async function handleSubmit(type: SuggestionType, context?: string) {
    try {
      await createSuggestion.mutateAsync({
        suggestion_type: type,
        patient_id: patientId,
        extra_context: context ? { text: context } : undefined,
      })
      setShowForm(false)
    } catch {
      // Error handled by mutation
    }
  }

  if (isLoading) return <LoadingState lines={6} />
  if (isError) return <ErrorState message="Erro ao carregar sugestões de IA." onRetry={refetch} />

  const suggestions = suggestionsData?.items ?? []

  return (
    <div className="space-y-4">
      {showForm ? (
        <SuggestionForm
          patientId={patientId}
          onSubmit={handleSubmit}
          isSubmitting={createSuggestion.isPending}
          onCancel={() => setShowForm(false)}
        />
      ) : (
        <div className="flex justify-end">
          <Button onClick={() => setShowForm(true)}>
            <Plus className="h-4 w-4 mr-2" />
            Nova Sugestão
          </Button>
        </div>
      )}

      {createSuggestion.isError && (
        <p className="text-sm text-destructive">Erro ao solicitar sugestão. Tente novamente.</p>
      )}

      {suggestions.length === 0 && !showForm ? (
        <EmptyState
          icon={<Sparkles className="h-12 w-12" />}
          title="Nenhuma sugestão"
          description="Solicite uma sugestão da IA para auxiliar no atendimento deste paciente."
          action={
            <Button onClick={() => setShowForm(true)}>
              <Plus className="h-4 w-4 mr-2" />
              Nova Sugestão
            </Button>
          }
        />
      ) : (
        <SuggestionList
          suggestions={suggestions}
          onSelect={(s) => setSelectedSuggestion(s)}
        />
      )}

      <SuggestionDetail
        suggestion={selectedSuggestion}
        onClose={() => setSelectedSuggestion(null)}
      />
    </div>
  )
}
