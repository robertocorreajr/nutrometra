"use client"

import { useState } from "react"
import Link from "next/link"
import { useDocumentVersions } from "@nutrometra/api-client/hooks"
import { LoadingState } from "@nutrometra/ui"
import { documentStatusLabels, documentStatusColors } from "@/lib/schemas/document"
import { format } from "date-fns"
import { ptBR } from "date-fns/locale"
import { ChevronDown, ChevronRight, History } from "lucide-react"

interface DocumentVersionListProps {
  documentId: string
  patientId: string
}

export function DocumentVersionList({ documentId, patientId }: DocumentVersionListProps) {
  const [isOpen, setIsOpen] = useState(false)
  const { data: versions, isLoading } = useDocumentVersions(documentId)

  const versionList = versions ?? []

  if (versionList.length === 0 && !isLoading) return null

  return (
    <div className="border rounded-md">
      <button
        type="button"
        onClick={() => setIsOpen(!isOpen)}
        className="flex items-center justify-between w-full px-4 py-3 text-sm font-medium text-left hover:bg-muted/50 transition-colors"
      >
        <span className="flex items-center gap-2">
          <History className="h-4 w-4 text-muted-foreground" />
          Histórico de versões ({versionList.length})
        </span>
        {isOpen ? (
          <ChevronDown className="h-4 w-4 text-muted-foreground" />
        ) : (
          <ChevronRight className="h-4 w-4 text-muted-foreground" />
        )}
      </button>

      {isOpen && (
        <div className="px-4 pb-4 space-y-2">
          {isLoading ? (
            <LoadingState lines={3} />
          ) : (
            versionList.map((version) => {
              const statusLabel = documentStatusLabels[version.status] ?? version.status
              const statusColor =
                documentStatusColors[version.status] ?? "bg-gray-100 text-gray-700 border-gray-200"

              return (
                <Link
                  key={version.id}
                  href={`/pacientes/${patientId}/documentos/${version.id}`}
                  className="flex items-center justify-between rounded-md border px-3 py-2 hover:bg-muted/50 transition-colors"
                >
                  <div className="flex items-center gap-2">
                    <span className="text-sm font-medium">
                      v{version.version_number}
                    </span>
                    <span
                      className={`inline-flex items-center rounded-full border px-2 py-0.5 text-xs font-medium ${statusColor}`}
                    >
                      {statusLabel}
                    </span>
                  </div>
                  <span className="text-xs text-muted-foreground">
                    {format(new Date(version.created_at), "dd/MM/yyyy HH:mm", {
                      locale: ptBR,
                    })}
                  </span>
                </Link>
              )
            })
          )}
        </div>
      )}
    </div>
  )
}
