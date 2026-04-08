import type { UseFormRegister } from "react-hook-form"
import type { AnamnesisFormValues } from "@/lib/schemas/anamnesis"

interface Props {
  register: UseFormRegister<AnamnesisFormValues>
  disabled?: boolean
}

export function StepHistory({ register, disabled }: Props) {
  return (
    <div className="space-y-6">
      <div className="space-y-2">
        <label htmlFor="history_present_illness" className="text-sm font-medium">Historia da Doenca Atual</label>
        <textarea id="history_present_illness" {...register("history_present_illness")} disabled={disabled} rows={4} placeholder="Evolucao dos sintomas, tratamentos anteriores, exames realizados..." className="flex w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50" />
      </div>
      <div className="space-y-2">
        <label htmlFor="past_medical_history" className="text-sm font-medium">Antecedentes Pessoais</label>
        <textarea id="past_medical_history" {...register("past_medical_history")} disabled={disabled} rows={4} placeholder="Doencas previas, cirurgias, internacoes, uso de medicamentos..." className="flex w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50" />
      </div>
      <div className="space-y-2">
        <label htmlFor="family_history" className="text-sm font-medium">Historico Familiar</label>
        <textarea id="family_history" {...register("family_history")} disabled={disabled} rows={4} placeholder="Diabetes, hipertensao, obesidade, cancer, doencas cardiovasculares na familia..." className="flex w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50" />
      </div>
    </div>
  )
}
