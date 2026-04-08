"use client"

import { Card, CardContent, CardHeader, CardTitle } from "@nutrometra/ui"
import { Users, CalendarCheck, CalendarDays } from "lucide-react"

interface KpiCardsProps {
  totalPatients: number
  appointmentsToday: number
  appointmentsWeek: number
}

export function KpiCards({ totalPatients, appointmentsToday, appointmentsWeek }: KpiCardsProps) {
  const cards = [
    {
      title: "Total de Pacientes",
      value: totalPatients,
      icon: Users,
      description: "pacientes cadastrados",
    },
    {
      title: "Consultas Hoje",
      value: appointmentsToday,
      icon: CalendarCheck,
      description: "agendadas para hoje",
    },
    {
      title: "Consultas na Semana",
      value: appointmentsWeek,
      icon: CalendarDays,
      description: "nos próximos 7 dias",
    },
  ]

  return (
    <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
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
              <div className="text-2xl font-bold">{card.value}</div>
              <p className="text-xs text-muted-foreground mt-1">{card.description}</p>
            </CardContent>
          </Card>
        )
      })}
    </div>
  )
}
