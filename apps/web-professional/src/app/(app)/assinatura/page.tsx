"use client"

import { Button, LoadingState, EmptyState } from "@nutrometra/ui"
import { useSubscription, useActivateTrial } from "@nutrometra/api-client/hooks"
import { SubscriptionStatus } from "@/components/billing/subscription-status"
import { TrialBanner } from "@/components/billing/trial-banner"

export default function AssinaturaPage() {
  const { data: subscription, isLoading, isError } = useSubscription()
  const activateTrial = useActivateTrial()

  if (isLoading) return <LoadingState lines={4} />

  if (isError || !subscription) {
    return (
      <div className="space-y-4">
        <EmptyState
          title="Sem assinatura ativa"
          description="Inicie um período de avaliação gratuito para explorar todos os recursos."
        />
        <div className="flex justify-center">
          <Button
            onClick={() => activateTrial.mutate()}
            disabled={activateTrial.isPending}
          >
            {activateTrial.isPending ? "Ativando..." : "Iniciar Trial Gratuito"}
          </Button>
        </div>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      {subscription.status === "trialing" && subscription.trial_ends_at && (
        <TrialBanner trialEndsAt={subscription.trial_ends_at} />
      )}
      <SubscriptionStatus />
    </div>
  )
}
