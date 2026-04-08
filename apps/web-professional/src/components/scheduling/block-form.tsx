"use client"

import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { Button, Input, Card, CardContent } from "@nutrometra/ui"
import { useCreateBlock } from "@nutrometra/api-client/hooks"
import { blockSchema, type BlockFormValues } from "@/lib/schemas/block"

interface BlockFormProps {
  professionalId: string
  onSuccess: () => void
  onCancel?: () => void
}

export function BlockForm({
  professionalId,
  onSuccess,
  onCancel,
}: BlockFormProps) {
  const createBlock = useCreateBlock(professionalId)

  const {
    register,
    handleSubmit,
    watch,
    reset,
    formState: { errors },
  } = useForm<BlockFormValues>({
    resolver: zodResolver(blockSchema),
    defaultValues: {
      start_at: "",
      end_at: "",
      reason: "",
      all_day: false,
    },
  })

  const allDay = watch("all_day")

  async function handleFormSubmit(data: BlockFormValues) {
    try {
      let startAt = data.start_at
      let endAt = data.end_at

      // If all_day, set times to full day
      if (data.all_day) {
        // Ensure date-only values get full day range
        if (startAt && !startAt.includes("T")) {
          startAt = `${startAt}T00:00:00`
        }
        if (endAt && !endAt.includes("T")) {
          endAt = `${endAt}T23:59:59`
        }
      }

      await createBlock.mutateAsync({
        start_at: startAt,
        end_at: endAt,
        reason: data.reason ?? "",
        all_day: data.all_day,
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
          {/* All day checkbox */}
          <div className="flex items-center gap-2">
            <input
              type="checkbox"
              id="all_day"
              {...register("all_day")}
              className="h-4 w-4 rounded border-input"
            />
            <label htmlFor="all_day" className="text-sm font-medium">
              Dia inteiro
            </label>
          </div>

          {/* Start */}
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div className="space-y-2">
              <label htmlFor="block-start" className="text-sm font-medium">
                {allDay ? "Data de Início *" : "Início *"}
              </label>
              <Input
                id="block-start"
                type={allDay ? "date" : "datetime-local"}
                {...register("start_at")}
              />
              {errors.start_at && (
                <p className="text-sm text-destructive">{errors.start_at.message}</p>
              )}
            </div>

            <div className="space-y-2">
              <label htmlFor="block-end" className="text-sm font-medium">
                {allDay ? "Data de Término *" : "Término *"}
              </label>
              <Input
                id="block-end"
                type={allDay ? "date" : "datetime-local"}
                {...register("end_at")}
              />
              {errors.end_at && (
                <p className="text-sm text-destructive">{errors.end_at.message}</p>
              )}
            </div>
          </div>

          {/* Reason */}
          <div className="space-y-2">
            <label htmlFor="block-reason" className="text-sm font-medium">
              Motivo
            </label>
            <Input
              id="block-reason"
              {...register("reason")}
              placeholder="Ex: Férias, Congresso, Feriado..."
            />
          </div>

          {/* Error message */}
          {createBlock.isError && (
            <p className="text-sm text-destructive">
              Erro ao criar bloqueio. Tente novamente.
            </p>
          )}

          {/* Actions */}
          <div className="flex gap-2 justify-end">
            {onCancel && (
              <Button type="button" variant="ghost" onClick={onCancel}>
                Cancelar
              </Button>
            )}
            <Button type="submit" disabled={createBlock.isPending}>
              {createBlock.isPending ? "Salvando..." : "Criar Bloqueio"}
            </Button>
          </div>
        </form>
      </CardContent>
    </Card>
  )
}
