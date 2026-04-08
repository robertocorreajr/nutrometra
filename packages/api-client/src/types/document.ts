import type { UUID } from "./common"

export type DocumentType = "exam_request" | "prescription" | "letter" | "other"
export type DocumentStatus = "draft" | "finalized" | "published"

export interface ClinicalDocument {
  id: UUID
  patient_id: UUID
  professional_id: UUID
  appointment_id?: string
  document_type: DocumentType
  title: string
  status: DocumentStatus
  content_json: unknown
  version_number: number
  previous_version_id?: string
  created_at: string
  updated_at: string
}

export interface CreateDocumentRequest {
  professional_id: string
  appointment_id?: string
  document_type: DocumentType
  title: string
  content_json: unknown
}

export interface UpdateDocumentRequest {
  title: string
  content_json: unknown
  appointment_id?: string
}
