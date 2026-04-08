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

export interface UpdateProfessionalRequest {
  user_id: string
  full_name: string
  registration_type: string
  registration_number: string
  registration_state?: string
  specialty?: string
  bio?: string
  phone?: string
  avatar_url?: string
}

export interface ProfessionalAddress {
  id: UUID
  tenant_id: UUID
  professional_id: UUID
  label: string
  street: string
  number?: string
  complement?: string
  neighborhood?: string
  city: string
  state: string
  zip_code: string
  country: string
  latitude?: number
  longitude?: number
  phone?: string
  notes?: string
  active: boolean
  created_at: string
  updated_at: string
}

export interface CreateAddressRequest {
  label: string
  street: string
  number?: string
  complement?: string
  neighborhood?: string
  city: string
  state: string
  zip_code: string
  country?: string
  phone?: string
  notes?: string
}

export type ServiceModeType = "onsite" | "online" | "home_visit"

export interface ProfessionalServiceMode {
  id: UUID
  tenant_id: UUID
  professional_id: UUID
  mode: ServiceModeType
  address_id?: string
  duration_min: number
  active: boolean
  created_at: string
  updated_at: string
}

export interface SetServiceModeRequest {
  mode: ServiceModeType
  address_id?: string
  duration_min: number
}
