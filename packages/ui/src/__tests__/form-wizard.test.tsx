import { describe, it, expect, vi, afterEach } from "vitest"
import { render, screen, cleanup } from "@testing-library/react"
import userEvent from "@testing-library/user-event"
import { FormWizard, type WizardStep } from "../components/form-wizard"

const steps: WizardStep[] = [
  { title: "Passo 1", description: "Descrição 1" },
  { title: "Passo 2", description: "Descrição 2" },
  { title: "Passo 3", description: "Descrição 3" },
]

describe("FormWizard", () => {
  afterEach(() => {
    cleanup()
  })

  it("renders current step info", () => {
    render(
      <FormWizard steps={steps} currentStep={0} onStepChange={vi.fn()} onSubmit={vi.fn()}>
        <div>Step content</div>
      </FormWizard>
    )
    expect(screen.getByText("Passo 1 de 3")).toBeTruthy()
    expect(screen.getByText("Step content")).toBeTruthy()
  })

  it("disables Previous on first step", () => {
    render(
      <FormWizard steps={steps} currentStep={0} onStepChange={vi.fn()} onSubmit={vi.fn()}>
        <div>Content</div>
      </FormWizard>
    )
    const prevBtn = screen.getByText("Anterior")
    expect(prevBtn).toBeDisabled()
  })

  it("shows submit label on last step", () => {
    render(
      <FormWizard steps={steps} currentStep={2} onStepChange={vi.fn()} onSubmit={vi.fn()} submitLabel="Salvar Anamnese">
        <div>Last step</div>
      </FormWizard>
    )
    expect(screen.getByText("Salvar Anamnese")).toBeTruthy()
  })

  it("calls onStepChange when clicking Next", async () => {
    const user = userEvent.setup()
    const onStepChange = vi.fn()
    render(
      <FormWizard steps={steps} currentStep={0} onStepChange={onStepChange} onSubmit={vi.fn()}>
        <div>Content</div>
      </FormWizard>
    )
    await user.click(screen.getByText("Próximo"))
    expect(onStepChange).toHaveBeenCalledWith(1)
  })

  it("calls onSubmit on last step", async () => {
    const user = userEvent.setup()
    const onSubmit = vi.fn()
    render(
      <FormWizard steps={steps} currentStep={2} onStepChange={vi.fn()} onSubmit={onSubmit}>
        <div>Content</div>
      </FormWizard>
    )
    await user.click(screen.getByText("Finalizar"))
    expect(onSubmit).toHaveBeenCalled()
  })
})
