import { Card, CardContent, CardHeader, CardTitle } from "@nutrometra/ui"
import type { BodyMeasurement } from "@nutrometra/api-client/types"
import { format, parseISO } from "date-fns"

interface LatestMeasurementProps {
  measurement: BodyMeasurement
}

export function LatestMeasurement({ measurement }: LatestMeasurementProps) {
  const fields = [
    { label: "Peso", value: measurement.weight_kg, unit: "kg" },
    { label: "Altura", value: measurement.height_cm, unit: "cm" },
    { label: "IMC", value: measurement.bmi, unit: "" },
    { label: "Gordura Corporal", value: measurement.body_fat_pct, unit: "%" },
    { label: "Massa Magra", value: measurement.lean_mass_kg, unit: "kg" },
    { label: "Massa Muscular", value: measurement.muscle_mass_kg, unit: "kg" },
    { label: "Agua Corporal", value: measurement.water_pct, unit: "%" },
    { label: "Gordura Visceral", value: measurement.visceral_fat, unit: "" },
    { label: "Taxa Metab. Basal", value: measurement.basal_metabolic_rate, unit: "kcal" },
    { label: "Cintura", value: measurement.waist_cm, unit: "cm" },
    { label: "Quadril", value: measurement.hip_cm, unit: "cm" },
  ].filter((f) => f.value != null)

  return (
    <Card>
      <CardHeader>
        <CardTitle className="text-base">Ultima Avaliacao</CardTitle>
        <p className="text-sm text-muted-foreground">{format(parseISO(measurement.measured_at), "dd/MM/yyyy")}</p>
      </CardHeader>
      <CardContent>
        <div className="grid grid-cols-2 gap-3">
          {fields.map((f) => (
            <div key={f.label}>
              <p className="text-xs text-muted-foreground">{f.label}</p>
              <p className="text-lg font-semibold">{f.value}{f.unit && <span className="text-sm text-muted-foreground ml-1">{f.unit}</span>}</p>
            </div>
          ))}
        </div>
      </CardContent>
    </Card>
  )
}
