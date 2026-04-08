import type { UseFormRegister } from "react-hook-form"
import type { AnamnesisFormValues } from "@/lib/schemas/anamnesis"

interface Props {
  register: UseFormRegister<AnamnesisFormValues>
  disabled?: boolean
}

export function StepAllergies({ register, disabled }: Props) {
  return (
    <div className="space-y-6">
      <div className="space-y-2">
        <label htmlFor="supplements" className="text-sm font-medium">Suplementos e Alergias</label>
        <p className="text-xs text-muted-foreground">Descreva suplementos em uso e alergias alimentares ou medicamentosas conhecidas. Alergias tambem podem ser registradas no perfil do paciente.</p>
        <textarea id="supplements" {...register("supplements")} disabled={disabled} rows={6} placeholder="Suplementos em uso (tipo, dose, frequencia), alergias alimentares, intolerancias, reacoes adversas a medicamentos..." className="flex w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50" />
      </div>
    </div>
  )
}
