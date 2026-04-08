import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query"
import { api } from "../client"
import type { FoodItem, CreateFoodItemRequest } from "../types/catalog"

export function useFoods(query?: string, group?: string, limit?: number) {
  const params = new URLSearchParams()
  if (query) params.set("q", query)
  if (group) params.set("group", group)
  if (limit) params.set("limit", String(limit))
  const qs = params.toString()
  return useQuery({
    queryKey: ["foods", { query, group, limit }],
    queryFn: () => api.get<FoodItem[]>(`/foods${qs ? `?${qs}` : ""}`),
  })
}

export function useFoodGroups() {
  return useQuery({
    queryKey: ["foods", "groups"],
    queryFn: () => api.get<string[]>("/foods/groups"),
  })
}

export function useFood(id: string) {
  return useQuery({
    queryKey: ["foods", id],
    queryFn: () => api.get<FoodItem>(`/foods/${id}`),
    enabled: !!id,
    staleTime: 5 * 60 * 1000,
  })
}

export function useCreateFood() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (data: CreateFoodItemRequest) => api.post<FoodItem>("/foods", data),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["foods"] })
    },
  })
}

export function useUpdateFood(id: string) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (data: CreateFoodItemRequest) => api.put<FoodItem>(`/foods/${id}`, data),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["foods"] })
      qc.invalidateQueries({ queryKey: ["foods", id] })
    },
  })
}

export function useDeleteFood() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id: string) => api.delete(`/foods/${id}`),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["foods"] })
    },
  })
}
