"use client"

import { useParams } from "next/navigation"
import { DataTable, LoadingState, ErrorState } from "@nutrometra/ui"
import type { Column } from "@nutrometra/ui"
import { useBackofficeInvoices } from "@nutrometra/api-client/hooks"
import type { Invoice } from "@nutrometra/api-client/types"
import { invoiceStatusLabels, invoiceStatusColors, formatCurrency } from "@/lib/schemas/backoffice"
import { format, parseISO } from "date-fns"

const columns: Column<Invoice>[] = [
  {
    key: "PeriodStart",
    header: "Periodo",
    render: (row) =>
      `${format(parseISO(row.PeriodStart), "dd/MM/yy")} - ${format(parseISO(row.PeriodEnd), "dd/MM/yy")}`,
    searchValue: (row) => row.PeriodStart,
  },
  {
    key: "Status",
    header: "Status",
    render: (row) => (
      <span
        className={`inline-flex items-center rounded-full px-2 py-1 text-xs font-medium ${
          invoiceStatusColors[row.Status] ?? "bg-gray-100 text-gray-700"
        }`}
      >
        {invoiceStatusLabels[row.Status] ?? row.Status}
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
    key: "DueDate",
    header: "Vencimento",
    render: (row) => (row.DueDate ? format(parseISO(row.DueDate), "dd/MM/yyyy") : "—"),
  },
]

export default function FaturasPage() {
  const { id } = useParams<{ id: string }>()
  const { data: invoices, isLoading, isError, refetch } = useBackofficeInvoices(id)

  if (isLoading) return <LoadingState />
  if (isError) return <ErrorState message="Nao foi possivel carregar as faturas." onRetry={refetch} />

  return (
    <DataTable
      columns={columns}
      data={invoices ?? []}
      keyExtractor={(row) => row.ID}
      searchPlaceholder="Buscar faturas..."
      emptyMessage="Nenhuma fatura encontrada."
    />
  )
}
