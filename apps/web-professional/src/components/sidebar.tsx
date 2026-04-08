"use client"

import Link from "next/link"
import { usePathname } from "next/navigation"
import {
  LayoutDashboard,
  Users,
  UsersRound,
  Calendar,
  UtensilsCrossed,
  FileText,
  CreditCard,
  Settings,
  type LucideIcon,
} from "lucide-react"
import { cn } from "@nutrometra/ui"

interface NavItem {
  label: string
  href: string
  icon: LucideIcon
  disabled?: boolean
}

const navItems: NavItem[] = [
  { label: "Dashboard", href: "/", icon: LayoutDashboard },
  { label: "Pacientes", href: "/pacientes", icon: Users },
  { label: "Equipe", href: "/equipe", icon: UsersRound },
  { label: "Agenda", href: "/agenda", icon: Calendar },
  { label: "Dietas", href: "/dietas", icon: UtensilsCrossed },
  { label: "Documentos", href: "/documentos", icon: FileText },
  { label: "Assinatura", href: "/assinatura", icon: CreditCard },
  { label: "Configurações", href: "/configuracoes", icon: Settings },
]

export function Sidebar() {
  const pathname = usePathname()

  function isActive(href: string): boolean {
    if (href === "/") return pathname === "/"
    return pathname.startsWith(href)
  }

  return (
    <nav className="flex flex-col h-full">
      {/* Logo */}
      <div className="p-4 border-b">
        <Link href="/" className="flex items-center gap-2">
          <div className="h-8 w-8 rounded-lg bg-primary flex items-center justify-center">
            <span className="text-primary-foreground font-bold text-sm">N</span>
          </div>
          <span className="font-semibold text-lg">Nutrometra</span>
        </Link>
      </div>

      {/* Menu items */}
      <div className="flex-1 py-4 px-3 space-y-1">
        {navItems.map((item) => {
          const Icon = item.icon
          const active = isActive(item.href)

          if (item.disabled) {
            return (
              <div
                key={item.href}
                className="flex items-center gap-3 px-3 py-2 rounded-md text-sm text-muted-foreground/50 cursor-not-allowed"
                title="Em breve"
              >
                <Icon className="h-4 w-4" />
                <span>{item.label}</span>
                <span className="ml-auto text-xs bg-muted px-1.5 py-0.5 rounded">Em breve</span>
              </div>
            )
          }

          return (
            <Link
              key={item.href}
              href={item.href}
              className={cn(
                "flex items-center gap-3 px-3 py-2 rounded-md text-sm font-medium transition-colors",
                active
                  ? "bg-primary/10 text-primary"
                  : "text-muted-foreground hover:bg-accent hover:text-accent-foreground"
              )}
            >
              <Icon className="h-4 w-4" />
              <span>{item.label}</span>
            </Link>
          )
        })}
      </div>

      {/* Footer */}
      <div className="p-4 border-t">
        <p className="text-xs text-muted-foreground">Nutrometra v1.0</p>
      </div>
    </nav>
  )
}
