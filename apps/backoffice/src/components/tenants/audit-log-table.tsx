"use client"

import { DataTable } from "@nutrometra/ui"
import type { Column } from "@nutrometra/ui"
import type { AuditLogEntry } from "@nutrometra/api-client/types"
import { format, parseISO } from "date-fns"

const columns: Column<AuditLogEntry>[] = [
  {
    key: "created_at",
    header: "Data",
    render: (row) => format(parseISO(row.created_at), "dd/MM/yyyy HH:mm"),
    searchValue: (row) => row.created_at,
  },
  {
    key: "actor_role",
    header: "Ator",
    render: (row) => (
      <span className="text-xs">
        <span className="font-medium">{row.actor_role}</span>
        <br />
        <span className="text-muted-foreground">{row.actor_id.slice(0, 8)}...</span>
      </span>
    ),
    searchValue: (row) => row.actor_role,
  },
  {
    key: "action",
    header: "Acao",
    render: (row) => row.action,
    searchValue: (row) => row.action,
  },
  {
    key: "resource_type",
    header: "Recurso",
    render: (row) => (
      <span className="font-mono text-xs">
        {row.resource_type}/{row.resource_id.slice(0, 8)}
      </span>
    ),
    searchValue: (row) => row.resource_type,
  },
]

interface AuditLogTableProps {
  entries: AuditLogEntry[]
}

export function AuditLogTable({ entries }: AuditLogTableProps) {
  return (
    <DataTable
      columns={columns}
      data={entries}
      keyExtractor={(row) => row.id}
      searchPlaceholder="Buscar auditoria..."
      emptyMessage="Nenhum registro de auditoria."
    />
  )
}
