"use client"

import { PageHeader, LoadingState, ErrorState } from "@nutrometra/ui"
import { useBackofficeTenants } from "@nutrometra/api-client/hooks"
import { KpiCards } from "@/components/dashboard/kpi-cards"

export default function DashboardPage() {
  const { data: tenants, isLoading, isError, refetch } = useBackofficeTenants()

  return (
    <>
      <PageHeader title="Dashboard" description="Visao geral do sistema" />
      {isLoading && <LoadingState />}
      {isError && <ErrorState message="Nao foi possivel carregar os dados." onRetry={refetch} />}
      {tenants && <KpiCards tenants={tenants} />}
    </>
  )
}
