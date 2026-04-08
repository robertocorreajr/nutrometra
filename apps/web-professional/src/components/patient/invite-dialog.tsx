"use client"

import { useState } from "react"
import { format } from "date-fns"
import { ptBR } from "date-fns/locale"
import { Copy, Check, X } from "lucide-react"
import { Button } from "@nutrometra/ui"
import { useGenerateInvite } from "@nutrometra/api-client/hooks"

interface InviteDialogProps {
  patientId: string
  open: boolean
  onClose: () => void
}

export function InviteDialog({ patientId, open, onClose }: InviteDialogProps) {
  const generateInvite = useGenerateInvite(patientId)
  const [copied, setCopied] = useState(false)

  if (!open) return null

  const invite = generateInvite.data

  function handleGenerate() {
    generateInvite.mutate()
  }

  function handleCopy() {
    if (invite?.code) {
      navigator.clipboard.writeText(invite.code)
      setCopied(true)
      setTimeout(() => setCopied(false), 2000)
    }
  }

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div className="absolute inset-0 bg-black/50" onClick={onClose} />
      <div className="relative z-10 w-full max-w-sm rounded-lg border bg-background shadow-lg p-6">
        <div className="flex items-start justify-between mb-4">
          <h3 className="text-lg font-semibold">Convite do Paciente</h3>
          <button
            type="button"
            onClick={onClose}
            className="text-muted-foreground hover:text-foreground"
            aria-label="Fechar"
          >
            <X className="h-5 w-5" />
          </button>
        </div>

        {!invite ? (
          <div className="space-y-4">
            <p className="text-sm text-muted-foreground">
              Gere um código de convite para o paciente acessar o portal.
            </p>
            <Button
              className="w-full"
              onClick={handleGenerate}
              disabled={generateInvite.isPending}
            >
              {generateInvite.isPending ? "Gerando..." : "Gerar Código de Convite"}
            </Button>
          </div>
        ) : (
          <div className="space-y-4">
            <div className="rounded-lg bg-muted p-4 text-center">
              <p className="text-xs text-muted-foreground mb-1">Código de Convite</p>
              <p className="text-2xl font-mono font-bold tracking-widest">
                {invite.code}
              </p>
            </div>

            <Button
              variant="outline"
              className="w-full"
              onClick={handleCopy}
            >
              {copied ? (
                <><Check className="h-4 w-4 mr-2" /> Copiado!</>
              ) : (
                <><Copy className="h-4 w-4 mr-2" /> Copiar Código</>
              )}
            </Button>

            <p className="text-xs text-muted-foreground text-center">
              Válido até {format(new Date(invite.expires_at), "dd/MM/yyyy 'às' HH:mm", { locale: ptBR })}
            </p>
          </div>
        )}
      </div>
    </div>
  )
}
