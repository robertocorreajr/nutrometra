"use client"

import Link from "next/link"
import { usePathname } from "next/navigation"
import { TabNav, type Tab } from "@nutrometra/ui"

export default function AssinaturaLayout({ children }: { children: React.ReactNode }) {
  const pathname = usePathname()
  const tabs: Tab[] = [
    { label: "Plano Atual", href: "/assinatura" },
    { label: "Planos", href: "/assinatura/planos" },
    { label: "Faturas", href: "/assinatura/faturas" },
    { label: "Pagamentos", href: "/assinatura/pagamentos" },
  ]
  function getActivePath(): string {
    const exactMatch = tabs.find((t) => t.href === pathname)
    if (exactMatch) return exactMatch.href
    const prefixMatch = tabs
      .filter((t) => pathname.startsWith(t.href + "/"))
      .sort((a, b) => b.href.length - a.href.length)[0]
    return prefixMatch?.href ?? "/assinatura"
  }
  return (
    <div>
      <h1 className="text-2xl font-bold tracking-tight mb-4">Assinatura</h1>
      <TabNav tabs={tabs} activePath={getActivePath()} LinkComponent={Link} />
      <div className="mt-6">{children}</div>
    </div>
  )
}
