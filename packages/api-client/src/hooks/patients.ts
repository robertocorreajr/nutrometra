import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query"
import { api } from "../client"
import type { Patient, CreatePatientRequest, PatientProfile, PatientInvite } from "../types/patient"

export function usePatients() {
  return useQuery({
    queryKey: ["patients"],
    queryFn: () => api.get<Patient[]>("/patients"),
  })
}

export function usePatient(id: string) {
  return useQuery({
    queryKey: ["patients", id],
    queryFn: () => api.get<Patient>(`/patients/${id}`),
    enabled: !!id,
  })
}

export function useCreatePatient() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (data: CreatePatientRequest) =>
      api.post<Patient>("/patients", data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["patients"] })
    },
  })
}

export function useUpdatePatient(id: string) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (data: CreatePatientRequest) =>
      api.put<Patient>(`/patients/${id}`, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["patients"] })
      queryClient.invalidateQueries({ queryKey: ["patients", id] })
    },
  })
}

export function usePatientProfile(patientId: string) {
  return useQuery({
    queryKey: ["patients", patientId, "profile"],
    queryFn: () => api.get<PatientProfile>(`/patients/${patientId}/profiles`),
    enabled: !!patientId,
  })
}

export function useGenerateInvite(patientId: string) {
  return useMutation({
    mutationFn: () =>
      api.post<PatientInvite>(`/patients/${patientId}/invites`),
  })
}
