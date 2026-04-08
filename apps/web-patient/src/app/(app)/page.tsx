"use client"

import { useMemo } from "react"
import { useAuth } from "@nutrometra/auth"
import { LoadingState } from "@nutrometra/ui"
import { useMyAppointments, useMyDiets } from "@nutrometra/api-client/hooks"
import { WelcomeCard } from "@/components/dashboard/welcome-card"
import { NextAppointmentCard } from "@/components/dashboard/next-appointment-card"
import { ActiveDietCard } from "@/components/dashboard/active-diet-card"

export default function DashboardPage() {
  const { user } = useAuth()
  // TODO: patientId should come from session/context. For now, show what we can.
  // Appointments are auto-scoped by the backend for authenticated patients.
  const now = useMemo(() => new Date().toISOString(), [])
  const { data: appointments, isLoading: loadingAppts } = useMyAppointments({ from: now })
  // Diets need patientId — placeholder until session provides it
  const patientId = "" // Will be populated from session in future
  const { data: diets, isLoading: loadingDiets } = useMyDiets(patientId)

  const nextAppointment = useMemo(() => {
    if (!appointments?.length) return undefined
    return appointments
      .filter((a) => a.status === "scheduled" || a.status === "confirmed")
      .sort((a, b) => a.start_at.localeCompare(b.start_at))[0]
  }, [appointments])

  const activeDiet = useMemo(() => {
    if (!diets?.length) return undefined
    return diets.find((d) => d.status === "published")
  }, [diets])

  if (loadingAppts && loadingDiets) return <LoadingState />

  return (
    <div className="space-y-4 max-w-lg mx-auto md:max-w-none">
      <WelcomeCard name={user?.name ?? undefined} />
      <div className="grid gap-4 md:grid-cols-2">
        <NextAppointmentCard appointment={nextAppointment} />
        <ActiveDietCard diet={activeDiet} />
      </div>
    </div>
  )
}
