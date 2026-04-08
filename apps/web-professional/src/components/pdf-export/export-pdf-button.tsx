"use client"

import { useState } from "react"
import { Button } from "@nutrometra/ui"
import { useRequestExport, useExport, downloadExport } from "@nutrometra/api-client/hooks"
import { Download, Loader2, FileDown, AlertCircle } from "lucide-react"

interface ExportPDFButtonProps {
  entityType: string
  entityId: string
  label?: string
}

type ExportState = "idle" | "requesting" | "polling" | "completed" | "failed"

export function ExportPDFButton({
  entityType,
  entityId,
  label = "Exportar PDF",
}: ExportPDFButtonProps) {
  const [state, setState] = useState<ExportState>("idle")
  const [exportId, setExportId] = useState<string | null>(null)
  const [errorMessage, setErrorMessage] = useState<string | null>(null)

  const requestExport = useRequestExport()

  const { data: exportData } = useExport(exportId ?? "", {
    refetchInterval:
      state === "polling" ? 2000 : false,
  })

  // React to export status changes
  if (exportData && state === "polling") {
    if (exportData.status === "completed") {
      setState("completed")
    } else if (exportData.status === "failed") {
      setState("failed")
      setErrorMessage(exportData.failure_reason ?? "Erro ao gerar PDF.")
    }
  }

  async function handleRequest() {
    setState("requesting")
    setErrorMessage(null)
    try {
      const result = await requestExport.mutateAsync({
        entity_type: entityType,
        entity_id: entityId,
      })
      setExportId(result.id)
      setState("polling")
    } catch {
      setState("failed")
      setErrorMessage("Erro ao solicitar exportação.")
    }
  }

  async function handleDownload() {
    if (!exportId) return
    try {
      await downloadExport(exportId)
    } catch {
      setErrorMessage("Erro ao baixar o arquivo.")
    }
  }

  function handleReset() {
    setState("idle")
    setExportId(null)
    setErrorMessage(null)
  }

  if (state === "idle") {
    return (
      <Button variant="outline" size="sm" onClick={handleRequest}>
        <FileDown className="h-4 w-4 mr-2" />
        {label}
      </Button>
    )
  }

  if (state === "requesting" || state === "polling") {
    return (
      <Button variant="outline" size="sm" disabled>
        <Loader2 className="h-4 w-4 mr-2 animate-spin" />
        {state === "requesting" ? "Solicitando..." : "Gerando PDF..."}
      </Button>
    )
  }

  if (state === "completed") {
    return (
      <div className="flex items-center gap-2">
        <Button variant="outline" size="sm" onClick={handleDownload}>
          <Download className="h-4 w-4 mr-2" />
          Baixar PDF
        </Button>
        <Button variant="ghost" size="sm" onClick={handleReset}>
          Novo export
        </Button>
      </div>
    )
  }

  // failed
  return (
    <div className="flex items-center gap-2">
      <span className="inline-flex items-center gap-1 text-sm text-destructive">
        <AlertCircle className="h-4 w-4" />
        {errorMessage}
      </span>
      <Button variant="ghost" size="sm" onClick={handleReset}>
        Tentar novamente
      </Button>
    </div>
  )
}
