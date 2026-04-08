import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query"
import { api } from "../client"
import type { Role, AssignRoleRequest } from "../types/rbac"

export function useRoles() {
  return useQuery({
    queryKey: ["roles"],
    queryFn: () => api.get<Role[]>("/roles"),
  })
}

export function useAssignRole(memberId: string) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (data: AssignRoleRequest) =>
      api.post<void>(`/members/${memberId}/roles`, data),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["professionals"] })
    },
  })
}

export function useRevokeRole(memberId: string) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (roleId: string) =>
      api.delete<void>(`/members/${memberId}/roles/${roleId}`),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["professionals"] })
    },
  })
}
