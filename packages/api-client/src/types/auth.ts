import type { UUID } from "./common"

export interface User {
  id: UUID
  external_id: string
  email: string
  display_name: string
  created_at: string
}

export interface TenantMembership {
  tenant_id: UUID
  tenant_name: string
  roles: string[]
}
