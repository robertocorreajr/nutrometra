"use client"

import { useState } from "react"
import { useFood, useDeleteMealItem } from "@nutrometra/api-client/hooks"
import type { DietMealItem } from "@nutrometra/api-client"
import { ConfirmDialog } from "@/components/shared/confirm-dialog"
import { SubstitutionList } from "@/components/diet/substitution-list"
import { Trash2, ChevronDown, ChevronRight } from "lucide-react"

interface MealItemRowProps {
  item: DietMealItem
  dietId: string
  isDraft: boolean
}

export function MealItemRow({ item, dietId, isDraft }: MealItemRowProps) {
  const { data: food } = useFood(item.food_item_id)
  const deleteMealItem = useDeleteMealItem(dietId)
  const [showDelete, setShowDelete] = useState(false)
  const [showSubs, setShowSubs] = useState(false)

  const subsCount = item.substitutions?.length ?? 0

  async function handleDelete() {
    try {
      await deleteMealItem.mutateAsync(item.id)
      setShowDelete(false)
    } catch {
      // Error handled by mutation state
    }
  }

  return (
    <div className="border-b last:border-b-0 py-2">
      <div className="flex items-start gap-2">
        {/* Expand substitutions toggle */}
        <button
          type="button"
          onClick={() => setShowSubs(!showSubs)}
          className="mt-0.5 text-muted-foreground hover:text-foreground shrink-0"
          aria-label={showSubs ? "Recolher substituicoes" : "Expandir substituicoes"}
        >
          {showSubs ? (
            <ChevronDown className="h-4 w-4" />
          ) : (
            <ChevronRight className="h-4 w-4" />
          )}
        </button>

        {/* Item info */}
        <div className="min-w-0 flex-1">
          <div className="flex items-center gap-2">
            <span className="text-sm font-medium">
              {food?.name ?? "Carregando..."}
            </span>
            <span className="text-sm text-muted-foreground">
              {item.quantity_value} {item.quantity_unit}
            </span>
          </div>
          {item.preparation_notes && (
            <p className="text-xs text-muted-foreground mt-0.5">
              {item.preparation_notes}
            </p>
          )}
          {subsCount > 0 && !showSubs && (
            <p className="text-xs text-muted-foreground mt-0.5">
              {subsCount} substituicao{subsCount > 1 ? "es" : ""}
            </p>
          )}
        </div>

        {/* Delete button */}
        {isDraft && (
          <button
            type="button"
            onClick={() => setShowDelete(true)}
            className="text-muted-foreground hover:text-destructive shrink-0 mt-0.5"
            aria-label="Remover alimento"
          >
            <Trash2 className="h-4 w-4" />
          </button>
        )}
      </div>

      {/* Substitutions */}
      {showSubs && (
        <SubstitutionList
          substitutions={item.substitutions ?? []}
          dietId={dietId}
          itemId={item.id}
          isDraft={isDraft}
        />
      )}

      <ConfirmDialog
        open={showDelete}
        onClose={() => setShowDelete(false)}
        onConfirm={handleDelete}
        title="Remover Alimento"
        description="Deseja remover este alimento da refeicao? As substituicoes tambem serao removidas."
        confirmLabel="Remover"
        variant="destructive"
        isLoading={deleteMealItem.isPending}
      />
    </div>
  )
}
