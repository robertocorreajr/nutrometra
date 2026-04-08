"use client"

import { useState } from "react"
import { format } from "date-fns"
import { ptBR } from "date-fns/locale"
import { Calendar } from "lucide-react"
import { Button } from "@nutrometra/ui"
import {
  useGoogleCalendarStatus,
  useGoogleCalendarAuthorize,
  useGoogleCalendarDisconnect,
} from "@nutrometra/api-client/hooks"
import { ConfirmDialog } from "@/components/shared/confirm-dialog"

export function GoogleCalendarCard() {
  const { data: status, isLoading } = useGoogleCalendarStatus()
  const authorize = useGoogleCalendarAuthorize()
  const disconnect = useGoogleCalendarDisconnect()
  const [showDisconnect, setShowDisconnect] = useState(false)

  if (isLoading) {
    return (
      <div className="rounded-lg border bg-card p-6 animate-pulse">
        <div className="h-6 w-48 bg-muted rounded" />
      </div>
    )
  }

  const connected = status?.connected ?? false

  function handleConnect() {
    authorize.mutate(undefined, {
      onSuccess: (data) => {
        window.location.href = data.authorize_url
      },
    })
  }

  return (
    <div className="rounded-lg border bg-card p-6">
      <div className="flex items-start gap-4">
        <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-primary/10">
          <Calendar className="h-5 w-5 text-primary" />
        </div>
        <div className="flex-1">
          <h3 className="font-semibold">Google Calendar</h3>
          <p className="text-sm text-muted-foreground mt-1">
            Sincronize suas consultas com o Google Calendar automaticamente.
          </p>

          {connected && status?.synced_at && (
            <p className="text-xs text-muted-foreground mt-2">
              Última sincronização: {format(new Date(status.synced_at), "dd/MM/yyyy 'às' HH:mm", { locale: ptBR })}
            </p>
          )}

          <div className="mt-4">
            {connected ? (
              <div className="flex items-center gap-3">
                <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">
                  Conectado
                </span>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => setShowDisconnect(true)}
                >
                  Desconectar
                </Button>
              </div>
            ) : (
              <Button onClick={handleConnect} disabled={authorize.isPending}>
                {authorize.isPending ? "Conectando..." : "Conectar Google Calendar"}
              </Button>
            )}
          </div>
        </div>
      </div>

      <ConfirmDialog
        open={showDisconnect}
        onClose={() => setShowDisconnect(false)}
        onConfirm={() => {
          disconnect.mutate(undefined, {
            onSuccess: () => setShowDisconnect(false),
          })
        }}
        title="Desconectar Google Calendar"
        description="Suas consultas não serão mais sincronizadas automaticamente. Você poderá reconectar a qualquer momento."
        confirmLabel="Desconectar"
        variant="destructive"
        isLoading={disconnect.isPending}
      />
    </div>
  )
}
