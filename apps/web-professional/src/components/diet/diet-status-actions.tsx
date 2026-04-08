"use client"

import { useState } from "react"
import { useRouter, useParams } from "next/navigation"
import { Button } from "@nutrometra/ui"
import {
  usePublishDiet,
  useArchiveDiet,
  useDeleteDiet,
  useNewDietVersion,
} from "@nutrometra/api-client/hooks"
import type { Diet } from "@nutrometra/api-client"
import { ConfirmDialog } from "@/components/shared/confirm-dialog"
import { ExportPDFButton } from "@/components/pdf-export/export-pdf-button"

interface DietStatusActionsProps {
  diet: Diet
  patientId: string
}

type DialogType = "publish" | "archive" | "delete" | "new-version" | null

export function DietStatusActions({ diet, patientId }: DietStatusActionsProps) {
  const router = useRouter()
  const params = useParams()
  const [dialog, setDialog] = useState<DialogType>(null)

  const publishDiet = usePublishDiet(diet.id, patientId)
  const archiveDiet = useArchiveDiet(diet.id, patientId)
  const deleteDiet = useDeleteDiet(patientId)
  const newVersion = useNewDietVersion(diet.id, patientId)

  async function handlePublish() {
    try {
      await publishDiet.mutateAsync()
      setDialog(null)
    } catch {
      // Error handled by mutation state
    }
  }

  async function handleArchive() {
    try {
      await archiveDiet.mutateAsync()
      setDialog(null)
    } catch {
      // Error handled by mutation state
    }
  }

  async function handleDelete() {
    try {
      await deleteDiet.mutateAsync(diet.id)
      setDialog(null)
      router.push(`/pacientes/${patientId}/dietas`)
    } catch {
      // Error handled by mutation state
    }
  }

  async function handleNewVersion() {
    try {
      const newDiet = await newVersion.mutateAsync()
      setDialog(null)
      router.push(`/pacientes/${patientId}/dietas/${newDiet.id}`)
    } catch {
      // Error handled by mutation state
    }
  }

  return (
    <div className="flex flex-wrap gap-2">
      {diet.status === "draft" && (
        <>
          <Button size="sm" onClick={() => setDialog("publish")}>
            Publicar
          </Button>
          <Button
            size="sm"
            variant="destructive"
            onClick={() => setDialog("delete")}
          >
            Excluir
          </Button>
        </>
      )}

      {diet.status === "published" && (
        <>
          <Button
            size="sm"
            variant="outline"
            onClick={() => setDialog("archive")}
          >
            Arquivar
          </Button>
          <Button size="sm" onClick={() => setDialog("new-version")}>
            Nova Versao
          </Button>
          <ExportPDFButton entityType="diet" entityId={diet.id} />
        </>
      )}

      {diet.status === "archived" && (
        <Button size="sm" onClick={() => setDialog("new-version")}>
          Nova Versao
        </Button>
      )}

      <ConfirmDialog
        open={dialog === "publish"}
        onClose={() => setDialog(null)}
        onConfirm={handlePublish}
        title="Publicar Dieta"
        description="Ao publicar, a dieta ficara visivel para o paciente. Deseja continuar?"
        confirmLabel="Publicar"
        isLoading={publishDiet.isPending}
      />

      <ConfirmDialog
        open={dialog === "archive"}
        onClose={() => setDialog(null)}
        onConfirm={handleArchive}
        title="Arquivar Dieta"
        description="A dieta sera arquivada e nao ficara mais visivel para o paciente. Voce pode criar uma nova versao depois."
        confirmLabel="Arquivar"
        isLoading={archiveDiet.isPending}
      />

      <ConfirmDialog
        open={dialog === "delete"}
        onClose={() => setDialog(null)}
        onConfirm={handleDelete}
        title="Excluir Dieta"
        description="Esta acao nao pode ser desfeita. A dieta e todas as suas refeicoes serao removidas."
        confirmLabel="Excluir"
        variant="destructive"
        isLoading={deleteDiet.isPending}
      />

      <ConfirmDialog
        open={dialog === "new-version"}
        onClose={() => setDialog(null)}
        onConfirm={handleNewVersion}
        title="Nova Versao"
        description="Uma nova versao da dieta sera criada como rascunho, com base nesta versao. Deseja continuar?"
        confirmLabel="Criar Nova Versao"
        isLoading={newVersion.isPending}
      />
    </div>
  )
}
