"use client"

import { PageHeader, LoadingState, ErrorState } from "@nutrometra/ui"
import { usePatients } from "@nutrometra/api-client/hooks"
import { KpiCards } from "@/components/dashboard/kpi-cards"
import { RecentPatients } from "@/components/dashboard/recent-patients"
import { UpcomingAppointments } from "@/components/dashboard/upcoming-appointments"

export default function DashboardPage() {
  const { data: patients, isLoading, isError, refetch } = usePatients()

  if (isLoading) return <LoadingState lines={6} />
  if (isError) return <ErrorState message="Não foi possível carregar o dashboard." onRetry={refetch} />

  const patientList = patients ?? []

  return (
    <div>
      <PageHeader title="Dashboard" description="Visão geral do seu consultório" />

      <div className="space-y-6">
        <KpiCards
          totalPatients={patientList.length}
          appointmentsToday={0}
          appointmentsWeek={0}
        />

        <div className="grid gap-6 lg:grid-cols-2">
          <RecentPatients patients={patientList} />
          <UpcomingAppointments />
        </div>
      </div>
    </div>
  )
}
