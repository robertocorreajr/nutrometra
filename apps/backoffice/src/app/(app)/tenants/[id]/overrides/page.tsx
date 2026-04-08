"use client"

import { useState } from "react"
import { useParams } from "next/navigation"
import { Button, LoadingState, ErrorState, EmptyState } from "@nutrometra/ui"
import { useBackofficeOverrides } from "@nutrometra/api-client/hooks"
import type { FeatureOverride } from "@nutrometra/api-client/types"
import { OverrideList } from "@/components/tenants/override-list"
import { OverrideForm } from "@/components/tenants/override-form"
import { Plus } from "lucide-react"

export default function OverridesPage() {
  const { id } = useParams<{ id: string }>()
  const { data: overrides, isLoading, isError, refetch } = useBackofficeOverrides(id)
  const [editing, setEditing] = useState<FeatureOverride | null>(null)
  const [creating, setCreating] = useState(false)

  if (isLoading) return <LoadingState />
  if (isError) return <ErrorState message="Nao foi possivel carregar os overrides." onRetry={refetch} />

  if (creating || editing) {
    return (
      <OverrideForm
        tenantId={id}
        existing={editing ?? undefined}
        onDone={() => {
          setCreating(false)
          setEditing(null)
        }}
      />
    )
  }

  return (
    <div className="space-y-4">
      <div className="flex justify-end">
        <Button onClick={() => setCreating(true)}>
          <Plus className="h-4 w-4 mr-2" />
          Novo Override
        </Button>
      </div>
      {overrides && overrides.length === 0 && (
        <EmptyState
          title="Nenhum override"
          description="Este tenant nao possui overrides de features."
        />
      )}
      {overrides && overrides.length > 0 && (
        <OverrideList tenantId={id} overrides={overrides} onEdit={setEditing} />
      )}
    </div>
  )
}
