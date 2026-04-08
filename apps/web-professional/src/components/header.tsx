"use client"

import { useAuth } from "@nutrometra/auth"
import { Button } from "@nutrometra/ui"
import { LogOut } from "lucide-react"
import { signOut } from "next-auth/react"

export function Header() {
  const { user } = useAuth()

  return (
    <div className="flex items-center justify-between w-full">
      <span className="font-semibold text-sm">Portal Profissional</span>
      <div className="flex items-center gap-3">
        {user && (
          <span className="text-sm text-muted-foreground hidden sm:inline">
            {user.name ?? user.email}
          </span>
        )}
        <Button
          variant="ghost"
          size="icon"
          onClick={() => signOut({ callbackUrl: "/auth/signin" })}
          aria-label="Sair"
        >
          <LogOut className="h-4 w-4" />
        </Button>
      </div>
    </div>
  )
}
