"use client"

import { useParams } from "next/navigation"
import { PageHeader, LoadingState, ErrorState, Card, CardContent } from "@nutrometra/ui"
import { useMyDiet } from "@nutrometra/api-client/hooks"
import { MealView } from "@/components/diet/meal-view"
import { format, parseISO } from "date-fns"

export default function DietDetailPage() {
  const { id } = useParams<{ id: string }>()
  const { data: diet, isLoading, isError, refetch } = useMyDiet(id)

  if (isLoading) return <LoadingState />
  if (isError) return <ErrorState message="Nao foi possivel carregar a dieta." onRetry={refetch} />
  if (!diet) return null

  return (
    <div className="max-w-lg mx-auto md:max-w-none space-y-4">
      <PageHeader title={diet.title} description={diet.objective} />
      <Card>
        <CardContent className="pt-4">
          <div className="flex flex-wrap gap-4 text-sm text-muted-foreground">
            {diet.published_at && <span>Publicada: {format(parseISO(diet.published_at), "dd/MM/yyyy")}</span>}
            {diet.valid_from && <span>De: {format(parseISO(diet.valid_from), "dd/MM/yyyy")}</span>}
            {diet.valid_until && <span>Ate: {format(parseISO(diet.valid_until), "dd/MM/yyyy")}</span>}
            <span>Versao: {diet.version_number}</span>
          </div>
        </CardContent>
      </Card>
      {diet.meals && diet.meals.length > 0 && <MealView meals={diet.meals} />}
    </div>
  )
}
