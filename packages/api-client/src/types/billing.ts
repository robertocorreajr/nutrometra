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
