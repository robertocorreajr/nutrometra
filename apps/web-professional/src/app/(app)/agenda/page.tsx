"use client"

import { useState, useMemo } from "react"
import { startOfWeek, endOfWeek, addWeeks, format } from "date-fns"
import { Plus, Calendar, Clock, Ban } from "lucide-react"
import {
  PageHeader,
  LoadingState,
  ErrorState,
  Button,
} from "@nutrometra/ui"
import {
  useProfessionalMe,
  useAppointments,
  useCreateAppointment,
} from "@nutrometra/api-client/hooks"
import type { Appointment } from "@nutrometra/api-client"
import { WeekCalendar } from "@/components/scheduling/week-calendar"
import { AppointmentForm } from "@/components/scheduling/appointment-form"
import { AppointmentDetailModal } from "@/components/scheduling/appointment-detail-modal"
import { AvailabilityList } from "@/components/scheduling/availability-list"
import { AvailabilityForm } from "@/components/scheduling/availability-form"
import { BlockList } from "@/components/scheduling/block-list"
import { BlockForm } from "@/components/scheduling/block-form"
import type { AppointmentFormValues } from "@/lib/schemas/appointment"

type Tab = "calendario" | "disponibilidade" | "bloqueios"

const tabs: { key: Tab; label: string; icon: React.ReactNode }[] = [
  { key: "calendario", label: "Calendário", icon: <Calendar className="h-4 w-4" /> },
  { key: "disponibilidade", label: "Disponibilidade", icon: <Clock className="h-4 w-4" /> },
  { key: "bloqueios", label: "Bloqueios", icon: <Ban className="h-4 w-4" /> },
]

export default function AgendaPage() {
  const { data: professional, isLoading: profLoading } = useProfessionalMe()

  const [activeTab, setActiveTab] = useState<Tab>("calendario")
  const [weekStart, setWeekStart] = useState<Date>(
    startOfWeek(new Date(), { weekStartsOn: 0 })
  )
  const [showForm, setShowForm] = useState(false)
  const [selectedAppointment, setSelectedAppointment] = useState<Appointment | null>(null)
  const [defaultSlotDate, setDefaultSlotDate] = useState<string | undefined>()

  const [showAvailabilityForm, setShowAvailabilityForm] = useState(false)
  const [showBlockForm, setShowBlockForm] = useState(false)

  const weekEnd = useMemo(() => endOfWeek(weekStart, { weekStartsOn: 0 }), [weekStart])

  const {
    data: appointments,
    isLoading: aptsLoading,
    isError: aptsError,
    refetch: aptsRefetch,
  } = useAppointments({
    professional_id: professional?.id,
    from: weekStart.toISOString(),
    to: weekEnd.toISOString(),
  })

  const createAppointment = useCreateAppointment()

  if (profLoading) return <LoadingState lines={6} />

  const professionalId = professional?.id ?? ""

  function handleNavigateWeek(dir: -1 | 1) {
    setWeekStart((prev) => addWeeks(prev, dir))
  }

  function handleSelectSlot(startTime: string) {
    setDefaultSlotDate(startTime)
    setShowForm(true)
  }

  function handleSelectAppointment(apt: Appointment) {
    setSelectedAppointment(apt)
  }

  async function handleCreateAppointment(data: AppointmentFormValues) {
    try {
      await createAppointment.mutateAsync({
        professional_id: professionalId,
        patient_id: data.patient_id,
        start_at: data.start_at,
        end_at: data.end_at,
        service_mode: data.service_mode,
        address_id: data.address_id,
        notes: data.notes,
      })
      setShowForm(false)
      setDefaultSlotDate(undefined)
    } catch {
      // Error handled by mutation
    }
  }

  return (
    <div>
      <PageHeader
        title="Agenda"
        description="Gerencie sua agenda e consultas"
        actions={
          activeTab === "calendario" ? (
            <Button onClick={() => { setDefaultSlotDate(undefined); setShowForm(true) }}>
              <Plus className="h-4 w-4 mr-2" />
              Nova Consulta
            </Button>
          ) : activeTab === "disponibilidade" ? (
            <Button onClick={() => setShowAvailabilityForm(true)}>
              <Plus className="h-4 w-4 mr-2" />
              Nova Disponibilidade
            </Button>
          ) : (
            <Button onClick={() => setShowBlockForm(true)}>
              <Plus className="h-4 w-4 mr-2" />
              Novo Bloqueio
            </Button>
          )
        }
      />

      {/* Tab buttons */}
      <div className="flex gap-1 mb-6 border-b">
        {tabs.map((tab) => (
          <button
            key={tab.key}
            type="button"
            onClick={() => setActiveTab(tab.key)}
            className={`flex items-center gap-1.5 px-4 py-2 text-sm font-medium border-b-2 transition-colors ${
              activeTab === tab.key
                ? "border-primary text-primary"
                : "border-transparent text-muted-foreground hover:text-foreground hover:border-muted-foreground/30"
            }`}
          >
            {tab.icon}
            <span className="hidden sm:inline">{tab.label}</span>
          </button>
        ))}
      </div>

      {/* Calendário tab */}
      {activeTab === "calendario" && (
        <div>
          {/* Create form */}
          {showForm && (
            <div className="mb-6">
              <AppointmentForm
                onSubmit={handleCreateAppointment}
                isSubmitting={createAppointment.isPending}
                onCancel={() => { setShowForm(false); setDefaultSlotDate(undefined) }}
                defaultDate={defaultSlotDate}
              />
              {createAppointment.isError && (
                <p className="text-sm text-destructive mt-2">
                  Erro ao criar consulta. Tente novamente.
                </p>
              )}
            </div>
          )}

          {/* Calendar */}
          {aptsLoading ? (
            <LoadingState lines={8} />
          ) : aptsError ? (
            <ErrorState
              message="Erro ao carregar consultas."
              onRetry={aptsRefetch}
            />
          ) : (
            <WeekCalendar
              weekStart={weekStart}
              appointments={appointments ?? []}
              onSelectAppointment={handleSelectAppointment}
              onSelectSlot={handleSelectSlot}
              onNavigateWeek={handleNavigateWeek}
            />
          )}

          {/* Detail modal */}
          <AppointmentDetailModal
            appointment={selectedAppointment}
            onClose={() => setSelectedAppointment(null)}
          />
        </div>
      )}

      {/* Disponibilidade tab */}
      {activeTab === "disponibilidade" && (
        <div className="space-y-6">
          {showAvailabilityForm && (
            <AvailabilityForm
              professionalId={professionalId}
              onSuccess={() => setShowAvailabilityForm(false)}
              onCancel={() => setShowAvailabilityForm(false)}
            />
          )}
          <AvailabilityList professionalId={professionalId} />
        </div>
      )}

      {/* Bloqueios tab */}
      {activeTab === "bloqueios" && (
        <div className="space-y-6">
          {showBlockForm && (
            <BlockForm
              professionalId={professionalId}
              onSuccess={() => setShowBlockForm(false)}
              onCancel={() => setShowBlockForm(false)}
            />
          )}
          <BlockList professionalId={professionalId} />
        </div>
      )}
    </div>
  )
}
