"use client"

import Link from "next/link"
import { PageHeader, DataTable, LoadingState, ErrorState, EmptyState } from "@nutrometra/ui"
import type { Column } from "@nutrometra/ui"
import { useBackofficeTenants } from "@nutrometra/api-client/hooks"
import type { BackofficeTenant } from "@nutrometra/api-client/types"
import { tenantStatusLabels, tenantStatusColors } from "@/lib/schemas/backoffice"
import { format, parseISO } from "date-fns"

const columns: Column<BackofficeTenant>[] = [
  {
    key: "tenant_name",
    header: "Nome",
    render: (row) => row.tenant_name,
    searchValue: (row) => row.tenant_name,
    sortable: true,
  },
  {
    key: "slug",
    header: "Slug",
    render: (row) => row.slug,
    searchValue: (row) => row.slug,
  },
  {
    key: "owner_email",
    header: "Email",
    render: (row) => row.owner_email,
    searchValue: (row) => row.owner_email,
  },
  {
    key: "status",
    header: "Status",
    render: (row) => (
      <span
        className={`inline-flex items-center rounded-full px-2 py-1 text-xs font-medium ${tenantStatusColors[row.status] ?? "bg-gray-100 text-gray-700"}`}
      >
        {tenantStatusLabels[row.status] ?? row.status}
      </span>
    ),
  },
  {
    key: "created_at",
    header: "Criado em",
    render: (row) => format(parseISO(row.created_at), "dd/MM/yyyy"),
  },
  {
    key: "actions",
    header: "",
    render: (row) => (
      <Link
        href={`/tenants/${row.id}`}
        className="text-sm font-medium text-primary hover:underline"
      >
        Ver
      </Link>
    ),
  },
]

export default function TenantsPage() {
  const { data: tenants, isLoading, isError, refetch } = useBackofficeTenants()

  return (
    <>
      <PageHeader title="Tenants" description="Gerenciar tenants do sistema" />
      {isLoading && <LoadingState />}
      {isError && (
        <ErrorState
          message="Nao foi possivel carregar os tenants."
          onRetry={refetch}
        />
      )}
      {tenants && tenants.length === 0 && (
        <EmptyState title="Nenhum tenant cadastrado" />
      )}
      {tenants && tenants.length > 0 && (
        <DataTable<BackofficeTenant>
          columns={columns}
          data={tenants}
          keyExtractor={(tenant) => tenant.id}
          searchPlaceholder="Buscar tenant..."
          emptyMessage="Nenhum tenant encontrado."
        />
      )}
    </>
  )
}
