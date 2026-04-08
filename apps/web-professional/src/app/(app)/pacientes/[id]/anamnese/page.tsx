"use client"

import { useParams } from "next/navigation"
import { LoadingState, ErrorState, EmptyState, Button, Card, CardContent, CardHeader, CardTitle } from "@nutrometra/ui"
import { useAnamneses } from "@nutrometra/api-client/hooks"
import type { Anamnesis } from "@nutrometra/api-client"
import { AnamnesisWizard } from "@/components/anamnesis/anamnesis-wizard"
import { Plus, FileText } from "lucide-react"
import { format } from "date-fns"
import { ptBR } from "date-fns/locale"
import { useState } from "react"

export default function AnamnesePage() {
  const params = useParams()
  const patientId = params.id as string
  const { data: anamneses, isLoading, isError, refetch } = useAnamneses(patientId)
  const [showWizard, setShowWizard] = useState(false)
  const [editingAnamnesis, setEditingAnamnesis] = useState<Anamnesis | null>(null)

  if (isLoading) return <LoadingState lines={6} />
  if (isError) return <ErrorState message="Erro ao carregar anamneses." onRetry={refetch} />

  const list = anamneses ?? []

  if (showWizard || editingAnamnesis) {
    return (
      <div>
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-lg font-semibold">
            {editingAnamnesis
              ? editingAnamnesis.status === "finalized"
                ? "Visualizar Anamnese"
                : "Editar Anamnese"
              : "Nova Anamnese"}
          </h2>
          <Button variant="ghost" size="sm" onClick={() => { setShowWizard(false); setEditingAnamnesis(null) }}>
            Cancelar
          </Button>
        </div>
        <AnamnesisWizard
          patientId={patientId}
          existing={editingAnamnesis}
          onComplete={() => { setShowWizard(false); setEditingAnamnesis(null); refetch() }}
        />
      </div>
    )
  }

  if (list.length === 0) {
    return (
      <EmptyState
        icon={<FileText className="h-12 w-12" />}
        title="Nenhuma anamnese"
        description="Crie a primeira anamnese para este paciente."
        action={
          <Button onClick={() => setShowWizard(true)}>
            <Plus className="h-4 w-4 mr-2" />
            Nova Anamnese
          </Button>
        }
      />
    )
  }

  return (
    <div className="space-y-4">
      <div className="flex justify-end">
        <Button onClick={() => setShowWizard(true)}>
          <Plus className="h-4 w-4 mr-2" />
          Nova Anamnese
        </Button>
      </div>

      <div className="space-y-3">
        {list.map((anamnesis) => (
          <Card
            key={anamnesis.id}
            className="cursor-pointer hover:bg-accent/50 transition-colors"
            onClick={() => setEditingAnamnesis(anamnesis)}
          >
            <CardHeader className="pb-2">
              <div className="flex items-center justify-between">
                <CardTitle className="text-sm font-medium">
                  Anamnese — {format(new Date(anamnesis.created_at), "dd/MM/yyyy", { locale: ptBR })}
                </CardTitle>
                <span className={`inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium ${anamnesis.status === "finalized" ? "bg-blue-100 text-blue-800" : "bg-yellow-100 text-yellow-800"}`}>
                  {anamnesis.status === "finalized" ? "Finalizada" : "Rascunho"}
                </span>
              </div>
            </CardHeader>
            <CardContent>
              <p className="text-sm text-muted-foreground line-clamp-2">{anamnesis.chief_complaint || "Sem queixa principal registrada."}</p>
              <p className="text-xs text-muted-foreground mt-2">
                {anamnesis.status === "draft" ? "Clique para continuar editando" : "Clique para visualizar"}
              </p>
            </CardContent>
          </Card>
        ))}
      </div>
    </div>
  )
}
