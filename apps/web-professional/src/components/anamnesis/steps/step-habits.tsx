import type { UseFormRegister } from "react-hook-form"
import type { AnamnesisFormValues } from "@/lib/schemas/anamnesis"

interface Props {
  register: UseFormRegister<AnamnesisFormValues>
  disabled?: boolean
}

export function StepHabits({ register, disabled }: Props) {
  const fields = [
    { name: "social_history" as const, label: "Historico Social", placeholder: "Tabagismo, etilismo, uso de substancias, rede de suporte..." },
    { name: "dietary_history" as const, label: "Historico Alimentar", placeholder: "Padrao alimentar atual, refeicoes por dia, preferencias, restricoes..." },
    { name: "physical_activity" as const, label: "Atividade Fisica", placeholder: "Tipo, frequencia, duracao, intensidade..." },
    { name: "sleep_pattern" as const, label: "Padrao de Sono", placeholder: "Horas de sono, qualidade, horarios, uso de medicamentos para dormir..." },
    { name: "bowel_habits" as const, label: "Habito Intestinal", placeholder: "Frequencia, consistencia (Escala de Bristol), desconfortos..." },
    { name: "water_intake" as const, label: "Ingestao Hidrica", placeholder: "Quantidade diaria estimada, tipos de liquidos consumidos..." },
  ]

  return (
    <div className="space-y-6">
      {fields.map((field) => (
        <div key={field.name} className="space-y-2">
          <label htmlFor={field.name} className="text-sm font-medium">{field.label}</label>
          <textarea id={field.name} {...register(field.name)} disabled={disabled} rows={3} placeholder={field.placeholder} className="flex w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50" />
        </div>
      ))}
    </div>
  )
}
