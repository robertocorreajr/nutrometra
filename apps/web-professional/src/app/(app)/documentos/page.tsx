"use client"

import { useState } from "react"
import Link from "next/link"
import { Search, FileText, ArrowRight } from "lucide-react"
import { Input, LoadingState, ErrorState, EmptyState } from "@nutrometra/ui"
import { usePatients } from "@nutrometra/api-client/hooks"

export default function DocumentosPage() {
  const { data: patients, isLoading, isError, refetch } = usePatients()
  const [search, setSearch] = useState("")

  if (isLoading) return <LoadingState lines={6} />
  if (isError) return <ErrorState message="Erro ao carregar pacientes." onRetry={refetch} />

  const filtered = (patients ?? []).filter((p) =>
    p.full_name.toLowerCase().includes(search.toLowerCase()),
  )

  return (
    <div className="space-y-4">
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <h1 className="text-2xl font-bold tracking-tight">Documentos</h1>
        <div className="relative w-full sm:max-w-xs">
          <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
          <Input
            type="text"
            placeholder="Buscar paciente..."
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            className="pl-9"
          />
        </div>
      </div>

      {filtered.length === 0 ? (
        <EmptyState
          icon={<FileText className="h-12 w-12" />}
          title="Nenhum paciente encontrado"
          description="Cadastre pacientes para gerenciar seus documentos."
        />
      ) : (
        <div className="grid gap-3">
          {filtered.map((patient) => (
            <Link
              key={patient.id}
              href={`/pacientes/${patient.id}/documentos`}
              className="flex items-center justify-between rounded-lg border p-4 hover:bg-muted/50 transition-colors"
            >
              <div className="min-w-0">
                <p className="font-medium truncate">{patient.full_name}</p>
                <p className="text-sm text-muted-foreground">Ver documentos do paciente</p>
              </div>
              <ArrowRight className="h-4 w-4 text-muted-foreground shrink-0" />
            </Link>
          ))}
        </div>
      )}
    </div>
  )
}
