import type { DefaultSession } from "next-auth"
import type { DefaultJWT } from "next-auth/jwt"

declare module "next-auth" {
  interface User {
    id: string
  }

  interface Session extends DefaultSession {
    user: DefaultSession["user"] & { id: string }
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
