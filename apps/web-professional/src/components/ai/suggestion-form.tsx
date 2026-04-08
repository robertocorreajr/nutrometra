"use client"

import { useState } from "react"
import { Sparkles, AlertTriangle } from "lucide-react"
import { Button, Card, CardContent } from "@nutrometra/ui"
import type { SuggestionType } from "@nutrometra/api-client"

interface SuggestionFormProps {
  patientId: string
  onSubmit: (type: SuggestionType, context?: string) => void
  isSubmitting: boolean
  onCancel?: () => void
}

const suggestionTypeOptions: { label: string; value: SuggestionType }[] = [
  { label: "Rascunho de Dieta", value: "diet_draft" },
  { label: "Estrutura de Refeições", value: "meal_structure" },
  { label: "Substituições Alimentares", value: "substitutions" },
  { label: "Resumo Clínico", value: "clinical_summary" },
  { label: "Checklist de Revisão", value: "review_checklist" },
]

export function SuggestionForm({
  patientId,
  onSubmit,
  isSubmitting,
  onCancel,
}: SuggestionFormProps) {
  const [selectedType, setSelectedType] = useState<SuggestionType>("diet_draft")
  const [extraContext, setExtraContext] = useState("")

  function handleFormSubmit(e: React.FormEvent) {
    e.preventDefault()
    onSubmit(selectedType, extraContext || undefined)
  }

  return (
    <Card>
      <CardContent className="pt-6">
        {/* AI disclaimer */}
        <div className="mb-6 flex items-start gap-3 rounded-md border border-yellow-300 bg-yellow-50 p-4">
          <AlertTriangle className="h-5 w-5 shrink-0 text-yellow-600 mt-0.5" />
          <p className="text-sm text-yellow-800">
            Sugestões geradas por IA devem ser revisadas pelo profissional antes de qualquer uso clínico.
          </p>
        </div>

        <form onSubmit={handleFormSubmit} className="space-y-4">
          <div className="space-y-2">
            <label htmlFor="suggestion-type" className="text-sm font-medium">
              Tipo de Sugestão *
            </label>
            <select
              id="suggestion-type"
              value={selectedType}
              onChange={(e) => setSelectedType(e.target.value as SuggestionType)}
              className="flex w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
            >
              {suggestionTypeOptions.map((opt) => (
                <option key={opt.value} value={opt.value}>
                  {opt.label}
                </option>
              ))}
            </select>
          </div>

          <div className="space-y-2">
            <label htmlFor="extra-context" className="text-sm font-medium">
              Contexto adicional (opcional)
            </label>
            <textarea
              id="extra-context"
              value={extraContext}
              onChange={(e) => setExtraContext(e.target.value)}
              rows={4}
              placeholder="Descreva informações adicionais que possam ajudar na sugestão..."
              className="flex w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
            />
          </div>

          <div className="flex gap-2 justify-end">
            {onCancel && (
              <Button type="button" variant="ghost" onClick={onCancel}>
                Cancelar
              </Button>
            )}
            <Button type="submit" disabled={isSubmitting}>
              <Sparkles className="h-4 w-4 mr-2" />
              {isSubmitting ? "Gerando..." : "Gerar Sugestão"}
            </Button>
          </div>
        </form>
      </CardContent>
    </Card>
  )
}
