"use client"

import { signIn } from "next-auth/react"
import { Button } from "@nutrometra/ui"

export default function SignInPage() {
  return (
    <div className="min-h-screen flex items-center justify-center bg-background">
      <div className="w-full max-w-sm mx-auto p-6">
        <div className="text-center mb-8">
          <div className="h-12 w-12 rounded-xl bg-primary flex items-center justify-center mx-auto mb-4">
            <span className="text-primary-foreground font-bold text-xl">N</span>
          </div>
          <h1 className="text-2xl font-bold">Nutrometra</h1>
          <p className="text-muted-foreground mt-1">Portal Profissional</p>
        </div>

        <Button
          className="w-full"
          size="lg"
          onClick={() => signIn("zitadel", { callbackUrl: "/" })}
        >
          Entrar com sua conta
        </Button>

        <p className="text-xs text-muted-foreground text-center mt-6">
          Ao entrar, você concorda com nossos termos de uso e política de privacidade.
        </p>
      </div>
    </div>
  )
}
