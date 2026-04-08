import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query"
import { api } from "../client"
import type {
  ClinicalDocument,
  CreateDocumentRequest,
  UpdateDocumentRequest,
} from "../types/document"

export function usePatientDocuments(patientId: string) {
  return useQuery({
    queryKey: ["patients", patientId, "documents"],
    queryFn: () => api.get<ClinicalDocument[]>(`/patients/${patientId}/documents`),
    enabled: !!patientId,
  })
}

export function useDocument(id: string) {
  return useQuery({
    queryKey: ["documents", id],
    queryFn: () => api.get<ClinicalDocument>(`/documents/${id}`),
    enabled: !!id,
  })
}

export function useCreateDocument(patientId: string) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (data: CreateDocumentRequest) =>
      api.post<ClinicalDocument>(`/patients/${patientId}/documents`, data),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["patients", patientId, "documents"] })
    },
  })
}

export function useUpdateDocument(id: string, patientId: string) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (data: UpdateDocumentRequest) =>
      api.put<ClinicalDocument>(`/documents/${id}`, data),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["documents", id] })
      qc.invalidateQueries({ queryKey: ["patients", patientId, "documents"] })
    },
  })
}

export function useFinalizeDocument(id: string, patientId: string) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: () => api.post<void>(`/documents/${id}/finalize`),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["documents", id] })
      qc.invalidateQueries({ queryKey: ["patients", patientId, "documents"] })
    },
  })
}

export function usePublishDocument(id: string, patientId: string) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: () => api.post<void>(`/documents/${id}/publish`),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["documents", id] })
      qc.invalidateQueries({ queryKey: ["patients", patientId, "documents"] })
    },
  })
}

export function useNewDocumentVersion(id: string, patientId: string) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: () => api.post<ClinicalDocument>(`/documents/${id}/new-version`),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["patients", patientId, "documents"] })
    },
  })
}

export function useDocumentVersions(id: string) {
  return useQuery({
    queryKey: ["documents", id, "versions"],
    queryFn: () => api.get<ClinicalDocument[]>(`/documents/${id}/versions`),
    enabled: !!id,
  })
}
