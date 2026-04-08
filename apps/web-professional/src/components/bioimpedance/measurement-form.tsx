"use client"

import { useState } from "react"
import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { Button, Input, Card, CardContent } from "@nutrometra/ui"
import { measurementSchema, type MeasurementFormValues } from "@/lib/schemas/bioimpedance"
import { ChevronDown, ChevronRight } from "lucide-react"

interface MeasurementFormProps {
  onSubmit: (data: MeasurementFormValues) => void
  isSubmitting: boolean
  onCancel?: () => void
}

interface FieldConfig {
  name: keyof MeasurementFormValues
  label: string
  step?: string
  placeholder?: string
}

const basicFields: FieldConfig[] = [
  { name: "weight_kg", label: "Peso (kg)", placeholder: "Ex: 72.5" },
  { name: "height_cm", label: "Altura (cm)", placeholder: "Ex: 170" },
]

const compositionFields: FieldConfig[] = [
  { name: "body_fat_pct", label: "Gordura corporal (%)", placeholder: "Ex: 22.5" },
  { name: "lean_mass_kg", label: "Massa magra (kg)", placeholder: "Ex: 55.0" },
  { name: "fat_mass_kg", label: "Massa gorda (kg)", placeholder: "Ex: 17.5" },
  { name: "muscle_mass_kg", label: "Massa muscular (kg)", placeholder: "Ex: 30.0" },
  { name: "bone_mass_kg", label: "Massa óssea (kg)", placeholder: "Ex: 3.2" },
  { name: "water_pct", label: "Água corporal (%)", placeholder: "Ex: 55.0" },
  { name: "visceral_fat", label: "Gordura visceral", placeholder: "Ex: 8" },
  { name: "basal_metabolic_rate", label: "Taxa metabólica basal (kcal)", placeholder: "Ex: 1650" },
]

const circumferenceFields: FieldConfig[] = [
  { name: "waist_cm", label: "Cintura (cm)", placeholder: "Ex: 80.0" },
  { name: "hip_cm", label: "Quadril (cm)", placeholder: "Ex: 95.0" },
  { name: "chest_cm", label: "Tórax (cm)", placeholder: "Ex: 98.0" },
  { name: "right_arm_cm", label: "Braço direito (cm)", placeholder: "Ex: 30.0" },
  { name: "left_arm_cm", label: "Braço esquerdo (cm)", placeholder: "Ex: 29.5" },
  { name: "right_thigh_cm", label: "Coxa direita (cm)", placeholder: "Ex: 55.0" },
  { name: "left_thigh_cm", label: "Coxa esquerda (cm)", placeholder: "Ex: 54.5" },
  { name: "right_calf_cm", label: "Panturrilha direita (cm)", placeholder: "Ex: 36.0" },
  { name: "left_calf_cm", label: "Panturrilha esquerda (cm)", placeholder: "Ex: 35.5" },
  { name: "neck_cm", label: "Pescoço (cm)", placeholder: "Ex: 38.0" },
  { name: "abdomen_cm", label: "Abdômen (cm)", placeholder: "Ex: 85.0" },
]

const skinfoldFields: FieldConfig[] = [
  { name: "triceps_sf_mm", label: "Tríceps (mm)", placeholder: "Ex: 12.0" },
  { name: "biceps_sf_mm", label: "Bíceps (mm)", placeholder: "Ex: 8.0" },
  { name: "subscapular_sf_mm", label: "Subescapular (mm)", placeholder: "Ex: 14.0" },
  { name: "suprailiac_sf_mm", label: "Suprailíaca (mm)", placeholder: "Ex: 18.0" },
  { name: "abdominal_sf_mm", label: "Abdominal (mm)", placeholder: "Ex: 20.0" },
  { name: "thigh_sf_mm", label: "Coxa (mm)", placeholder: "Ex: 16.0" },
  { name: "calf_sf_mm", label: "Panturrilha (mm)", placeholder: "Ex: 10.0" },
]

interface SectionProps {
  title: string
  isOpen: boolean
  onToggle: () => void
  children: React.ReactNode
}

function Section({ title, isOpen, onToggle, children }: SectionProps) {
  return (
    <div className="border rounded-md">
      <button
        type="button"
        onClick={onToggle}
        className="flex items-center justify-between w-full px-4 py-3 text-sm font-medium text-left hover:bg-muted/50 transition-colors"
      >
        {title}
        {isOpen ? (
          <ChevronDown className="h-4 w-4 text-muted-foreground" />
        ) : (
          <ChevronRight className="h-4 w-4 text-muted-foreground" />
        )}
      </button>
      {isOpen && <div className="px-4 pb-4">{children}</div>}
    </div>
  )
}

