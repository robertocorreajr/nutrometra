"use client"

import Link from "next/link"
import { format } from "date-fns"
import { ptBR } from "date-fns/locale"
import { Card, CardContent, CardHeader, CardTitle, LoadingState } from "@nutrometra/ui"
import { Calendar, ArrowRight } from "lucide-react"
import { useProfessionalMe, useAppointments } from "@nutrometra/api-client/hooks"
import { appointmentStatusLabels, appointmentStatusColors } from "@/lib/schemas/appointment"

export function UpcomingAppointments() {
  const { data: professional } = useProfessionalMe()
  const professionalId = professional?.id ?? ""

  const now = new Date()
  const from = now.toISOString()
  const to = new Date(now.getTime() + 7 * 24 * 60 * 60 * 1000).toISOString()

  const { data: appointments, isLoading, isError } = useAppointments({
    professional_id: professionalId,
    from,
    to,
  })

  const upcoming = (appointments ?? [])
    .filter((a) => a.status !== "cancelled" && a.status !== "no_show")
    .sort((a, b) => new Date(a.start_at).getTime() - new Date(b.start_at).getTime())
    .slice(0, 5)

  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between">
        <CardTitle className="text-base">Próximas Consultas</CardTitle>
        <Link
          href="/agenda"
          className="text-sm text-primary hover:underline inline-flex items-center gap-1"
        >
          Ver agenda
          <ArrowRight className="h-3 w-3" />
        </Link>
      </CardHeader>
      <CardContent>
        {isLoading ? (
          <LoadingState lines={3} />
        ) : isError ? (
          <p className="text-sm text-destructive">Erro ao carregar consultas.</p>
        ) : upcoming.length === 0 ? (
          <div className="flex flex-col items-center justify-center py-8 text-center">
            <Calendar className="h-8 w-8 text-muted-foreground mb-3" />
            <p className="text-sm text-muted-foreground">
              Nenhuma consulta agendada para os próximos 7 dias.
            </p>
          </div>
        ) : (
          <div className="space-y-3">
            {upcoming.map((appt) => (
              <div
                key={appt.id}
                className="flex items-center justify-between gap-3 py-2 border-b last:border-0"
              >
                <div className="min-w-0">
                  <p className="text-sm font-medium">
                    {format(new Date(appt.start_at), "EEEE, dd/MM", { locale: ptBR })}
                  </p>
                  <p className="text-xs text-muted-foreground">
                    {format(new Date(appt.start_at), "HH:mm")} –{" "}
                    {format(new Date(appt.end_at), "HH:mm")}
                  </p>
                </div>
                <span
                  className={`shrink-0 inline-flex items-center rounded-full px-2 py-0.5 text-xs font-medium border ${appointmentStatusColors[appt.status]}`}
                >
                  {appointmentStatusLabels[appt.status]}
                </span>
              </div>
            ))}
          </div>
        )}
      </CardContent>
    </Card>
  )
}
