"use client"

import Link from "next/link"
import { Card, CardContent, CardHeader, CardTitle, Button } from "@nutrometra/ui"
import type { Patient } from "@nutrometra/api-client"
import { formatDistanceToNow } from "date-fns"
import { ptBR } from "date-fns/locale"
import { ArrowRight } from "lucide-react"

interface RecentPatientsProps {
  patients: Patient[]
}

export function RecentPatients({ patients }: RecentPatientsProps) {
  const recent = [...patients]
    .sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime())
    .slice(0, 5)

  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between">
        <CardTitle className="text-base">Pacientes Recentes</CardTitle>
        <Button variant="ghost" size="sm" asChild>
          <Link href="/pacientes">
            Ver todos <ArrowRight className="h-4 w-4 ml-1" />
          </Link>
        </Button>
      </CardHeader>
      <CardContent>
        {recent.length === 0 ? (
          <p className="text-sm text-muted-foreground">Nenhum paciente cadastrado.</p>
        ) : (
          <div className="space-y-3">
            {recent.map((patient) => (
              <Link
                key={patient.id}
                href={`/pacientes/${patient.id}`}
                className="flex items-center justify-between p-2 rounded-md hover:bg-accent transition-colors"
              >
                <div>
                  <p className="text-sm font-medium">{patient.full_name}</p>
                  <p className="text-xs text-muted-foreground">
                    {patient.email || "Sem e-mail"}
                  </p>
                </div>
                <span className="text-xs text-muted-foreground">
                  {formatDistanceToNow(new Date(patient.created_at), {
                    addSuffix: true,
                    locale: ptBR,
                  })}
                </span>
              </Link>
            ))}
          </div>
        )}
      </CardContent>
    </Card>
  )
}
