"use client"

import { useAuth } from "@nutrometra/auth"

export function Header() {
  const { user } = useAuth()

  return (
    <div className="flex items-center justify-between w-full">
      <div className="flex items-center gap-2">
        <div className="h-7 w-7 rounded-lg bg-primary flex items-center justify-center">
          <span className="text-primary-foreground font-bold text-xs">N</span>
        </div>
        <span className="font-semibold text-sm">Nutrometra</span>
      </div>
      {user && (
        <span className="text-sm text-muted-foreground truncate max-w-[150px]">
          {user.name ?? user.email}
        </span>
      )}
    </div>
  )
}
