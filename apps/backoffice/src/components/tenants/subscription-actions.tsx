"use client"

import { useState } from "react"
import { Button, Input } from "@nutrometra/ui"
import {
  useBackofficeOverridePlan,
  useBackofficeCancelSubscription,
  useBackofficeReactivateSubscription,
  usePlans,
} from "@nutrometra/api-client/hooks"
import { ConfirmDialog } from "@/components/shared/confirm-dialog"
import type { Subscription } from "@nutrometra/api-client/types"

interface SubscriptionActionsProps {
  tenantId: string
  subscription: Subscription
}

export function SubscriptionActions({ tenantId, subscription }: SubscriptionActionsProps) {
  const { data: plans } = usePlans()
  const overridePlan = useBackofficeOverridePlan(tenantId)
  const cancelSub = useBackofficeCancelSubscription(tenantId)
  const reactivateSub = useBackofficeReactivateSubscription(tenantId)

  const [showOverride, setShowOverride] = useState(false)
  const [showCancel, setShowCancel] = useState(false)
  const [showReactivate, setShowReactivate] = useState(false)
  const [selectedPlanId, setSelectedPlanId] = useState("")
  const [reason, setReason] = useState("")

  const isCancelled = subscription.status === "cancelled" || subscription.status === "expired"

  return (
    <div className="flex flex-wrap gap-2">
      {!isCancelled && (
        <>
          <Button variant="outline" onClick={() => setShowOverride(true)}>
            Alterar Plano
          </Button>
          <Button variant="destructive" onClick={() => setShowCancel(true)}>
            Cancelar
          </Button>
        </>
      )}
      {isCancelled && (
        <Button onClick={() => setShowReactivate(true)}>Reativar</Button>
      )}

      {/* Override Plan Dialog */}
      {showOverride && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
          <div className="absolute inset-0 bg-black/50" onClick={() => setShowOverride(false)} />
          <div className="relative z-10 w-full max-w-md rounded-lg border bg-background shadow-lg p-6">
            <h3 className="text-lg font-semibold mb-4">Alterar Plano</h3>
            <div className="space-y-3">
              <div>
                <label className="text-sm font-medium">Plano</label>
                <select
                  className="w-full mt-1 rounded-md border px-3 py-2 text-sm"
                  value={selectedPlanId}
                  onChange={(e) => setSelectedPlanId(e.target.value)}
                >
                  <option value="">Selecione...</option>
                  {plans?.map((p) => (
                    <option key={p.id} value={p.id}>
                      {p.name}
                    </option>
                  ))}
                </select>
              </div>
              <div>
                <label className="text-sm font-medium">Motivo</label>
                <Input
                  className="mt-1"
                  value={reason}
                  onChange={(e) => setReason(e.target.value)}
                  placeholder="Motivo da alteracao"
                />
              </div>
            </div>
            <div className="flex justify-end gap-3 mt-6">
              <Button variant="outline" onClick={() => setShowOverride(false)}>
                Cancelar
              </Button>
              <Button
                disabled={!selectedPlanId || overridePlan.isPending}
                onClick={() => {
                  overridePlan.mutate(
                    { plan_id: selectedPlanId, reason: reason || undefined },
                    {
                      onSuccess: () => {
                        setShowOverride(false)
                        setSelectedPlanId("")
                        setReason("")
                      },
                    }
                  )
                }}
              >
                {overridePlan.isPending ? "Aguarde..." : "Confirmar"}
              </Button>
            </div>
          </div>
        </div>
      )}

      {/* Cancel Dialog */}
      <ConfirmDialog
        open={showCancel}
        onClose={() => setShowCancel(false)}
        onConfirm={() =>
          cancelSub.mutate(undefined, { onSuccess: () => setShowCancel(false) })
        }
        title="Cancelar Assinatura"
        description="Tem certeza que deseja cancelar a assinatura deste tenant? Esta acao pode ser revertida com a reativacao."
        confirmLabel="Cancelar Assinatura"
        variant="destructive"
        isLoading={cancelSub.isPending}
      />

      {/* Reactivate Dialog */}
      <ConfirmDialog
        open={showReactivate}
        onClose={() => setShowReactivate(false)}
        onConfirm={() =>
          reactivateSub.mutate(undefined, { onSuccess: () => setShowReactivate(false) })
        }
        title="Reativar Assinatura"
        description="Deseja reativar a assinatura deste tenant?"
        confirmLabel="Reativar"
        isLoading={reactivateSub.isPending}
      />
    </div>
  )
}
