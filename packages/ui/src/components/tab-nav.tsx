"use client"

import * as React from "react"
import { cn } from "../lib/utils"

export interface Tab {
  label: string
  href: string
  disabled?: boolean
}

interface TabNavProps {
  tabs: Tab[]
  activePath: string
  LinkComponent?: React.ElementType
}

export function TabNav({ tabs, activePath, LinkComponent = "a" }: TabNavProps) {
  function isActive(href: string): boolean {
    return href === activePath
  }

  return (
    <nav
      aria-label="Navegação por abas"
      className="overflow-x-auto border-b border-border"
    >
      <ul className="flex min-w-max gap-0" role="tablist">
        {tabs.map((tab) => {
          const active = isActive(tab.href)

          if (tab.disabled) {
            return (
              <li key={tab.href} role="presentation">
                <span
                  aria-disabled="true"
                  className={cn(
                    "inline-flex items-center px-4 py-3 text-sm font-medium border-b-2 border-transparent",
                    "cursor-not-allowed text-muted-foreground/50 select-none"
                  )}
                >
                  {tab.label}
                </span>
              </li>
            )
          }

          return (
            <li key={tab.href} role="presentation">
              <LinkComponent
                href={tab.href}
                role="tab"
                aria-selected={active}
                className={cn(
                  "inline-flex items-center px-4 py-3 text-sm font-medium border-b-2 transition-colors",
                  active
                    ? "border-primary text-primary"
                    : "border-transparent text-muted-foreground hover:text-foreground hover:border-border"
                )}
              >
                {tab.label}
              </LinkComponent>
            </li>
          )
        })}
      </ul>
    </nav>
  )
}
