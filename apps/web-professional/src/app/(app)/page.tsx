"use client"

import { useMemo } from "react"
import { startOfDay, endOfDay, addDays } from "date-fns"
import { PageHeader, LoadingState, ErrorState } from "@nutrometra/ui"
import { usePatients, useProfessionalMe, useAppointments } from "@nutrometra/api-client/hooks"
import { KpiCards } from "@/components/dashboard/kpi-cards"
import { RecentPatients } from "@/components/dashboard/recent-patients"
import { UpcomingAppointments } from "@/components/dashboard/upcoming-appointments"

const CANCELLED_STATUSES = ["cancelled", "no_show"] as const

export default function DashboardPage() {
  const { data: patients, isLoading, isError, refetch } = usePatients()
  const { data: professional } = useProfessionalMe()

  const professionalId = professional?.id ?? ""

  const now = useMemo(() => new Date(), [])
  const todayFrom = useMemo(() => startOfDay(now).toISOString(), [now])
  const todayTo = useMemo(() => endOfDay(now).toISOString(), [now])
  const weekTo = useMemo(() => endOfDay(addDays(now, 6)).toISOString(), [now])

  const { data: todayAppointments } = useAppointments({
    professional_id: professionalId,
    from: todayFrom,
    to: todayTo,
  })

  const { data: weekAppointments } = useAppointments({
    professional_id: professionalId,
    from: todayFrom,
    to: weekTo,
  })

  const appointmentsToday = useMemo(
    () =>
      (todayAppointments ?? []).filter(
        (a) => !CANCELLED_STATUSES.includes(a.status as typeof CANCELLED_STATUSES[number]),
      ).length,
    [todayAppointments],
  )

  const appointmentsWeek = useMemo(
    () =>
      (weekAppointments ?? []).filter(
        (a) => !CANCELLED_STATUSES.includes(a.status as typeof CANCELLED_STATUSES[number]),
      ).length,
    [weekAppointments],
  )

  if (isLoading) return <LoadingState lines={6} />
  if (isError) return <ErrorState message="Não foi possível carregar o dashboard." onRetry={refetch} />

  const patientList = patients ?? []

  return (
    <div>
      <PageHeader title="Dashboard" description="Visão geral do seu consultório" />

      <div className="space-y-6">
        <KpiCards
          totalPatients={patientList.length}
          appointmentsToday={appointmentsToday}
          appointmentsWeek={appointmentsWeek}
        />

        <div className="grid gap-6 lg:grid-cols-2">
          <RecentPatients patients={patientList} />
          <UpcomingAppointments />
        </div>
      </div>
    </div>
  )
}
