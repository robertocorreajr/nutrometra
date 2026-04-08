"use client"

import { useState } from "react"
import { useRouter } from "next/navigation"
import Link from "next/link"
import { ArrowLeft } from "lucide-react"
import { useCreatePatient } from "@nutrometra/api-client/hooks"
import { useAuth } from "@nutrometra/auth"
import { PageHeader, Card, CardContent } from "@nutrometra/ui"
import { PatientForm } from "@/components/patient/patient-form"
import type { PatientFormValues } from "@/lib/schemas/patient"

export default function NovoPacientePage() {
  const router = useRouter()
  const { user } = useAuth()
  const { mutateAsync: createPatient, isPending } = useCreatePatient()
  const [errorMessage, setErrorMessage] = useState<string | null>(null)

  async function handleSubmit(values: PatientFormValues) {
    setErrorMessage(null)
    try {
      const result = await createPatient({
        professional_id: user?.id ?? "",
        full_name: values.full_name,
        email: values.email ?? "",
        phone: values.phone ?? "",
        cpf: values.cpf ?? "",
        date_of_birth: values.date_of_birth || undefined,
        gender: values.gender ?? "prefer_not_to_say",
        notes: values.notes ?? "",
      })
      router.push(`/pacientes/${result.id}`)
    } catch {
      setErrorMessage(
        "Não foi possível cadastrar o paciente. Verifique os dados e tente novamente."
      )
    }
  }

  return (
    <div>
      <div className="mb-4">
        <Link
          href="/pacientes"
          className="inline-flex items-center gap-1 text-sm text-muted-foreground hover:text-foreground"
        >
          <ArrowLeft className="h-4 w-4" />
          Voltar para pacientes
        </Link>
      </div>

      <PageHeader
        title="Novo Paciente"
        description="Preencha os dados cadastrais do paciente."
      />

      <Card className="max-w-xl">
        <CardContent className="pt-6">
          {errorMessage && (
            <div
              role="alert"
              className="mb-4 rounded-md bg-destructive/10 px-4 py-3 text-sm text-destructive"
            >
              {errorMessage}
            </div>
          )}
          <PatientForm
            onSubmit={handleSubmit}
            isSubmitting={isPending}
            submitLabel="Cadastrar Paciente"
          />
        </CardContent>
      </Card>
    </div>
  )
}
