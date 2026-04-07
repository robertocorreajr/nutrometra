import { withAuth } from "next-auth/middleware"

export function createAuthMiddleware(publicPaths: string[] = []) {
  return withAuth({
    pages: {
      signIn: "/auth/signin",
    },
    callbacks: {
      authorized({ token }) {
        return !!token
      },
    },
  })
}

export function createMiddlewareMatcher(publicPaths: string[] = []) {
  return {
    matcher: [
      "/((?!api/auth|_next/static|_next/image|favicon.ico|auth).*)",
      ...publicPaths,
    ],
  }
}
