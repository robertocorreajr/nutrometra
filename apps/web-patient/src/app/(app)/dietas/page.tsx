"use client"

import { useAuth } from "@nutrometra/auth"
import { PageHeader, LoadingState, ErrorState, EmptyState } from "@nutrometra/ui"
import { useMyDiets } from "@nutrometra/api-client/hooks"
import { DietCard } from "@/components/diet/diet-card"

export default function DietasPage() {
  const { patientId } = useAuth()
  const { data: diets, isLoading, isError, refetch } = useMyDiets(patientId ?? "")

  return (
    <div className="max-w-lg mx-auto md:max-w-none">
      <PageHeader title="Minhas Dietas" description="Dietas prescritas pelo seu nutricionista" />
      {isLoading && <LoadingState />}
      {isError && <ErrorState message="Nao foi possivel carregar as dietas." onRetry={refetch} />}
      {diets && diets.length === 0 && <EmptyState title="Nenhuma dieta" description="Seu nutricionista ainda nao publicou dietas para voce." />}
      {diets && diets.length > 0 && (
        <div className="space-y-3">
          {diets.filter((d) => d.status === "published").map((diet) => (
            <DietCard key={diet.id} diet={diet} />
          ))}
        </div>
      )}
    </div>
  )
}
