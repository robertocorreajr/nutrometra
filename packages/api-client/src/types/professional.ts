import type { UUID } from "./common"

/** Espelha professionalResponse do handler Go (professional/handler.go) */
export interface Professional {
  id: UUID
  tenant_id: UUID
  user_id: UUID
  full_name: string
  registration_type: string
  registration_number: string
  registration_state?: string
  specialty?: string
  bio?: string
  phone?: string
  avatar_url?: string
  active: boolean
  created_at: string
  updated_at: string
}
