import type { UUID } from "./common"

export type ExportStatus = "pending" | "processing" | "completed" | "failed"
export type ExportType = "pdf" | "print_job"

export interface ExportedFile {
  id: UUID
  tenant_id: UUID
  related_entity_type: string
  related_entity_id: UUID
  export_type: ExportType
  file_key?: string
  status: ExportStatus
  requested_by_user_id: UUID
  created_at: string
  completed_at?: string
  failure_reason?: string
}

export interface RequestExportRequest {
  entity_type: string
  entity_id: string
}
