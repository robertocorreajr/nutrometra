import { useQuery } from "@tanstack/react-query"
import { api } from "../client"
import type { Professional } from "../types/professional"

export function useProfessionalMe() {
  return useQuery({
    queryKey: ["professionals", "me"],
    queryFn: () => api.get<Professional>("/professionals/me"),
  })
}
