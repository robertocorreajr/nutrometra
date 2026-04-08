"use client"

import { CheckCircle, XCircle, X, Sparkles } from "lucide-react"
import { Button, Card, CardContent, CardHeader, CardTitle } from "@nutrometra/ui"
import {
  useAISuggestion,
  useAcceptAISuggestion,
  useRejectAISuggestion,
} from "@nutrometra/api-client/hooks"
import type { AISuggestion, SuggestionType, SuggestionStatus } from "@nutrometra/api-client"

interface SuggestionDetailProps {
  suggestion: AISuggestion | null
  onClose: () => void
}

const suggestionTypeLabels: Record<SuggestionType, string> = {
  diet_draft: "Rascunho de Dieta",
  meal_structure: "Estrutura de Refeições",
  substitutions: "Substituições Alimentares",
  clinical_summary: "Resumo Clínico",
  review_checklist: "Checklist de Revisão",
}

const statusBadgeStyles: Record<SuggestionStatus, string> = {
  pending: "bg-yellow-100 text-yellow-800",
  generating: "bg-blue-100 text-blue-800",
  completed: "bg-green-100 text-green-800",
  failed: "bg-red-100 text-red-800",
  accepted: "border border-green-500 text-green-700",
  rejected: "bg-gray-100 text-gray-600",
}

const statusLabels: Record<SuggestionStatus, string> = {
  pending: "Pendente",
  generating: "Gerando",
  completed: "Concluída",
  failed: "Falhou",
  accepted: "Aceita",
  rejected: "Rejeitada",
}

export function SuggestionDetail({ suggestion, onClose }: SuggestionDetailProps) {
  if (!suggestion) return null

  const isPolling = suggestion.status === "pending" || suggestion.status === "generating"

  const { data: liveSuggestion } = useAISuggestion(suggestion.id, {
    refetchInterval: isPolling ? 3000 : false,
  })

  const current = liveSuggestion ?? suggestion

  const acceptMutation = useAcceptAISuggestion(current.id)
  const rejectMutation = useRejectAISuggestion(current.id)

  const isCompleted = current.status === "completed"
  const isMutating = acceptMutation.isPending || rejectMutation.isPending

  async function handleAccept() {
    try {
      await acceptMutation.mutateAsync()
    } catch {
      // Error handled by mutation
    }
  }

  async function handleReject() {
    try {
      await rejectMutation.mutateAsync()
    } catch {
      // Error handled by mutation
    }
  }

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
      {/* Backdrop */}
      <div className="absolute inset-0 bg-black/50" onClick={onClose} />

      {/* Modal */}
      <Card className="relative z-10 w-full max-w-lg max-h-[90vh] overflow-y-auto">
        <CardHeader>
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-2">
              <Sparkles className="h-5 w-5 text-primary" />
              <CardTitle className="text-lg">
                {suggestionTypeLabels[current.suggestion_type] ?? current.suggestion_type}
              </CardTitle>
            </div>
            <button
              type="button"
              onClick={onClose}
              className="text-muted-foreground hover:text-foreground"
            >
              <X className="h-5 w-5" />
            </button>
          </div>

          {/* Status badge */}
          <div className="mt-2">
            <span
              className={`inline-flex items-center rounded-full px-2.5 py-1 text-xs font-medium ${statusBadgeStyles[current.status]}`}
            >
              {statusLabels[current.status]}
              {current.status === "generating" && (
                <span className="ml-1 animate-pulse">...</span>
              )}
            </span>
          </div>
        </CardHeader>

        <CardContent>
          {/* Response text */}
          {current.response_text ? (
            <div className="prose prose-sm max-w-none mb-6">
              <div className="whitespace-pre-wrap rounded-md bg-muted p-4 text-sm">
                {current.response_text}
              </div>
            </div>
          ) : (
            <div className="mb-6">
              {current.status === "pending" || current.status === "generating" ? (
                <div className="flex items-center gap-2 text-sm text-muted-foreground">
                  <span className="animate-spin">
                    <Sparkles className="h-4 w-4" />
                  </span>
                  Aguardando resposta da IA...
                </div>
              ) : current.status === "failed" ? (
                <p className="text-sm text-destructive">
                  A geração da sugestão falhou. Tente novamente.
                </p>
              ) : null}
            </div>
          )}

          {/* Metadata */}
          <div className="space-y-1 text-xs text-muted-foreground mb-6">
            {current.model_id && (
              <p>
                <span className="font-medium">Modelo:</span> {current.model_id}
              </p>
            )}
            {current.input_tokens != null && (
              <p>
                <span className="font-medium">Tokens entrada:</span> {current.input_tokens}
              </p>
            )}
            {current.output_tokens != null && (
              <p>
                <span className="font-medium">Tokens saída:</span> {current.output_tokens}
              </p>
            )}
          </div>

          {/* Error messages */}
          {(acceptMutation.isError || rejectMutation.isError) && (
            <p className="text-sm text-destructive mb-4">
              Erro ao atualizar sugestão. Tente novamente.
            </p>
          )}

          {/* Action buttons — only when completed */}
          {isCompleted && (
            <div className="flex gap-2 justify-end">
              <Button
                variant="outline"
                size="sm"
                disabled={isMutating}
                onClick={handleReject}
              >
                <XCircle className="h-4 w-4 mr-2" />
                {rejectMutation.isPending ? "Rejeitando..." : "Rejeitar"}
              </Button>
              <Button
                size="sm"
                disabled={isMutating}
                onClick={handleAccept}
              >
                <CheckCircle className="h-4 w-4 mr-2" />
                {acceptMutation.isPending ? "Aceitando..." : "Aceitar"}
              </Button>
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  )
}
