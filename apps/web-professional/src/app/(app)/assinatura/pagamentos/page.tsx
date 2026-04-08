"use client"

import { format } from "date-fns"
import { ptBR } from "date-fns/locale"
import { LoadingState, EmptyState } from "@nutrometra/ui"
import { usePayments } from "@nutrometra/api-client/hooks"
import { paymentStatusLabels, paymentStatusColors, formatCurrency } from "@/lib/schemas/billing"

export default function PagamentosPage() {
  const { data: payments, isLoading } = usePayments()

  if (isLoading) return <LoadingState lines={4} />

  const list = payments ?? []
  if (list.length === 0) {
    return <EmptyState title="Nenhum pagamento" description="Nenhum pagamento registrado." />
  }

  return (
    <div className="rounded-lg border">
      <table className="w-full text-sm">
        <thead>
          <tr className="border-b bg-muted/50">
            <th className="text-left p-3 font-medium">Data</th>
            <th className="text-left p-3 font-medium">Status</th>
            <th className="text-right p-3 font-medium">Valor</th>
            <th className="text-left p-3 font-medium">Motivo</th>
          </tr>
        </thead>
        <tbody>
          {list.map((pmt) => {
            const statusLabel = paymentStatusLabels[pmt.Status] ?? pmt.Status
            const statusColor = paymentStatusColors[pmt.Status] ?? "bg-gray-100 text-gray-800"
            return (
              <tr key={pmt.ID} className="border-b last:border-0">
                <td className="p-3">
                  {pmt.PaidAt
                    ? format(new Date(pmt.PaidAt), "dd/MM/yyyy HH:mm", { locale: ptBR })
                    : format(new Date(pmt.CreatedAt), "dd/MM/yyyy HH:mm", { locale: ptBR })}
                </td>
                <td className="p-3">
                  <span className={`inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium ${statusColor}`}>
                    {statusLabel}
                  </span>
                </td>
                <td className="p-3 text-right">{formatCurrency(pmt.AmountCents, pmt.Currency)}</td>
                <td className="p-3 text-muted-foreground">{pmt.FailureReason ?? "—"}</td>
              </tr>
            )
          })}
        </tbody>
      </table>
    </div>
  )
}
