import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query"
import { api } from "../client"
import type { ExportedFile, RequestExportRequest } from "../types/export"

export function useRequestExport() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (data: RequestExportRequest) =>
      api.post<ExportedFile>("/exports", data),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["exports"] })
    },
  })
}

export function useExports() {
  return useQuery({
    queryKey: ["exports"],
    queryFn: () => api.get<ExportedFile[]>("/exports"),
  })
}

export function useExport(
  id: string,
  options?: { refetchInterval?: number | false },
) {
  return useQuery({
    queryKey: ["exports", id],
    queryFn: () => api.get<ExportedFile>(`/exports/${id}`),
    enabled: !!id,
    refetchInterval: options?.refetchInterval,
  })
}

export async function downloadExport(id: string) {
  await api.download(`/exports/${id}/download`, "export.pdf")
}
