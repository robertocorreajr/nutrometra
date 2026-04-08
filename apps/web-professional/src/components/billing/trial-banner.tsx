"use client"

import { differenceInDays } from "date-fns"

interface TrialBannerProps {
  trialEndsAt: string
}

export function TrialBanner({ trialEndsAt }: TrialBannerProps) {
  const daysLeft = differenceInDays(new Date(trialEndsAt), new Date())
  if (daysLeft < 0) return null

  return (
    <div className="rounded-lg border border-blue-200 bg-blue-50 px-4 py-3 text-sm text-blue-800">
      <p className="font-medium">
        Período de avaliação: {daysLeft === 0 ? "último dia" : `${daysLeft} dia${daysLeft !== 1 ? "s" : ""} restante${daysLeft !== 1 ? "s" : ""}`}
      </p>
      <p className="text-blue-600 mt-1">
        Escolha um plano para continuar usando todos os recursos.
      </p>
    </div>
  )
}
