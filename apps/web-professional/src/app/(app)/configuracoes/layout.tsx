"use client"

import Link from "next/link"
import { usePathname } from "next/navigation"
import { TabNav, type Tab } from "@nutrometra/ui"

export default function ConfiguracoesLayout({ children }: { children: React.ReactNode }) {
  const pathname = usePathname()

  const tabs: Tab[] = [
    { label: "Perfil", href: "/configuracoes" },
    { label: "Endereços", href: "/configuracoes/enderecos" },
    { label: "Modalidades", href: "/configuracoes/modalidades" },
  ]

  function getActivePath(): string {
    const exactMatch = tabs.find((t) => t.href === pathname)
    if (exactMatch) return exactMatch.href

    const prefixMatch = tabs
      .filter((t) => pathname.startsWith(t.href + "/"))
      .sort((a, b) => b.href.length - a.href.length)[0]

    return prefixMatch?.href ?? "/configuracoes"
  }

  return (
    <div>
      <h1 className="text-2xl font-bold tracking-tight mb-4">Configurações</h1>
      <TabNav tabs={tabs} activePath={getActivePath()} LinkComponent={Link} />
      <div className="mt-6">{children}</div>
    </div>
  )
}
