"use client"

import { useState } from "react"
import { useRouter } from "next/navigation"
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@nutrometra/ui"
import { useActivateInvite } from "@nutrometra/api-client/hooks"
import { ActivateForm } from "@/components/invite/activate-form"

export default function AtivarPage() {
  const router = useRouter()
  const activateInvite = useActivateInvite()
  const [error, setError] = useState<string>()
  const [success, setSuccess] = useState<{ professional_name: string } | null>(null)

  function handleActivate(code: string) {
    setError(undefined)
    activateInvite.mutate(
      { code },
      {
        onSuccess: (data) => {
          setSuccess({ professional_name: data.professional_name })
          setTimeout(() => router.push("/"), 2000)
        },
        onError: () => {
          setError("Codigo invalido ou expirado. Verifique com seu nutricionista.")
        },
      }
    )
  }

  return (
    <div className="flex items-center justify-center min-h-[60vh] px-4">
      <Card className="w-full max-w-sm">
        <CardHeader className="text-center">
          <CardTitle>Ativar Convite</CardTitle>
          <CardDescription>
            Insira o codigo de convite fornecido pelo seu nutricionista
          </CardDescription>
        </CardHeader>
        <CardContent>
          {success ? (
            <div className="text-center space-y-2">
              <div className="h-12 w-12 rounded-full bg-green-100 flex items-center justify-center mx-auto">
                <svg className="h-6 w-6 text-green-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                </svg>
              </div>
              <p className="font-medium">Convite ativado!</p>
              <p className="text-sm text-muted-foreground">
                Voce foi vinculado ao profissional {success.professional_name}
              </p>
              <p className="text-xs text-muted-foreground">Redirecionando...</p>
            </div>
          ) : (
            <ActivateForm
              onActivate={handleActivate}
              isPending={activateInvite.isPending}
              error={error}
            />
          )}
        </CardContent>
      </Card>
    </div>
  )
}
