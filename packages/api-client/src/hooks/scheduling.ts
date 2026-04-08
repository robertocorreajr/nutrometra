import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query"
import { api } from "../client"
import type {
  AvailabilityRule,
  CreateAvailabilityRuleRequest,
  ScheduleBlock,
  CreateBlockRequest,
  TimeSlot,
  Appointment,
  CreateAppointmentRequest,
  UpdateStatusRequest,
  RescheduleRequest,
} from "../types/scheduling"

// --- Availability Rules ---

export function useAvailabilityRules(professionalId: string) {
  return useQuery({
    queryKey: ["availability-rules", professionalId],
    queryFn: () =>
      api.get<AvailabilityRule[]>(
        `/professionals/${professionalId}/availability`,
      ),
    enabled: !!professionalId,
  })
}

export function useCreateAvailabilityRule(professionalId: string) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (data: CreateAvailabilityRuleRequest) =>
      api.post<AvailabilityRule>(
        `/professionals/${professionalId}/availability`,
        data,
      ),
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: ["availability-rules", professionalId],
      })
    },
  })
}

export function useUpdateAvailabilityRule(professionalId: string) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ id, ...data }: { id: string } & CreateAvailabilityRuleRequest) =>
      api.put<AvailabilityRule>(`/availability/${id}`, data),
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: ["availability-rules", professionalId],
      })
    },
  })
}

export function useDeleteAvailabilityRule(professionalId: string) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (id: string) => api.delete<void>(`/availability/${id}`),
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: ["availability-rules", professionalId],
      })
    },
  })
}

// --- Available Slots ---

export function useAvailableSlots(
  professionalId: string,
  date: string,
  duration?: number,
) {
  const params = new URLSearchParams({ date })
  if (duration) params.set("duration", String(duration))

  return useQuery({
    queryKey: ["available-slots", professionalId, date, duration],
    queryFn: () =>
      api.get<TimeSlot[]>(
        `/professionals/${professionalId}/slots?${params.toString()}`,
      ),
    enabled: !!professionalId && !!date,
  })
}

// --- Schedule Blocks ---

export function useBlocks(professionalId: string, from?: string, to?: string) {
  const params = new URLSearchParams()
  if (from) params.set("from", from)
  if (to) params.set("to", to)
  const qs = params.toString()

  return useQuery({
    queryKey: ["blocks", professionalId, from, to],
    queryFn: () =>
      api.get<ScheduleBlock[]>(
        `/professionals/${professionalId}/blocks${qs ? `?${qs}` : ""}`,
      ),
    enabled: !!professionalId,
  })
}

export function useCreateBlock(professionalId: string) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (data: CreateBlockRequest) =>
      api.post<ScheduleBlock>(
        `/professionals/${professionalId}/blocks`,
        data,
      ),
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: ["blocks", professionalId],
      })
    },
  })
}

export function useDeleteBlock(professionalId: string) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (id: string) => api.delete<void>(`/blocks/${id}`),
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: ["blocks", professionalId],
      })
    },
  })
}

// --- Appointments ---

export interface UseAppointmentsParams {
  professional_id?: string
  from?: string
  to?: string
  status?: string
}

export function useAppointments(params: UseAppointmentsParams) {
  const searchParams = new URLSearchParams()
  if (params.professional_id)
    searchParams.set("professional_id", params.professional_id)
  if (params.from) searchParams.set("from", params.from)
  if (params.to) searchParams.set("to", params.to)
  if (params.status) searchParams.set("status", params.status)
  const qs = searchParams.toString()

  return useQuery({
    queryKey: ["appointments", params],
    queryFn: () =>
      api.get<Appointment[]>(`/appointments${qs ? `?${qs}` : ""}`),
  })
}

export function useAppointment(id: string) {
  return useQuery({
    queryKey: ["appointments", id],
    queryFn: () => api.get<Appointment>(`/appointments/${id}`),
    enabled: !!id,
  })
}

export function useCreateAppointment() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (data: CreateAppointmentRequest) =>
      api.post<Appointment>("/appointments", data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["appointments"] })
    },
  })
}

export function useUpdateAppointmentStatus(id: string) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (data: UpdateStatusRequest) =>
      api.patch<Appointment>(`/appointments/${id}/status`, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["appointments"] })
      queryClient.invalidateQueries({ queryKey: ["appointments", id] })
    },
  })
}

export function useRescheduleAppointment(id: string) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (data: RescheduleRequest) =>
      api.patch<Appointment>(`/appointments/${id}/reschedule`, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["appointments"] })
      queryClient.invalidateQueries({ queryKey: ["appointments", id] })
    },
  })
}
