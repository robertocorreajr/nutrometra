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
