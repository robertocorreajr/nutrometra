"use client"

import { SessionProvider, getSession } from "next-auth/react"
import { ApiProvider, configureApiClient } from "@nutrometra/api-client"

let cachedTenantId: string | null = null

// Configure API client to use Next.js proxy (avoids CORS) with lazy session fetching
configureApiClient({
  baseUrl: "/api/v1",
  getAccessToken: async () => {
    const session = await getSession()
    cachedTenantId = (session as any)?.tenantId ?? null
    return (session as any)?.accessToken ?? null
  },
  getTenantId: () => cachedTenantId,
})

export function Providers({ children }: { children: React.ReactNode }) {
  return (
    <SessionProvider>
      <ApiProvider>{children}</ApiProvider>
    </SessionProvider>
  )
}
