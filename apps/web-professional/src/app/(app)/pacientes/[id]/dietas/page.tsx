"use client"

import { useState } from "react"
import { useParams, useRouter } from "next/navigation"
import { Button, LoadingState, ErrorState, EmptyState, Card, CardContent } from "@nutrometra/ui"
import { usePatientDiets } from "@nutrometra/api-client/hooks"
import { CreateDietForm } from "@/components/diet/create-diet-form"
import { dietStatusLabels, dietStatusColors } from "@/lib/schemas/diet"
import { Plus, UtensilsCrossed } from "lucide-react"
import type { Diet } from "@nutrometra/api-client"

function formatDate(dateStr?: string): string {
  if (!dateStr) return ""
  try {
    return new Date(dateStr).toLocaleDateString("pt-BR")
  } catch {
    return dateStr
  }
}

function DietCard({ diet, onClick }: { diet: Diet; onClick: () => void }) {
  return (
    <Card className="cursor-pointer hover:shadow-md transition-shadow" onClick={onClick}>
      <CardContent className="p-4">
        <div className="flex items-start justify-between gap-2">
          <div className="min-w-0 flex-1">
            <h3 className="font-medium text-sm truncate">{diet.title}</h3>
            {diet.objective && (
              <p className="text-xs text-muted-foreground mt-1 line-clamp-2">
                {diet.objective}
              </p>
            )}
          </div>
          <div className="flex flex-col items-end gap-1 shrink-0">
            <span
              className={`inline-flex items-center rounded-full border px-2 py-0.5 text-xs font-medium ${dietStatusColors[diet.status]}`}
            >
              {dietStatusLabels[diet.status]}
            </span>
            <span className="text-xs text-muted-foreground">v{diet.version_number}</span>
          </div>
        </div>
        {(diet.valid_from || diet.valid_until) && (
          <p className="text-xs text-muted-foreground mt-2">
            {diet.valid_from && `De ${formatDate(diet.valid_from)}`}
            {diet.valid_from && diet.valid_until && " "}
            {diet.valid_until && `ate ${formatDate(diet.valid_until)}`}
          </p>
        )}
      </CardContent>
    </Card>
  )
}

export default function DietasPage() {
  const params = useParams()
  const router = useRouter()
  const patientId = params.id as string

  const { data: diets, isLoading, isError, refetch } = usePatientDiets(patientId)
  const [showForm, setShowForm] = useState(false)

  function handleDietCreated(diet: Diet) {
    setShowForm(false)
    router.push(`/pacientes/${patientId}/dietas/${diet.id}`)
  }

  if (isLoading) return <LoadingState lines={4} />
  if (isError) {
    return (
      <ErrorState
        message="Erro ao carregar dietas do paciente."
        onRetry={refetch}
      />
    )
  }

  const dietList = diets ?? []

  return (
    <div className="space-y-4">
      {showForm ? (
        <CreateDietForm
          patientId={patientId}
          onSuccess={handleDietCreated}
          onCancel={() => setShowForm(false)}
        />
      ) : (
        <div className="flex justify-end">
          <Button onClick={() => setShowForm(true)}>
            <Plus className="h-4 w-4 mr-2" />
            Nova Dieta
          </Button>
        </div>
      )}

      {dietList.length === 0 && !showForm ? (
        <EmptyState
          icon={<UtensilsCrossed className="h-12 w-12" />}
          title="Nenhuma dieta"
          description="Crie a primeira dieta para este paciente."
          action={
            <Button onClick={() => setShowForm(true)}>
              <Plus className="h-4 w-4 mr-2" />
              Nova Dieta
            </Button>
          }
        />
      ) : (
        <div className="grid gap-3">
          {dietList.map((diet) => (
            <DietCard
              key={diet.id}
              diet={diet}
              onClick={() => router.push(`/pacientes/${patientId}/dietas/${diet.id}`)}
            />
          ))}
        </div>
      )}
    </div>
  )
}
