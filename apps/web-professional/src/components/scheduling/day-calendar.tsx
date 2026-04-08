"use client"

import { format, isSameDay } from "date-fns"
import { ptBR } from "date-fns/locale"
import type { Appointment } from "@nutrometra/api-client"
import { AppointmentCard } from "./appointment-card"

interface DayCalendarProps {
  date: Date
  appointments: Appointment[]
  onSelectAppointment: (apt: Appointment) => void
  onSelectSlot: (startTime: string) => void
}

const HOUR_START = 6
const HOUR_END = 22
const SLOT_MINUTES = 30

function generateTimeSlots(): string[] {
  const slots: string[] = []
  for (let h = HOUR_START; h < HOUR_END; h++) {
    slots.push(`${String(h).padStart(2, "0")}:00`)
    slots.push(`${String(h).padStart(2, "0")}:30`)
  }
  return slots
}

function getAppointmentsForSlot(
  appointments: Appointment[],
  date: Date,
  slotTime: string
): Appointment[] {
  return appointments.filter((apt) => {
    const start = new Date(apt.start_at)
    if (!isSameDay(start, date)) return false
    const aptTime = format(start, "HH:mm")
    return aptTime === slotTime
  })
}

export function DayCalendar({
  date,
  appointments,
  onSelectAppointment,
  onSelectSlot,
}: DayCalendarProps) {
  const timeSlots = generateTimeSlots()

  function handleSlotClick(slotTime: string) {
    const [hours, minutes] = slotTime.split(":").map(Number)
    const slotDate = new Date(date)
    slotDate.setHours(hours, minutes, 0, 0)
    onSelectSlot(slotDate.toISOString())
  }

  return (
    <div>
      <h3 className="text-sm font-semibold mb-3 text-center">
        {format(date, "EEEE, dd 'de' MMMM", { locale: ptBR })}
      </h3>
      <div className="border rounded-lg overflow-hidden">
        {timeSlots.map((slot) => {
          const slotAppointments = getAppointmentsForSlot(appointments, date, slot)
          return (
            <div
              key={slot}
              className="flex border-b last:border-b-0 min-h-[40px] hover:bg-accent/30 transition-colors"
            >
              <div className="w-14 shrink-0 flex items-start justify-end pr-2 pt-1 text-xs text-muted-foreground border-r bg-muted/30">
                {slot}
              </div>
              <button
                type="button"
                className="flex-1 p-1 text-left cursor-pointer"
                onClick={() => {
                  if (slotAppointments.length === 0) {
                    handleSlotClick(slot)
                  }
                }}
              >
                <div className="space-y-1">
                  {slotAppointments.map((apt) => (
                    <AppointmentCard
                      key={apt.id}
                      appointment={apt}
                      onSelect={onSelectAppointment}
                    />
                  ))}
                </div>
              </button>
            </div>
          )
        })}
      </div>
    </div>
  )
}
