"use client"

import { useState, useCallback } from "react"
import { Button } from "@nutrometra/ui"
import { useUpdateDocument } from "@nutrometra/api-client/hooks"
import type { ClinicalDocument } from "@nutrometra/api-client"
import { Save } from "lucide-react"

interface DocumentContentEditorProps {
  document: ClinicalDocument
  patientId: string
}

export function DocumentContentEditor({ document, patientId }: DocumentContentEditorProps) {
  const content = document.content_json as { body: string } | null
  const [body, setBody] = useState(content?.body ?? "")
  const [isDirty, setIsDirty] = useState(false)

  const updateDocument = useUpdateDocument(document.id, patientId)
  const isReadonly = document.status !== "draft"

  const handleSave = useCallback(async () => {
    if (!isDirty) return

    try {
      await updateDocument.mutateAsync({
        title: document.title,
        content_json: { body },
      })
      setIsDirty(false)
    } catch {
      // Error handled by mutation state
    }
  }, [body, isDirty, document.title, updateDocument])

  function handleChange(value: string) {
    setBody(value)
    setIsDirty(true)
  }

  return (
    <div className="space-y-2">
      <div className="flex items-center justify-between">
        <label htmlFor="document-body" className="text-sm font-medium">
          Conteúdo
        </label>
        {!isReadonly && isDirty && (
          <Button
            size="sm"
            variant="outline"
            onClick={handleSave}
            disabled={updateDocument.isPending}
          >
            <Save className="h-3.5 w-3.5 mr-1.5" />
            {updateDocument.isPending ? "Salvando..." : "Salvar"}
          </Button>
        )}
      </div>

      <textarea
        id="document-body"
        value={body}
        onChange={(e) => handleChange(e.target.value)}
        onBlur={() => {
          if (!isReadonly && isDirty) {
            handleSave()
          }
        }}
        readOnly={isReadonly}
        rows={12}
        placeholder={
          isReadonly
            ? "Este documento não pode mais ser editado."
            : "Digite o conteúdo do documento..."
        }
        className={`flex w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 ${
          isReadonly ? "bg-muted cursor-not-allowed" : ""
        }`}
      />

      {updateDocument.isError && (
        <p className="text-sm text-destructive">
          Erro ao salvar conteúdo. Tente novamente.
        </p>
      )}

      {!isReadonly && isDirty && (
        <p className="text-xs text-muted-foreground">
          Alterações não salvas. O conteúdo será salvo automaticamente ao sair do campo.
        </p>
      )}

      {isReadonly && (
        <p className="text-xs text-muted-foreground">
          Documento {document.status === "finalized" ? "finalizado" : "publicado"} — não é possível editar. Crie uma nova versão para fazer alterações.
        </p>
      )}
    </div>
  )
}
