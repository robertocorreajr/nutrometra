"use client"

import { useState, useRef, useEffect } from "react"
import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { FormWizard, type WizardStep, Button } from "@nutrometra/ui"
import {
  useCreateAnamnesis,
  useUpdateAnamnesis,
  useFinalizeAnamnesis,
} from "@nutrometra/api-client/hooks"
import type { Anamnesis, AnamnesisRequest } from "@nutrometra/api-client"
import { useAuth } from "@nutrometra/auth"
import { anamnesisSchema, type AnamnesisFormValues } from "@/lib/schemas/anamnesis"
import { StepClinicalData } from "./steps/step-clinical-data"
import { StepHistory } from "./steps/step-history"
import { StepHabits } from "./steps/step-habits"
import { StepAllergies } from "./steps/step-allergies"
import { StepReview } from "./steps/step-review"

const wizardSteps: WizardStep[] = [
  { title: "Dados Clinicos", description: "Queixa principal e motivo da consulta" },
  { title: "Historico", description: "Historico medico pessoal e familiar" },
  { title: "Habitos", description: "Habitos alimentares, sono, atividade fisica e intestino" },
  { title: "Alergias", description: "Suplementos, alergias e intolerancias" },
  { title: "Revisao", description: "Revise e finalize a anamnese" },
]

interface AnamnesisWizardProps {
  patientId: string
  existing?: Anamnesis | null
  onComplete?: () => void
}

export function AnamnesisWizard({ patientId, existing, onComplete }: AnamnesisWizardProps) {
  const { user } = useAuth()
  const [currentStep, setCurrentStep] = useState(0)
  const [anamnesisId, setAnamnesisId] = useState<string | null>(existing?.id ?? null)

  const createAnamnesis = useCreateAnamnesis(patientId)
  const updateAnamnesis = useUpdateAnamnesis(anamnesisId ?? "", patientId)
  const finalizeAnamnesis = useFinalizeAnamnesis(anamnesisId ?? "", patientId)

  const isFinalized = existing?.status === "finalized"

  const { register, getValues, watch } = useForm<AnamnesisFormValues>({
    resolver: zodResolver(anamnesisSchema),
    defaultValues: existing
      ? {
          chief_complaint: existing.chief_complaint ?? "",
          history_present_illness: existing.history_present_illness ?? "",
          past_medical_history: existing.past_medical_history ?? "",
          family_history: existing.family_history ?? "",
          social_history: existing.social_history ?? "",
          dietary_history: existing.dietary_history ?? "",
          physical_activity: existing.physical_activity ?? "",
          sleep_pattern: existing.sleep_pattern ?? "",
          bowel_habits: existing.bowel_habits ?? "",
          water_intake: existing.water_intake ?? "",
          supplements: existing.supplements ?? "",
          observations: existing.observations ?? "",
        }
      : undefined,
  })

  const autosaveTimer = useRef<ReturnType<typeof setTimeout> | null>(null)

  useEffect(() => {
    if (!anamnesisId || isFinalized) return

    const subscription = watch(() => {
      if (autosaveTimer.current) clearTimeout(autosaveTimer.current)
      autosaveTimer.current = setTimeout(() => {
        const values = getValues()
        const payload = buildPayload(values)
        updateAnamnesis.mutate(payload)
      }, 2000)
    })

    return () => {
      subscription.unsubscribe()
      if (autosaveTimer.current) clearTimeout(autosaveTimer.current)
    }
  }, [anamnesisId, isFinalized, watch, getValues, updateAnamnesis])

  function buildPayload(values: AnamnesisFormValues): AnamnesisRequest {
    return {
      patient_id: patientId,
      professional_id: user?.id ?? "",
      chief_complaint: values.chief_complaint ?? "",
      history_present_illness: values.history_present_illness ?? "",
      past_medical_history: values.past_medical_history ?? "",
      family_history: values.family_history ?? "",
      social_history: values.social_history ?? "",
      dietary_history: values.dietary_history ?? "",
      physical_activity: values.physical_activity ?? "",
      sleep_pattern: values.sleep_pattern ?? "",
      bowel_habits: values.bowel_habits ?? "",
      water_intake: values.water_intake ?? "",
      supplements: values.supplements ?? "",
      observations: values.observations ?? "",
    }
  }

  async function handleFinalSubmit() {
    const values = getValues()
    const payload = buildPayload(values)

    try {
      if (!anamnesisId) {
        const created = await createAnamnesis.mutateAsync(payload)
        setAnamnesisId(created.id)
        await finalizeAnamnesis.mutateAsync()
      } else {
        await updateAnamnesis.mutateAsync(payload)
        await finalizeAnamnesis.mutateAsync()
      }
      onComplete?.()
    } catch {
      // Error handled by mutation
    }
  }

  async function handleSaveDraft() {
    const values = getValues()
    const payload = buildPayload(values)

    try {
      if (!anamnesisId) {
        const created = await createAnamnesis.mutateAsync(payload)
        setAnamnesisId(created.id)
      } else {
        await updateAnamnesis.mutateAsync(payload)
      }
    } catch {
      // Error handled by mutation
    }
  }

  const isPending = createAnamnesis.isPending || updateAnamnesis.isPending || finalizeAnamnesis.isPending

  return (
    <div>
      {anamnesisId && !isFinalized && (
        <div className="flex items-center gap-2 mb-4 text-xs text-muted-foreground">
          <span className="h-2 w-2 rounded-full bg-green-400" />
          {updateAnamnesis.isPending ? "Salvando..." : "Rascunho salvo automaticamente"}
        </div>
      )}

      <FormWizard
        steps={wizardSteps}
        currentStep={currentStep}
        onStepChange={setCurrentStep}
        onSubmit={handleFinalSubmit}
        isSubmitting={isPending}
        submitLabel="Finalizar Anamnese"
      >
        {currentStep === 0 && <StepClinicalData register={register} disabled={isFinalized} />}
        {currentStep === 1 && <StepHistory register={register} disabled={isFinalized} />}
        {currentStep === 2 && <StepHabits register={register} disabled={isFinalized} />}
        {currentStep === 3 && <StepAllergies register={register} disabled={isFinalized} />}
        {currentStep === 4 && <StepReview getValues={getValues} register={register} disabled={isFinalized} />}
      </FormWizard>

      {!isFinalized && (
        <div className="mt-4 flex justify-end">
          <Button variant="outline" size="sm" onClick={handleSaveDraft} disabled={isPending}>
            {isPending ? "Salvando..." : "Salvar Rascunho"}
          </Button>
        </div>
      )}
    </div>
  )
}
