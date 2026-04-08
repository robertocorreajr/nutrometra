"use client"

import { useParams } from "next/navigation"
import { DataTable, LoadingState, ErrorState } from "@nutrometra/ui"
import type { Column } from "@nutrometra/ui"
import { useBackofficePayments } from "@nutrometra/api-client/hooks"
import type { Payment } from "@nutrometra/api-client/types"
import { paymentStatusLabels, paymentStatusColors, formatCurrency } from "@/lib/schemas/backoffice"
import { format, parseISO } from "date-fns"

const columns: Column<Payment>[] = [
  {
    key: "CreatedAt",
    header: "Data",
    render: (row) => format(parseISO(row.CreatedAt), "dd/MM/yyyy HH:mm"),
    searchValue: (row) => row.CreatedAt,
  },
  {
    key: "Status",
    header: "Status",
    render: (row) => (
      <span
        className={`inline-flex items-center rounded-full px-2 py-1 text-xs font-medium ${
          paymentStatusColors[row.Status] ?? "bg-gray-100 text-gray-700"
        }`}
      >
        {paymentStatusLabels[row.Status] ?? row.Status}
      </span>
    ),
    searchValue: (row) => row.Status,
  },
  {
    key: "AmountCents",
    header: "Valor",
    render: (row) => formatCurrency(row.AmountCents, row.Currency),
  },
  {
    key: "FailureReason",
    header: "Motivo Falha",
    render: (row) => row.FailureReason ?? "—",
  },
]

export default function PagamentosPage() {
  const { id } = useParams<{ id: string }>()
  const { data: payments, isLoading, isError, refetch } = useBackofficePayments(id)

  if (isLoading) return <LoadingState />
  if (isError)
    return <ErrorState message="Nao foi possivel carregar os pagamentos." onRetry={refetch} />

  return (
    <DataTable
      columns={columns}
      data={payments ?? []}
      keyExtractor={(row) => row.ID}
      searchPlaceholder="Buscar pagamentos..."
      emptyMessage="Nenhum pagamento encontrado."
    />
  )
}
