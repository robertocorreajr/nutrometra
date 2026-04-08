"use client"

import * as React from "react"
import { cn } from "../lib/utils"

export interface WizardStep {
  title: string
  description?: string
}

interface FormWizardProps {
  steps: WizardStep[]
  currentStep: number
  onStepChange: (step: number) => void
  onSubmit: () => void
  validateStep?: (step: number) => Promise<boolean> | boolean
  isSubmitting?: boolean
  submitLabel?: string
  children: React.ReactNode
}

export function FormWizard({
  steps,
  currentStep,
  onStepChange,
  onSubmit,
  validateStep,
  isSubmitting = false,
  submitLabel = "Finalizar",
  children,
}: FormWizardProps) {
  const isFirst = currentStep === 0
  const isLast = currentStep === steps.length - 1

  async function handleNext() {
    if (validateStep) {
      const valid = await validateStep(currentStep)
      if (!valid) return
    }
    if (isLast) {
      onSubmit()
    } else {
      onStepChange(currentStep + 1)
    }
  }

  function handlePrevious() {
    if (!isFirst) {
      onStepChange(currentStep - 1)
    }
  }

  return (
    <div className="space-y-6">
      <div className="space-y-2">
        <div className="flex justify-between text-sm text-muted-foreground">
          <span>Passo {currentStep + 1} de {steps.length}</span>
          <span>{steps[currentStep].title}</span>
        </div>
        <div className="w-full bg-muted rounded-full h-2">
          <div
            className="bg-primary h-2 rounded-full transition-all duration-300"
            style={{ width: `${((currentStep + 1) / steps.length) * 100}%` }}
          />
        </div>
        <div className="flex gap-1">
          {steps.map((step, i) => (
            <button
              key={i}
              type="button"
              onClick={() => { if (i < currentStep) onStepChange(i) }}
              disabled={i > currentStep}
              className={cn(
                "flex-1 text-xs py-1 rounded transition-colors text-center",
                i === currentStep
                  ? "bg-primary/10 text-primary font-medium"
                  : i < currentStep
                    ? "text-muted-foreground hover:bg-accent cursor-pointer"
                    : "text-muted-foreground/40 cursor-not-allowed"
              )}
              title={step.title}
            >
              <span className="hidden sm:inline">{step.title}</span>
              <span className="sm:hidden">{i + 1}</span>
            </button>
          ))}
        </div>
      </div>

      {steps[currentStep].description && (
        <p className="text-sm text-muted-foreground">{steps[currentStep].description}</p>
      )}

      <div>{children}</div>

      <div className="flex justify-between pt-4 border-t">
        <button
          type="button"
          onClick={handlePrevious}
          disabled={isFirst}
          className={cn(
            "px-4 py-2 text-sm font-medium rounded-md border transition-colors",
            isFirst ? "opacity-50 cursor-not-allowed" : "hover:bg-accent"
          )}
        >
          Anterior
        </button>
        <button
          type="button"
          onClick={handleNext}
          disabled={isSubmitting}
          className="px-4 py-2 text-sm font-medium rounded-md bg-primary text-primary-foreground hover:bg-primary/90 disabled:opacity-50 transition-colors"
        >
          {isSubmitting ? "Salvando..." : isLast ? submitLabel : "Próximo"}
        </button>
      </div>
    </div>
  )
}
