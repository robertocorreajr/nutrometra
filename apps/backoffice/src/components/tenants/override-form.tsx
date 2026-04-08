"use client"

import { useState } from "react"
import { Button, Input, Card, CardContent, CardHeader, CardTitle } from "@nutrometra/ui"
import {
  useBackofficeCreateOverride,
  useBackofficeUpdateOverride,
} from "@nutrometra/api-client/hooks"
import type { FeatureOverride } from "@nutrometra/api-client/types"

interface OverrideFormProps {
  tenantId: string
  existing?: FeatureOverride
  onDone: () => void
}

export function OverrideForm({ tenantId, existing, onDone }: OverrideFormProps) {
  const [featureKey, setFeatureKey] = useState(existing?.feature_key ?? "")
  const [enabled, setEnabled] = useState(existing?.enabled ?? true)
  const [limit, setLimit] = useState(existing?.limit?.toString() ?? "")
  const [reason, setReason] = useState(existing?.reason ?? "")

  const createOverride = useBackofficeCreateOverride(tenantId)
  const updateOverride = useBackofficeUpdateOverride(tenantId)
  const isPending = createOverride.isPending || updateOverride.isPending

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    const limitValue = limit ? parseInt(limit, 10) : undefined
    if (existing) {
      updateOverride.mutate(
        { featureKey: existing.feature_key, enabled, limit: limitValue, reason: reason || undefined },
        { onSuccess: onDone }
      )
    } else {
      createOverride.mutate(
        { feature_key: featureKey, enabled, limit: limitValue, reason: reason || undefined },
        { onSuccess: onDone }
      )
    }
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>{existing ? "Editar Override" : "Novo Override"}</CardTitle>
      </CardHeader>
      <CardContent>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="text-sm font-medium">Feature Key</label>
            <Input
              className="mt-1"
              value={featureKey}
              onChange={(e) => setFeatureKey(e.target.value)}
              disabled={!!existing}
              placeholder="ex: professionals:create"
              required
            />
          </div>
          <div className="flex items-center gap-2">
            <input
              type="checkbox"
              id="enabled"
              checked={enabled}
              onChange={(e) => setEnabled(e.target.checked)}
              className="rounded"
            />
            <label htmlFor="enabled" className="text-sm font-medium">
              Habilitado
            </label>
          </div>
          <div>
            <label className="text-sm font-medium">Limite (opcional)</label>
            <Input
              className="mt-1"
              type="number"
              value={limit}
              onChange={(e) => setLimit(e.target.value)}
              placeholder="Sem limite"
            />
          </div>
          <div>
            <label className="text-sm font-medium">Motivo</label>
            <Input
              className="mt-1"
              value={reason}
              onChange={(e) => setReason(e.target.value)}
              placeholder="Motivo do override"
            />
          </div>
          <div className="flex gap-2">
            <Button type="submit" disabled={isPending || (!existing && !featureKey)}>
              {isPending ? "Salvando..." : existing ? "Salvar" : "Criar"}
            </Button>
            <Button type="button" variant="outline" onClick={onDone}>
              Cancelar
            </Button>
          </div>
        </form>
      </CardContent>
    </Card>
  )
}
