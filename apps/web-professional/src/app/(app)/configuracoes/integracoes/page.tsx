"use client"

import { GoogleCalendarCard } from "@/components/settings/google-calendar-card"

export default function IntegracoesPage() {
  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-lg font-semibold mb-1">Integrações</h2>
        <p className="text-sm text-muted-foreground">
          Conecte serviços externos para automatizar seu fluxo de trabalho.
        </p>
      </div>

      <GoogleCalendarCard />
    </div>
  )
}
