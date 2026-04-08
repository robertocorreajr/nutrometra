"use client"

import { useState } from "react"
import { UserCircle } from "lucide-react"
import { Button } from "@nutrometra/ui"
import { useProfessionals } from "@nutrometra/api-client/hooks"
import type { Professional } from "@nutrometra/api-client"
import { RoleManager } from "./role-manager"

export function MemberList() {
  const { data: professionals, isLoading, isError, refetch } = useProfessionals()
  const [selectedMember, setSelectedMember] = useState<Professional | null>(null)

  if (isLoading) {
    return (
      <div className="space-y-3">
        {[1, 2, 3].map((i) => (
          <div key={i} className="rounded-lg border bg-card p-4 animate-pulse">
            <div className="h-5 w-48 bg-muted rounded" />
          </div>
        ))}
      </div>
    )
  }

  if (isError) {
    return (
      <div className="rounded-lg border bg-card p-6 text-center">
        <p className="text-sm text-muted-foreground mb-2">Não foi possível carregar a equipe.</p>
        <Button variant="outline" size="sm" onClick={() => refetch()}>Tentar novamente</Button>
      </div>
    )
  }

  const list = professionals ?? []

  if (list.length === 0) {
    return (
      <div className="rounded-lg border bg-card p-6 text-center">
        <p className="text-sm text-muted-foreground">Nenhum membro na equipe.</p>
      </div>
    )
  }

  return (
    <>
      <div className="space-y-3">
        {list.map((member) => (
          <div
            key={member.id}
            className="rounded-lg border bg-card p-4 flex items-center gap-4"
          >
            <div className="flex h-10 w-10 items-center justify-center rounded-full bg-muted">
              <UserCircle className="h-6 w-6 text-muted-foreground" />
            </div>
            <div className="flex-1 min-w-0">
              <p className="font-medium truncate">{member.full_name}</p>
              <div className="flex items-center gap-2 text-xs text-muted-foreground mt-0.5">
                {member.specialty && <span>{member.specialty}</span>}
                {member.registration_number && (
                  <span>
                    {member.registration_number}
                    {member.registration_state ? `/${member.registration_state}` : ""}
                  </span>
                )}
              </div>
            </div>
            <Button
              variant="outline"
              size="sm"
              onClick={() => setSelectedMember(member)}
            >
              Gerenciar Roles
            </Button>
          </div>
        ))}
      </div>

      {selectedMember && (
        <RoleManager
          memberId={selectedMember.id}
          memberName={selectedMember.full_name}
          open={!!selectedMember}
          onClose={() => setSelectedMember(null)}
        />
      )}
    </>
  )
}
