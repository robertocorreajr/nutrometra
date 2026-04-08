"use client"

import { useState } from "react"
import { Button, Card, CardContent } from "@nutrometra/ui"
import { useDeleteMeal } from "@nutrometra/api-client/hooks"
import type { DietMeal } from "@nutrometra/api-client"
import { MealForm } from "@/components/diet/meal-form"
import { MealItemRow } from "@/components/diet/meal-item-row"
import { AddMealItemForm } from "@/components/diet/add-meal-item-form"
import { ConfirmDialog } from "@/components/shared/confirm-dialog"
import { ChevronDown, ChevronRight, Pencil, Trash2, Plus } from "lucide-react"

interface MealCardProps {
  meal: DietMeal
  dietId: string
  isDraft: boolean
}

export function MealCard({ meal, dietId, isDraft }: MealCardProps) {
  const [expanded, setExpanded] = useState(true)
  const [showEdit, setShowEdit] = useState(false)
  const [showAddItem, setShowAddItem] = useState(false)
  const [showDelete, setShowDelete] = useState(false)

  const deleteMeal = useDeleteMeal(dietId)

  async function handleDelete() {
    try {
      await deleteMeal.mutateAsync(meal.id)
      setShowDelete(false)
    } catch {
      // Error handled by mutation state
    }
  }

  if (showEdit) {
    return (
      <MealForm
        dietId={dietId}
        meal={meal}
        onClose={() => setShowEdit(false)}
      />
    )
  }

  const items = meal.items ?? []

  return (
    <Card>
      <CardContent className="p-0">
        {/* Header */}
        <div className="flex items-center justify-between px-4 py-3 border-b">
          <button
            type="button"
            onClick={() => setExpanded(!expanded)}
            className="flex items-center gap-2 min-w-0 flex-1 text-left"
          >
            {expanded ? (
              <ChevronDown className="h-4 w-4 shrink-0 text-muted-foreground" />
            ) : (
              <ChevronRight className="h-4 w-4 shrink-0 text-muted-foreground" />
            )}
            <span className="font-medium text-sm truncate">{meal.meal_name}</span>
            <span className="text-xs text-muted-foreground shrink-0">
              Refeicao {meal.meal_order}
            </span>
          </button>

          {isDraft && (
            <div className="flex items-center gap-1 shrink-0 ml-2">
              <button
                type="button"
                onClick={() => setShowEdit(true)}
                className="text-muted-foreground hover:text-foreground p-1"
                aria-label="Editar refeicao"
              >
                <Pencil className="h-4 w-4" />
              </button>
              <button
                type="button"
                onClick={() => setShowDelete(true)}
                className="text-muted-foreground hover:text-destructive p-1"
                aria-label="Excluir refeicao"
              >
                <Trash2 className="h-4 w-4" />
              </button>
            </div>
          )}
        </div>

        {/* Notes */}
        {expanded && meal.notes && (
          <div className="px-4 py-2 bg-muted/30">
            <p className="text-xs text-muted-foreground">{meal.notes}</p>
          </div>
        )}

        {/* Items */}
        {expanded && (
          <div className="px-4">
            {items.length === 0 && !showAddItem && (
              <p className="text-sm text-muted-foreground py-3">
                Nenhum alimento adicionado.
              </p>
            )}

            {items.map((item) => (
              <MealItemRow
                key={item.id}
                item={item}
                dietId={dietId}
                isDraft={isDraft}
              />
            ))}

            {/* Add item form */}
            {showAddItem && (
              <div className="py-2">
                <AddMealItemForm
                  dietId={dietId}
                  mealId={meal.id}
                  onClose={() => setShowAddItem(false)}
                />
              </div>
            )}

            {/* Add item button */}
            {isDraft && !showAddItem && (
              <div className="py-2">
                <Button
                  type="button"
                  variant="ghost"
                  size="sm"
                  className="text-xs"
                  onClick={() => setShowAddItem(true)}
                >
                  <Plus className="h-3.5 w-3.5 mr-1" />
                  Adicionar Alimento
                </Button>
              </div>
            )}
          </div>
        )}
      </CardContent>

      <ConfirmDialog
        open={showDelete}
        onClose={() => setShowDelete(false)}
        onConfirm={handleDelete}
        title="Excluir Refeicao"
        description="Deseja excluir esta refeicao e todos os seus alimentos? Esta acao nao pode ser desfeita."
        confirmLabel="Excluir"
        variant="destructive"
        isLoading={deleteMeal.isPending}
      />
    </Card>
  )
}
