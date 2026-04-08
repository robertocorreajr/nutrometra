"use client"

import { useState, useEffect } from "react"
import { Search, X, Apple } from "lucide-react"
import { Button, Input, LoadingState, EmptyState } from "@nutrometra/ui"
import { useFoods, useFoodGroups } from "@nutrometra/api-client/hooks"
import type { FoodItem } from "@nutrometra/api-client"

interface FoodSearchModalProps {
  open: boolean
  onClose: () => void
  onSelect: (food: FoodItem) => void
}

const sourceLabels: Record<string, string> = {
  system: "Sistema",
  taco: "TACO",
  ibge: "IBGE",
  tenant: "Personalizado",
}

function formatNutrition(food: FoodItem): string {
  const nf = food.nutrition_facts
  if (!nf) return ""
  const parts: string[] = []
  if (nf.calories_kcal != null) parts.push(`Cal ${nf.calories_kcal}`)
  if (nf.protein_g != null) parts.push(`Prot ${nf.protein_g}g`)
  if (nf.carbs_g != null) parts.push(`Carb ${nf.carbs_g}g`)
  if (nf.total_fat_g != null) parts.push(`Gord ${nf.total_fat_g}g`)
  return parts.join(" | ")
}

export function FoodSearchModal({ open, onClose, onSelect }: FoodSearchModalProps) {
  const [searchInput, setSearchInput] = useState("")
  const [debouncedQuery, setDebouncedQuery] = useState("")
  const [selectedGroup, setSelectedGroup] = useState("")

  // Debounce search input
  useEffect(() => {
    const timer = setTimeout(() => setDebouncedQuery(searchInput), 300)
    return () => clearTimeout(timer)
  }, [searchInput])

  // Reset state when modal opens
  useEffect(() => {
    if (open) {
      setSearchInput("")
      setDebouncedQuery("")
      setSelectedGroup("")
    }
  }, [open])

  const { data: foods, isLoading } = useFoods(
    debouncedQuery || undefined,
    selectedGroup || undefined,
  )
  const { data: groups } = useFoodGroups()

  if (!open) return null

  function handleSelect(food: FoodItem) {
    onSelect(food)
    onClose()
  }

  return (
    <div className="fixed inset-0 z-50 flex items-start justify-center p-4 pt-[10vh]">
      {/* Backdrop */}
      <div
        className="absolute inset-0 bg-black/50"
        onClick={onClose}
      />

      {/* Modal */}
      <div className="relative z-10 w-full max-w-lg rounded-lg border bg-background shadow-lg flex flex-col max-h-[80vh]">
        {/* Header */}
        <div className="flex items-center justify-between border-b px-4 py-3">
          <h2 className="text-lg font-semibold">Buscar Alimento</h2>
          <button
            type="button"
            onClick={onClose}
            className="text-muted-foreground hover:text-foreground"
            aria-label="Fechar"
          >
            <X className="h-5 w-5" />
          </button>
        </div>

        {/* Search & Filter */}
        <div className="space-y-3 border-b px-4 py-3">
          <div className="relative">
            <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
            <Input
              type="text"
              placeholder="Buscar por nome do alimento..."
              value={searchInput}
              onChange={(e) => setSearchInput(e.target.value)}
              className="pl-9"
              autoFocus
            />
          </div>

          {groups && groups.length > 0 && (
            <select
              value={selectedGroup}
              onChange={(e) => setSelectedGroup(e.target.value)}
              className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
              aria-label="Filtrar por grupo alimentar"
            >
              <option value="">Todos os grupos</option>
              {groups.map((group) => (
                <option key={group} value={group}>
                  {group}
                </option>
              ))}
            </select>
          )}
        </div>

        {/* Results */}
        <div className="flex-1 overflow-y-auto px-4 py-2">
          {isLoading && <LoadingState lines={4} />}

          {!isLoading && (!foods || foods.length === 0) && (
            <EmptyState
              icon={<Apple className="h-10 w-10" />}
              title="Nenhum alimento encontrado"
              description={
                debouncedQuery
                  ? "Tente ajustar os termos de busca ou o grupo alimentar."
                  : "Digite para buscar alimentos no catálogo."
              }
            />
          )}

          {!isLoading && foods && foods.length > 0 && (
            <ul className="divide-y" role="listbox" aria-label="Resultados da busca">
              {foods.map((food) => (
                <li key={food.id}>
                  <button
                    type="button"
                    onClick={() => handleSelect(food)}
                    className="w-full text-left px-2 py-3 hover:bg-muted/50 rounded-md transition-colors"
                    role="option"
                    aria-selected={false}
                  >
                    <div className="flex items-start justify-between gap-2">
                      <div className="min-w-0 flex-1">
                        <p className="font-medium text-sm truncate">{food.name}</p>
                        <div className="flex flex-wrap items-center gap-1.5 mt-1">
                          <span className="inline-flex items-center rounded-full bg-secondary px-2 py-0.5 text-xs text-secondary-foreground">
                            {food.food_group}
                          </span>
                          {food.source && food.source !== "system" && (
                            <span className="inline-flex items-center rounded-full border px-2 py-0.5 text-xs text-muted-foreground">
                              {sourceLabels[food.source] ?? food.source}
                            </span>
                          )}
                        </div>
                        {food.nutrition_facts && (
                          <p className="text-xs text-muted-foreground mt-1">
                            {formatNutrition(food)}
                          </p>
                        )}
                      </div>
                    </div>
                  </button>
                </li>
              ))}
            </ul>
          )}
        </div>

        {/* Footer */}
        <div className="border-t px-4 py-3 flex justify-end">
          <Button variant="outline" size="sm" onClick={onClose}>
            Cancelar
          </Button>
        </div>
      </div>
    </div>
  )
}
