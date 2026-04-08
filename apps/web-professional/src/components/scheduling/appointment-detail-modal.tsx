"use client"

import { useState } from "react"
import { format } from "date-fns"
import { ptBR } from "date-fns/locale"
import { X, Clock, MapPin, Video, Home } from "lucide-react"
import { Button, Card, CardContent } from "@nutrometra/ui"
import { useUpdateAppointmentStatus } from "@nutrometra/api-client/hooks"
import type { Appointment, AppointmentStatus } from "@nutrometra/api-client"
import {
  appointmentStatusLabels,
  appointmentStatusColors,
  serviceModeOptions,
} from "@/lib/schemas/appointment"

interface AppointmentDetailModalProps {
  appointment: Appointment | null
  onClose: () => void
}

const serviceModeIcons: Record<string, React.ReactNode> = {
  onsite: <MapPin className="h-4 w-4" />,
  online: <Video className="h-4 w-4" />,
  home_visit: <Home className="h-4 w-4" />,
}

export function AppointmentDetailModal({
  appointment,
  onClose,
}: AppointmentDetailModalProps) {
  const [showCancelReason, setShowCancelReason] = useState(false)
  const [cancelReason, setCancelReason] = useState("")

  const updateStatus = useUpdateAppointmentStatus(appointment?.id ?? "")

  if (!appointment) return null

  const startDate = new Date(appointment.start_at)
  const endDate = new Date(appointment.end_at)
  const statusLabel = appointmentStatusLabels[appointment.status] ?? appointment.status
  const statusColor = appointmentStatusColors[appointment.status] ?? ""
  const serviceModeLabel =
    serviceModeOptions.find((o) => o.value === appointment.service_mode)?.label ??
    appointment.service_mode

  async function handleStatusChange(status: AppointmentStatus, reason?: string) {
    try {
      await updateStatus.mutateAsync({
        status,
        cancellation_reason: reason,
      })
      setShowCancelReason(false)
      setCancelReason("")
      onClose()
    } catch {
      // Error handled by mutation
    }
  }

  const isScheduled = appointment.status === "scheduled"
  const isConfirmed = appointment.status === "confirmed"

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
      {/* Backdrop */}
      <div
        className="absolute inset-0 bg-black/50"
        onClick={onClose}
      />

      {/* Modal */}
      <Card className="relative z-10 w-full max-w-md">
        <CardContent className="pt-6">
          {/* Close button */}
          <button
            type="button"
            onClick={onClose}
            className="absolute top-4 right-4 text-muted-foreground hover:text-foreground"
          >
            <X className="h-5 w-5" />
          </button>

          <h2 className="text-lg font-semibold mb-4">Detalhes da Consulta</h2>

          {/* Status badge */}
          <div className="mb-4">
            <span
              className={`inline-flex items-center rounded-full px-2.5 py-1 text-xs font-medium border ${statusColor}`}
            >
              {statusLabel}
            </span>
          </div>

          {/* Date/Time */}
          <div className="space-y-3 mb-6">
            <div className="flex items-center gap-2 text-sm">
              <Clock className="h-4 w-4 text-muted-foreground" />
              <span>
                {format(startDate, "EEEE, dd 'de' MMMM 'de' yyyy", { locale: ptBR })}
              </span>
            </div>
            <div className="flex items-center gap-2 text-sm pl-6">
              <span className="font-medium">
                {format(startDate, "HH:mm")} - {format(endDate, "HH:mm")}
              </span>
            </div>

            {/* Service mode */}
            <div className="flex items-center gap-2 text-sm">
              {serviceModeIcons[appointment.service_mode]}
              <span>{serviceModeLabel}</span>
            </div>

            {/* Patient ID */}
            {appointment.patient_id && (
              <div className="text-sm">
                <span className="text-muted-foreground">Paciente: </span>
                <span className="font-mono text-xs">{appointment.patient_id}</span>
              </div>
            )}

            {/* Notes */}
            {appointment.notes && (
              <div className="text-sm">
                <span className="text-muted-foreground">Observações: </span>
                <p className="mt-1 whitespace-pre-wrap">{appointment.notes}</p>
              </div>
            )}

            {/* Cancellation reason */}
            {appointment.cancellation_reason && (
              <div className="text-sm">
                <span className="text-muted-foreground">Motivo do cancelamento: </span>
                <p className="mt-1">{appointment.cancellation_reason}</p>
              </div>
            )}
          </div>

          {/* Error message */}
          {updateStatus.isError && (
            <p className="text-sm text-destructive mb-4">
              Erro ao atualizar status. Tente novamente.
            </p>
          )}

          {/* Cancel reason input */}
          {showCancelReason && (
            <div className="mb-4 space-y-2">
              <label htmlFor="cancel-reason" className="text-sm font-medium">
                Motivo do cancelamento
              </label>
              <textarea
                id="cancel-reason"
                value={cancelReason}
                onChange={(e) => setCancelReason(e.target.value)}
                rows={2}
                placeholder="Descreva o motivo..."
                className="flex w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
              />
              <div className="flex gap-2 justify-end">
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={() => {
                    setShowCancelReason(false)
                    setCancelReason("")
                  }}
                >
                  Voltar
                </Button>
                <Button
                  variant="destructive"
                  size="sm"
                  disabled={updateStatus.isPending}
                  onClick={() => handleStatusChange("cancelled", cancelReason)}
                >
                  {updateStatus.isPending ? "Cancelando..." : "Confirmar Cancelamento"}
                </Button>
              </div>
            </div>
          )}

          {/* Action buttons */}
          {!showCancelReason && (isScheduled || isConfirmed) && (
            <div className="flex flex-wrap gap-2">
              {isScheduled && (
                <Button
                  size="sm"
                  disabled={updateStatus.isPending}
                  onClick={() => handleStatusChange("confirmed")}
                >
                  Confirmar
                </Button>
              )}
              {isConfirmed && (
                <Button
                  size="sm"
                  disabled={updateStatus.isPending}
                  onClick={() => handleStatusChange("completed")}
                >
                  Completar
                </Button>
              )}
              <Button
                variant="destructive"
                size="sm"
                disabled={updateStatus.isPending}
                onClick={() => setShowCancelReason(true)}
              >
                Cancelar
              </Button>
              <Button
                variant="outline"
                size="sm"
                disabled={updateStatus.isPending}
                onClick={() => handleStatusChange("no_show")}
              >
                Não Compareceu
              </Button>
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  )
}
