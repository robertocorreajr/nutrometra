"use client"

import { useParams } from "next/navigation"
import { LoadingState, ErrorState } from "@nutrometra/ui"
import { useBackofficeAuditLog } from "@nutrometra/api-client/hooks"
import { AuditLogTable } from "@/components/tenants/audit-log-table"

export default function AuditoriaPage() {
  const { id } = useParams<{ id: string }>()
  const { data: entries, isLoading, isError, refetch } = useBackofficeAuditLog(id)

  if (isLoading) return <LoadingState />
  if (isError) return <ErrorState message="Nao foi possivel carregar a auditoria." onRetry={refetch} />

  return <AuditLogTable entries={entries ?? []} />
}
