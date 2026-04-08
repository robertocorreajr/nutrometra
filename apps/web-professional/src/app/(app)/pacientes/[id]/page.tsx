"use client"

import { useParams } from "next/navigation"
import { format } from "date-fns"
import { ptBR } from "date-fns/locale"
import { usePatient, usePatientProfile } from "@nutrometra/api-client/hooks"
import { LoadingState, ErrorState, Card, CardHeader, CardTitle, CardContent } from "@nutrometra/ui"

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

const genderLabels: Record<string, string> = {
  male: "Masculino",
  female: "Feminino",
  other: "Outro",
  prefer_not_to_say: "Prefiro não dizer",
}

function formatDate(dateStr: string | undefined): string {
  if (!dateStr) return "—"
  try {
    return format(new Date(dateStr), "dd/MM/yyyy", { locale: ptBR })
  } catch {
    return dateStr
  }
}

function joinList(items: string[] | undefined): string {
  if (!items || items.length === 0) return "—"
  return items.join(", ")
}

// ---------------------------------------------------------------------------
// InfoRow
// ---------------------------------------------------------------------------

interface InfoRowProps {
  label: string
  value: string | undefined | null
}

function InfoRow({ label, value }: InfoRowProps) {
  return (
    <div className="flex flex-col gap-0.5 py-2 sm:flex-row sm:gap-4">
      <dt className="text-sm font-medium text-muted-foreground sm:w-40 shrink-0">
        {label}
      </dt>
      <dd className="text-sm text-foreground break-words">{value || "—"}</dd>
    </div>
  )
}

// ---------------------------------------------------------------------------
// Page
// ---------------------------------------------------------------------------

export default function PatientResumePage() {
  const params = useParams<{ id: string }>()
  const patientId = params.id

  const {
    data: patient,
    isLoading: patientLoading,
    isError: patientError,
    refetch: refetchPatient,
  } = usePatient(patientId)

  const {
    data: profile,
    isLoading: profileLoading,
    isError: profileError,
    refetch: refetchProfile,
  } = usePatientProfile(patientId)

  if (patientLoading || profileLoading) return <LoadingState lines={6} />

  if (patientError) {
    return (
      <ErrorState
        message="Não foi possível carregar os dados cadastrais."
        onRetry={refetchPatient}
      />
    )
  }

  if (profileError) {
    return (
      <ErrorState
        message="Não foi possível carregar o perfil clínico."
        onRetry={refetchProfile}
      />
    )
  }

  return (
    <div className="grid gap-6 md:grid-cols-2">
      {/* Card 1: Dados Cadastrais */}
      <Card>
        <CardHeader>
          <CardTitle className="text-base">Dados Cadastrais</CardTitle>
        </CardHeader>
        <CardContent>
          <dl className="divide-y divide-border">
            <InfoRow label="Nome" value={patient?.full_name} />
            <InfoRow label="E-mail" value={patient?.email} />
            <InfoRow label="Telefone" value={patient?.phone} />
            <InfoRow label="CPF" value={patient?.cpf} />
            <InfoRow
              label="Data de nascimento"
              value={formatDate(patient?.date_of_birth)}
            />
            <InfoRow
              label="Gênero"
              value={
                patient?.gender ? genderLabels[patient.gender] ?? patient.gender : undefined
              }
            />
          </dl>
        </CardContent>
      </Card>

      {/* Card 2: Perfil Clínico */}
      <Card>
        <CardHeader>
          <CardTitle className="text-base">Perfil Clínico</CardTitle>
        </CardHeader>
        <CardContent>
          <dl className="divide-y divide-border">
            <InfoRow label="Ocupação" value={profile?.occupation} />
            <InfoRow label="Estado civil" value={profile?.marital_status} />
            <InfoRow label="Tipo sanguíneo" value={profile?.blood_type} />
            <InfoRow
              label="Alergias"
              value={joinList(profile?.allergies)}
            />
            <InfoRow
              label="Condições"
              value={joinList(profile?.chronic_conditions)}
            />
            <InfoRow
              label="Medicamentos"
              value={joinList(profile?.medications)}
            />
            <InfoRow
              label="Contato de emergência"
              value={
                profile?.emergency_contact_name
                  ? `${profile.emergency_contact_name}${profile.emergency_contact_phone ? ` — ${profile.emergency_contact_phone}` : ""}`
                  : undefined
              }
            />
          </dl>
        </CardContent>
      </Card>
    </div>
  )
}
