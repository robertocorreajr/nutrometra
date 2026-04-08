"use client"

import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { Button, Input, Card, CardContent } from "@nutrometra/ui"
import { useAddMeal, useUpdateMeal } from "@nutrometra/api-client/hooks"
import { mealSchema, type MealFormValues } from "@/lib/schemas/diet"
import type { DietMeal } from "@nutrometra/api-client"

interface MealFormProps {
  dietId: string
  meal?: DietMeal
  onClose: () => void
}

export function MealForm({ dietId, meal, onClose }: MealFormProps) {
  const addMeal = useAddMeal(dietId)
  const updateMeal = useUpdateMeal(dietId)
  const isEdit = !!meal

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<MealFormValues>({
    resolver: zodResolver(mealSchema),
    defaultValues: {
      meal_name: meal?.meal_name ?? "",
      meal_order: meal?.meal_order ?? 1,
      notes: meal?.notes ?? "",
    },
  })

  async function onSubmit(data: MealFormValues) {
    try {
      if (isEdit && meal) {
        await updateMeal.mutateAsync({
          mealId: meal.id,
          meal_name: data.meal_name,
          meal_order: data.meal_order,
          notes: data.notes || undefined,
        })
      } else {
        await addMeal.mutateAsync({
          meal_name: data.meal_name,
          meal_order: data.meal_order,
          notes: data.notes || undefined,
        })
      }
      onClose()
    } catch {
      // Error handled by mutation state
    }
  }

  const isPending = isEdit ? updateMeal.isPending : addMeal.isPending
  const isError = isEdit ? updateMeal.isError : addMeal.isError

  return (
    <Card>
      <CardContent className="p-4">
        <h4 className="text-sm font-semibold mb-3">
          {isEdit ? "Editar Refeicao" : "Nova Refeicao"}
        </h4>
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-3">
          <div>
            <label htmlFor="meal_name" className="block text-sm font-medium mb-1">
              Nome *
            </label>
            <Input
              id="meal_name"
              placeholder="Ex: Cafe da manha"
              {...register("meal_name")}
            />
            {errors.meal_name && (
              <p className="text-sm text-destructive mt-1">{errors.meal_name.message}</p>
            )}
          </div>

          <div>
            <label htmlFor="meal_order" className="block text-sm font-medium mb-1">
              Ordem *
            </label>
            <Input
              id="meal_order"
              type="number"
              min={1}
              {...register("meal_order")}
            />
            {errors.meal_order && (
              <p className="text-sm text-destructive mt-1">{errors.meal_order.message}</p>
            )}
          </div>

          <div>
            <label htmlFor="meal_notes" className="block text-sm font-medium mb-1">
              Observacoes
            </label>
            <textarea
              id="meal_notes"
              rows={2}
              placeholder="Observacoes opcionais..."
              className="flex w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
              {...register("notes")}
            />
          </div>

          {isError && (
            <p className="text-sm text-destructive">Erro ao salvar refeicao. Tente novamente.</p>
          )}

          <div className="flex justify-end gap-2">
            <Button type="button" variant="outline" size="sm" onClick={onClose}>
              Cancelar
            </Button>
            <Button type="submit" size="sm" disabled={isPending}>
              {isPending ? "Salvando..." : isEdit ? "Salvar" : "Adicionar"}
            </Button>
          </div>
        </form>
      </CardContent>
    </Card>
  )
}
