export interface ActivateInviteRequest {
  code: string
}

export interface ActivateInviteResponse {
  patient_id: string
  tenant_id: string
  professional_name: string
}
