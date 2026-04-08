"use client"

import { useState } from "react"
import {
  useProfessionalMe,
  useAddresses,
  useCreateAddress,
  useUpdateAddress,
} from "@nutrometra/api-client/hooks"
import type { ProfessionalAddress } from "@nutrometra/api-client"
import { Button, LoadingState, ErrorState, EmptyState } from "@nutrometra/ui"
import { MapPin, Plus } from "lucide-react"
import { AddressList } from "@/components/settings/address-list"
import { AddressForm } from "@/components/settings/address-form"
import type { AddressFormValues } from "@/lib/schemas/professional"

type ViewMode = "list" | "create" | "edit"

export default function EnderecosPage() {
  const { data: professional, isLoading: profLoading } = useProfessionalMe()
  const professionalId = professional?.id ?? ""
  const {
    data: addresses,
    isLoading,
    isError,
    refetch,
  } = useAddresses(professionalId)
  const createMutation = useCreateAddress(professionalId)
  const updateMutation = useUpdateAddress(professionalId)

  const [viewMode, setViewMode] = useState<ViewMode>("list")
  const [editingAddress, setEditingAddress] = useState<ProfessionalAddress | null>(null)

  if (profLoading || isLoading) return <LoadingState lines={4} />
  if (isError || !professional) {
    return <ErrorState message="Erro ao carregar endereços." onRetry={refetch} />
  }

  function handleCreate(data: AddressFormValues) {
    createMutation.mutate(data, {
      onSuccess: () => {
        setViewMode("list")
      },
    })
  }

  function handleEdit(address: ProfessionalAddress) {
    setEditingAddress(address)
    setViewMode("edit")
  }

  function handleUpdate(data: AddressFormValues) {
    if (!editingAddress) return
    updateMutation.mutate(
      { addressId: editingAddress.id, ...data },
      {
        onSuccess: () => {
          setEditingAddress(null)
          setViewMode("list")
        },
      },
    )
  }

  function handleCancel() {
    setEditingAddress(null)
    setViewMode("list")
  }

  if (viewMode === "create") {
    return (
      <AddressForm
        onSubmit={handleCreate}
        isSubmitting={createMutation.isPending}
        onCancel={handleCancel}
      />
    )
  }

  if (viewMode === "edit" && editingAddress) {
    return (
      <AddressForm
        onSubmit={handleUpdate}
        isSubmitting={updateMutation.isPending}
        onCancel={handleCancel}
        defaultValues={{
          label: editingAddress.label,
          street: editingAddress.street,
          number: editingAddress.number ?? "",
          complement: editingAddress.complement ?? "",
          neighborhood: editingAddress.neighborhood ?? "",
          city: editingAddress.city,
          state: editingAddress.state,
          zip_code: editingAddress.zip_code,
          country: editingAddress.country,
          phone: editingAddress.phone ?? "",
          notes: editingAddress.notes ?? "",
        }}
      />
    )
  }

  const addressList = addresses ?? []

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <p className="text-sm text-muted-foreground">
          Gerencie os endereços onde você atende.
        </p>
        <Button size="sm" onClick={() => setViewMode("create")}>
          <Plus className="h-4 w-4 mr-1" />
          Novo Endereço
        </Button>
      </div>

      {addressList.length === 0 ? (
        <EmptyState
          icon={<MapPin className="h-12 w-12" />}
          title="Nenhum endereço cadastrado"
          description="Adicione seus endereços de atendimento para começar."
          action={
            <Button size="sm" onClick={() => setViewMode("create")}>
              <Plus className="h-4 w-4 mr-1" />
              Adicionar Endereço
            </Button>
          }
        />
      ) : (
        <AddressList addresses={addressList} onEdit={handleEdit} />
      )}
    </div>
  )
}
