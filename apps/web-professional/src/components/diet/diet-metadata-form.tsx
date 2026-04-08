"use client"

import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { Button, Input, Card, CardContent } from "@nutrometra/ui"
import { useUpdateDiet } from "@nutrometra/api-client/hooks"
import { createDietSchema, type CreateDietFormValues } from "@/lib/schemas/diet"
import type { Diet } from "@nutrometra/api-client"

interface DietMetadataFormProps {
  diet: Diet
  patientId: string
}

export function DietMetadataForm({ diet, patientId }: DietMetadataFormProps) {
  const updateDiet = useUpdateDiet(diet.id, patientId)

  const {
    register,
    handleSubmit,
    formState: { errors, isDirty },
  } = useForm<CreateDietFormValues>({
    resolver: zodResolver(createDietSchema),
    defaultValues: {
      title: diet.title,
      objective: diet.objective ?? "",
      valid_from: diet.valid_from?.slice(0, 10) ?? "",
      valid_until: diet.valid_until?.slice(0, 10) ?? "",
    },
  })

  if (diet.status !== "draft") return null

  async function onSubmit(data: CreateDietFormValues) {
    try {
      await updateDiet.mutateAsync({
        title: data.title,
        objective: data.objective || undefined,
        valid_from: data.valid_from || undefined,
        valid_until: data.valid_until || undefined,
      })
    } catch {
      // Error handled by mutation state
    }
  }

  return (
    <Card>
      <CardContent className="p-4 sm:p-6">
        <h3 className="text-sm font-semibold mb-3">Dados da Dieta</h3>
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <div>
            <label htmlFor="meta-title" className="block text-sm font-medium mb-1">
              Titulo *
            </label>
            <Input id="meta-title" {...register("title")} />
            {errors.title && (
              <p className="text-sm text-destructive mt-1">{errors.title.message}</p>
            )}
          </div>

          <div>
            <label htmlFor="meta-objective" className="block text-sm font-medium mb-1">
              Objetivo
            </label>
            <textarea
              id="meta-objective"
              rows={2}
              className="flex w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
              {...register("objective")}
            />
          </div>

          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div>
              <label htmlFor="meta-valid-from" className="block text-sm font-medium mb-1">
                Inicio
              </label>
              <Input id="meta-valid-from" type="date" {...register("valid_from")} />
            </div>
            <div>
              <label htmlFor="meta-valid-until" className="block text-sm font-medium mb-1">
                Fim
              </label>
              <Input id="meta-valid-until" type="date" {...register("valid_until")} />
            </div>
          </div>

          {updateDiet.isError && (
            <p className="text-sm text-destructive">Erro ao atualizar. Tente novamente.</p>
          )}

          {updateDiet.isSuccess && (
            <p className="text-sm text-green-600">Salvo com sucesso.</p>
          )}

          <div className="flex justify-end">
            <Button type="submit" size="sm" disabled={!isDirty || updateDiet.isPending}>
              {updateDiet.isPending ? "Salvando..." : "Salvar"}
            </Button>
          </div>
        </form>
      </CardContent>
    </Card>
  )
}
