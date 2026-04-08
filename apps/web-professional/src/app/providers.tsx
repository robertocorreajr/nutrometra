"use client"

import { SessionProvider, useSession } from "next-auth/react"
import { ApiProvider, configureApiClient } from "@nutrometra/api-client"
import { useEffect } from "react"

function ApiClientConfigurator({ children }: { children: React.ReactNode }) {
  const { data: session } = useSession()

  useEffect(() => {
    configureApiClient({
      baseUrl: process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080",
      getAccessToken: async () => session?.accessToken ?? null,
      getTenantId: () => session?.tenantId ?? null,
    })
  }, [session])

  return <>{children}</>
}

export function Providers({ children }: { children: React.ReactNode }) {
  return (
    <SessionProvider>
      <ApiProvider>
        <ApiClientConfigurator>{children}</ApiClientConfigurator>
      </ApiProvider>
    </SessionProvider>
  )
}
