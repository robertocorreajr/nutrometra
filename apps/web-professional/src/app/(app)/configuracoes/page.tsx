"use client"

import { PageHeader } from "@nutrometra/ui"
import { Settings } from "lucide-react"

export default function ConfiguracoesPage() {
  return (
    <div>
      <PageHeader title="Configurações" description="Gerencie as configurações da sua conta" />
      <div className="flex flex-col items-center justify-center py-12 text-center">
        <Settings className="h-12 w-12 text-muted-foreground mb-4" />
        <p className="text-muted-foreground">Configurações avançadas serão disponibilizadas em versão futura.</p>
      </div>
    </div>
  )
}
