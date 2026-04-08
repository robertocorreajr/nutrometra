"use client"

import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { format } from "date-fns"
import { Button, Input, Card, CardContent } from "@nutrometra/ui"
import { usePatients } from "@nutrometra/api-client/hooks"
import {
  appointmentSchema,
  serviceModeOptions,
  type AppointmentFormValues,
} from "@/lib/schemas/appointment"

interface AppointmentFormProps {
  onSubmit: (data: AppointmentFormValues) => void
  isSubmitting: boolean
  onCancel: () => void
  defaultDate?: string
  defaultTime?: string
}

export function AppointmentForm({
  onSubmit,
  isSubmitting,
  onCancel,
  defaultDate,
  defaultTime,
}: AppointmentFormProps) {
  const { data: patients } = usePatients()

  // Parse defaults from ISO string or explicit values
  const parsedDate = defaultDate
    ? defaultDate.includes("T")
      ? format(new Date(defaultDate), "yyyy-MM-dd")
      : defaultDate
    : ""
  const parsedTime = defaultTime
    ? defaultTime
    : defaultDate && defaultDate.includes("T")
      ? format(new Date(defaultDate), "HH:mm")
      : ""

  const {
    register,
    handleSubmit,
    watch,
    setValue,
    formState: { errors },
  } = useForm<AppointmentFormValues & { date: string; start_time: string; end_time: string }>({
    resolver: zodResolver(appointmentSchema),
    defaultValues: {
      patient_id: "",
      start_at: "",
      end_at: "",
      service_mode: "onsite",
      notes: "",
    },
  })

  // Local state for date/time decomposition
  const dateValue = watch("date" as never) as string | undefined
  const startTimeValue = watch("start_time" as never) as string | undefined

  function handleFormSubmit(data: AppointmentFormValues) {
    onSubmit(data)
  }

  function buildDateTime(date: string, time: string): string {
    if (!date || !time) return ""
    return `${date}T${time}:00`
  }

  return (
    <Card>
      <CardContent className="pt-6">
        <form
          onSubmit={handleSubmit(handleFormSubmit)}
          className="space-y-4"
        >
          {/* Patient selector */}
          <div className="space-y-2">
            <label htmlFor="patient_id" className="text-sm font-medium">
              Paciente *
            </label>
            <select
              id="patient_id"
              {...register("patient_id")}
              className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
            >
              <option value="">Selecione um paciente</option>
              {(patients ?? []).map((p) => (
                <option key={p.id} value={p.id}>
                  {p.full_name}
                </option>
              ))}
            </select>
            {errors.patient_id && (
              <p className="text-sm text-destructive">{errors.patient_id.message}</p>
            )}
          </div>

          {/* Date */}
          <div className="space-y-2">
            <label htmlFor="apt-date" className="text-sm font-medium">
              Data *
            </label>
            <Input
              id="apt-date"
              type="date"
              defaultValue={parsedDate}
              onChange={(e) => {
                const date = e.target.value
                setValue("date" as never, date as never)
                const time = startTimeValue ?? parsedTime
                if (date && time) {
                  setValue("start_at", buildDateTime(date, time))
                }
              }}
            />
          </div>

          {/* Start time */}
          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <label htmlFor="apt-start-time" className="text-sm font-medium">
                Hora Início *
              </label>
              <Input
                id="apt-start-time"
                type="time"
                defaultValue={parsedTime}
                onChange={(e) => {
                  const time = e.target.value
                  setValue("start_time" as never, time as never)
                  const date = dateValue ?? parsedDate
                  if (date && time) {
                    setValue("start_at", buildDateTime(date, time))
                  }
                }}
              />
              {errors.start_at && (
                <p className="text-sm text-destructive">{errors.start_at.message}</p>
              )}
            </div>

            {/* End time */}
            <div className="space-y-2">
              <label htmlFor="apt-end-time" className="text-sm font-medium">
                Hora Término *
              </label>
              <Input
                id="apt-end-time"
                type="time"
                onChange={(e) => {
                  const time = e.target.value
                  setValue("end_time" as never, time as never)
                  const date = dateValue ?? parsedDate
                  if (date && time) {
                    setValue("end_at", buildDateTime(date, time))
                  }
                }}
              />
              {errors.end_at && (
                <p className="text-sm text-destructive">{errors.end_at.message}</p>
              )}
            </div>
          </div>

          {/* Service mode */}
          <div className="space-y-2">
            <label className="text-sm font-medium">Modalidade *</label>
            <div className="flex gap-4">
              {serviceModeOptions.map((opt) => (
                <label key={opt.value} className="flex items-center gap-2 text-sm cursor-pointer">
                  <input
                    type="radio"
                    value={opt.value}
                    {...register("service_mode")}
                    className="h-4 w-4"
                  />
                  {opt.label}
                </label>
              ))}
            </div>
          </div>

          {/* Notes */}
          <div className="space-y-2">
            <label htmlFor="apt-notes" className="text-sm font-medium">
              Observações
            </label>
            <textarea
              id="apt-notes"
              {...register("notes")}
              rows={3}
              placeholder="Observações sobre a consulta..."
              className="flex w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
            />
          </div>

          {/* Actions */}
          <div className="flex gap-2 justify-end">
            <Button type="button" variant="ghost" onClick={onCancel}>
              Cancelar
            </Button>
            <Button type="submit" disabled={isSubmitting}>
              {isSubmitting ? "Agendando..." : "Agendar Consulta"}
            </Button>
          </div>
        </form>
      </CardContent>
    </Card>
  )
}
