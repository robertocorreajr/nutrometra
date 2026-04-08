import { useQuery, useMutation } from "@tanstack/react-query"
import { api } from "../client"
import type { ActivateInviteRequest, ActivateInviteResponse } from "../types/patient-portal"
import type { Appointment } from "../types/scheduling"
import type { Diet } from "../types/diet"
import type { ClinicalDocument } from "../types/document"
import type { BodyMeasurement } from "../types/bioimpedance"
import type { PatientProfile } from "../types/patient"

// --- Invite Activation ---

export function useActivateInvite() {
  return useMutation({
    mutationFn: (data: ActivateInviteRequest) =>
      api.post<ActivateInviteResponse>("/invites/activate", data),
  })
}

// --- Appointments ---

export interface MyAppointmentsParams {
  from?: string
  to?: string
}

export function useMyAppointments(params?: MyAppointmentsParams) {
  const qs = new URLSearchParams()
  if (params?.from) qs.set("from", params.from)
  if (params?.to) qs.set("to", params.to)
  const query = qs.toString()
  return useQuery({
    queryKey: ["my-appointments", params],
    queryFn: () =>
      api.get<Appointment[]>(`/appointments${query ? `?${query}` : ""}`),
  })
}

// --- Diets ---

export function useMyDiets(patientId: string) {
  return useQuery({
    queryKey: ["my-diets", patientId],
    queryFn: () => api.get<Diet[]>(`/patients/${patientId}/diets`),
    enabled: !!patientId,
  })
}

export function useMyDiet(id: string) {
  return useQuery({
    queryKey: ["my-diet", id],
    queryFn: () => api.get<Diet>(`/diets/${id}`),
    enabled: !!id,
  })
}

// --- Documents ---

export function useMyDocuments(patientId: string) {
  return useQuery({
    queryKey: ["my-documents", patientId],
    queryFn: () =>
      api.get<ClinicalDocument[]>(`/patients/${patientId}/documents`),
    enabled: !!patientId,
  })
}

// --- Measurements ---

export function useMyMeasurements(patientId: string) {
  return useQuery({
    queryKey: ["my-measurements", patientId],
    queryFn: () =>
      api.get<BodyMeasurement[]>(`/patients/${patientId}/measurements`),
    enabled: !!patientId,
  })
}

// --- Profile ---

export function useMyProfile(patientId: string) {
  return useQuery({
    queryKey: ["my-profile", patientId],
    queryFn: () =>
      api.get<PatientProfile>(`/patients/${patientId}/profiles`),
    enabled: !!patientId,
  })
}
