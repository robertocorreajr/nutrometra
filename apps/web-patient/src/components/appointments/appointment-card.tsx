import { Card, CardContent } from "@nutrometra/ui"
import { Calendar, Clock, MapPin, Video } from "lucide-react"
import { format, parseISO } from "date-fns"
import { ptBR } from "date-fns/locale"
import type { Appointment } from "@nutrometra/api-client/types"

interface AppointmentCardProps {
  appointment: Appointment
}

const statusLabels: Record<string, string> = {
  scheduled: "Agendada", confirmed: "Confirmada", completed: "Realizada",
  cancelled: "Cancelada", no_show: "Faltou",
}
const statusColors: Record<string, string> = {
  scheduled: "bg-blue-100 text-blue-700", confirmed: "bg-green-100 text-green-700",
  completed: "bg-gray-100 text-gray-700", cancelled: "bg-red-100 text-red-700",
  no_show: "bg-yellow-100 text-yellow-700",
}

export function AppointmentCard({ appointment }: AppointmentCardProps) {
  const ModeIcon = appointment.service_mode === "online" ? Video : MapPin

  return (
    <Card>
      <CardContent className="pt-4">
        <div className="flex items-start justify-between">
          <div className="flex gap-3">
            <div className="flex flex-col items-center pt-1">
              <Calendar className="h-5 w-5 text-primary" />
            </div>
            <div>
              <p className="font-semibold">
                {format(parseISO(appointment.start_at), "EEEE, dd 'de' MMMM", { locale: ptBR })}
              </p>
              <div className="flex items-center gap-1 text-sm text-muted-foreground mt-1">
                <Clock className="h-3.5 w-3.5" />
                <span>{format(parseISO(appointment.start_at), "HH:mm")} - {format(parseISO(appointment.end_at), "HH:mm")}</span>
              </div>
              <div className="flex items-center gap-1 text-sm text-muted-foreground mt-1">
                <ModeIcon className="h-3.5 w-3.5" />
                <span>{appointment.service_mode === "online" ? "Online" : appointment.service_mode === "onsite" ? "Presencial" : "Domiciliar"}</span>
              </div>
            </div>
          </div>
          <span className={`inline-flex items-center rounded-full px-2 py-1 text-xs font-medium ${statusColors[appointment.status] ?? "bg-gray-100 text-gray-700"}`}>
            {statusLabels[appointment.status] ?? appointment.status}
          </span>
        </div>
        {appointment.notes && (
          <p className="text-sm text-muted-foreground mt-3 pl-8">{appointment.notes}</p>
        )}
      </CardContent>
    </Card>
  )
}
