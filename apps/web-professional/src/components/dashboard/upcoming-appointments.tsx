import { Card, CardContent, CardHeader, CardTitle } from "@nutrometra/ui"
import { Calendar } from "lucide-react"

export function UpcomingAppointments() {
  return (
    <Card>
      <CardHeader>
        <CardTitle className="text-base">Próximas Consultas</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="flex flex-col items-center justify-center py-8 text-center">
          <Calendar className="h-8 w-8 text-muted-foreground mb-3" />
          <p className="text-sm text-muted-foreground">
            O módulo de agenda será habilitado em breve.
          </p>
        </div>
      </CardContent>
    </Card>
  )
}
