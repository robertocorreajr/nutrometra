"use client"

import type { ClinicalDocument } from "@nutrometra/api-client"
import {
  documentTypeLabels,
  documentStatusLabels,
  documentStatusColors,
} from "@/lib/schemas/document"

interface DocumentHeaderProps {
  document: ClinicalDocument
}

export function DocumentHeader({ document }: DocumentHeaderProps) {
  const typeLabel = documentTypeLabels[document.document_type] ?? document.document_type
  const statusLabel = documentStatusLabels[document.status] ?? document.status
  const statusColor = documentStatusColors[document.status] ?? "bg-gray-100 text-gray-700 border-gray-200"

  return (
    <div className="space-y-2">
      <div className="flex flex-wrap items-center gap-2">
        <span className="inline-flex items-center rounded-full border px-2.5 py-0.5 text-xs font-medium bg-gray-100 text-gray-700 border-gray-200">
          {typeLabel}
        </span>
        <span
          className={`inline-flex items-center rounded-full border px-2.5 py-0.5 text-xs font-medium ${statusColor}`}
        >
          {statusLabel}
        </span>
        <span className="text-xs text-muted-foreground">
          Versão {document.version_number}
        </span>
      </div>
      <h2 className="text-xl font-semibold tracking-tight">{document.title}</h2>
    </div>
  )
}
