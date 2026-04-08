"use client"

import { format } from "date-fns"
import { ptBR } from "date-fns/locale"
import { ExternalLink } from "lucide-react"
import { LoadingState, ErrorState, EmptyState } from "@nutrometra/ui"
import { useInvoices } from "@nutrometra/api-client/hooks"
import { invoiceStatusLabels, invoiceStatusColors, formatCurrency } from "@/lib/schemas/billing"

export default function FaturasPage() {
  const { data: invoices, isLoading, isError, refetch } = useInvoices()

  if (isLoading) return <LoadingState lines={4} />
  if (isError) return <ErrorState message="Não foi possível carregar as faturas." onRetry={refetch} />

  const list = invoices ?? []
  if (list.length === 0) {
    return <EmptyState title="Nenhuma fatura" description="Nenhuma fatura disponível." />
  }

  return (
    <div className="rounded-lg border">
      <table className="w-full text-sm">
        <thead>
          <tr className="border-b bg-muted/50">
            <th className="text-left p-3 font-medium">Período</th>
            <th className="text-left p-3 font-medium">Status</th>
            <th className="text-right p-3 font-medium">Valor</th>
            <th className="text-right p-3 font-medium"></th>
          </tr>
        </thead>
        <tbody>
          {list.map((inv) => {
            const statusLabel = invoiceStatusLabels[inv.Status] ?? inv.Status
            const statusColor = invoiceStatusColors[inv.Status] ?? "bg-gray-100 text-gray-800"
            return (
              <tr key={inv.ID} className="border-b last:border-0">
                <td className="p-3">
                  {format(new Date(inv.PeriodStart), "dd/MM/yyyy", { locale: ptBR })} — {format(new Date(inv.PeriodEnd), "dd/MM/yyyy", { locale: ptBR })}
                </td>
                <td className="p-3">
                  <span className={`inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium ${statusColor}`}>
                    {statusLabel}
                  </span>
                </td>
                <td className="p-3 text-right">{formatCurrency(inv.AmountCents, inv.Currency)}</td>
                <td className="p-3 text-right">
                  {inv.HostedURL && (
                    <a href={inv.HostedURL} target="_blank" rel="noopener noreferrer" className="text-primary hover:underline inline-flex items-center gap-1">
                      Ver <ExternalLink className="h-3 w-3" />
                    </a>
                  )}
                </td>
              </tr>
            )
          })}
        </tbody>
      </table>
    </div>
  )
}
