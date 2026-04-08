"use client"

import { PageHeader } from "@nutrometra/ui"
import { MemberList } from "@/components/team/member-list"

export default function EquipePage() {
  return (
    <div>
      <PageHeader
        title="Equipe"
        description="Gerencie os profissionais e seus papéis no consultório."
      />
      <div className="mt-6">
        <MemberList />
      </div>
    </div>
  )
}
