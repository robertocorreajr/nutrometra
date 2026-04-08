import type { UUID } from "./common"

export interface Plan {
  id: UUID
  code: string
  name: string
  billing_cycle: string
  currency: string
  price_cents: number
}

export interface Subscription {
  id: UUID
  tenant_id: UUID
  plan_id: UUID
  status: "trialing" | "active" | "past_due" | "cancelled" | "expired"
  started_at: string
  trial_ends_at?: string
  renews_at?: string
}

export interface Entitlement {
  feature_key: string
  enabled: boolean
  limit?: number
  source: "override" | "plan" | "default"
}

export type InvoiceStatus = "draft" | "open" | "paid" | "void" | "uncollectible"

export interface Invoice {
  ID: string
  TenantID: string
  SubscriptionID: string
  ProviderInvoiceID?: string
  Status: InvoiceStatus
  AmountCents: number
  Currency: string
  HostedURL?: string
  PeriodStart: string
  PeriodEnd: string
  DueDate?: string
  PaidAt?: string
  CreatedAt: string
  UpdatedAt: string
}

export type PaymentStatus = "pending" | "succeeded" | "failed" | "refunded"

export interface Payment {
  ID: string
  TenantID: string
  InvoiceID: string
  ProviderPaymentID?: string
  Status: PaymentStatus
  AmountCents: number
  Currency: string
  FailureReason?: string
  PaidAt?: string
  CreatedAt: string
}

export interface CheckoutRequest {
  plan_id: string
  email: string
  name: string
}

export interface ChangePlanRequest {
  plan_id: string
}
