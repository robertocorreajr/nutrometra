import type { UUID } from "./common"

export interface Role {
  id: UUID
  code: string
  name: string
  description?: string
  application_scope: "tenant" | "backoffice"
}

export interface AssignRoleRequest {
  role_code: string
}
