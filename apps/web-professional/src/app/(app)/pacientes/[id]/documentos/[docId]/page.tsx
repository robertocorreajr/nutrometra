"use client"

import { useParams } from "next/navigation"
import Link from "next/link"
import { useDocument } from "@nutrometra/api-client/hooks"
import { LoadingState, ErrorState } from "@nutrometra/ui"
import { DocumentHeader } from "@/components/documents/document-header"
import { DocumentContentEditor } from "@/components/documents/document-content-editor"
import { DocumentStatusActions } from "@/components/documents/document-status-actions"
import { DocumentVersionList } from "@/components/documents/document-version-list"
import { ArrowLeft } from "lucide-react"

export default function DocumentoDetalhe() {
  const params = useParams()
  const patientId = params.id as string
  const docId = params.docId as string

  const { data: document, isLoading, isError, refetch } = useDocument(docId)

  if (isLoading) return <LoadingState lines={8} />
  if (isError || !document) {
    return (
      <ErrorState
        message="Não foi possível carregar o documento."
        onRetry={refetch}
      />
    )
  }

  return (
    <div className="space-y-6">
      {/* Back link */}
      <Link
        href={`/pacientes/${patientId}/documentos`}
        className="inline-flex items-center gap-1 text-sm text-muted-foreground hover:text-foreground"
      >
        <ArrowLeft className="h-4 w-4" />
        Voltar para documentos
      </Link>

      {/* Header */}
      <DocumentHeader document={document} />

      {/* Content editor */}
      <DocumentContentEditor document={document} patientId={patientId} />

      {/* Status actions */}
      <DocumentStatusActions document={document} patientId={patientId} />

      {/* Version history */}
      <DocumentVersionList documentId={document.id} patientId={patientId} />
    </div>
  )
}
