"use client"

import { signIn } from "next-auth/react"
import { Button, Card, CardContent, CardHeader, CardTitle, CardDescription } from "@nutrometra/ui"

export default function SignInPage() {
  return (
    <div className="min-h-screen flex items-center justify-center bg-muted/30 p-4">
      <Card className="w-full max-w-sm">
        <CardHeader className="text-center">
          <div className="mx-auto mb-4 h-12 w-12 rounded-xl bg-primary flex items-center justify-center">
            <span className="text-primary-foreground font-bold text-xl">N</span>
          </div>
          <CardTitle>Backoffice</CardTitle>
          <CardDescription>Painel administrativo Nutrometra</CardDescription>
        </CardHeader>
        <CardContent>
          <Button className="w-full" onClick={() => signIn("zitadel", { callbackUrl: "/" })}>
            Entrar com SSO
          </Button>
        </CardContent>
      </Card>
    </div>
  )
}
