import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query"
import { api } from "../client"
import type {
  BodyMeasurement,
  CreateMeasurementRequest,
} from "../types/bioimpedance"

export function useMeasurements(patientId: string) {
  return useQuery({
    queryKey: ["patients", patientId, "measurements"],
    queryFn: () =>
      api.get<BodyMeasurement[]>(`/patients/${patientId}/measurements`),
    enabled: !!patientId,
  })
}

export function useMeasurement(id: string) {
  return useQuery({
    queryKey: ["measurements", id],
    queryFn: () => api.get<BodyMeasurement>(`/measurements/${id}`),
    enabled: !!id,
  })
}

export function useCreateMeasurement(patientId: string) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (data: CreateMeasurementRequest) =>
      api.post<BodyMeasurement>("/measurements", data),
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: ["patients", patientId, "measurements"],
      })
    },
  })
}

export function usePublishMeasurement(id: string, patientId: string) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (data: { message?: string }) =>
      api.post<void>(`/measurements/${id}/publish`, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["measurements", id] })
      queryClient.invalidateQueries({
        queryKey: ["patients", patientId, "measurements"],
      })
    },
  })
}
