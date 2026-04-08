"use client"

import { useState } from "react"
import { format } from "date-fns"
import { ptBR } from "date-fns/locale"
import { Trash2, Calendar } from "lucide-react"
import { Button, Card, CardContent, LoadingState, ErrorState } from "@nutrometra/ui"
import { useBlocks, useDeleteBlock } from "@nutrometra/api-client/hooks"

interface BlockListProps {
  professionalId: string
}

export function BlockList({ professionalId }: BlockListProps) {
  const { data: blocks, isLoading, isError, refetch } = useBlocks(professionalId)
  const deleteBlock = useDeleteBlock(professionalId)
  const [confirmDeleteId, setConfirmDeleteId] = useState<string | null>(null)

  if (isLoading) return <LoadingState lines={4} />
  if (isError) return <ErrorState message="Erro ao carregar bloqueios." onRetry={refetch} />

  const blockList = blocks ?? []

  if (blockList.length === 0) {
    return (
      <p className="text-sm text-muted-foreground text-center py-8">
        Nenhum bloqueio de agenda cadastrado.
      </p>
    )
  }

  // Sort chronologically
  const sorted = [...blockList].sort(
    (a, b) => new Date(a.start_at).getTime() - new Date(b.start_at).getTime()
  )

  async function handleDelete(id: string) {
    try {
      await deleteBlock.mutateAsync(id)
      setConfirmDeleteId(null)
    } catch {
      // Error handled by mutation
    }
  }

  return (
    <div className="space-y-3">
      {sorted.map((block) => {
        const startDate = new Date(block.start_at)
        const endDate = new Date(block.end_at)

        return (
          <Card key={block.id}>
            <CardContent className="py-4">
              <div className="flex items-center justify-between gap-2">
                <div className="flex items-start gap-3 min-w-0">
                  <Calendar className="h-4 w-4 text-muted-foreground mt-0.5 shrink-0" />
                  <div className="min-w-0">
                    <div className="flex items-center gap-2 flex-wrap">
                      <span className="text-sm font-medium">
                        {block.all_day
                          ? format(startDate, "dd/MM/yyyy", { locale: ptBR })
                          : format(startDate, "dd/MM/yyyy HH:mm", { locale: ptBR })}
                        {" - "}
                        {block.all_day
                          ? format(endDate, "dd/MM/yyyy", { locale: ptBR })
                          : format(endDate, "dd/MM/yyyy HH:mm", { locale: ptBR })}
                      </span>
                      {block.all_day && (
                        <span className="inline-flex items-center rounded-full bg-purple-100 text-purple-800 px-2 py-0.5 text-[10px] font-medium">
                          Dia inteiro
                        </span>
                      )}
                    </div>
                    {block.reason && (
                      <p className="text-sm text-muted-foreground mt-0.5 truncate">
                        {block.reason}
                      </p>
                    )}
                  </div>
                </div>
                <div className="shrink-0">
                  {confirmDeleteId === block.id ? (
                    <div className="flex items-center gap-1">
                      <Button
                        variant="destructive"
                        size="sm"
                        disabled={deleteBlock.isPending}
                        onClick={() => handleDelete(block.id)}
                      >
                        {deleteBlock.isPending ? "..." : "Confirmar"}
                      </Button>
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => setConfirmDeleteId(null)}
                      >
                        Não
                      </Button>
                    </div>
                  ) : (
                    <Button
                      variant="ghost"
                      size="icon"
                      onClick={() => setConfirmDeleteId(block.id)}
                    >
                      <Trash2 className="h-4 w-4 text-muted-foreground" />
                    </Button>
                  )}
                </div>
              </div>
            </CardContent>
          </Card>
        )
      })}

      {deleteBlock.isError && (
        <p className="text-sm text-destructive">
          Erro ao excluir bloqueio. Tente novamente.
        </p>
      )}
    </div>
  )
}
