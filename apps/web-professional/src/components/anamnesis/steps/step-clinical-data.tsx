import type { UseFormRegister } from "react-hook-form"
import type { AnamnesisFormValues } from "@/lib/schemas/anamnesis"

interface Props {
  register: UseFormRegister<AnamnesisFormValues>
  disabled?: boolean
}

export function StepClinicalData({ register, disabled }: Props) {
  return (
    <div className="space-y-4">
      <div className="space-y-2">
        <label htmlFor="chief_complaint" className="text-sm font-medium">
          Queixa Principal / Motivo da Consulta
        </label>
        <textarea
          id="chief_complaint"
          {...register("chief_complaint")}
          disabled={disabled}
          rows={6}
          placeholder="Descreva a queixa principal do paciente, motivo da consulta, sintomas relatados, expectativas e objetivos..."
          className="flex w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
        />
      </div>
    </div>
  )
}
