"use client"

import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { Button, Input, Card, CardContent } from "@nutrometra/ui"
import { useCreateDiet, useProfessionalMe } from "@nutrometra/api-client/hooks"
import { createDietSchema, type CreateDietFormValues } from "@/lib/schemas/diet"
import type { Diet } from "@nutrometra/api-client"

interface CreateDietFormProps {
  patientId: string
  onSuccess: (diet: Diet) => void
  onCancel: () => void
}

export function CreateDietForm({ patientId, onSuccess, onCancel }: CreateDietFormProps) {
  const { data: professional } = useProfessionalMe()
  const createDiet = useCreateDiet(patientId)

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<CreateDietFormValues>({
    resolver: zodResolver(createDietSchema),
    defaultValues: {
      title: "",
      objective: "",
      valid_from: "",
      valid_until: "",
    },
  })

  async function onSubmit(data: CreateDietFormValues) {
    if (!professional?.id) return

    try {
      const result = await createDiet.mutateAsync({
        professional_id: professional.id,
        title: data.title,
        objective: data.objective || undefined,
        valid_from: data.valid_from || undefined,
        valid_until: data.valid_until || undefined,
      })
      onSuccess(result)
    } catch {
      // Error handled by mutation state
    }
  }

  return (
    <Card>
      <CardContent className="p-4 sm:p-6">
        <h3 className="text-lg font-semibold mb-4">Nova Dieta</h3>
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <div>
            <label htmlFor="title" className="block text-sm font-medium mb-1">
              Titulo *
            </label>
            <Input
              id="title"
              placeholder="Ex: Plano alimentar para emagrecimento"
              {...register("title")}
            />
            {errors.title && (
              <p className="text-sm text-destructive mt-1">{errors.title.message}</p>
            )}
          </div>

          <div>
            <label htmlFor="objective" className="block text-sm font-medium mb-1">
              Objetivo
            </label>
            <textarea
              id="objective"
              rows={3}
              placeholder="Descreva o objetivo desta dieta..."
              className="flex w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
              {...register("objective")}
            />
          </div>

          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div>
              <label htmlFor="valid_from" className="block text-sm font-medium mb-1">
                Validade - Inicio
              </label>
              <Input id="valid_from" type="date" {...register("valid_from")} />
            </div>
            <div>
              <label htmlFor="valid_until" className="block text-sm font-medium mb-1">
                Validade - Fim
              </label>
              <Input id="valid_until" type="date" {...register("valid_until")} />
            </div>
          </div>

          {createDiet.isError && (
            <p className="text-sm text-destructive">Erro ao criar dieta. Tente novamente.</p>
          )}

          <div className="flex justify-end gap-3">
            <Button type="button" variant="outline" onClick={onCancel}>
              Cancelar
            </Button>
            <Button type="submit" disabled={createDiet.isPending}>
              {createDiet.isPending ? "Criando..." : "Criar Dieta"}
            </Button>
          </div>
        </form>
      </CardContent>
    </Card>
  )
}
