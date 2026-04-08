import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query"
import { api } from "../client"
import type {
  AISuggestion,
  CreateSuggestionRequest,
} from "../types/ai"
import type { PaginatedResponse } from "../types/common"

interface AISuggestionsParams {
  limit?: number
  offset?: number
}

export function useAISuggestions(params?: AISuggestionsParams) {
  const searchParams = new URLSearchParams()
  if (params?.limit != null) searchParams.set("limit", String(params.limit))
  if (params?.offset != null) searchParams.set("offset", String(params.offset))
  const qs = searchParams.toString()

  return useQuery({
    queryKey: ["ai-suggestions", params],
    queryFn: () =>
      api.get<PaginatedResponse<AISuggestion>>(
        `/ai/suggestions${qs ? `?${qs}` : ""}`
      ),
  })
}

export function useAISuggestion(id: string) {
  return useQuery({
    queryKey: ["ai-suggestions", id],
    queryFn: () => api.get<AISuggestion>(`/ai/suggestions/${id}`),
    enabled: !!id,
  })
}

export function useCreateAISuggestion() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (data: CreateSuggestionRequest) =>
      api.post<AISuggestion>("/ai/suggestions", data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["ai-suggestions"] })
    },
  })
}

export function useAcceptAISuggestion(id: string) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: () =>
      api.post<void>(`/ai/suggestions/${id}/accept`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["ai-suggestions", id] })
      queryClient.invalidateQueries({ queryKey: ["ai-suggestions"] })
    },
  })
}

export function useRejectAISuggestion(id: string) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: () =>
      api.post<void>(`/ai/suggestions/${id}/reject`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["ai-suggestions", id] })
      queryClient.invalidateQueries({ queryKey: ["ai-suggestions"] })
    },
  })
}
