"use client"

import { LoadingState, ErrorState } from "@nutrometra/ui"
import { usePlans, useSubscription, useProfessionalMe, useCheckout, useChangePlan } from "@nutrometra/api-client/hooks"
import { PlanCard } from "@/components/billing/plan-card"

export default function PlanosPage() {
  const { data: plans, isLoading, isError, refetch } = usePlans()
  const { data: subscription } = useSubscription()
  const { data: professional } = useProfessionalMe()
  const checkout = useCheckout()
  const changePlan = useChangePlan()

  if (isLoading) return <LoadingState lines={4} />
  if (isError) return <ErrorState message="Não foi possível carregar os planos." onRetry={refetch} />

  const planList = plans ?? []
  const hasSubscription = !!subscription

  function handleSelect(planId: string) {
    if (hasSubscription) {
      changePlan.mutate({ plan_id: planId })
    } else {
      checkout.mutate({
        plan_id: planId,
        email: professional?.phone ?? "",
        name: professional?.full_name ?? "",
      })
    }
  }

  return (
    <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
      {planList.map((plan) => (
        <PlanCard
          key={plan.id}
          name={plan.name}
          priceCents={plan.price_cents}
          currency={plan.currency}
          billingCycle={plan.billing_cycle}
          isCurrent={subscription?.plan_id === plan.id}
          onSelect={() => handleSelect(plan.id)}
          isLoading={checkout.isPending || changePlan.isPending}
        />
      ))}
    </div>
  )
}
