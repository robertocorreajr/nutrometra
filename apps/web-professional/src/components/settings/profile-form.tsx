"use client"

import { useState } from "react"
import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { Button, Input, Card, CardContent, CardHeader, CardTitle } from "@nutrometra/ui"
import { useUpdateProfessional } from "@nutrometra/api-client/hooks"
import type { Professional } from "@nutrometra/api-client"
import {
  profileSchema,
  type ProfileFormValues,
  registrationTypes,
  brazilianStates,
} from "@/lib/schemas/professional"

interface ProfileFormProps {
  professional: Professional
}

export function ProfileForm({ professional }: ProfileFormProps) {
  const [successMessage, setSuccessMessage] = useState<string | null>(null)
  const updateMutation = useUpdateProfessional(professional.id)

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<ProfileFormValues>({
    resolver: zodResolver(profileSchema),
    defaultValues: {
      full_name: professional.full_name,
      registration_type: professional.registration_type,
      registration_number: professional.registration_number,
      registration_state: professional.registration_state ?? "",
      specialty: professional.specialty ?? "",
      bio: professional.bio ?? "",
      phone: professional.phone ?? "",
    },
  })

  function onSubmit(data: ProfileFormValues) {
    setSuccessMessage(null)
    updateMutation.mutate(
      {
        user_id: professional.user_id,
        full_name: data.full_name,
        registration_type: data.registration_type,
        registration_number: data.registration_number,
        registration_state: data.registration_state || undefined,
        specialty: data.specialty || undefined,
        bio: data.bio || undefined,
        phone: data.phone || undefined,
      },
      {
        onSuccess: () => {
          setSuccessMessage("Perfil atualizado com sucesso.")
        },
      },
    )
  }

  const selectClasses =
    "flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"

  return (
    <Card>
      <CardHeader>
        <CardTitle>Dados do Profissional</CardTitle>
      </CardHeader>
      <CardContent>
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <div className="space-y-1">
            <label htmlFor="full_name" className="text-sm font-medium">
              Nome completo
            </label>
            <Input id="full_name" {...register("full_name")} placeholder="Seu nome completo" />
            {errors.full_name && (
              <p className="text-sm text-destructive">{errors.full_name.message}</p>
            )}
          </div>

          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div className="space-y-1">
              <label htmlFor="registration_type" className="text-sm font-medium">
                Tipo de registro
              </label>
              <select id="registration_type" {...register("registration_type")} className={selectClasses}>
                <option value="">Selecione...</option>
                {registrationTypes.map((rt) => (
                  <option key={rt.value} value={rt.value}>
                    {rt.label}
                  </option>
                ))}
              </select>
              {errors.registration_type && (
                <p className="text-sm text-destructive">{errors.registration_type.message}</p>
              )}
            </div>

            <div className="space-y-1">
              <label htmlFor="registration_number" className="text-sm font-medium">
                Número de registro
              </label>
              <Input
                id="registration_number"
                {...register("registration_number")}
                placeholder="Ex: 12345"
              />
              {errors.registration_number && (
                <p className="text-sm text-destructive">{errors.registration_number.message}</p>
              )}
            </div>
          </div>

          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div className="space-y-1">
              <label htmlFor="registration_state" className="text-sm font-medium">
                Estado do registro
              </label>
              <select id="registration_state" {...register("registration_state")} className={selectClasses}>
                <option value="">Selecione...</option>
                {brazilianStates.map((uf) => (
                  <option key={uf} value={uf}>
                    {uf}
                  </option>
                ))}
              </select>
              {errors.registration_state && (
                <p className="text-sm text-destructive">{errors.registration_state.message}</p>
              )}
            </div>

            <div className="space-y-1">
              <label htmlFor="specialty" className="text-sm font-medium">
                Especialidade
              </label>
              <Input
                id="specialty"
                {...register("specialty")}
                placeholder="Ex: Nutrição Esportiva"
              />
              {errors.specialty && (
                <p className="text-sm text-destructive">{errors.specialty.message}</p>
              )}
            </div>
          </div>

          <div className="space-y-1">
            <label htmlFor="phone" className="text-sm font-medium">
              Telefone
            </label>
            <Input id="phone" {...register("phone")} placeholder="(11) 99999-9999" />
            {errors.phone && (
              <p className="text-sm text-destructive">{errors.phone.message}</p>
            )}
          </div>

          <div className="space-y-1">
            <label htmlFor="bio" className="text-sm font-medium">
              Bio
            </label>
            <textarea
              id="bio"
              {...register("bio")}
              rows={4}
              placeholder="Uma breve descrição sobre você..."
              className="flex w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
            />
            {errors.bio && (
              <p className="text-sm text-destructive">{errors.bio.message}</p>
            )}
          </div>

          {successMessage && (
            <p className="text-sm text-green-600 font-medium">{successMessage}</p>
          )}

          {updateMutation.isError && (
            <p className="text-sm text-destructive">
              Erro ao salvar perfil. Tente novamente.
            </p>
          )}

          <div className="flex justify-end">
            <Button type="submit" disabled={updateMutation.isPending}>
              {updateMutation.isPending ? "Salvando..." : "Salvar Perfil"}
            </Button>
          </div>
        </form>
      </CardContent>
    </Card>
  )
}
