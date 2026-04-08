"use client"

import Link from "next/link"
import { usePathname } from "next/navigation"
import { Home, UtensilsCrossed, Activity, FileText, Calendar, User, type LucideIcon } from "lucide-react"
import { cn } from "@nutrometra/ui"

interface NavItem {
  label: string
  href: string
  icon: LucideIcon
}

const navItems: NavItem[] = [
  { label: "Inicio", href: "/", icon: Home },
  { label: "Dietas", href: "/dietas", icon: UtensilsCrossed },
  { label: "Medidas", href: "/medidas", icon: Activity },
  { label: "Docs", href: "/documentos", icon: FileText },
  { label: "Agenda", href: "/agenda", icon: Calendar },
  { label: "Perfil", href: "/perfil", icon: User },
]

export function BottomNav() {
  const pathname = usePathname()

  function isActive(href: string): boolean {
    if (href === "/") return pathname === "/"
    return pathname.startsWith(href)
  }

  return (
    <nav className="fixed bottom-0 left-0 right-0 z-50 border-t bg-background md:hidden">
      <div className="flex justify-around items-center h-16 max-w-lg mx-auto">
        {navItems.map((item) => {
          const Icon = item.icon
          const active = isActive(item.href)

          return (
            <Link
              key={item.href}
              href={item.href}
              className={cn(
                "flex flex-col items-center justify-center gap-1 flex-1 py-2 text-xs transition-colors",
                active
                  ? "text-primary"
                  : "text-muted-foreground"
              )}
            >
              <Icon className="h-5 w-5" />
              <span>{item.label}</span>
            </Link>
          )
        })}
      </div>
    </nav>
  )
}
