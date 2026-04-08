"use client"

import { useParams } from "next/navigation"
import { PageHeader, Card, CardContent, LoadingState, ErrorState } from "@nutrometra/ui"
import { useMyDocument } from "@nutrometra/api-client/hooks"
import { format, parseISO } from "date-fns"

const typeLabels: Record<string, string> = {
  prescription: "Prescricao", referral: "Encaminhamento",
  medical_certificate: "Atestado", lab_request: "Solicitacao de Exames",
  clinical_report: "Relatorio Clinico",
}

export default function DocumentDetailPage() {
  const { id } = useParams<{ id: string }>()
  const { data: doc, isLoading, isError, refetch } = useMyDocument(id)

  if (isLoading) return <LoadingState />
  if (isError) return <ErrorState message="Nao foi possivel carregar o documento." onRetry={refetch} />
  if (!doc) return <ErrorState message="Documento nao encontrado." />

  return (
    <div className="max-w-lg mx-auto md:max-w-none space-y-4">
      <PageHeader title={doc.title} description={typeLabels[doc.document_type] ?? doc.document_type} />
      <Card>
        <CardContent className="pt-4 space-y-3">
          <div className="text-sm text-muted-foreground">
            Criado em: {format(parseISO(doc.created_at), "dd/MM/yyyy HH:mm")}
          </div>
          <div className="text-sm text-muted-foreground">Versao: {doc.version_number}</div>
          {doc.content_json != null && (
            <div className="mt-4 prose prose-sm max-w-none">
              <pre className="whitespace-pre-wrap text-sm bg-muted/50 p-4 rounded-md">
                {typeof doc.content_json === "string"
                  ? doc.content_json
                  : JSON.stringify(doc.content_json as Record<string, unknown>, null, 2)}
              </pre>
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  )
}
