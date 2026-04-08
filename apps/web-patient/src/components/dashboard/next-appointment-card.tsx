import { Card, CardContent, CardHeader, CardTitle } from "@nutrometra/ui"
import { Calendar } from "lucide-react"
import { format, parseISO } from "date-fns"
import { ptBR } from "date-fns/locale"
import type { Appointment } from "@nutrometra/api-client/types"

interface NextAppointmentCardProps {
  appointment?: Appointment
}

export function NextAppointmentCard({ appointment }: NextAppointmentCardProps) {
  return (
    <Card>
      <CardHeader className="pb-2">
        <CardTitle className="text-sm font-medium text-muted-foreground flex items-center gap-2">
          <Calendar className="h-4 w-4" />
          Proxima Consulta
        </CardTitle>
      </CardHeader>
      <CardContent>
        {appointment ? (
          <div>
            <p className="text-lg font-semibold">
              {format(parseISO(appointment.start_at), "dd 'de' MMMM", { locale: ptBR })}
            </p>
            <p className="text-sm text-muted-foreground">
              {format(parseISO(appointment.start_at), "HH:mm")} - {format(parseISO(appointment.end_at), "HH:mm")}
            </p>
            <p className="text-xs text-muted-foreground mt-1 capitalize">
              {appointment.service_mode === "onsite" ? "Presencial" : "Online"}
            </p>
          </div>
        ) : (
          <p className="text-sm text-muted-foreground">Nenhuma consulta agendada</p>
        )}
      </CardContent>
    </Card>
  )
}
