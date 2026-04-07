"use client"

import { useSession } from "next-auth/react"

export function useAuth() {
  const { data: session, status } = useSession()
  return {
    user: session?.user ?? null,
    accessToken: session?.accessToken ?? null,
    isAuthenticated: status === "authenticated",
    isLoading: status === "loading",
    tenantId: session?.tenantId ?? null,
    roles: session?.roles ?? [],
  }
}

export function useTenant() {
  const { tenantId } = useAuth()
  return { tenantId }
}

export function usePermission(code: string): boolean {
  // For now, return true — will be wired to useEntitlements in portal implementation
  return true
}
