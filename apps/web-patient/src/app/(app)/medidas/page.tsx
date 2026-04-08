"use client"

import { useAuth } from "@nutrometra/auth"
import { PageHeader, LoadingState, ErrorState, EmptyState } from "@nutrometra/ui"
import { useMyMeasurements } from "@nutrometra/api-client/hooks"
import { LatestMeasurement } from "@/components/measurements/latest-measurement"
import { MeasurementHistory } from "@/components/measurements/measurement-history"

export default function MedidasPage() {
  const { patientId } = useAuth()
  const { data: measurements, isLoading, isError, refetch } = useMyMeasurements(patientId ?? "")

  return (
    <div className="max-w-lg mx-auto md:max-w-none space-y-4">
      <PageHeader title="Medidas" description="Sua composicao corporal" />
      {isLoading && <LoadingState />}
      {isError && <ErrorState message="Nao foi possivel carregar as medidas." onRetry={refetch} />}
      {measurements && measurements.length === 0 && <EmptyState title="Nenhuma medida" description="Nenhuma avaliacao registrada ainda." />}
      {measurements && measurements.length > 0 && (
        <>
          <LatestMeasurement measurement={[...measurements].sort((a, b) => b.measured_at.localeCompare(a.measured_at))[0]} />
          {measurements.length > 1 && <MeasurementHistory measurements={measurements} />}
        </>
      )}
    </div>
  )
}
