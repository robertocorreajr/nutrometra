import type { UUID } from "./common"

export type ServiceMode = "onsite" | "online" | "home_visit"
export type AppointmentStatus = "scheduled" | "confirmed" | "completed" | "cancelled" | "no_show"
export type AppointmentSource = "professional" | "patient" | "system"

export interface AvailabilityRule {
  id: UUID
  professional_id: UUID
  day_of_week: number
  start_time: string
  end_time: string
  service_mode: ServiceMode
  address_id?: UUID
  active: boolean
}

export interface CreateAvailabilityRuleRequest {
  day_of_week: number
  start_time: string
  end_time: string
  service_mode: ServiceMode
  address_id?: string
}

export interface ScheduleBlock {
  id: UUID
  professional_id: UUID
  start_at: string
  end_at: string
  reason?: string
  all_day: boolean
}

export interface CreateBlockRequest {
  start_at: string
  end_at: string
  reason: string
  all_day: boolean
}

export interface TimeSlot {
  start: string
  end: string
  service_mode: ServiceMode
  address_id?: UUID
}

export interface Appointment {
  id: UUID
  tenant_id: UUID
  professional_id: UUID
  patient_id?: UUID
  start_at: string
  end_at: string
  service_mode: ServiceMode
  address_id?: UUID
  status: AppointmentStatus
  source: AppointmentSource
  notes?: string
  cancellation_reason?: string
  created_at: string
  updated_at: string
}

export interface CreateAppointmentRequest {
  professional_id: string
  patient_id?: string
  start_at: string
  end_at: string
  service_mode: ServiceMode
  address_id?: string
  source?: string
  notes?: string
}

export interface UpdateStatusRequest {
  status: AppointmentStatus
  cancellation_reason?: string
}

export interface RescheduleRequest {
  start_at: string
  end_at: string
}
