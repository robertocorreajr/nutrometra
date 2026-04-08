"use client"

import { useParams } from "next/navigation"
import { Card, CardContent, CardHeader, CardTitle } from "@nutrometra/ui"
import { useBackofficeTenant } from "@nutrometra/api-client/hooks"
import { tenantStatusLabels, tenantStatusColors } from "@/lib/schemas/backoffice"
import { format, parseISO } from "date-fns"

export default function TenantResumoPage() {
  const { id } = useParams<{ id: string }>()
  const { data: tenant } = useBackofficeTenant(id)

  if (!tenant) return null

  const fields = [
    { label: "Nome", value: tenant.tenant_name },
    { label: "Slug", value: tenant.slug },
    { label: "Email", value: tenant.owner_email },
    {
      label: "Status",
      value: (
        <span
          className={`inline-flex items-center rounded-full px-2 py-1 text-xs font-medium ${
            tenantStatusColors[tenant.status] ?? "bg-gray-100 text-gray-700"
          }`}
        >
          {tenantStatusLabels[tenant.status] ?? tenant.status}
        </span>
      ),
    },
    { label: "Criado em", value: format(parseISO(tenant.created_at), "dd/MM/yyyy HH:mm") },
    { label: "Atualizado em", value: format(parseISO(tenant.updated_at), "dd/MM/yyyy HH:mm") },
  ]

  return (
    <Card>
      <CardHeader>
        <CardTitle>Dados do Tenant</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        {fields.map((f) => (
          <div key={f.label} className="flex flex-col sm:flex-row sm:items-center gap-1">
            <span className="text-sm font-medium text-muted-foreground w-40">{f.label}</span>
            <span className="text-sm">{f.value}</span>
          </div>
        ))}
      </CardContent>
    </Card>
  )
}
