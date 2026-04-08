"use client"

import { X } from "lucide-react"
import { Button } from "@nutrometra/ui"
import { useRoles, useAssignRole, useRevokeRole } from "@nutrometra/api-client/hooks"
import { roleLabels, roleColors } from "@/lib/schemas/team"

interface RoleManagerProps {
  memberId: string
  memberName: string
  open: boolean
  onClose: () => void
}

export function RoleManager({ memberId, memberName, open, onClose }: RoleManagerProps) {
  const { data: roles, isLoading } = useRoles()
  const assignRole = useAssignRole(memberId)
  const revokeRole = useRevokeRole(memberId)

  if (!open) return null

  const tenantRoles = (roles ?? []).filter((r) => r.application_scope === "tenant")

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div className="absolute inset-0 bg-black/50" onClick={onClose} />
      <div className="relative z-10 w-full max-w-md rounded-lg border bg-background shadow-lg p-6">
        <div className="flex items-start justify-between mb-4">
          <div>
            <h3 className="text-lg font-semibold">Gerenciar Roles</h3>
            <p className="text-sm text-muted-foreground">{memberName}</p>
          </div>
          <button
            type="button"
            onClick={onClose}
            className="text-muted-foreground hover:text-foreground"
            aria-label="Fechar"
          >
            <X className="h-5 w-5" />
          </button>
        </div>

        {isLoading ? (
          <div className="space-y-2">
            {[1, 2, 3].map((i) => (
              <div key={i} className="h-10 bg-muted rounded animate-pulse" />
            ))}
          </div>
        ) : (
          <div className="space-y-2">
            {tenantRoles.map((role) => {
              const label = roleLabels[role.code] ?? role.name
              const color = roleColors[role.code] ?? "bg-gray-100 text-gray-800"

              return (
                <div
                  key={role.id}
                  className="flex items-center justify-between rounded-lg border p-3"
                >
                  <div>
                    <span className={`inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium ${color}`}>
                      {label}
                    </span>
                    {role.description && (
                      <p className="text-xs text-muted-foreground mt-1">{role.description}</p>
                    )}
                  </div>
                  <div className="flex gap-2">
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => assignRole.mutate({ role_code: role.code })}
                      disabled={assignRole.isPending}
                    >
                      Atribuir
                    </Button>
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => revokeRole.mutate(role.id)}
                      disabled={revokeRole.isPending}
                      className="text-red-600 hover:text-red-700 hover:bg-red-50"
                    >
                      Remover
                    </Button>
                  </div>
                </div>
              )
            })}
          </div>
        )}
      </div>
    </div>
  )
}
