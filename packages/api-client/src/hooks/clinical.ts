import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query"
import { api } from "../client"
import type {
  Anamnesis,
  AnamnesisRequest,
  ProgressNote,
  ProgressNoteRequest,
  ClinicalAttachment,
  AttachmentRequest,
} from "../types/clinical"

// --- Anamnesis ---

export function useAnamneses(patientId: string) {
  return useQuery({
    queryKey: ["patients", patientId, "anamneses"],
    queryFn: () => api.get<Anamnesis[]>(`/patients/${patientId}/anamneses`),
    enabled: !!patientId,
  })
}

export function useAnamnesis(id: string) {
  return useQuery({
    queryKey: ["anamneses", id],
    queryFn: () => api.get<Anamnesis>(`/anamneses/${id}`),
    enabled: !!id,
  })
}

export function useCreateAnamnesis(patientId: string) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (data: AnamnesisRequest) =>
      api.post<Anamnesis>(`/patients/${patientId}/anamneses`, data),
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: ["patients", patientId, "anamneses"],
      })
    },
  })
}

export function useUpdateAnamnesis(id: string, patientId: string) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (data: AnamnesisRequest) =>
      api.put<Anamnesis>(`/anamneses/${id}`, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["anamneses", id] })
      queryClient.invalidateQueries({
        queryKey: ["patients", patientId, "anamneses"],
      })
    },
  })
}

export function useFinalizeAnamnesis(id: string, patientId: string) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: () => api.post<void>(`/anamneses/${id}/finalize`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["anamneses", id] })
      queryClient.invalidateQueries({
        queryKey: ["patients", patientId, "anamneses"],
      })
    },
  })
}

// --- Progress Notes ---

export function useProgressNotes(patientId: string) {
  return useQuery({
    queryKey: ["patients", patientId, "notes"],
    queryFn: () =>
      api.get<ProgressNote[]>(`/patients/${patientId}/notes`),
    enabled: !!patientId,
  })
}

export function useCreateProgressNote(patientId: string) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (data: ProgressNoteRequest) =>
      api.post<ProgressNote>(`/patients/${patientId}/notes`, data),
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: ["patients", patientId, "notes"],
      })
    },
  })
}

// --- Attachments ---

export function useAttachments(patientId: string) {
  return useQuery({
    queryKey: ["patients", patientId, "attachments"],
    queryFn: () =>
      api.get<ClinicalAttachment[]>(`/patients/${patientId}/attachments`),
    enabled: !!patientId,
  })
}

export function useCreateAttachment(patientId: string) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (data: AttachmentRequest) =>
      api.post<ClinicalAttachment>(
        `/patients/${patientId}/attachments`,
        data
      ),
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: ["patients", patientId, "attachments"],
      })
    },
  })
}
