import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query"
import { api } from "../client"
import type {
  Professional,
  UpdateProfessionalRequest,
  ProfessionalAddress,
  CreateAddressRequest,
  ProfessionalServiceMode,
  SetServiceModeRequest,
} from "../types/professional"

export function useProfessionalMe() {
  return useQuery({
    queryKey: ["professionals", "me"],
    queryFn: () => api.get<Professional>("/professionals/me"),
  })
}

export function useUpdateProfessional(id: string) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (data: UpdateProfessionalRequest) =>
      api.put<Professional>(`/professionals/${id}`, data),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["professionals", "me"] })
    },
  })
}

export function useAddresses(professionalId: string) {
  return useQuery({
    queryKey: ["addresses", professionalId],
    queryFn: () =>
      api.get<ProfessionalAddress[]>(`/professionals/${professionalId}/addresses`),
    enabled: !!professionalId,
  })
}

export function useCreateAddress(professionalId: string) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (data: CreateAddressRequest) =>
      api.post<ProfessionalAddress>(
        `/professionals/${professionalId}/addresses`,
        data,
      ),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["addresses", professionalId] })
    },
  })
}

export function useUpdateAddress(professionalId: string) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({
      addressId,
      ...data
    }: CreateAddressRequest & { addressId: string }) =>
      api.put<ProfessionalAddress>(
        `/professionals/${professionalId}/addresses/${addressId}`,
        data,
      ),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["addresses", professionalId] })
    },
  })
}

export function useServiceModes(professionalId: string) {
  return useQuery({
    queryKey: ["service-modes", professionalId],
    queryFn: () =>
      api.get<ProfessionalServiceMode[]>(
        `/professionals/${professionalId}/service-modes`,
      ),
    enabled: !!professionalId,
  })
}

export function useSetServiceMode(professionalId: string) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (data: SetServiceModeRequest) =>
      api.put<ProfessionalServiceMode>(
        `/professionals/${professionalId}/service-modes`,
        data,
      ),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["service-modes", professionalId] })
    },
  })
}

export function useProfessionals() {
  return useQuery({
    queryKey: ["professionals"],
    queryFn: () => api.get<Professional[]>("/professionals"),
  })
}

export function useCreateProfessional() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (data: Partial<Professional>) =>
      api.post<Professional>("/professionals", data),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["professionals"] })
    },
  })
}
