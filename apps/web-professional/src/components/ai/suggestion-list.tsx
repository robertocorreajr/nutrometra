"use client"

import { format } from "date-fns"
import { ptBR } from "date-fns/locale"
import { Card, CardContent } from "@nutrometra/ui"
import type { AISuggestion, SuggestionType, SuggestionStatus } from "@nutrometra/api-client"

interface SuggestionListProps {
  suggestions: AISuggestion[]
  onSelect: (suggestion: AISuggestion) => void
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

export function SuggestionList({ suggestions, onSelect }: SuggestionListProps) {
  const sorted = [...suggestions].sort(
    (a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
  )

  return (
    <div className="space-y-3">
      {sorted.map((suggestion) => (
        <Card
          key={suggestion.id}
          className="cursor-pointer transition-colors hover:bg-muted/50"
          onClick={() => onSelect(suggestion)}
        >
          <CardContent className="py-4">
            <div className="flex items-center justify-between gap-3">
              <div className="min-w-0">
                <p className="font-medium truncate">
                  {suggestionTypeLabels[suggestion.suggestion_type] ?? suggestion.suggestion_type}
                </p>
                <p className="text-sm text-muted-foreground mt-1">
                  {format(new Date(suggestion.created_at), "dd 'de' MMM 'de' yyyy, HH:mm", {
                    locale: ptBR,
                  })}
                </p>
              </div>
              <span
                className={`inline-flex items-center shrink-0 rounded-full px-2.5 py-1 text-xs font-medium ${statusBadgeStyles[suggestion.status]}`}
              >
                {statusLabels[suggestion.status]}
                {suggestion.status === "generating" && (
                  <span className="ml-1 animate-pulse">...</span>
                )}
              </span>
            </div>
          </CardContent>
        </Card>
      ))}
    </div>
  )
}
