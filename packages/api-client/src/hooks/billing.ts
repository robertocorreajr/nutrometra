import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query"
import { api } from "../client"
import type {
  Plan,
  Subscription,
  Entitlement,
  Invoice,
  Payment,
  CheckoutRequest,
  ChangePlanRequest,
} from "../types/billing"

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

export function useActivateTrial() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: () => api.post<Subscription>("/subscription/trial"),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["subscription"] })
      qc.invalidateQueries({ queryKey: ["entitlements"] })
    },
  })
}

export function useCheckout() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (data: CheckoutRequest) =>
      api.post<Subscription>("/subscription/checkout", data),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["subscription"] })
    },
  })
}

export function useChangePlan() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (data: ChangePlanRequest) =>
      api.patch<Subscription>("/subscription/plan", data),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["subscription"] })
      qc.invalidateQueries({ queryKey: ["entitlements"] })
    },
  })
}

export function useCancelSubscription() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: () => api.post<void>("/subscription/cancel"),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["subscription"] })
      qc.invalidateQueries({ queryKey: ["entitlements"] })
    },
  })
}

export function useInvoices() {
  return useQuery({
    queryKey: ["invoices"],
    queryFn: () => api.get<Invoice[]>("/invoices"),
  })
}

export function usePayments() {
  return useQuery({
    queryKey: ["payments"],
    queryFn: () => api.get<Payment[]>("/payments"),
  })
}
