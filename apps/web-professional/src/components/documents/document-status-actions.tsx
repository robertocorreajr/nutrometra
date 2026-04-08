"use client"

import { useState } from "react"
import { useRouter } from "next/navigation"
import { Button } from "@nutrometra/ui"
import {
  useFinalizeDocument,
  usePublishDocument,
  useNewDocumentVersion,
} from "@nutrometra/api-client/hooks"
import type { ClinicalDocument } from "@nutrometra/api-client"
import { CheckCircle, Send, FilePlus, AlertTriangle } from "lucide-react"

type ConfirmAction = "finalize" | "publish" | "new-version" | null

interface DocumentStatusActionsProps {
  document: ClinicalDocument
  patientId: string
}

export function DocumentStatusActions({ document, patientId }: DocumentStatusActionsProps) {
  const router = useRouter()
  const [confirmAction, setConfirmAction] = useState<ConfirmAction>(null)

  const finalizeDocument = useFinalizeDocument(document.id, patientId)
  const publishDocument = usePublishDocument(document.id, patientId)
  const newVersion = useNewDocumentVersion(document.id, patientId)

  const isProcessing =
    finalizeDocument.isPending || publishDocument.isPending || newVersion.isPending

  async function handleFinalize() {
    try {
      await finalizeDocument.mutateAsync()
      setConfirmAction(null)
    } catch {
      // Error handled by mutation state
    }
  }

  async function handlePublish() {
    try {
      await publishDocument.mutateAsync()
      setConfirmAction(null)
    } catch {
      // Error handled by mutation state
    }
  }

  async function handleNewVersion() {
    try {
      const newDoc = await newVersion.mutateAsync()
      setConfirmAction(null)
      router.push(`/pacientes/${patientId}/documentos/${newDoc.id}`)
    } catch {
      // Error handled by mutation state
    }
  }

  const errorMessage =
    finalizeDocument.isError
      ? "Erro ao finalizar documento."
      : publishDocument.isError
        ? "Erro ao publicar documento."
        : newVersion.isError
          ? "Erro ao criar nova versão."
          : null

  return (
    <div className="space-y-3">
      {/* Confirmation inline area */}
      {confirmAction && (
        <div className="rounded-md border border-amber-200 bg-amber-50 p-4 space-y-3">
          <div className="flex items-start gap-2">
            <AlertTriangle className="h-5 w-5 text-amber-600 mt-0.5 shrink-0" />
            <div className="space-y-1">
              <p className="text-sm font-medium text-amber-800">
                {confirmAction === "finalize" && "Finalizar documento?"}
                {confirmAction === "publish" && "Publicar documento?"}
                {confirmAction === "new-version" && "Criar nova versão?"}
              </p>
              <p className="text-sm text-amber-700">
                {confirmAction === "finalize" &&
                  "Após finalizar, o conteúdo não poderá mais ser editado."}
                {confirmAction === "publish" &&
                  "O documento será visível para o paciente após a publicação."}
                {confirmAction === "new-version" &&
                  "Uma cópia editável será criada como rascunho com o conteúdo atual."}
              </p>
            </div>
          </div>
          <div className="flex gap-2 justify-end">
            <Button
              size="sm"
              variant="ghost"
              onClick={() => setConfirmAction(null)}
              disabled={isProcessing}
            >
              Cancelar
            </Button>
            <Button
              size="sm"
              onClick={() => {
                if (confirmAction === "finalize") handleFinalize()
                if (confirmAction === "publish") handlePublish()
                if (confirmAction === "new-version") handleNewVersion()
              }}
              disabled={isProcessing}
            >
              {isProcessing ? "Processando..." : "Confirmar"}
            </Button>
          </div>
        </div>
      )}

      {errorMessage && (
        <p className="text-sm text-destructive">{errorMessage}</p>
      )}

      {/* Action buttons based on status */}
      {!confirmAction && (
        <div className="flex flex-wrap gap-2">
          {document.status === "draft" && (
            <Button
              variant="outline"
              onClick={() => setConfirmAction("finalize")}
              disabled={isProcessing}
            >
              <CheckCircle className="h-4 w-4 mr-2" />
              Finalizar
            </Button>
          )}

          {document.status === "finalized" && (
            <>
              <Button
                variant="outline"
                onClick={() => setConfirmAction("publish")}
                disabled={isProcessing}
              >
                <Send className="h-4 w-4 mr-2" />
                Publicar
              </Button>
              <Button
                variant="ghost"
                onClick={() => setConfirmAction("new-version")}
                disabled={isProcessing}
              >
                <FilePlus className="h-4 w-4 mr-2" />
                Nova Versão
              </Button>
            </>
          )}

          {document.status === "published" && (
            <Button
              variant="ghost"
              onClick={() => setConfirmAction("new-version")}
              disabled={isProcessing}
            >
              <FilePlus className="h-4 w-4 mr-2" />
              Nova Versão
            </Button>
          )}
        </div>
      )}
    </div>
  )
}
