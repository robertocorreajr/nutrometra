"use client"

import { format } from "date-fns"
import { MapPin, Video, Home } from "lucide-react"
import type { Appointment } from "@nutrometra/api-client"
import {
  appointmentStatusLabels,
  appointmentStatusColors,
} from "@/lib/schemas/appointment"

interface AppointmentCardProps {
  appointment: Appointment
  onSelect: (apt: Appointment) => void
}

const serviceModeIcons: Record<string, React.ReactNode> = {
  onsite: <MapPin className="h-3 w-3" />,
  online: <Video className="h-3 w-3" />,
  home_visit: <Home className="h-3 w-3" />,
}

const statusBorderColors: Record<string, string> = {
  scheduled: "border-l-blue-500",
  confirmed: "border-l-green-500",
  completed: "border-l-gray-400",
  cancelled: "border-l-red-500",
  no_show: "border-l-orange-500",
}

export function AppointmentCard({ appointment, onSelect }: AppointmentCardProps) {
  const startTime = format(new Date(appointment.start_at), "HH:mm")
  const endTime = format(new Date(appointment.end_at), "HH:mm")
  const statusLabel = appointmentStatusLabels[appointment.status] ?? appointment.status
  const statusColor = appointmentStatusColors[appointment.status] ?? ""
  const borderColor = statusBorderColors[appointment.status] ?? "border-l-gray-300"

  return (
    <button
      type="button"
      onClick={() => onSelect(appointment)}
      className={`w-full text-left rounded-md border border-l-4 ${borderColor} bg-card p-2 text-xs shadow-sm hover:shadow-md transition-shadow cursor-pointer`}
    >
      <div className="flex items-center justify-between gap-1">
        <span className="font-medium">
          {startTime} - {endTime}
        </span>
        <span className="flex items-center gap-0.5 text-muted-foreground">
          {serviceModeIcons[appointment.service_mode]}
        </span>
      </div>
      <div className="mt-1">
        <span
          className={`inline-flex items-center rounded-full px-1.5 py-0.5 text-[10px] font-medium border ${statusColor}`}
        >
          {statusLabel}
        </span>
      </div>
      {appointment.notes && (
        <p className="mt-1 text-muted-foreground truncate">{appointment.notes}</p>
      )}
    </button>
  )
}
