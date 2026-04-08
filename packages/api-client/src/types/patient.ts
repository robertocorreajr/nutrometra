import type { UUID } from "./common"

/** Mirrors patientResponse from Go handler (patient/handler.go) */
export interface Patient {
  id: UUID
  tenant_id: UUID
  professional_id: UUID
  full_name: string
  email?: string
  phone?: string
  cpf?: string
  date_of_birth?: string // YYYY-MM-DD
  gender?: string
  notes?: string
  active: boolean
  created_at: string // RFC3339
  updated_at: string // RFC3339
}

/** Mirrors createPatientRequest from Go handler */
export interface CreatePatientRequest {
  professional_id: string
  full_name: string
  email: string
  phone: string
  cpf: string
  date_of_birth?: string // YYYY-MM-DD
  gender: string
  notes: string
}

/** Mirrors profileResponse from Go handler (patient/handler.go) */
export interface PatientProfile {
  id: UUID
  patient_id: UUID
  occupation?: string
  marital_status?: string
  ethnicity?: string
  blood_type?: string
  allergies: string[]
  chronic_conditions: string[]
  medications: string[]
  emergency_contact_name?: string
  emergency_contact_phone?: string
  created_at: string
  updated_at: string
}

/** Mirrors inviteResponse from Go handler */
export interface PatientInvite {
  invite_id: UUID
  code: string
  expires_at: string
}

export type Gender = "male" | "female" | "other" | "prefer_not_to_say"
