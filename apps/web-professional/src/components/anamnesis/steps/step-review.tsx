import type { UseFormGetValues, UseFormRegister } from "react-hook-form"
import type { AnamnesisFormValues } from "@/lib/schemas/anamnesis"

interface Props {
  getValues: UseFormGetValues<AnamnesisFormValues>
  register: UseFormRegister<AnamnesisFormValues>
  disabled?: boolean
}

const fieldLabels: Record<keyof AnamnesisFormValues, string> = {
  chief_complaint: "Queixa Principal",
  history_present_illness: "História da Doença Atual",
  past_medical_history: "Antecedentes Pessoais",
  family_history: "Histórico Familiar",
  social_history: "Histórico Social",
  dietary_history: "Histórico Alimentar",
  physical_activity: "Atividade Física",
  sleep_pattern: "Padrão de Sono",
  bowel_habits: "Hábito Intestinal",
  water_intake: "Ingestão Hídrica",
  supplements: "Suplementos e Alergias",
  observations: "Observações Gerais",
}

export function StepReview({ getValues, register, disabled }: Props) {
  const values = getValues()
  const filledFields = Object.entries(values).filter(([_, v]) => v && v.trim().length > 0)

  return (
    <div className="space-y-6">
      <div className="rounded-md border p-4 bg-muted/30">
        <h3 className="text-sm font-medium mb-3">Resumo do preenchimento</h3>
        <div className="grid gap-2 sm:grid-cols-2">
          {Object.entries(fieldLabels).map(([key, label]) => {
            const value = values[key as keyof AnamnesisFormValues]
            const filled = value && value.trim().length > 0
            return (
              <div key={key} className="flex items-center gap-2 text-sm">
                <span className={`h-2 w-2 rounded-full ${filled ? "bg-green-500" : "bg-muted-foreground/30"}`} />
                <span className={filled ? "" : "text-muted-foreground"}>{label}</span>
              </div>
            )
          })}
        </div>
      </div>

      <div className="space-y-2">
        <label htmlFor="observations" className="text-sm font-medium">Observações Gerais</label>
        <textarea id="observations" {...register("observations")} disabled={disabled} rows={4} placeholder="Observações adicionais, impressão clínica, plano terapêutico inicial..." className="flex w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50" />
      </div>

      {filledFields.length > 0 && (
        <div className="space-y-4">
          <h3 className="text-sm font-medium">Detalhes preenchidos</h3>
          {filledFields.filter(([key]) => key !== "observations").map(([key, value]) => (
            <div key={key} className="space-y-1">
              <p className="text-xs font-medium text-muted-foreground uppercase tracking-wide">{fieldLabels[key as keyof AnamnesisFormValues]}</p>
              <p className="text-sm whitespace-pre-wrap border-l-2 border-muted pl-3">{value as string}</p>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
