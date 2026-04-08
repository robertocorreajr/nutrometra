"use client"

import { SessionProvider } from "next-auth/react"
import { ApiProvider } from "@nutrometra/api-client"

export function Providers({ children }: { children: React.ReactNode }) {
  return (
    <SessionProvider>
      <ApiProvider>
        {children}
      </ApiProvider>
    </SessionProvider>
  )
}
