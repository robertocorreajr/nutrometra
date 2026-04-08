"use client"

import Link from "next/link"
import { useAuth } from "@nutrometra/auth"
import { PageHeader, Card, CardContent, LoadingState, ErrorState, EmptyState } from "@nutrometra/ui"
import { useMyDocuments } from "@nutrometra/api-client/hooks"
import { format, parseISO } from "date-fns"
import { FileText } from "lucide-react"

const typeLabels: Record<string, string> = {
  prescription: "Prescricao", referral: "Encaminhamento",
  medical_certificate: "Atestado", lab_request: "Solicitacao de Exames",
  clinical_report: "Relatorio Clinico",
}

export default function DocumentosPage() {
  const { patientId } = useAuth()
  const { data: documents, isLoading, isError, refetch } = useMyDocuments(patientId ?? "")

  return (
    <div className="max-w-lg mx-auto md:max-w-none">
      <PageHeader title="Documentos" description="Documentos clinicos" />
      {isLoading && <LoadingState />}
      {isError && <ErrorState message="Nao foi possivel carregar os documentos." onRetry={refetch} />}
      {documents && documents.length === 0 && <EmptyState title="Nenhum documento" description="Nenhum documento publicado." />}
      {documents && documents.length > 0 && (
        <div className="space-y-3">
          {documents.filter((d) => d.status === "published").map((doc) => (
            <Link key={doc.id} href={`/documentos/${doc.id}`}>
              <Card className="hover:border-primary/50 transition-colors">
                <CardContent className="pt-4 flex items-start gap-3">
                  <FileText className="h-5 w-5 text-muted-foreground mt-0.5" />
                  <div>
                    <h3 className="font-medium">{doc.title}</h3>
                    <p className="text-sm text-muted-foreground">
                      {typeLabels[doc.document_type] ?? doc.document_type} · {format(parseISO(doc.created_at), "dd/MM/yyyy")}
                    </p>
                  </div>
                </CardContent>
              </Card>
            </Link>
          ))}
        </div>
      )}
    </div>
  )
}
