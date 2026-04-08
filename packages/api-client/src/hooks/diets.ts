import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query"
import { api } from "../client"
import type {
  Diet,
  CreateDietRequest,
  UpdateDietRequest,
  CreateMealRequest,
  DietMeal,
  CreateMealItemRequest,
  DietMealItem,
  CreateSubstitutionRequest,
  DietSubstitution,
} from "../types/diet"

export function usePatientDiets(patientId: string) {
  return useQuery({
    queryKey: ["patients", patientId, "diets"],
    queryFn: () => api.get<Diet[]>(`/patients/${patientId}/diets`),
    enabled: !!patientId,
  })
}

export function useDiet(id: string) {
  return useQuery({
    queryKey: ["diets", id],
    queryFn: () => api.get<Diet>(`/diets/${id}`),
    enabled: !!id,
  })
}

export function useCreateDiet(patientId: string) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (data: CreateDietRequest) =>
      api.post<Diet>(`/patients/${patientId}/diets`, data),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["patients", patientId, "diets"] })
    },
  })
}

export function useUpdateDiet(id: string, patientId: string) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (data: UpdateDietRequest) => api.put<Diet>(`/diets/${id}`, data),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["diets", id] })
      qc.invalidateQueries({ queryKey: ["patients", patientId, "diets"] })
    },
  })
}

export function useDeleteDiet(patientId: string) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id: string) => api.delete(`/diets/${id}`),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["patients", patientId, "diets"] })
    },
  })
}

export function usePublishDiet(id: string, patientId: string) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: () => api.post<Diet>(`/diets/${id}/publish`),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["diets", id] })
      qc.invalidateQueries({ queryKey: ["patients", patientId, "diets"] })
    },
  })
}

export function useArchiveDiet(id: string, patientId: string) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: () => api.post<Diet>(`/diets/${id}/archive`),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["diets", id] })
      qc.invalidateQueries({ queryKey: ["patients", patientId, "diets"] })
    },
  })
}

export function useNewDietVersion(id: string, patientId: string) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: () => api.post<Diet>(`/diets/${id}/new-version`),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["patients", patientId, "diets"] })
    },
  })
}

export function useAddMeal(dietId: string) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (data: CreateMealRequest) =>
      api.post<DietMeal>(`/diets/${dietId}/meals`, data),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["diets", dietId] })
    },
  })
}

export function useUpdateMeal(dietId: string) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ mealId, ...data }: CreateMealRequest & { mealId: string }) =>
      api.put<DietMeal>(`/diets/${dietId}/meals/${mealId}`, data),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["diets", dietId] })
    },
  })
}

export function useDeleteMeal(dietId: string) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (mealId: string) =>
      api.delete(`/diets/${dietId}/meals/${mealId}`),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["diets", dietId] })
    },
  })
}

export function useAddMealItem(dietId: string, mealId: string) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (data: CreateMealItemRequest) =>
      api.post<DietMealItem>(`/diets/${dietId}/meals/${mealId}/items`, data),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["diets", dietId] })
    },
  })
}

export function useUpdateMealItem(dietId: string) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ itemId, ...data }: CreateMealItemRequest & { itemId: string }) =>
      api.put<DietMealItem>(`/diet-items/${itemId}`, data),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["diets", dietId] })
    },
  })
}

export function useDeleteMealItem(dietId: string) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (itemId: string) => api.delete(`/diet-items/${itemId}`),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["diets", dietId] })
    },
  })
}

export function useAddSubstitution(dietId: string) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ itemId, ...data }: CreateSubstitutionRequest & { itemId: string }) =>
      api.post<DietSubstitution>(`/diet-items/${itemId}/substitutions`, data),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["diets", dietId] })
    },
  })
}

export function useDeleteSubstitution(dietId: string) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (subId: string) => api.delete(`/diet-substitutions/${subId}`),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["diets", dietId] })
    },
  })
}
