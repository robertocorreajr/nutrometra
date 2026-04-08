import { Card, CardContent, CardHeader, CardTitle } from "@nutrometra/ui"
import type { DietMeal } from "@nutrometra/api-client/types"

interface MealViewProps {
  meals: DietMeal[]
}

export function MealView({ meals }: MealViewProps) {
  const sorted = [...meals].sort((a, b) => a.meal_order - b.meal_order)

  return (
    <div className="space-y-4">
      {sorted.map((meal) => (
        <Card key={meal.id}>
          <CardHeader className="pb-2">
            <CardTitle className="text-base">{meal.meal_name}</CardTitle>
            {meal.notes && <p className="text-sm text-muted-foreground">{meal.notes}</p>}
          </CardHeader>
          <CardContent>
            {meal.items && meal.items.length > 0 ? (
              <ul className="space-y-2">
                {[...meal.items].sort((a, b) => a.sort_order - b.sort_order).map((item) => (
                  <li key={item.id} className="flex justify-between items-start text-sm border-b last:border-0 pb-2 last:pb-0">
                    <div className="flex-1">
                      <span className="font-medium">{item.amount_description ?? `${item.quantity_value} ${item.quantity_unit}`}</span>
                      {item.preparation_notes && (
                        <p className="text-xs text-muted-foreground mt-0.5">{item.preparation_notes}</p>
                      )}
                      {item.substitutions && item.substitutions.length > 0 && (
                        <div className="mt-1 pl-3 border-l-2 border-muted">
                          <p className="text-xs text-muted-foreground font-medium">Substituicoes:</p>
                          {item.substitutions.map((sub) => (
                            <p key={sub.id} className="text-xs text-muted-foreground">
                              {sub.notes ?? `${sub.quantity_value ?? ""} ${sub.quantity_unit ?? ""}`}
                            </p>
                          ))}
                        </div>
                      )}
                    </div>
                  </li>
                ))}
              </ul>
            ) : (
              <p className="text-sm text-muted-foreground">Nenhum item</p>
            )}
          </CardContent>
        </Card>
      ))}
    </div>
  )
}
