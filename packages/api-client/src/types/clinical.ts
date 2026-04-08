import type { UUID } from "./common"

export interface Anamnesis {
  id: UUID
  patient_id: UUID
  professional_id: UUID
  status: AnamnesisStatus
  chief_complaint?: string
  history_present_illness?: string
  past_medical_history?: string
  family_history?: string
  social_history?: string
  dietary_history?: string
  physical_activity?: string
  sleep_pattern?: string
  bowel_habits?: string
  water_intake?: string
  supplements?: string
  observations?: string
  finalized_at?: string
  created_at: string
  updated_at: string
}

export type AnamnesisStatus = "draft" | "finalized"

export interface AnamnesisRequest {
  patient_id: string
  professional_id: string
  chief_complaint: string
  history_present_illness: string
  past_medical_history: string
  family_history: string
  social_history: string
  dietary_history: string
  physical_activity: string
  sleep_pattern: string
  bowel_habits: string
  water_intake: string
  supplements: string
  observations: string
}

export interface ProgressNote {
  id: UUID
  patient_id: UUID
  professional_id: UUID
  appointment_id?: UUID
  title: string
  content: string
  visible_to_patient: boolean
  created_at: string
  updated_at: string
}

export interface ProgressNoteRequest {
  professional_id: string
  appointment_id?: string
  title: string
  content: string
  visible_to_patient: boolean
}

export type AttachmentCategory =
  | "general"
  | "exam"
  | "lab_result"
  | "prescription"
  | "photo"
  | "other"

export interface ClinicalAttachment {
  id: UUID
  patient_id: UUID
  professional_id: UUID
  file_name: string
  file_type: string
  file_size_bytes: number
  storage_key: string
  category: AttachmentCategory
  description?: string
  created_at: string
}

export interface AttachmentRequest {
  professional_id: string
  file_name: string
  file_type: string
  file_size_bytes: number
  storage_key: string
  category: AttachmentCategory
  description: string
}
