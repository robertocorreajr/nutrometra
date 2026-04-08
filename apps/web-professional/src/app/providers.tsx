"use client"

import { SessionProvider, getSession } from "next-auth/react"
import { ApiProvider, configureApiClient } from "@nutrometra/api-client"

// Configure API client with lazy session fetching (no React context dependency)
configureApiClient({
  baseUrl: process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080",
  getAccessToken: async () => {
    const session = await getSession()
    return (session as any)?.accessToken ?? null
  },
  getTenantId: () => null,
})

export function Providers({ children }: { children: React.ReactNode }) {
  return (
    <SessionProvider>
      <ApiProvider>{children}</ApiProvider>
    </SessionProvider>
  )
}
