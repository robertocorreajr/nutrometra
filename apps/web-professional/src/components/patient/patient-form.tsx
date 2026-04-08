"use client"

import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { Button, Input } from "@nutrometra/ui"
import {
  patientSchema,
  genderOptions,
  type PatientFormValues,
} from "@/lib/schemas/patient"

interface PatientFormProps {
  defaultValues?: Partial<PatientFormValues>
  onSubmit: (values: PatientFormValues) => void | Promise<void>
  isSubmitting?: boolean
  submitLabel?: string
}

export function PatientForm({
  defaultValues,
  onSubmit,
  isSubmitting = false,
  submitLabel = "Salvar",
}: PatientFormProps) {
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<PatientFormValues>({
    resolver: zodResolver(patientSchema),
    defaultValues: {
      full_name: "",
      email: "",
      phone: "",
      cpf: "",
      date_of_birth: "",
      gender: "prefer_not_to_say",
      notes: "",
      ...defaultValues,
    },
  })

  return (
    <form onSubmit={handleSubmit(onSubmit)} noValidate className="space-y-5">
      {/* full_name */}
      <div className="space-y-1">
        <label
          htmlFor="full_name"
          className="block text-sm font-medium text-foreground"
        >
          Nome completo <span className="text-destructive">*</span>
        </label>
        <Input
          id="full_name"
          {...register("full_name")}
          placeholder="Ex.: Maria da Silva"
          aria-invalid={!!errors.full_name}
        />
        {errors.full_name && (
          <p className="text-sm text-destructive">{errors.full_name.message}</p>
        )}
      </div>

      {/* email */}
      <div className="space-y-1">
        <label
          htmlFor="email"
          className="block text-sm font-medium text-foreground"
        >
          E-mail
        </label>
        <Input
          id="email"
          type="email"
          {...register("email")}
          placeholder="exemplo@email.com"
          aria-invalid={!!errors.email}
        />
        {errors.email && (
          <p className="text-sm text-destructive">{errors.email.message}</p>
        )}
      </div>

      {/* phone */}
      <div className="space-y-1">
        <label
          htmlFor="phone"
          className="block text-sm font-medium text-foreground"
        >
          Telefone
        </label>
        <Input
          id="phone"
          type="tel"
          {...register("phone")}
          placeholder="(11) 99999-9999"
          aria-invalid={!!errors.phone}
        />
        {errors.phone && (
          <p className="text-sm text-destructive">{errors.phone.message}</p>
        )}
      </div>

      {/* cpf */}
      <div className="space-y-1">
        <label
          htmlFor="cpf"
          className="block text-sm font-medium text-foreground"
        >
          CPF
        </label>
        <Input
          id="cpf"
          {...register("cpf")}
          placeholder="Somente números (11 dígitos)"
          maxLength={11}
          aria-invalid={!!errors.cpf}
        />
        {errors.cpf && (
          <p className="text-sm text-destructive">{errors.cpf.message}</p>
        )}
      </div>

      {/* date_of_birth */}
      <div className="space-y-1">
        <label
          htmlFor="date_of_birth"
          className="block text-sm font-medium text-foreground"
        >
          Data de nascimento
        </label>
        <Input
          id="date_of_birth"
          type="date"
          {...register("date_of_birth")}
          aria-invalid={!!errors.date_of_birth}
        />
        {errors.date_of_birth && (
          <p className="text-sm text-destructive">
            {errors.date_of_birth.message}
          </p>
        )}
      </div>

      {/* gender */}
      <div className="space-y-1">
        <label
          htmlFor="gender"
          className="block text-sm font-medium text-foreground"
        >
          Gênero
        </label>
        <select
          id="gender"
          {...register("gender")}
          className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
          aria-invalid={!!errors.gender}
        >
          {genderOptions.map((opt) => (
            <option key={opt.value} value={opt.value}>
              {opt.label}
            </option>
          ))}
        </select>
        {errors.gender && (
          <p className="text-sm text-destructive">{errors.gender.message}</p>
        )}
      </div>

      {/* notes */}
      <div className="space-y-1">
        <label
          htmlFor="notes"
          className="block text-sm font-medium text-foreground"
        >
          Observações
        </label>
        <textarea
          id="notes"
          {...register("notes")}
          rows={3}
          placeholder="Informações adicionais relevantes..."
          className="flex min-h-[80px] w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
          aria-invalid={!!errors.notes}
        />
        {errors.notes && (
          <p className="text-sm text-destructive">{errors.notes.message}</p>
        )}
      </div>

      <Button type="submit" disabled={isSubmitting} className="w-full sm:w-auto">
        {isSubmitting ? "Salvando..." : submitLabel}
      </Button>
    </form>
  )
}
