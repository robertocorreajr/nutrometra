"use client"

import { useState } from "react"
import { MapPin, Video, Home, Trash2 } from "lucide-react"
import { Button, Card, CardContent, LoadingState, ErrorState } from "@nutrometra/ui"
import {
  useAvailabilityRules,
  useDeleteAvailabilityRule,
} from "@nutrometra/api-client/hooks"
import type { AvailabilityRule } from "@nutrometra/api-client"
import { dayOfWeekOptions } from "@/lib/schemas/availability"
import { serviceModeOptions } from "@/lib/schemas/appointment"

interface AvailabilityListProps {
  professionalId: string
}

const serviceModeIcons: Record<string, React.ReactNode> = {
  onsite: <MapPin className="h-3.5 w-3.5" />,
  online: <Video className="h-3.5 w-3.5" />,
  home_visit: <Home className="h-3.5 w-3.5" />,
}

function getDayLabel(dayOfWeek: number): string {
  return dayOfWeekOptions.find((d) => d.value === dayOfWeek)?.label ?? String(dayOfWeek)
}

function getServiceModeLabel(mode: string): string {
  return serviceModeOptions.find((o) => o.value === mode)?.label ?? mode
}

export function AvailabilityList({ professionalId }: AvailabilityListProps) {
  const { data: rules, isLoading, isError, refetch } = useAvailabilityRules(professionalId)
  const deleteRule = useDeleteAvailabilityRule(professionalId)
  const [confirmDeleteId, setConfirmDeleteId] = useState<string | null>(null)

  if (isLoading) return <LoadingState lines={4} />
  if (isError) return <ErrorState message="Erro ao carregar disponibilidade." onRetry={refetch} />

  const ruleList = rules ?? []

  if (ruleList.length === 0) {
    return (
      <p className="text-sm text-muted-foreground text-center py-8">
        Nenhuma regra de disponibilidade cadastrada.
      </p>
    )
  }

  // Group by day_of_week
  const grouped = ruleList.reduce<Record<number, AvailabilityRule[]>>((acc, rule) => {
    if (!acc[rule.day_of_week]) acc[rule.day_of_week] = []
    acc[rule.day_of_week].push(rule)
    return acc
  }, {})

  const sortedDays = Object.keys(grouped)
    .map(Number)
    .sort((a, b) => a - b)

  async function handleDelete(id: string) {
    try {
      await deleteRule.mutateAsync(id)
      setConfirmDeleteId(null)
    } catch {
      // Error handled by mutation
    }
  }

  return (
    <div className="space-y-4">
      {sortedDays.map((day) => (
        <Card key={day}>
          <CardContent className="pt-4 pb-4">
            <h4 className="text-sm font-semibold mb-2">{getDayLabel(day)}</h4>
            <div className="space-y-2">
              {grouped[day].map((rule) => (
                <div
                  key={rule.id}
                  className="flex items-center justify-between gap-2 rounded-md border p-2 text-sm"
                >
                  <div className="flex items-center gap-3">
                    <span className="font-medium">
                      {rule.start_time} - {rule.end_time}
                    </span>
                    <span className="flex items-center gap-1 text-muted-foreground">
                      {serviceModeIcons[rule.service_mode]}
                      <span className="text-xs">{getServiceModeLabel(rule.service_mode)}</span>
                    </span>
                  </div>
                  <div>
                    {confirmDeleteId === rule.id ? (
                      <div className="flex items-center gap-1">
                        <Button
                          variant="destructive"
                          size="sm"
                          disabled={deleteRule.isPending}
                          onClick={() => handleDelete(rule.id)}
                        >
                          {deleteRule.isPending ? "..." : "Confirmar"}
                        </Button>
                        <Button
                          variant="ghost"
                          size="sm"
                          onClick={() => setConfirmDeleteId(null)}
                        >
                          Não
                        </Button>
                      </div>
                    ) : (
                      <Button
                        variant="ghost"
                        size="icon"
                        onClick={() => setConfirmDeleteId(rule.id)}
                      >
                        <Trash2 className="h-4 w-4 text-muted-foreground" />
                      </Button>
                    )}
                  </div>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      ))}

      {deleteRule.isError && (
        <p className="text-sm text-destructive">
          Erro ao excluir regra. Tente novamente.
        </p>
      )}
    </div>
  )
}
