"use client"

import { PageHeader, LoadingState, ErrorState, EmptyState } from "@nutrometra/ui"
import { useMyAppointments } from "@nutrometra/api-client/hooks"
import { AppointmentCard } from "@/components/appointments/appointment-card"

export default function AgendaPage() {
  const { data: appointments, isLoading, isError, refetch } = useMyAppointments()

  return (
    <div className="max-w-lg mx-auto md:max-w-none space-y-4">
      <PageHeader title="Agenda" description="Suas consultas agendadas" />
      {isLoading && <LoadingState />}
      {isError && <ErrorState message="Nao foi possivel carregar a agenda." onRetry={refetch} />}
      {appointments && appointments.length === 0 && <EmptyState title="Nenhuma consulta" description="Voce nao possui consultas agendadas." />}
      {appointments && appointments.length > 0 && (
        <div className="space-y-3">
          {[...appointments]
            .sort((a, b) => a.start_at.localeCompare(b.start_at))
            .map((apt) => (
              <AppointmentCard key={apt.id} appointment={apt} />
            ))}
        </div>
      )}
    </div>
  )
}
