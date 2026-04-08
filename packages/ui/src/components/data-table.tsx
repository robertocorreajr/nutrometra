"use client"

import * as React from "react"
import { cn } from "../lib/utils"
import { Input } from "./input"

export interface Column<T> {
  key: string
  header: string
  render: (item: T) => React.ReactNode
  sortable?: boolean
  searchValue?: (item: T) => string
  className?: string
}

interface DataTableProps<T> {
  data: T[]
  columns: Column<T>[]
  keyExtractor: (item: T) => string
  searchPlaceholder?: string
  renderMobileCard?: (item: T) => React.ReactNode
  pageSize?: number
  actions?: React.ReactNode
  emptyMessage?: string
}

type SortDir = "asc" | "desc"

export function DataTable<T>({
  data,
  columns,
  keyExtractor,
  searchPlaceholder = "Buscar...",
  renderMobileCard,
  pageSize = 10,
  actions,
  emptyMessage = "Nenhum item encontrado.",
}: DataTableProps<T>) {
  const [search, setSearch] = React.useState("")
  const [sortKey, setSortKey] = React.useState<string | null>(null)
  const [sortDir, setSortDir] = React.useState<SortDir>("asc")
  const [page, setPage] = React.useState(0)

  // Reset page when search changes
  React.useEffect(() => {
    setPage(0)
  }, [search])

  const filteredData = React.useMemo(() => {
    if (!search.trim()) return data
    const term = search.toLowerCase()
    return data.filter((item) =>
      columns.some((col) => {
        if (!col.searchValue) return false
        return col.searchValue(item).toLowerCase().includes(term)
      })
    )
  }, [data, search, columns])

  const sortedData = React.useMemo(() => {
    if (!sortKey) return filteredData
    const col = columns.find((c) => c.key === sortKey)
    if (!col?.searchValue) return filteredData
    return [...filteredData].sort((a, b) => {
      const av = col.searchValue!(a).toLowerCase()
      const bv = col.searchValue!(b).toLowerCase()
      if (av < bv) return sortDir === "asc" ? -1 : 1
      if (av > bv) return sortDir === "asc" ? 1 : -1
      return 0
    })
  }, [filteredData, sortKey, sortDir, columns])

  const totalPages = Math.max(1, Math.ceil(sortedData.length / pageSize))
  const safePage = Math.min(page, totalPages - 1)
  const paginatedData = sortedData.slice(safePage * pageSize, safePage * pageSize + pageSize)

  function handleSort(key: string) {
    if (sortKey === key) {
      setSortDir((d) => (d === "asc" ? "desc" : "asc"))
    } else {
      setSortKey(key)
      setSortDir("asc")
    }
    setPage(0)
  }

  const isEmpty = sortedData.length === 0

  return (
    <div className="space-y-4">
      {/* Toolbar */}
      <div className="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
        <Input
          placeholder={searchPlaceholder}
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          className="max-w-sm"
        />
        {actions && <div className="flex gap-2">{actions}</div>}
      </div>

      {/* Desktop table */}
      <div className="hidden md:block rounded-md border">
        <table className="w-full text-sm">
          <thead>
            <tr className="border-b bg-muted/50">
              {columns.map((col) => (
                <th
                  key={col.key}
                  className={cn(
                    "px-4 py-3 text-left font-medium text-muted-foreground",
                    col.sortable && "cursor-pointer select-none hover:text-foreground",
                    col.className
                  )}
                  onClick={() => col.sortable && handleSort(col.key)}
                >
                  <span className="inline-flex items-center gap-1">
                    {col.header}
                    {col.sortable && sortKey === col.key && (
                      <span className="text-xs">{sortDir === "asc" ? "↑" : "↓"}</span>
                    )}
                  </span>
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {isEmpty ? (
              <tr>
                <td
                  colSpan={columns.length}
                  className="px-4 py-8 text-center text-muted-foreground"
                >
                  {emptyMessage}
                </td>
              </tr>
            ) : (
              paginatedData.map((item) => (
                <tr
                  key={keyExtractor(item)}
                  className="border-b last:border-0 hover:bg-muted/30 transition-colors"
                >
                  {columns.map((col) => (
                    <td key={col.key} className={cn("px-4 py-3", col.className)}>
                      {col.render(item)}
                    </td>
                  ))}
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>

      {/* Mobile cards */}
      <div className="md:hidden space-y-3">
        {isEmpty ? (
          <p className="py-8 text-center text-muted-foreground text-sm">{emptyMessage}</p>
        ) : (
          paginatedData.map((item) => (
            <div key={keyExtractor(item)}>
              {renderMobileCard ? (
                renderMobileCard(item)
              ) : (
                <div className="rounded-md border p-4 space-y-2">
                  {columns.map((col) => (
                    <div key={col.key} className="flex justify-between gap-2 text-sm">
                      <span className="text-muted-foreground font-medium">{col.header}</span>
                      <span className="text-right">{col.render(item)}</span>
                    </div>
                  ))}
                </div>
              )}
            </div>
          ))
        )}
      </div>

      {/* Pagination */}
      {totalPages > 1 && (
        <div className="flex items-center justify-between text-sm">
          <span className="text-muted-foreground">
            Página {safePage + 1} de {totalPages}
          </span>
          <div className="flex gap-2">
            <button
              onClick={() => setPage((p) => Math.max(0, p - 1))}
              disabled={safePage === 0}
              className="px-3 py-1.5 rounded-md border text-sm disabled:opacity-50 disabled:cursor-not-allowed hover:bg-muted/50 transition-colors"
            >
              Anterior
            </button>
            <button
              onClick={() => setPage((p) => Math.min(totalPages - 1, p + 1))}
              disabled={safePage >= totalPages - 1}
              className="px-3 py-1.5 rounded-md border text-sm disabled:opacity-50 disabled:cursor-not-allowed hover:bg-muted/50 transition-colors"
            >
              Próxima
            </button>
          </div>
        </div>
      )}
    </div>
  )
}
