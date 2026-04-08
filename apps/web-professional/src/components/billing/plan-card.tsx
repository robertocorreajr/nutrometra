"use client"

import { Button } from "@nutrometra/ui"
import { cn } from "@nutrometra/ui"
import { formatCurrency } from "@/lib/schemas/billing"

interface PlanCardProps {
  name: string
  priceCents: number
  currency: string
  billingCycle: string
  isCurrent: boolean
  onSelect?: () => void
  isLoading?: boolean
}

export function PlanCard({
  name,
  priceCents,
  currency,
  billingCycle,
  isCurrent,
  onSelect,
  isLoading,
}: PlanCardProps) {
  const cycleLabel = billingCycle === "monthly" ? "/mês" : billingCycle === "yearly" ? "/ano" : ""

  return (
    <div
      className={cn(
        "rounded-lg border p-6 flex flex-col",
        isCurrent ? "border-primary bg-primary/5 ring-2 ring-primary" : "bg-card"
      )}
    >
      <h3 className="text-lg font-semibold">{name}</h3>
      <p className="text-2xl font-bold mt-2">
        {priceCents === 0 ? "Grátis" : formatCurrency(priceCents, currency)}
        {priceCents > 0 && <span className="text-sm font-normal text-muted-foreground">{cycleLabel}</span>}
      </p>
      <div className="mt-auto pt-4">
        {isCurrent ? (
          <Button variant="outline" disabled className="w-full">Plano Atual</Button>
        ) : (
          <Button className="w-full" onClick={onSelect} disabled={isLoading}>
            {isLoading ? "Aguarde..." : "Selecionar"}
          </Button>
        )}
      </div>
    </div>
  )
}
