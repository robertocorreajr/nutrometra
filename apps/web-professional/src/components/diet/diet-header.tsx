"use client"

import type { Diet } from "@nutrometra/api-client"
import { dietStatusLabels, dietStatusColors } from "@/lib/schemas/diet"

interface DietHeaderProps {
  diet: Diet
}

function formatDate(dateStr?: string): string {
  if (!dateStr) return ""
  try {
    return new Date(dateStr).toLocaleDateString("pt-BR")
  } catch {
    return dateStr
  }
}

export function DietHeader({ diet }: DietHeaderProps) {
  return (
    <div className="space-y-2">
      <div className="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
        <h2 className="text-xl font-bold tracking-tight">{diet.title}</h2>
        <div className="flex items-center gap-2">
          <span
            className={`inline-flex items-center rounded-full border px-2.5 py-0.5 text-xs font-medium ${dietStatusColors[diet.status]}`}
          >
            {dietStatusLabels[diet.status]}
          </span>
          <span className="text-xs text-muted-foreground">
            v{diet.version_number}
          </span>
        </div>
      </div>

      {diet.objective && (
        <p className="text-sm text-muted-foreground">{diet.objective}</p>
      )}

      {(diet.valid_from || diet.valid_until) && (
        <p className="text-xs text-muted-foreground">
          {diet.valid_from && `De ${formatDate(diet.valid_from)}`}
          {diet.valid_from && diet.valid_until && " "}
          {diet.valid_until && `ate ${formatDate(diet.valid_until)}`}
        </p>
      )}
    </div>
  )
}
