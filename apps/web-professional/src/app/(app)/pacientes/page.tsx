"use client"

import Link from "next/link"
import { Plus, Users } from "lucide-react"
import {
  PageHeader,
  LoadingState,
  ErrorState,
  EmptyState,
  Button,
  DataTable,
  Card,
  CardContent,
} from "@nutrometra/ui"
import type { Column } from "@nutrometra/ui"
import { usePatients } from "@nutrometra/api-client/hooks"
import type { Patient } from "@nutrometra/api-client"

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

const columns: Column<Patient>[] = [
  {
    key: "full_name",
    header: "Nome",
    render: (p) => (
      <Link
        href={`/pacientes/${p.id}`}
        className="font-medium text-primary hover:underline"
      >
        {p.full_name}
      </Link>
    ),
    searchValue: (p) => p.full_name,
    sortable: true,
  },
  {
    key: "email",
    header: "E-mail",
    render: (p) => p.email ?? "—",
    searchValue: (p) => p.email ?? "",
  },
  {
    key: "phone",
    header: "Telefone",
    render: (p) => p.phone ?? "—",
    className: "hidden lg:table-cell",
  },
  {
    key: "active",
    header: "Status",
    render: (p) => <StatusBadge active={p.active} />,
    searchValue: (p) => (p.active ? "ativo" : "inativo"),
  },
]

function MobilePatientCard({ patient }: { patient: Patient }) {
  return (
    <Card>
      <CardContent className="p-4">
        <div className="flex items-start justify-between gap-2">
          <div className="min-w-0">
            <Link
              href={`/pacientes/${patient.id}`}
              className="font-medium text-primary hover:underline truncate block"
            >
              {patient.full_name}
            </Link>
            {patient.email && (
              <p className="text-sm text-muted-foreground mt-0.5 truncate">
                {patient.email}
              </p>
            )}
          </div>
          <StatusBadge active={patient.active} />
        </div>
      </CardContent>
    </Card>
  )
}

export default function PacientesPage() {
  const { data: patients, isLoading, isError, refetch } = usePatients()

  if (isLoading) return <LoadingState lines={6} />
  if (isError)
    return (
      <ErrorState
        message="Não foi possível carregar a lista de pacientes."
        onRetry={refetch}
      />
    )

  const patientList = patients ?? []

  if (patientList.length === 0) {
    return (
      <div>
        <PageHeader title="Pacientes" description="Gerencie seus pacientes" />
        <EmptyState
          icon={<Users className="h-10 w-10 text-muted-foreground" />}
          title="Nenhum paciente cadastrado"
          description="Comece cadastrando seu primeiro paciente."
          action={
            <Button asChild>
              <Link href="/pacientes/novo">
                <Plus className="mr-2 h-4 w-4" />
                Novo Paciente
              </Link>
            </Button>
          }
        />
      </div>
    )
  }

  return (
    <div>
      <PageHeader title="Pacientes" description="Gerencie seus pacientes" />
      <DataTable
        data={patientList}
        columns={columns}
        keyExtractor={(p) => p.id}
        searchPlaceholder="Buscar paciente..."
        renderMobileCard={(p) => <MobilePatientCard patient={p} />}
        emptyMessage="Nenhum paciente encontrado para esta busca."
        actions={
          <Button asChild>
            <Link href="/pacientes/novo">
              <Plus className="mr-2 h-4 w-4" />
              Novo Paciente
            </Link>
          </Button>
        }
      />
    </div>
  )
}
