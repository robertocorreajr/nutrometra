import type { DefaultSession, DefaultJWT } from "next-auth"

declare module "next-auth" {
  interface Session extends DefaultSession {
    accessToken?: string
    tenantId?: string
    roles?: string[]
  }
}

declare module "next-auth/jwt" {
  interface JWT extends DefaultJWT {
    accessToken?: string
    refreshToken?: string
    expiresAt?: number
    tenantId?: string
    roles?: string[]
  }
}

export type {} // ensure this is treated as a module
