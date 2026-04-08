"use client"

import { SessionProvider, getSession } from "next-auth/react"
import { ApiProvider, configureApiClient } from "@nutrometra/api-client"

// Configure API client to use Next.js proxy (avoids CORS) with lazy session fetching
configureApiClient({
  baseUrl: "/api/v1",
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
