"use client"

import { useState } from "react"
import { Button } from "@nutrometra/ui"
import { useBackofficeDeleteOverride } from "@nutrometra/api-client/hooks"
import type { FeatureOverride } from "@nutrometra/api-client/types"
import { ConfirmDialog } from "@/components/shared/confirm-dialog"
import { Pencil, Trash2 } from "lucide-react"

interface OverrideListProps {
  tenantId: string
  overrides: FeatureOverride[]
  onEdit: (override: FeatureOverride) => void
}

export function OverrideList({ tenantId, overrides, onEdit }: OverrideListProps) {
  const deleteOverride = useBackofficeDeleteOverride(tenantId)
  const [deleteTarget, setDeleteTarget] = useState<string | null>(null)

  return (
    <>
      <div className="rounded-md border">
        <table className="w-full text-sm">
          <thead>
            <tr className="border-b bg-muted/50">
              <th className="px-4 py-3 text-left font-medium text-muted-foreground">Feature</th>
              <th className="px-4 py-3 text-left font-medium text-muted-foreground">Status</th>
              <th className="px-4 py-3 text-left font-medium text-muted-foreground">Limite</th>
              <th className="px-4 py-3 text-left font-medium text-muted-foreground">Motivo</th>
              <th className="px-4 py-3 text-right font-medium text-muted-foreground">Acoes</th>
            </tr>
          </thead>
          <tbody>
            {overrides.map((o) => (
              <tr key={o.id} className="border-b last:border-0">
                <td className="px-4 py-3 font-mono text-xs">{o.feature_key}</td>
                <td className="px-4 py-3">
                  <span
                    className={`inline-flex items-center rounded-full px-2 py-1 text-xs font-medium ${
                      o.enabled ? "bg-green-100 text-green-700" : "bg-red-100 text-red-700"
                    }`}
                  >
                    {o.enabled ? "Ativo" : "Inativo"}
                  </span>
                </td>
                <td className="px-4 py-3 text-sm">{o.limit ?? "Ilimitado"}</td>
                <td className="px-4 py-3 text-sm text-muted-foreground">{o.reason ?? "—"}</td>
                <td className="px-4 py-3 text-right">
                  <div className="flex justify-end gap-1">
                    <Button variant="ghost" size="icon" onClick={() => onEdit(o)}>
                      <Pencil className="h-4 w-4" />
                    </Button>
                    <Button
                      variant="ghost"
                      size="icon"
                      onClick={() => setDeleteTarget(o.feature_key)}
                    >
                      <Trash2 className="h-4 w-4 text-destructive" />
                    </Button>
                  </div>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      <ConfirmDialog
        open={!!deleteTarget}
        onClose={() => setDeleteTarget(null)}
        onConfirm={() => {
          if (deleteTarget) {
            deleteOverride.mutate(deleteTarget, { onSuccess: () => setDeleteTarget(null) })
          }
        }}
        title="Remover Override"
        description={`Tem certeza que deseja remover o override "${deleteTarget}"?`}
        confirmLabel="Remover"
        variant="destructive"
        isLoading={deleteOverride.isPending}
      />
    </>
  )
}
