import type { UUID } from "./common"

export interface BackofficeTenant {
  id: UUID
  tenant_name: string
  slug: string
  status: string
  owner_email: string
  created_at: string
  updated_at: string
}

export interface FeatureOverride {
  id: UUID
  tenant_id: UUID
  feature_key: string
  enabled: boolean
  limit?: number
  reason?: string
  created_by: string
  created_at: string
  updated_at: string
}

export interface CreateOverrideRequest {
  feature_key: string
  enabled: boolean
  limit?: number
  reason?: string
}

export interface UpdateOverrideRequest {
  enabled: boolean
  limit?: number
  reason?: string
}

export interface OverridePlanRequest {
  plan_id: string
  reason?: string
}

export interface AuditLogEntry {
  id: UUID
  tenant_id: UUID
  actor_id: string
  actor_role: string
  action: string
  resource_type: string
  resource_id: string
  details?: Record<string, unknown>
  ip_address?: string
  created_at: string
}
