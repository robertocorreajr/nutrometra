import { useQuery } from "@tanstack/react-query"
import { api } from "../client"
import type { Plan, Subscription, Entitlement } from "../types/billing"

export function usePlans() {
  return useQuery({
    queryKey: ["plans"],
    queryFn: () => api.get<Plan[]>("/plans"),
  })
}

export function useSubscription() {
  return useQuery({
    queryKey: ["subscription"],
    queryFn: () => api.get<Subscription>("/subscription"),
  })
}

export function useEntitlements() {
  return useQuery({
    queryKey: ["entitlements"],
    queryFn: () => api.get<Record<string, Entitlement>>("/entitlements"),
    staleTime: 5 * 60 * 1000, // 5 min — aligned with backend cache TTL
  })
}
