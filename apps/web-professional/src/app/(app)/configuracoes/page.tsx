"use client"

import { useProfessionalMe } from "@nutrometra/api-client/hooks"
import { LoadingState, ErrorState } from "@nutrometra/ui"
import { ProfileForm } from "@/components/settings/profile-form"

export default function PerfilPage() {
  const { data: professional, isLoading, isError, refetch } = useProfessionalMe()

  if (isLoading) return <LoadingState lines={6} />
  if (isError || !professional) return <ErrorState message="Erro ao carregar perfil." onRetry={refetch} />

  return <ProfileForm professional={professional} />
}
