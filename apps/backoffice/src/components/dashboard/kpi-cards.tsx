"use client"

import { Card, CardContent, CardHeader, CardTitle } from "@nutrometra/ui"
import { Building2, CreditCard, Clock } from "lucide-react"
import type { BackofficeTenant } from "@nutrometra/api-client/types"

interface KpiCardsProps {
  tenants: BackofficeTenant[]
}

export function KpiCards({ tenants }: KpiCardsProps) {
  const total = tenants.length
  const active = tenants.filter((t) => t.status === "active").length
  const trialing = tenants.filter((t) => t.status === "trialing").length

  const cards = [
    { title: "Total Tenants", value: total, icon: Building2 },
    { title: "Assinaturas Ativas", value: active, icon: CreditCard },
    { title: "Em Trial", value: trialing, icon: Clock },
  ]

  return (
    <div className="grid gap-4 sm:grid-cols-3">
      {cards.map((card) => {
        const Icon = card.icon
        return (
          <Card key={card.title}>
            <CardHeader className="flex flex-row items-center justify-between pb-2">
              <CardTitle className="text-sm font-medium text-muted-foreground">
                {card.title}
              </CardTitle>
              <Icon className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <p className="text-2xl font-bold">{card.value}</p>
            </CardContent>
          </Card>
        )
      })}
    </div>
  )
}
