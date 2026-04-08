"use client"

import { useState } from "react"
import { Button, Input, Card, CardContent } from "@nutrometra/ui"
import { useSetServiceMode } from "@nutrometra/api-client/hooks"
import type { ProfessionalAddress, ProfessionalServiceMode } from "@nutrometra/api-client"
import { serviceModeLabels } from "@/lib/schemas/professional"
import type { ServiceModeType } from "@nutrometra/api-client"
import { Check } from "lucide-react"

interface ServiceModesFormProps {
  professionalId: string
  serviceModes: ProfessionalServiceMode[]
  addresses: ProfessionalAddress[]
}

interface ModeConfig {
  mode: ServiceModeType
  label: string
  description: string
  requiresAddress: boolean
}

const modeConfigs: ModeConfig[] = [
  {
    mode: "onsite",
    label: serviceModeLabels.onsite,
    description: "Atendimento no consultório ou clínica",
    requiresAddress: true,
  },
  {
    mode: "online",
    label: serviceModeLabels.online,
    description: "Atendimento por videochamada",
    requiresAddress: false,
  },
  {
    mode: "home_visit",
    label: serviceModeLabels.home_visit,
    description: "Atendimento na residência do paciente",
    requiresAddress: false,
  },
]

export function ServiceModesForm({
  professionalId,
  serviceModes,
  addresses,
}: ServiceModesFormProps) {
  const setModeMutation = useSetServiceMode(professionalId)
  const [savedMode, setSavedMode] = useState<string | null>(null)

  function getExistingMode(mode: ServiceModeType): ProfessionalServiceMode | undefined {
    return serviceModes.find((sm) => sm.mode === mode)
  }

  function handleSave(
    mode: ServiceModeType,
    durationMin: number,
    addressId?: string,
  ) {
    setSavedMode(null)
    setModeMutation.mutate(
      {
        mode,
        duration_min: durationMin,
        address_id: addressId || undefined,
      },
      {
        onSuccess: () => {
          setSavedMode(mode)
        },
      },
    )
  }

  const selectClasses =
    "flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"

  return (
    <div className="space-y-4">
      <p className="text-sm text-muted-foreground">
        Configure as modalidades de atendimento que você oferece e a duração padrão de cada uma.
      </p>

      {modeConfigs.map((config) => (
        <ServiceModeCard
          key={config.mode}
          config={config}
          existing={getExistingMode(config.mode)}
          addresses={addresses}
          onSave={handleSave}
          isSaving={setModeMutation.isPending}
          isSaved={savedMode === config.mode}
          selectClasses={selectClasses}
        />
      ))}

      {setModeMutation.isError && (
        <p className="text-sm text-destructive">
          Erro ao salvar modalidade. Tente novamente.
        </p>
      )}
    </div>
  )
}

interface ServiceModeCardProps {
  config: ModeConfig
  existing?: ProfessionalServiceMode
  addresses: ProfessionalAddress[]
  onSave: (mode: ServiceModeType, durationMin: number, addressId?: string) => void
  isSaving: boolean
  isSaved: boolean
  selectClasses: string
}

function ServiceModeCard({
  config,
  existing,
  addresses,
  onSave,
  isSaving,
  isSaved,
  selectClasses,
}: ServiceModeCardProps) {
  const [enabled, setEnabled] = useState(existing?.active ?? false)
  const [durationMin, setDurationMin] = useState(existing?.duration_min ?? 50)
  const [addressId, setAddressId] = useState(existing?.address_id ?? "")

  function handleToggle() {
    setEnabled(!enabled)
  }

  function handleSave() {
    onSave(config.mode, durationMin, addressId || undefined)
  }

  const activeAddresses = addresses.filter((a) => a.active)

  return (
    <Card>
      <CardContent className="pt-4 pb-4">
        <div className="space-y-3">
          <div className="flex items-start justify-between gap-3">
            <div className="flex items-start gap-3 min-w-0">
              <button
                type="button"
                role="switch"
                aria-checked={enabled}
                onClick={handleToggle}
                className={`
                  relative inline-flex h-6 w-11 shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors
                  ${enabled ? "bg-primary" : "bg-gray-200"}
                `}
              >
                <span
                  className={`
                    pointer-events-none inline-block h-5 w-5 rounded-full bg-white shadow-sm ring-0 transition-transform
                    ${enabled ? "translate-x-5" : "translate-x-0"}
                  `}
                />
              </button>
              <div>
                <p className="font-semibold text-sm">{config.label}</p>
                <p className="text-xs text-muted-foreground">{config.description}</p>
              </div>
            </div>

            {isSaved && (
              <span className="inline-flex items-center gap-1 text-xs text-green-600 shrink-0">
                <Check className="h-3 w-3" />
                Salvo
              </span>
            )}
          </div>

          {enabled && (
            <div className="pl-14 space-y-3">
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <div className="space-y-1">
                  <label
                    htmlFor={`duration-${config.mode}`}
                    className="text-sm font-medium"
                  >
                    Duração padrão (minutos)
                  </label>
                  <Input
                    id={`duration-${config.mode}`}
                    type="number"
                    min={1}
                    value={durationMin}
                    onChange={(e) => setDurationMin(Number(e.target.value))}
                  />
                </div>

                {config.requiresAddress && (
                  <div className="space-y-1">
                    <label
                      htmlFor={`address-${config.mode}`}
                      className="text-sm font-medium"
                    >
                      Endereço de atendimento
                    </label>
                    <select
                      id={`address-${config.mode}`}
                      value={addressId}
                      onChange={(e) => setAddressId(e.target.value)}
                      className={selectClasses}
                    >
                      <option value="">Selecione um endereço</option>
                      {activeAddresses.map((addr) => (
                        <option key={addr.id} value={addr.id}>
                          {addr.label}
                        </option>
                      ))}
                    </select>
                    {activeAddresses.length === 0 && (
                      <p className="text-xs text-muted-foreground">
                        Nenhum endereço ativo. Cadastre um na aba Endereços.
                      </p>
                    )}
                  </div>
                )}
              </div>

              <div className="flex justify-end">
                <Button
                  type="button"
                  size="sm"
                  onClick={handleSave}
                  disabled={isSaving}
                >
                  {isSaving ? "Salvando..." : "Salvar"}
                </Button>
              </div>
            </div>
          )}
        </div>
      </CardContent>
    </Card>
  )
}
