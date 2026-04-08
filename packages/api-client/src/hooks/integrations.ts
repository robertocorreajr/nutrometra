import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query"
import { api } from "../client"
import type { GoogleCalendarStatus, GoogleAuthorizeResponse } from "../types/integration"

export function useGoogleCalendarStatus() {
  return useQuery({
    queryKey: ["google-calendar-status"],
    queryFn: () => api.get<GoogleCalendarStatus>("/integrations/google/status"),
  })
}

export function useGoogleCalendarAuthorize() {
  return useMutation({
    mutationFn: () => api.get<GoogleAuthorizeResponse>("/integrations/google/authorize"),
  })
}

export function useGoogleCalendarDisconnect() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: () => api.post<void>("/integrations/google/disconnect"),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["google-calendar-status"] })
    },
  })
}
