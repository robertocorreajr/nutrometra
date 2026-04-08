"use client"

import { Button, Card, CardContent } from "@nutrometra/ui"
import type { ProfessionalAddress } from "@nutrometra/api-client"
import { Pencil } from "lucide-react"

interface AddressListProps {
  addresses: ProfessionalAddress[]
  onEdit: (address: ProfessionalAddress) => void
}

function StatusBadge({ active }: { active: boolean }) {
  return (
    <span
      className={
        active
          ? "inline-flex items-center rounded-full px-2 py-0.5 text-xs font-medium bg-green-100 text-green-700"
          : "inline-flex items-center rounded-full px-2 py-0.5 text-xs font-medium bg-gray-100 text-gray-600"
      }
    >
      {active ? "Ativo" : "Inativo"}
    </span>
  )
}

function formatAddress(addr: ProfessionalAddress): string {
  const parts: string[] = []

  let streetLine = addr.street
  if (addr.number) streetLine += `, ${addr.number}`
  if (addr.complement) streetLine += ` - ${addr.complement}`
  parts.push(streetLine)

  if (addr.neighborhood) parts.push(addr.neighborhood)

  parts.push(`${addr.city}/${addr.state}`)

  if (addr.zip_code) parts.push(`CEP: ${addr.zip_code}`)

  return parts.join(" - ")
}

export function AddressList({ addresses, onEdit }: AddressListProps) {
  return (
    <div className="space-y-3">
      {addresses.map((address) => (
        <Card key={address.id}>
          <CardContent className="pt-4 pb-4">
            <div className="flex flex-col gap-2 sm:flex-row sm:items-start sm:justify-between">
              <div className="min-w-0 flex-1">
                <div className="flex items-center gap-2 mb-1">
                  <span className="font-semibold text-sm">{address.label}</span>
                  <StatusBadge active={address.active} />
                </div>
                <p className="text-sm text-muted-foreground break-words">
                  {formatAddress(address)}
                </p>
                {address.phone && (
                  <p className="text-sm text-muted-foreground mt-1">
                    Tel: {address.phone}
                  </p>
                )}
              </div>
              <Button
                variant="ghost"
                size="sm"
                onClick={() => onEdit(address)}
                className="self-start shrink-0"
              >
                <Pencil className="h-4 w-4 mr-1" />
                Editar
              </Button>
            </div>
          </CardContent>
        </Card>
      ))}
    </div>
  )
}
