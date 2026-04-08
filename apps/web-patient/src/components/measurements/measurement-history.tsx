import { Card, CardContent, CardHeader, CardTitle } from "@nutrometra/ui"
import type { BodyMeasurement } from "@nutrometra/api-client/types"
import { format, parseISO } from "date-fns"

interface MeasurementHistoryProps {
  measurements: BodyMeasurement[]
}

export function MeasurementHistory({ measurements }: MeasurementHistoryProps) {
  const sorted = [...measurements].sort((a, b) => b.measured_at.localeCompare(a.measured_at))

  return (
    <Card>
      <CardHeader>
        <CardTitle className="text-base">Historico</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="space-y-3">
          {sorted.map((m) => (
            <div key={m.id} className="flex items-center justify-between text-sm border-b last:border-0 pb-2 last:pb-0">
              <span className="text-muted-foreground">{format(parseISO(m.measured_at), "dd/MM/yyyy")}</span>
              <div className="flex gap-4">
                {m.weight_kg != null && <span>{m.weight_kg} kg</span>}
                {m.body_fat_pct != null && <span>{m.body_fat_pct}%</span>}
                {m.bmi != null && <span>IMC {m.bmi}</span>}
              </div>
            </div>
          ))}
        </div>
      </CardContent>
    </Card>
  )
}
