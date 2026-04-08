"use client"

import { useState } from "react"
import { Button } from "@nutrometra/ui"
import { useFood, useAddSubstitution, useDeleteSubstitution } from "@nutrometra/api-client/hooks"
import type { DietSubstitution, FoodItem } from "@nutrometra/api-client"
import { FoodSearchModal } from "@/components/catalog/food-search-modal"
import { ConfirmDialog } from "@/components/shared/confirm-dialog"
import { Plus, Trash2 } from "lucide-react"

interface SubstitutionListProps {
  substitutions: DietSubstitution[]
  dietId: string
  itemId: string
  isDraft: boolean
}

function SubstitutionRow({
  sub,
  dietId,
  isDraft,
}: {
  sub: DietSubstitution
  dietId: string
  isDraft: boolean
}) {
  const { data: food } = useFood(sub.substitute_food_item_id)
  const deleteSub = useDeleteSubstitution(dietId)
  const [showDelete, setShowDelete] = useState(false)

  async function handleDelete() {
    try {
      await deleteSub.mutateAsync(sub.id)
      setShowDelete(false)
    } catch {
      // Error handled by mutation state
    }
  }

  return (
    <div className="flex items-center justify-between py-1.5 px-2 rounded-md hover:bg-muted/30 text-sm">
      <div className="min-w-0 flex-1">
        <span className="font-medium">{food?.name ?? "Carregando..."}</span>
        {sub.quantity_value != null && sub.quantity_unit && (
          <span className="text-muted-foreground ml-2">
            {sub.quantity_value} {sub.quantity_unit}
          </span>
        )}
        {sub.notes && (
          <span className="text-muted-foreground ml-2">({sub.notes})</span>
        )}
      </div>
      {isDraft && (
        <button
          type="button"
          onClick={() => setShowDelete(true)}
          className="text-muted-foreground hover:text-destructive ml-2 shrink-0"
          aria-label="Remover substituicao"
        >
          <Trash2 className="h-3.5 w-3.5" />
        </button>
      )}

      <ConfirmDialog
        open={showDelete}
        onClose={() => setShowDelete(false)}
        onConfirm={handleDelete}
        title="Remover Substituicao"
        description="Deseja remover esta substituicao?"
        confirmLabel="Remover"
        variant="destructive"
        isLoading={deleteSub.isPending}
      />
    </div>
  )
}

export function SubstitutionList({ substitutions, dietId, itemId, isDraft }: SubstitutionListProps) {
  const [showFoodSearch, setShowFoodSearch] = useState(false)
  const addSub = useAddSubstitution(dietId)

  async function handleFoodSelect(food: FoodItem) {
    try {
      await addSub.mutateAsync({
        itemId,
        substitute_food_item_id: food.id,
        quantity_value: undefined,
        quantity_unit: undefined,
        notes: undefined,
        sort_order: substitutions.length,
      })
    } catch {
      // Error handled by mutation state
    }
  }

  return (
    <div className="mt-2 ml-4 border-l-2 border-muted pl-3">
      <div className="flex items-center justify-between mb-1">
        <span className="text-xs font-medium text-muted-foreground uppercase tracking-wide">
          Substituicoes
        </span>
        {isDraft && (
          <Button
            type="button"
            variant="ghost"
            size="sm"
            className="h-6 text-xs"
            onClick={() => setShowFoodSearch(true)}
          >
            <Plus className="h-3 w-3 mr-1" />
            Substituicao
          </Button>
        )}
      </div>

      {substitutions.length === 0 && (
        <p className="text-xs text-muted-foreground py-1">
          Nenhuma substituicao cadastrada.
        </p>
      )}

      {substitutions.map((sub) => (
        <SubstitutionRow
          key={sub.id}
          sub={sub}
          dietId={dietId}
          isDraft={isDraft}
        />
      ))}

      {addSub.isError && (
        <p className="text-xs text-destructive mt-1">Erro ao adicionar substituicao.</p>
      )}

      <FoodSearchModal
        open={showFoodSearch}
        onClose={() => setShowFoodSearch(false)}
        onSelect={handleFoodSelect}
      />
    </div>
  )
}
