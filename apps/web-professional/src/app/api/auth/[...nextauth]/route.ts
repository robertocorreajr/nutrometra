import NextAuth from "next-auth"
import { createAuthOptions } from "@nutrometra/auth"

const handler = NextAuth(createAuthOptions())

export { handler as GET, handler as POST }
