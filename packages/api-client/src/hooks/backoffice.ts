import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query"
import { api } from "../client"
import type {
  BackofficeTenant,
  FeatureOverride,
  CreateOverrideRequest,
  UpdateOverrideRequest,
  OverridePlanRequest,
  AuditLogEntry,
} from "../types/backoffice"
import type { Subscription, Invoice, Payment } from "../types/billing"

// --- Tenants ---

export function useBackofficeTenants() {
  return useQuery({
    queryKey: ["backoffice", "tenants"],
    queryFn: () => api.get<BackofficeTenant[]>("/backoffice/tenants"),
  })
}

export function useBackofficeTenant(id: string) {
  return useQuery({
    queryKey: ["backoffice", "tenants", id],
    queryFn: () => api.get<BackofficeTenant>(`/backoffice/tenants/${id}`),
    enabled: !!id,
  })
}

// --- Subscription ---

export function useBackofficeTenantSubscription(tenantId: string) {
  return useQuery({
    queryKey: ["backoffice", "subscription", tenantId],
    queryFn: () =>
      api.get<Subscription>(`/backoffice/tenants/${tenantId}/subscription`),
    enabled: !!tenantId,
  })
}

export function useBackofficeOverridePlan(tenantId: string) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (data: OverridePlanRequest) =>
      api.patch<Subscription>(
        `/backoffice/tenants/${tenantId}/subscription/plan`,
        data,
      ),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["backoffice", "subscription", tenantId] })
    },
  })
}

export function useBackofficeCancelSubscription(tenantId: string) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: () =>
      api.post<void>(`/backoffice/tenants/${tenantId}/subscription/cancel`),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["backoffice", "subscription", tenantId] })
    },
  })
}

export function useBackofficeReactivateSubscription(tenantId: string) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: () =>
      api.post<Subscription>(
        `/backoffice/tenants/${tenantId}/subscription/reactivate`,
      ),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["backoffice", "subscription", tenantId] })
    },
  })
}

// --- Invoices & Payments ---

export function useBackofficeInvoices(tenantId: string) {
  return useQuery({
    queryKey: ["backoffice", "invoices", tenantId],
    queryFn: () =>
      api.get<Invoice[]>(`/backoffice/tenants/${tenantId}/invoices`),
    enabled: !!tenantId,
  })
}

export function useBackofficePayments(tenantId: string) {
  return useQuery({
    queryKey: ["backoffice", "payments", tenantId],
    queryFn: () =>
      api.get<Payment[]>(`/backoffice/tenants/${tenantId}/payments`),
    enabled: !!tenantId,
  })
}

// --- Feature Overrides ---

export function useBackofficeOverrides(tenantId: string) {
  return useQuery({
    queryKey: ["backoffice", "overrides", tenantId],
    queryFn: () =>
      api.get<FeatureOverride[]>(`/backoffice/tenants/${tenantId}/overrides`),
    enabled: !!tenantId,
  })
}

export function useBackofficeCreateOverride(tenantId: string) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (data: CreateOverrideRequest) =>
      api.post<FeatureOverride>(
        `/backoffice/tenants/${tenantId}/overrides`,
        data,
      ),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["backoffice", "overrides", tenantId] })
    },
  })
}

export function useBackofficeUpdateOverride(tenantId: string) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ featureKey, ...data }: { featureKey: string } & UpdateOverrideRequest) =>
      api.put<FeatureOverride>(
        `/backoffice/tenants/${tenantId}/overrides/${featureKey}`,
        data,
      ),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["backoffice", "overrides", tenantId] })
    },
  })
}

export function useBackofficeDeleteOverride(tenantId: string) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (featureKey: string) =>
      api.delete<void>(
        `/backoffice/tenants/${tenantId}/overrides/${featureKey}`,
      ),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["backoffice", "overrides", tenantId] })
    },
  })
}

// --- Audit Log ---

export function useBackofficeAuditLog(tenantId: string) {
  return useQuery({
    queryKey: ["backoffice", "audit", tenantId],
    queryFn: () =>
      api.get<AuditLogEntry[]>(`/backoffice/tenants/${tenantId}/audit`),
    enabled: !!tenantId,
  })
}
