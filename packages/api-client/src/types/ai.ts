import type { UUID } from "./common"

export type SuggestionType =
  | "diet_draft"
  | "meal_structure"
  | "substitutions"
  | "clinical_summary"
  | "review_checklist"

export type SuggestionStatus =
  | "pending"
  | "generating"
  | "completed"
  | "failed"
  | "accepted"
  | "rejected"

/** Mirrors aiSuggestionResponse from Go handler */
export interface AISuggestion {
  id: UUID
  tenant_id: UUID
  user_id: UUID
  suggestion_type: SuggestionType
  status: SuggestionStatus
  response_text?: string
  model_id?: string
  input_tokens?: number
  output_tokens?: number
  created_at: string
  completed_at?: string
  reviewed_at?: string
  reviewed_by?: UUID
}

/** Mirrors createSuggestionRequest from Go handler */
export interface CreateSuggestionRequest {
  suggestion_type: SuggestionType
  patient_id?: string
  extra_context?: Record<string, unknown>
}
