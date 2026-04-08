"use client"

import { useState } from "react"
import { useParams, useRouter } from "next/navigation"
import { LoadingState, ErrorState, EmptyState, Button } from "@nutrometra/ui"
import { usePatientDocuments } from "@nutrometra/api-client/hooks"
import type { ClinicalDocument } from "@nutrometra/api-client"
import { DocumentForm } from "@/components/documents/document-form"
import {
  documentTypeLabels,
  documentStatusLabels,
  documentStatusColors,
} from "@/lib/schemas/document"
import { format } from "date-fns"
import { ptBR } from "date-fns/locale"
import { Plus, FileText } from "lucide-react"

const filterOptions = [
  { value: "", label: "Todos" },
  { value: "exam_request", label: "Solicitação de Exame" },
  { value: "prescription", label: "Prescrição" },
  { value: "letter", label: "Carta" },
  { value: "other", label: "Outro" },
] as const

export default function DocumentosPage() {
  const params = useParams()
  const router = useRouter()
  const patientId = params.id as string

  const { data: documents, isLoading, isError, refetch } = usePatientDocuments(patientId)

  const [showForm, setShowForm] = useState(false)
  const [typeFilter, setTypeFilter] = useState("")

  function handleSuccess(doc: ClinicalDocument) {
    setShowForm(false)
    router.push(`/pacientes/${patientId}/documentos/${doc.id}`)
  }

  if (isLoading) return <LoadingState lines={6} />
  if (isError) {
    return (
      <ErrorState
        message="Erro ao carregar documentos."
        onRetry={refetch}
      />
    )
  }

  const allDocuments = documents ?? []
  const filteredDocuments = typeFilter
    ? allDocuments.filter((d) => d.document_type === typeFilter)
    : allDocuments

  return (
    <div className="space-y-4">
      {showForm ? (
        <DocumentForm
          patientId={patientId}
          onSuccess={handleSuccess}
          onCancel={() => setShowForm(false)}
        />
      ) : (
        <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
          <div className="flex items-center gap-2">
            <select
              value={typeFilter}
              onChange={(e) => setTypeFilter(e.target.value)}
              className="flex h-9 rounded-md border border-input bg-background px-3 py-1 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
              aria-label="Filtrar por tipo"
            >
              {filterOptions.map((opt) => (
                <option key={opt.value} value={opt.value}>
                  {opt.label}
                </option>
              ))}
            </select>
          </div>
          <Button onClick={() => setShowForm(true)}>
            <Plus className="h-4 w-4 mr-2" />
            Novo Documento
          </Button>
        </div>
      )}

      {filteredDocuments.length === 0 && !showForm ? (
        <EmptyState
          icon={<FileText className="h-12 w-12" />}
          title="Nenhum documento"
          description={
            typeFilter
              ? "Nenhum documento encontrado com este filtro."
              : "Crie solicitações de exame, prescrições, cartas e outros documentos clínicos."
          }
          action={
            !typeFilter ? (
              <Button onClick={() => setShowForm(true)}>
                <Plus className="h-4 w-4 mr-2" />
                Novo Documento
              </Button>
            ) : undefined
          }
        />
      ) : (
        <div className="grid gap-3">
          {filteredDocuments.map((doc) => {
            const typeLabel = documentTypeLabels[doc.document_type] ?? doc.document_type
            const statusLabel = documentStatusLabels[doc.status] ?? doc.status
            const statusColor =
              documentStatusColors[doc.status] ?? "bg-gray-100 text-gray-700 border-gray-200"

            return (
              <button
                key={doc.id}
                type="button"
                onClick={() => router.push(`/pacientes/${patientId}/documentos/${doc.id}`)}
                className="flex flex-col gap-2 rounded-lg border p-4 text-left hover:bg-muted/50 transition-colors sm:flex-row sm:items-center sm:justify-between"
              >
                <div className="flex flex-col gap-1.5 min-w-0">
                  <div className="flex flex-wrap items-center gap-2">
                    <span className="inline-flex items-center rounded-full border px-2 py-0.5 text-xs font-medium bg-gray-100 text-gray-700 border-gray-200">
                      {typeLabel}
                    </span>
                    <span
                      className={`inline-flex items-center rounded-full border px-2 py-0.5 text-xs font-medium ${statusColor}`}
                    >
                      {statusLabel}
                    </span>
                    <span className="text-xs text-muted-foreground">
                      v{doc.version_number}
                    </span>
                  </div>
                  <span className="text-sm font-medium truncate">{doc.title}</span>
                </div>
                <span className="text-xs text-muted-foreground whitespace-nowrap">
                  {format(new Date(doc.created_at), "dd/MM/yyyy", { locale: ptBR })}
                </span>
              </button>
            )
          })}
        </div>
      )}
    </div>
  )
}
