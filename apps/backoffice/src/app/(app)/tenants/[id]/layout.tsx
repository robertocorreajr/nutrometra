"use client"

import Link from "next/link"
import { usePathname, useParams } from "next/navigation"
import { TabNav, type Tab, PageHeader, LoadingState, ErrorState } from "@nutrometra/ui"
import { useBackofficeTenant } from "@nutrometra/api-client/hooks"

export default function TenantDetailLayout({ children }: { children: React.ReactNode }) {
  const { id } = useParams<{ id: string }>()
  const pathname = usePathname()
  const { data: tenant, isLoading, isError } = useBackofficeTenant(id)
  const basePath = `/tenants/${id}`
  const tabs: Tab[] = [
    { label: "Resumo", href: basePath },
    { label: "Assinatura", href: `${basePath}/assinatura` },
    { label: "Faturas", href: `${basePath}/faturas` },
    { label: "Pagamentos", href: `${basePath}/pagamentos` },
    { label: "Overrides", href: `${basePath}/overrides` },
    { label: "Auditoria", href: `${basePath}/auditoria` },
  ]

  function getActivePath(): string {
    const exactMatch = tabs.find((t) => t.href === pathname)
    if (exactMatch) return exactMatch.href
    const prefixMatch = tabs
      .filter((t) => pathname.startsWith(t.href + "/"))
      .sort((a, b) => b.href.length - a.href.length)[0]
    return prefixMatch?.href ?? basePath
  }

  if (isLoading) return <LoadingState />
  if (isError || !tenant) return <ErrorState message="Nao foi possivel carregar o tenant." />

  return (
    <div>
      <PageHeader
        title={tenant.tenant_name}
        description={`Slug: ${tenant.slug} · ${tenant.owner_email}`}
      />
      <TabNav tabs={tabs} activePath={getActivePath()} LinkComponent={Link} />
      <div className="mt-6">{children}</div>
    </div>
  )
}
