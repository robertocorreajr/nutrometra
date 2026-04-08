import { withAuth } from "next-auth/middleware"

export default withAuth({
  pages: {
    signIn: "/signin",
  },
})

export const config = {
  matcher: ["/((?!api/auth|api/v1|_next/static|_next/image|favicon.ico|signin).*)"],
}
