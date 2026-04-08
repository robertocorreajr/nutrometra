"use client"

import { useParams } from "next/navigation"
import { Card, CardContent, CardHeader, CardTitle, LoadingState, ErrorState } from "@nutrometra/ui"
import { useBackofficeTenantSubscription, usePlans } from "@nutrometra/api-client/hooks"
import {
  subscriptionStatusLabels,
  subscriptionStatusColors,
  formatCurrency,
} from "@/lib/schemas/backoffice"
import { SubscriptionActions } from "@/components/tenants/subscription-actions"
import { format, parseISO } from "date-fns"

export default function AssinaturaPage() {
  const { id } = useParams<{ id: string }>()
  const { data: subscription, isLoading, isError, refetch } = useBackofficeTenantSubscription(id)
  const { data: plans } = usePlans()

  if (isLoading) return <LoadingState />
  if (isError) return <ErrorState message="Nao foi possivel carregar a assinatura." onRetry={refetch} />
  if (!subscription) return <ErrorState message="Nenhuma assinatura encontrada." />

  const plan = plans?.find((p) => p.id === subscription.plan_id)

  return (
    <div className="space-y-6">
      <Card>
        <CardHeader>
          <CardTitle>Assinatura</CardTitle>
        </CardHeader>
        <CardContent className="space-y-3">
          <div className="flex items-center gap-2">
            <span className="text-sm text-muted-foreground">Status:</span>
            <span
              className={`inline-flex items-center rounded-full px-2 py-1 text-xs font-medium ${
                subscriptionStatusColors[subscription.status] ?? "bg-gray-100 text-gray-700"
              }`}
            >
              {subscriptionStatusLabels[subscription.status] ?? subscription.status}
            </span>
          </div>
          {plan && (
            <div className="flex items-center gap-2">
              <span className="text-sm text-muted-foreground">Plano:</span>
              <span className="text-sm font-medium">
                {plan.name} — {formatCurrency(plan.price_cents, plan.currency)}/{plan.billing_cycle}
              </span>
            </div>
          )}
          <div className="flex items-center gap-2">
            <span className="text-sm text-muted-foreground">Inicio:</span>
            <span className="text-sm">{format(parseISO(subscription.started_at), "dd/MM/yyyy")}</span>
          </div>
          {subscription.trial_ends_at && (
            <div className="flex items-center gap-2">
              <span className="text-sm text-muted-foreground">Trial ate:</span>
              <span className="text-sm">
                {format(parseISO(subscription.trial_ends_at), "dd/MM/yyyy")}
              </span>
            </div>
          )}
          {subscription.renews_at && (
            <div className="flex items-center gap-2">
              <span className="text-sm text-muted-foreground">Renovacao:</span>
              <span className="text-sm">
                {format(parseISO(subscription.renews_at), "dd/MM/yyyy")}
              </span>
            </div>
          )}
        </CardContent>
      </Card>
      <SubscriptionActions tenantId={id} subscription={subscription} />
    </div>
  )
}
