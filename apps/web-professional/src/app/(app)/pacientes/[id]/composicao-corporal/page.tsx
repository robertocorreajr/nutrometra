"use client"

import { useParams } from "next/navigation"
import { useState } from "react"
import { LoadingState, ErrorState, EmptyState, Button } from "@nutrometra/ui"
import { useMeasurements, useCreateMeasurement, useProfessionalMe } from "@nutrometra/api-client/hooks"
import { MeasurementForm } from "@/components/bioimpedance/measurement-form"
import { MeasurementList } from "@/components/bioimpedance/measurement-list"
import type { MeasurementFormValues } from "@/lib/schemas/bioimpedance"
import { Plus, Scale } from "lucide-react"

export default function ComposicaoCorporalPage() {
  const params = useParams()
  const patientId = params.id as string
  const { data: professional } = useProfessionalMe()

  const { data: measurements, isLoading, isError, refetch } = useMeasurements(patientId)
  const createMeasurement = useCreateMeasurement(patientId)

  const [showForm, setShowForm] = useState(false)

  async function handleSubmit(data: MeasurementFormValues) {
    // Filter out undefined and NaN values before sending
    const cleaned: Record<string, unknown> = {
      patient_id: patientId,
      professional_id: professional?.id ?? "",
      source: data.source,
    }

    const numericKeys = [
      "weight_kg", "height_cm", "body_fat_pct", "lean_mass_kg", "fat_mass_kg",
      "muscle_mass_kg", "bone_mass_kg", "water_pct", "visceral_fat", "basal_metabolic_rate",
      "waist_cm", "hip_cm", "chest_cm", "right_arm_cm", "left_arm_cm",
      "right_thigh_cm", "left_thigh_cm", "right_calf_cm", "left_calf_cm",
      "neck_cm", "abdomen_cm", "triceps_sf_mm", "biceps_sf_mm", "subscapular_sf_mm",
      "suprailiac_sf_mm", "abdominal_sf_mm", "thigh_sf_mm", "calf_sf_mm",
    ] as const

    for (const key of numericKeys) {
      const value = data[key]
      if (value !== undefined && !Number.isNaN(value)) {
        cleaned[key] = value
      }
    }

    if (data.measured_at) {
      cleaned.measured_at = data.measured_at
    }
    if (data.device_model) {
      cleaned.device_model = data.device_model
    }
    if (data.notes) {
      cleaned.notes = data.notes
    }

    try {
      await createMeasurement.mutateAsync(cleaned as Parameters<typeof createMeasurement.mutateAsync>[0])
      setShowForm(false)
    } catch {
      // Error handled by mutation
    }
  }

  if (isLoading) return <LoadingState lines={6} />
  if (isError) return <ErrorState message="Erro ao carregar medicoes." onRetry={refetch} />

  const measurementList = measurements ?? []

  return (
    <div className="space-y-4">
      {showForm ? (
        <MeasurementForm onSubmit={handleSubmit} isSubmitting={createMeasurement.isPending} onCancel={() => setShowForm(false)} />
      ) : (
        <div className="flex justify-end">
          <Button onClick={() => setShowForm(true)}>
            <Plus className="h-4 w-4 mr-2" />
            Nova Medicao
          </Button>
        </div>
      )}

      {createMeasurement.isError && (
        <p className="text-sm text-destructive">Erro ao salvar medicao. Tente novamente.</p>
      )}

      {measurementList.length === 0 && !showForm ? (
        <EmptyState
          icon={<Scale className="h-12 w-12" />}
          title="Nenhuma medicao"
          description="Registre a primeira medicao corporal deste paciente."
          action={
            <Button onClick={() => setShowForm(true)}>
              <Plus className="h-4 w-4 mr-2" />
              Nova Medicao
            </Button>
          }
        />
      ) : (
        <MeasurementList measurements={measurementList} />
      )}
    </div>
  )
}
