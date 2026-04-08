"use client"

import { useState } from "react"
import { format } from "date-fns"
import { ptBR } from "date-fns/locale"
import { Button } from "@nutrometra/ui"
import { useSubscription, usePlans, useCancelSubscription } from "@nutrometra/api-client/hooks"
import { ConfirmDialog } from "@/components/shared/confirm-dialog"
import {
  subscriptionStatusLabels,
  subscriptionStatusColors,
} from "@/lib/schemas/billing"

export function SubscriptionStatus() {
  const { data: subscription } = useSubscription()
  const { data: plans } = usePlans()
  const cancelMutation = useCancelSubscription()
  const [showCancel, setShowCancel] = useState(false)

  if (!subscription) return null

  const plan = (plans ?? []).find((p) => p.id === subscription.plan_id)
  const statusLabel = subscriptionStatusLabels[subscription.status] ?? subscription.status
  const statusColor = subscriptionStatusColors[subscription.status] ?? "bg-gray-100 text-gray-800"

  return (
    <div className="rounded-lg border bg-card p-6">
      <div className="flex items-center justify-between mb-4">
        <h2 className="text-lg font-semibold">Plano Atual</h2>
        <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${statusColor}`}>
          {statusLabel}
        </span>
      </div>

      <div className="space-y-2 text-sm">
        <p><span className="text-muted-foreground">Plano:</span> {plan?.name ?? "—"}</p>
        <p><span className="text-muted-foreground">Início:</span> {format(new Date(subscription.started_at), "dd/MM/yyyy", { locale: ptBR })}</p>
        {subscription.status === "trialing" && subscription.trial_ends_at && (
          <p><span className="text-muted-foreground">Trial até:</span> {format(new Date(subscription.trial_ends_at), "dd/MM/yyyy", { locale: ptBR })}</p>
        )}
        {subscription.status === "active" && subscription.renews_at && (
          <p><span className="text-muted-foreground">Renova em:</span> {format(new Date(subscription.renews_at), "dd/MM/yyyy", { locale: ptBR })}</p>
        )}
      </div>

      {(subscription.status === "trialing" || subscription.status === "active") && (
        <div className="mt-4">
          <Button variant="destructive" size="sm" onClick={() => setShowCancel(true)}>
            Cancelar Assinatura
          </Button>
        </div>
      )}

      <ConfirmDialog
        open={showCancel}
        onClose={() => setShowCancel(false)}
        onConfirm={() => {
          cancelMutation.mutate(undefined, { onSuccess: () => setShowCancel(false) })
        }}
        title="Cancelar Assinatura"
        description="Tem certeza que deseja cancelar sua assinatura? Você perderá acesso aos recursos do plano atual."
        confirmLabel="Sim, Cancelar"
        variant="destructive"
        isLoading={cancelMutation.isPending}
      />
    </div>
  )
}
