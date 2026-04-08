import type { UUID } from "./common"

export type MeasurementSource = "manual" | "device" | "import"

/** Mirrors bodyMeasurementResponse from Go handler */
export interface BodyMeasurement {
  id: UUID
  patient_id: UUID
  professional_id: UUID
  measured_at: string
  weight_kg?: number
  height_cm?: number
  bmi?: number
  body_fat_pct?: number
  lean_mass_kg?: number
  fat_mass_kg?: number
  muscle_mass_kg?: number
  bone_mass_kg?: number
  water_pct?: number
  visceral_fat?: number
  basal_metabolic_rate?: number
  waist_cm?: number
  hip_cm?: number
  source: MeasurementSource
  device_model?: string
  notes?: string
  created_at: string
}

/** Mirrors createMeasurementRequest from Go handler */
export interface CreateMeasurementRequest {
  patient_id: string
  professional_id: string
  measured_at?: string
  weight_kg?: number
  height_cm?: number
  body_fat_pct?: number
  lean_mass_kg?: number
  fat_mass_kg?: number
  muscle_mass_kg?: number
  bone_mass_kg?: number
  water_pct?: number
  visceral_fat?: number
  basal_metabolic_rate?: number
  waist_cm?: number
  hip_cm?: number
  chest_cm?: number
  right_arm_cm?: number
  left_arm_cm?: number
  right_thigh_cm?: number
  left_thigh_cm?: number
  right_calf_cm?: number
  left_calf_cm?: number
  neck_cm?: number
  abdomen_cm?: number
  triceps_sf_mm?: number
  biceps_sf_mm?: number
  subscapular_sf_mm?: number
  suprailiac_sf_mm?: number
  abdominal_sf_mm?: number
  thigh_sf_mm?: number
  calf_sf_mm?: number
  source: MeasurementSource
  device_model?: string
  notes?: string
}
