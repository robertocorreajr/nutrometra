import { useState } from "react"
import { Card, CardContent, CardHeader, CardTitle } from "@nutrometra/ui"
import type { BodyMeasurement } from "@nutrometra/api-client"
import { format } from "date-fns"
import { ptBR } from "date-fns/locale"
import { Scale, Ruler, Activity, ChevronDown, ChevronRight } from "lucide-react"

interface MeasurementListProps {
  measurements: BodyMeasurement[]
}

function formatValue(value: number | undefined, unit: string): string | null {
  if (value === undefined || value === null) return null
  return `${value}${unit}`
}

function MetricItem({ label, value }: { label: string; value: string | null }) {
  if (!value) return null
  return (
    <div className="flex justify-between text-sm">
      <span className="text-muted-foreground">{label}</span>
      <span className="font-medium">{value}</span>
    </div>
  )
}

export function MeasurementList({ measurements }: MeasurementListProps) {
  const [expandedIds, setExpandedIds] = useState<Set<string>>(new Set())

  const sorted = [...measurements].sort(
    (a, b) => new Date(b.measured_at).getTime() - new Date(a.measured_at).getTime()
  )

  function toggleExpand(id: string) {
    setExpandedIds((prev) => {
      const next = new Set(prev)
      if (next.has(id)) {
        next.delete(id)
      } else {
        next.add(id)
      }
      return next
    })
  }

  function computeBMI(weight?: number, height?: number): string | null {
    if (!weight || !height) return null
    const heightM = height / 100
    const bmi = weight / (heightM * heightM)
    return bmi.toFixed(1)
  }

  return (
    <div className="space-y-3">
      {sorted.map((m) => {
        const isExpanded = expandedIds.has(m.id)
        const bmi = m.bmi?.toFixed(1) ?? computeBMI(m.weight_kg, m.height_cm)

        return (
          <Card key={m.id}>
            <CardHeader className="pb-2">
              <button
                type="button"
                onClick={() => toggleExpand(m.id)}
                className="flex items-center justify-between w-full text-left"
              >
                <CardTitle className="text-sm font-medium">
                  {format(new Date(m.measured_at), "dd/MM/yyyy", { locale: ptBR })}
                </CardTitle>
                <div className="flex items-center gap-3">
                  {m.weight_kg !== undefined && (
                    <span className="flex items-center gap-1 text-xs text-muted-foreground">
                      <Scale className="h-3 w-3" />
                      {m.weight_kg} kg
                    </span>
                  )}
                  {bmi && (
                    <span className="flex items-center gap-1 text-xs text-muted-foreground">
                      <Ruler className="h-3 w-3" />
                      IMC {bmi}
                    </span>
                  )}
                  {m.body_fat_pct !== undefined && (
                    <span className="flex items-center gap-1 text-xs text-muted-foreground">
                      <Activity className="h-3 w-3" />
                      {m.body_fat_pct}% GC
                    </span>
                  )}
                  {isExpanded ? (
                    <ChevronDown className="h-4 w-4 text-muted-foreground" />
                  ) : (
                    <ChevronRight className="h-4 w-4 text-muted-foreground" />
                  )}
                </div>
              </button>
            </CardHeader>

            {isExpanded && (
              <CardContent>
                <div className="grid grid-cols-1 sm:grid-cols-2 gap-x-8 gap-y-1">
                  <MetricItem label="Peso" value={formatValue(m.weight_kg, " kg")} />
                  <MetricItem label="Altura" value={formatValue(m.height_cm, " cm")} />
                  <MetricItem label="IMC" value={bmi ? `${bmi}` : null} />
                  <MetricItem label="Gordura corporal" value={formatValue(m.body_fat_pct, "%")} />
                  <MetricItem label="Massa magra" value={formatValue(m.lean_mass_kg, " kg")} />
                  <MetricItem label="Massa gorda" value={formatValue(m.fat_mass_kg, " kg")} />
                  <MetricItem label="Massa muscular" value={formatValue(m.muscle_mass_kg, " kg")} />
                  <MetricItem label="Massa ossea" value={formatValue(m.bone_mass_kg, " kg")} />
                  <MetricItem label="Agua corporal" value={formatValue(m.water_pct, "%")} />
                  <MetricItem label="Gordura visceral" value={formatValue(m.visceral_fat, "")} />
                  <MetricItem label="Taxa metabolica basal" value={formatValue(m.basal_metabolic_rate, " kcal")} />
                  <MetricItem label="Cintura" value={formatValue(m.waist_cm, " cm")} />
                  <MetricItem label="Quadril" value={formatValue(m.hip_cm, " cm")} />
                </div>

                <div className="mt-3 flex items-center gap-4 text-xs text-muted-foreground">
                  <span>Origem: {m.source === "manual" ? "Manual" : m.source === "device" ? "Dispositivo" : "Importacao"}</span>
                  {m.device_model && <span>Dispositivo: {m.device_model}</span>}
                </div>

                {m.notes && (
                  <p className="mt-2 text-sm text-muted-foreground whitespace-pre-wrap">
                    {m.notes}
                  </p>
                )}
              </CardContent>
            )}
          </Card>
        )
      })}
    </div>
  )
}
