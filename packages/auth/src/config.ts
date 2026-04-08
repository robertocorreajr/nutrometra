import type { NextAuthOptions } from "next-auth"

export function createAuthOptions(overrides?: Partial<NextAuthOptions>): NextAuthOptions {
  return {
    providers: [
      {
        id: "zitadel",
        name: "Zitadel",
        type: "oauth",
        wellKnown: `${process.env.ZITADEL_ISSUER}/.well-known/openid-configuration`,
        clientId: process.env.ZITADEL_CLIENT_ID,
        clientSecret: process.env.ZITADEL_CLIENT_SECRET || "",
        authorization: {
          params: {
            scope: "openid profile email",
          },
        },
        idToken: true,
        checks: ["pkce", "state"],
        profile(profile) {
          return {
            id: profile.sub,
            name: profile.name ?? profile.preferred_username,
            email: profile.email,
          }
        },
      },
    ],
    callbacks: {
      async jwt({ token, account }) {
        if (account) {
          token.accessToken = account.access_token
          token.refreshToken = account.refresh_token
          token.expiresAt = account.expires_at

          // Resolve tenant on first login by calling the Go API
          const apiUrl = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8081"
          try {
            const meRes = await fetch(`${apiUrl}/auth/me`, {
              headers: { Authorization: `Bearer ${account.access_token}` },
            })
            if (meRes.ok) {
              const me = await meRes.json()
              token.userId = me.user_id

              // Fetch user's tenants directly
              const tenantsRes = await fetch(`${apiUrl}/auth/me/tenants`, {
                headers: { Authorization: `Bearer ${account.access_token}` },
              })
              if (tenantsRes.ok) {
                const tenants = await tenantsRes.json()
                if (Array.isArray(tenants) && tenants.length > 0) {
                  token.tenantId = tenants[0].tenant_id
                }
              }
            }
          } catch {
            // API may not be available — tenant will be resolved later
          }
        }
        return token
      },
      async session({ session, token }) {
        session.accessToken = token.accessToken as string | undefined
        session.tenantId = token.tenantId as string | undefined
        session.roles = token.roles as string[] | undefined
        return session
      },
    },
    pages: {
      signIn: "/signin",
    },
    session: {
      strategy: "jwt",
    },
    ...overrides,
  }
}
