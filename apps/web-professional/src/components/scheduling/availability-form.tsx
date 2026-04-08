"use client"

import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { Button, Input, Card, CardContent } from "@nutrometra/ui"
import { useCreateAvailabilityRule } from "@nutrometra/api-client/hooks"
import {
  availabilitySchema,
  dayOfWeekOptions,
  type AvailabilityFormValues,
} from "@/lib/schemas/availability"
import { serviceModeOptions } from "@/lib/schemas/appointment"

interface AvailabilityFormProps {
  professionalId: string
  onSuccess: () => void
  onCancel?: () => void
}

export function AvailabilityForm({
  professionalId,
  onSuccess,
  onCancel,
}: AvailabilityFormProps) {
  const createRule = useCreateAvailabilityRule(professionalId)

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<AvailabilityFormValues>({
    resolver: zodResolver(availabilitySchema),
    defaultValues: {
      day_of_week: 1,
      start_time: "08:00",
      end_time: "18:00",
      service_mode: "onsite",
    },
  })

  async function handleFormSubmit(data: AvailabilityFormValues) {
    try {
      await createRule.mutateAsync({
        day_of_week: data.day_of_week,
        start_time: data.start_time,
        end_time: data.end_time,
        service_mode: data.service_mode,
        address_id: data.address_id,
      })
      reset()
      onSuccess()
    } catch {
      // Error handled by mutation
    }
  }

  return (
    <Card>
      <CardContent className="pt-6">
        <form onSubmit={handleSubmit(handleFormSubmit)} className="space-y-4">
          {/* Day of week */}
          <div className="space-y-2">
            <label htmlFor="day_of_week" className="text-sm font-medium">
              Dia da Semana *
            </label>
            <select
              id="day_of_week"
              {...register("day_of_week", { valueAsNumber: true })}
              className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
            >
              {dayOfWeekOptions.map((opt) => (
                <option key={opt.value} value={opt.value}>
                  {opt.label}
                </option>
              ))}
            </select>
            {errors.day_of_week && (
              <p className="text-sm text-destructive">{errors.day_of_week.message}</p>
            )}
          </div>

          {/* Time range */}
          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <label htmlFor="avail-start" className="text-sm font-medium">
                Início *
              </label>
              <Input id="avail-start" type="time" {...register("start_time")} />
              {errors.start_time && (
                <p className="text-sm text-destructive">{errors.start_time.message}</p>
              )}
            </div>
            <div className="space-y-2">
              <label htmlFor="avail-end" className="text-sm font-medium">
                Término *
              </label>
              <Input id="avail-end" type="time" {...register("end_time")} />
              {errors.end_time && (
                <p className="text-sm text-destructive">{errors.end_time.message}</p>
              )}
            </div>
          </div>

          {/* Service mode */}
          <div className="space-y-2">
            <label htmlFor="avail-service-mode" className="text-sm font-medium">
              Modalidade *
            </label>
            <select
              id="avail-service-mode"
              {...register("service_mode")}
              className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
            >
              {serviceModeOptions.map((opt) => (
                <option key={opt.value} value={opt.value}>
                  {opt.label}
                </option>
              ))}
            </select>
          </div>

          {/* Error message */}
          {createRule.isError && (
            <p className="text-sm text-destructive">
              Erro ao salvar regra de disponibilidade. Tente novamente.
            </p>
          )}

          {/* Actions */}
          <div className="flex gap-2 justify-end">
            {onCancel && (
              <Button type="button" variant="ghost" onClick={onCancel}>
                Cancelar
              </Button>
            )}
            <Button type="submit" disabled={createRule.isPending}>
              {createRule.isPending ? "Salvando..." : "Salvar Disponibilidade"}
            </Button>
          </div>
        </form>
      </CardContent>
    </Card>
  )
}
