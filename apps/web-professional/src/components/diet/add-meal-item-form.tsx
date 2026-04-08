"use client"

import { useState } from "react"
import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { Button, Input, Card, CardContent } from "@nutrometra/ui"
import { useAddMealItem } from "@nutrometra/api-client/hooks"
import { mealItemSchema, type MealItemFormValues } from "@/lib/schemas/diet"
import type { FoodItem } from "@nutrometra/api-client"
import { FoodSearchModal } from "@/components/catalog/food-search-modal"
import { Search } from "lucide-react"

interface AddMealItemFormProps {
  dietId: string
  mealId: string
  onClose: () => void
}

const unitOptions = [
  { value: "g", label: "Gramas (g)" },
  { value: "ml", label: "Mililitros (ml)" },
  { value: "unidade", label: "Unidade" },
]

export function AddMealItemForm({ dietId, mealId, onClose }: AddMealItemFormProps) {
  const [showFoodSearch, setShowFoodSearch] = useState(false)
  const [selectedFood, setSelectedFood] = useState<FoodItem | null>(null)
  const addMealItem = useAddMealItem(dietId, mealId)

  const {
    register,
    handleSubmit,
    setValue,
    formState: { errors },
  } = useForm<MealItemFormValues>({
    resolver: zodResolver(mealItemSchema),
    defaultValues: {
      food_item_id: "",
      quantity_value: 100,
      quantity_unit: "g",
      sort_order: 0,
      preparation_notes: "",
    },
  })

  function handleFoodSelect(food: FoodItem) {
    setSelectedFood(food)
    setValue("food_item_id", food.id, { shouldValidate: true })
  }

  async function onSubmit(data: MealItemFormValues) {
    try {
      await addMealItem.mutateAsync({
        food_item_id: data.food_item_id,
        quantity_value: data.quantity_value,
        quantity_unit: data.quantity_unit,
        sort_order: data.sort_order,
        preparation_notes: data.preparation_notes || undefined,
      })
      onClose()
    } catch {
      // Error handled by mutation state
    }
  }

  return (
    <>
      <Card>
        <CardContent className="p-4">
          <h4 className="text-sm font-semibold mb-3">Adicionar Alimento</h4>
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-3">
            {/* Food selection */}
            <div>
              <label className="block text-sm font-medium mb-1">Alimento *</label>
              <input type="hidden" {...register("food_item_id")} />
              {selectedFood ? (
                <div className="flex items-center justify-between rounded-md border p-2">
                  <span className="text-sm font-medium">{selectedFood.name}</span>
                  <Button
                    type="button"
                    variant="outline"
                    size="sm"
                    onClick={() => setShowFoodSearch(true)}
                  >
                    Trocar
                  </Button>
                </div>
              ) : (
                <Button
                  type="button"
                  variant="outline"
                  className="w-full justify-start text-muted-foreground"
                  onClick={() => setShowFoodSearch(true)}
                >
                  <Search className="h-4 w-4 mr-2" />
                  Buscar alimento...
                </Button>
              )}
              {errors.food_item_id && (
                <p className="text-sm text-destructive mt-1">{errors.food_item_id.message}</p>
              )}
            </div>

            {/* Quantity and unit */}
            <div className="grid grid-cols-2 gap-3">
              <div>
                <label htmlFor="quantity_value" className="block text-sm font-medium mb-1">
                  Quantidade *
                </label>
                <Input
                  id="quantity_value"
                  type="number"
                  step="0.1"
                  min="0.1"
                  {...register("quantity_value")}
                />
                {errors.quantity_value && (
                  <p className="text-sm text-destructive mt-1">{errors.quantity_value.message}</p>
                )}
              </div>
              <div>
                <label htmlFor="quantity_unit" className="block text-sm font-medium mb-1">
                  Unidade *
                </label>
                <select
                  id="quantity_unit"
                  className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
                  {...register("quantity_unit")}
                >
                  {unitOptions.map((opt) => (
                    <option key={opt.value} value={opt.value}>
                      {opt.label}
                    </option>
                  ))}
                </select>
              </div>
            </div>

            {/* Preparation notes */}
            <div>
              <label htmlFor="prep_notes" className="block text-sm font-medium mb-1">
                Modo de preparo
              </label>
              <textarea
                id="prep_notes"
                rows={2}
                placeholder="Ex: cozido, grelhado, cru..."
                className="flex w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
                {...register("preparation_notes")}
              />
            </div>

            {addMealItem.isError && (
              <p className="text-sm text-destructive">Erro ao adicionar alimento. Tente novamente.</p>
            )}

            <div className="flex justify-end gap-2">
              <Button type="button" variant="outline" size="sm" onClick={onClose}>
                Cancelar
              </Button>
              <Button type="submit" size="sm" disabled={addMealItem.isPending}>
                {addMealItem.isPending ? "Adicionando..." : "Adicionar"}
              </Button>
            </div>
          </form>
        </CardContent>
      </Card>

      <FoodSearchModal
        open={showFoodSearch}
        onClose={() => setShowFoodSearch(false)}
        onSelect={handleFoodSelect}
      />
    </>
  )
}
