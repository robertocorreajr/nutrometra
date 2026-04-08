"use client"

import Link from "next/link"
import { useParams, usePathname } from "next/navigation"
import { ArrowLeft } from "lucide-react"
import { usePatient } from "@nutrometra/api-client/hooks"
import { TabNav, LoadingState, ErrorState, type Tab } from "@nutrometra/ui"

function StatusBadge({ active }: { active: boolean }) {
  return (
    <span
      className={
        active
          ? "inline-flex items-center rounded-full px-2 py-0.5 text-xs font-medium bg-green-100 text-green-700"
          : "inline-flex items-center rounded-full px-2 py-0.5 text-xs font-medium bg-gray-100 text-gray-600"
      }
    >
      {active ? "Ativo" : "Inativo"}
    </span>
  )
}

interface PatientLayoutProps {
  children: React.ReactNode
}

export default function PatientLayout({ children }: PatientLayoutProps) {
  const params = useParams<{ id: string }>()
  const pathname = usePathname()
  const patientId = params.id

  const { data: patient, isLoading, isError, refetch } = usePatient(patientId)

  if (isLoading) return <LoadingState lines={4} />
  if (isError || !patient) {
    return (
      <ErrorState
        message="Não foi possível carregar os dados do paciente."
        onRetry={refetch}
      />
    )
  }

  const basePath = `/pacientes/${patientId}`

  const tabs: Tab[] = [
    { label: "Resumo", href: basePath },
    { label: "Anamnese", href: `${basePath}/anamnese` },
    { label: "Evoluções", href: `${basePath}/evolucoes` },
    { label: "Anexos", href: `${basePath}/anexos` },
    { label: "Composição Corporal", href: `${basePath}/composicao-corporal` },
    { label: "IA Assistiva", href: `${basePath}/ia` },
  ]

  // Active detection: exact match OR prefix match for sub-routes
  function getActivePath(): string {
    const exactMatch = tabs.find((t) => t.href === pathname)
    if (exactMatch) return exactMatch.href

    // Find the longest prefix match (most specific sub-tab)
    const prefixMatch = tabs
      .filter((t) => pathname.startsWith(t.href + "/"))
      .sort((a, b) => b.href.length - a.href.length)[0]

    return prefixMatch?.href ?? basePath
  }

  return (
    <div>
      {/* Back button */}
      <div className="mb-4">
        <Link
          href="/pacientes"
          className="inline-flex items-center gap-1 text-sm text-muted-foreground hover:text-foreground"
        >
          <ArrowLeft className="h-4 w-4" />
          Voltar para pacientes
        </Link>
      </div>

      {/* Patient header */}
      <div className="flex flex-col gap-1 sm:flex-row sm:items-center sm:justify-between mb-4">
        <div className="flex items-center gap-3 min-w-0">
          <h1 className="text-2xl font-bold tracking-tight truncate">
            {patient.full_name}
          </h1>
          <StatusBadge active={patient.active} />
        </div>
      </div>

      {/* Tab navigation */}
      <TabNav tabs={tabs} activePath={getActivePath()} LinkComponent={Link} />

      {/* Tab content */}
      <div className="mt-6">{children}</div>
    </div>
  )
}
