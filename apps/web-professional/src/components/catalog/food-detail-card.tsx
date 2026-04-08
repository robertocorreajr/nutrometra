"use client"

import { Card, CardContent, CardHeader, CardTitle } from "@nutrometra/ui"
import type { FoodItem } from "@nutrometra/api-client"

interface FoodDetailCardProps {
  food: FoodItem
}

const sourceLabels: Record<string, string> = {
  system: "Sistema",
  taco: "TACO",
  ibge: "IBGE",
  tenant: "Personalizado",
}

interface NutrientRow {
  label: string
  value: number | undefined
  unit: string
}

function getNutrientRows(food: FoodItem): NutrientRow[] {
  const nf = food.nutrition_facts
  if (!nf) return []
  return [
    { label: "Calorias", value: nf.calories_kcal, unit: "kcal" },
    { label: "Proteína", value: nf.protein_g, unit: "g" },
    { label: "Carboidratos", value: nf.carbs_g, unit: "g" },
    { label: "Fibra", value: nf.fiber_g, unit: "g" },
    { label: "Açúcar", value: nf.sugar_g, unit: "g" },
    { label: "Gordura Total", value: nf.total_fat_g, unit: "g" },
    { label: "Gordura Saturada", value: nf.saturated_fat_g, unit: "g" },
    { label: "Gordura Trans", value: nf.trans_fat_g, unit: "g" },
    { label: "Colesterol", value: nf.cholesterol_mg, unit: "mg" },
    { label: "Sódio", value: nf.sodium_mg, unit: "mg" },
    { label: "Potássio", value: nf.potassium_mg, unit: "mg" },
    { label: "Cálcio", value: nf.calcium_mg, unit: "mg" },
    { label: "Ferro", value: nf.iron_mg, unit: "mg" },
  ]
}

export function FoodDetailCard({ food }: FoodDetailCardProps) {
  const nutrientRows = getNutrientRows(food).filter((r) => r.value != null)
  const measures = food.household_measures ?? []

  return (
    <Card>
      <CardHeader>
        <CardTitle className="text-base">{food.name}</CardTitle>
        <div className="flex flex-wrap items-center gap-1.5 mt-1">
          <span className="inline-flex items-center rounded-full bg-secondary px-2 py-0.5 text-xs text-secondary-foreground">
            {food.food_group}
          </span>
          <span className="inline-flex items-center rounded-full border px-2 py-0.5 text-xs text-muted-foreground">
            {sourceLabels[food.source] ?? food.source}
          </span>
          {food.brand && (
            <span className="text-xs text-muted-foreground">
              {food.brand}
            </span>
          )}
        </div>
      </CardHeader>

      <CardContent className="space-y-4">
        {/* Serving info */}
        <div className="text-sm text-muted-foreground">
          Porção: {food.serving_size_g}g ({food.serving_label})
        </div>

        {/* Nutrition facts grid */}
        {nutrientRows.length > 0 && (
          <div>
            <h4 className="text-sm font-medium mb-2">Informação Nutricional</h4>
            <div className="grid grid-cols-2 gap-x-4 gap-y-1.5">
              {nutrientRows.map((row) => (
                <div
                  key={row.label}
                  className="flex items-center justify-between text-sm border-b border-dashed py-1"
                >
                  <span className="text-muted-foreground">{row.label}</span>
                  <span className="font-medium tabular-nums">
                    {row.value} {row.unit}
                  </span>
                </div>
              ))}
            </div>
          </div>
        )}

        {/* Household measures */}
        {measures.length > 0 && (
          <div>
            <h4 className="text-sm font-medium mb-2">Medidas Caseiras</h4>
            <ul className="space-y-1">
              {measures.map((m) => (
                <li key={m.id} className="text-sm text-muted-foreground">
                  1 {m.label} = {m.grams}g
                </li>
              ))}
            </ul>
          </div>
        )}
      </CardContent>
    </Card>
  )
}