export function MeasurementForm({ onSubmit, isSubmitting, onCancel }: MeasurementFormProps) {
  const [openSections, setOpenSections] = useState<Record<string, boolean>>({
    basic: true,
    composition: false,
    circumferences: false,
    skinfolds: false,
    metadata: false,
  })

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<MeasurementFormValues>({
    resolver: zodResolver(measurementSchema),
    defaultValues: {
      source: "manual",
      device_model: "",
      notes: "",
    },
  })

  function toggleSection(key: string) {
    setOpenSections((prev) => ({ ...prev, [key]: !prev[key] }))
  }

  function handleFormSubmit(data: MeasurementFormValues) {
    onSubmit(data)
    reset()
  }

  function renderNumberFields(fields: FieldConfig[]) {
    return (
      <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
        {fields.map((field) => (
          <div key={field.name} className="space-y-1">
            <label htmlFor={field.name} className="text-sm font-medium">
              {field.label}
            </label>
            <Input
              id={field.name}
              type="number"
              step="0.1"
              placeholder={field.placeholder}
              {...register(field.name)}
            />
            {errors[field.name] && (
              <p className="text-sm text-destructive">
                {errors[field.name]?.message}
              </p>
            )}
          </div>
        ))}
      </div>
    )
  }

  return (
    <Card>
      <CardContent className="pt-6">
        <form onSubmit={handleSubmit(handleFormSubmit)} className="space-y-4">
          {/* Dados Basicos - always open */}
          <Section
            title="Dados Basicos"
            isOpen={openSections.basic ?? true}
            onToggle={() => toggleSection("basic")}
          >
            {renderNumberFields(basicFields)}
          </Section>

          {/* Composicao Corporal */}
          <Section
            title="Composicao Corporal"
            isOpen={openSections.composition ?? false}
            onToggle={() => toggleSection("composition")}
          >
            {renderNumberFields(compositionFields)}
          </Section>

          {/* Circunferencias */}
          <Section
            title="Circunferencias"
            isOpen={openSections.circumferences ?? false}
            onToggle={() => toggleSection("circumferences")}
          >
            {renderNumberFields(circumferenceFields)}
          </Section>

          {/* Dobras Cutaneas */}
          <Section
            title="Dobras Cutaneas"
            isOpen={openSections.skinfolds ?? false}
            onToggle={() => toggleSection("skinfolds")}
          >
            {renderNumberFields(skinfoldFields)}
          </Section>

          {/* Metadados */}
          <Section
            title="Metadados"
            isOpen={openSections.metadata ?? false}
            onToggle={() => toggleSection("metadata")}
          >
            <div className="space-y-4">
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <div className="space-y-1">
                  <label htmlFor="source" className="text-sm font-medium">
                    Origem
                  </label>
                  <select
                    id="source"
                    {...register("source")}
                    className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
                  >
                    <option value="manual">Manual</option>
                    <option value="device">Dispositivo</option>
                    <option value="import">Importacao</option>
                  </select>
                </div>

                <div className="space-y-1">
                  <label htmlFor="device_model" className="text-sm font-medium">
                    Modelo do dispositivo
                  </label>
                  <Input
                    id="device_model"
                    {...register("device_model")}
                    placeholder="Ex: InBody 270, Tanita BC-545N"
                  />
                </div>
              </div>

              <div className="space-y-1">
                <label htmlFor="measured_at" className="text-sm font-medium">
                  Data da medicao
                </label>
                <Input
                  id="measured_at"
                  type="date"
                  {...register("measured_at")}
                />
              </div>

              <div className="space-y-1">
                <label htmlFor="notes" className="text-sm font-medium">
                  Observacoes
                </label>
                <textarea
                  id="notes"
                  {...register("notes")}
                  rows={3}
                  placeholder="Observacoes sobre a medicao..."
                  className="flex w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
                />
              </div>
            </div>
          </Section>

          <div className="flex gap-2 justify-end">
            {onCancel && (
              <Button type="button" variant="ghost" onClick={onCancel}>
                Cancelar
              </Button>
            )}
            <Button type="submit" disabled={isSubmitting}>
              {isSubmitting ? "Salvando..." : "Salvar Medicao"}
            </Button>
          </div>
        </form>
      </CardContent>
    </Card>
  )
}
