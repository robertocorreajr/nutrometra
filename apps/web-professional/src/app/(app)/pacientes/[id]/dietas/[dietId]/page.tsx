"use client"

import { useState } from "react"
import { useParams } from "next/navigation"
import Link from "next/link"
import { ArrowLeft, Plus } from "lucide-react"
import { Button, LoadingState, ErrorState } from "@nutrometra/ui"
import { useDiet } from "@nutrometra/api-client/hooks"
import { DietHeader } from "@/components/diet/diet-header"
import { DietMetadataForm } from "@/components/diet/diet-metadata-form"
import { DietStatusActions } from "@/components/diet/diet-status-actions"
import { MealCard } from "@/components/diet/meal-card"
import { MealForm } from "@/components/diet/meal-form"

export default function DietBuilderPage() {
  const params = useParams()
  const patientId = params.id as string
  const dietId = params.dietId as string

  const { data: diet, isLoading, isError, refetch } = useDiet(dietId)
  const [showAddMeal, setShowAddMeal] = useState(false)

  if (isLoading) return <LoadingState lines={6} />
  if (isError || !diet) {
    return (
      <ErrorState
        message="Erro ao carregar a dieta."
        onRetry={refetch}
      />
    )
  }

  const isDraft = diet.status === "draft"
  const meals = [...(diet.meals ?? [])].sort((a, b) => a.meal_order - b.meal_order)

  return (
    <div className="space-y-4">
      {/* Back link */}
      <Link
        href={`/pacientes/${patientId}/dietas`}
        className="inline-flex items-center gap-1 text-sm text-muted-foreground hover:text-foreground"
      >
        <ArrowLeft className="h-4 w-4" />
        Voltar para dietas
      </Link>

      {/* Diet header */}
      <DietHeader diet={diet} />

      {/* Metadata form (draft only) */}
      <DietMetadataForm diet={diet} patientId={patientId} />

      {/* Status actions */}
      <DietStatusActions diet={diet} patientId={patientId} />

      {/* Meals */}
      <div className="space-y-3">
        <h3 className="text-sm font-semibold text-muted-foreground uppercase tracking-wide">
          Refeicoes ({meals.length})
        </h3>

        {meals.map((meal) => (
          <MealCard
            key={meal.id}
            meal={meal}
            dietId={dietId}
            isDraft={isDraft}
          />
        ))}

        {/* Add meal form */}
        {showAddMeal && (
          <MealForm
            dietId={dietId}
            onClose={() => setShowAddMeal(false)}
          />
        )}

        {/* Add meal button (draft only) */}
        {isDraft && !showAddMeal && (
          <Button
            variant="outline"
            className="w-full"
            onClick={() => setShowAddMeal(true)}
          >
            <Plus className="h-4 w-4 mr-2" />
            Adicionar Refeicao
          </Button>
        )}
      </div>
    </div>
  )
}
