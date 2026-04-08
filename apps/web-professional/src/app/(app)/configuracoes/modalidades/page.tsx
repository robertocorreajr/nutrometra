"use client"

import {
  useProfessionalMe,
  useServiceModes,
  useAddresses,
} from "@nutrometra/api-client/hooks"
import { LoadingState, ErrorState } from "@nutrometra/ui"
import { ServiceModesForm } from "@/components/settings/service-modes-form"

export default function ModalidadesPage() {
  const { data: professional, isLoading: profLoading } = useProfessionalMe()
  const professionalId = professional?.id ?? ""

  const {
    data: serviceModes,
    isLoading: modesLoading,
    isError: modesError,
    refetch: refetchModes,
  } = useServiceModes(professionalId)

  const {
    data: addresses,
    isLoading: addressesLoading,
  } = useAddresses(professionalId)

  if (profLoading || modesLoading || addressesLoading) {
    return <LoadingState lines={4} />
  }

  if (modesError || !professional) {
    return (
      <ErrorState
        message="Erro ao carregar modalidades de atendimento."
        onRetry={refetchModes}
      />
    )
  }

  return (
    <ServiceModesForm
      professionalId={professional.id}
      serviceModes={serviceModes ?? []}
      addresses={addresses ?? []}
    />
  )
}
