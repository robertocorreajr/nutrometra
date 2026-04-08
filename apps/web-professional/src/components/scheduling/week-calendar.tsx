"use client"

import { useState } from "react"
import { format, addDays, isToday, isSameDay, startOfWeek } from "date-fns"
import { ptBR } from "date-fns/locale"
import { ChevronLeft, ChevronRight } from "lucide-react"
import { Button } from "@nutrometra/ui"
import type { Appointment } from "@nutrometra/api-client"
import { AppointmentCard } from "./appointment-card"
import { DayCalendar } from "./day-calendar"

interface WeekCalendarProps {
  weekStart: Date
  appointments: Appointment[]
  onSelectAppointment: (apt: Appointment) => void
  onSelectSlot: (startTime: string) => void
  onNavigateWeek: (dir: -1 | 1) => void
}

const DAY_NAMES_SHORT = ["Dom", "Seg", "Ter", "Qua", "Qui", "Sex", "Sáb"]
const HOUR_START = 6
const HOUR_END = 22

function generateHourSlots(): string[] {
  const slots: string[] = []
  for (let h = HOUR_START; h < HOUR_END; h++) {
    slots.push(`${String(h).padStart(2, "0")}:00`)
    slots.push(`${String(h).padStart(2, "0")}:30`)
  }
  return slots
}

function getWeekDays(weekStart: Date): Date[] {
  return Array.from({ length: 7 }, (_, i) => addDays(weekStart, i))
}

function getAppointmentsForDaySlot(
  appointments: Appointment[],
  day: Date,
  slotTime: string
): Appointment[] {
  return appointments.filter((apt) => {
    const start = new Date(apt.start_at)
    if (!isSameDay(start, day)) return false
    return format(start, "HH:mm") === slotTime
  })
}

export function WeekCalendar({
  weekStart,
  appointments,
  onSelectAppointment,
  onSelectSlot,
  onNavigateWeek,
}: WeekCalendarProps) {
  const [selectedDay, setSelectedDay] = useState<Date>(new Date())
  const days = getWeekDays(weekStart)
  const timeSlots = generateHourSlots()

  function handleTodayClick() {
    const today = new Date()
    const currentWeekStart = startOfWeek(today, { weekStartsOn: 0 })
    if (currentWeekStart.getTime() !== weekStart.getTime()) {
      // Navigate to today's week by calculating direction
      const diff = currentWeekStart.getTime() - weekStart.getTime()
      onNavigateWeek(diff > 0 ? 1 : -1)
    }
    setSelectedDay(today)
  }

  function handleSlotClick(day: Date, slotTime: string) {
    const [hours, minutes] = slotTime.split(":").map(Number)
    const slotDate = new Date(day)
    slotDate.setHours(hours, minutes, 0, 0)
    onSelectSlot(slotDate.toISOString())
  }

  return (
    <div>
      {/* Navigation header */}
      <div className="flex items-center justify-between mb-4">
        <div className="flex items-center gap-2">
          <Button variant="outline" size="icon" onClick={() => onNavigateWeek(-1)}>
            <ChevronLeft className="h-4 w-4" />
          </Button>
          <Button variant="outline" size="icon" onClick={() => onNavigateWeek(1)}>
            <ChevronRight className="h-4 w-4" />
          </Button>
          <Button variant="outline" size="sm" onClick={handleTodayClick}>
            Hoje
          </Button>
        </div>
        <h3 className="text-sm font-semibold">
          {format(weekStart, "dd MMM", { locale: ptBR })} -{" "}
          {format(addDays(weekStart, 6), "dd MMM yyyy", { locale: ptBR })}
        </h3>
      </div>

      {/* Mobile: day selector + DayCalendar */}
      <div className="md:hidden">
        <div className="flex gap-1 mb-3 overflow-x-auto">
          {days.map((day, i) => (
            <button
              key={i}
              type="button"
              onClick={() => setSelectedDay(day)}
              className={`flex flex-col items-center px-2 py-1.5 rounded-md text-xs min-w-[44px] transition-colors ${
                isSameDay(day, selectedDay)
                  ? "bg-primary text-primary-foreground"
                  : isToday(day)
                    ? "bg-accent text-accent-foreground"
                    : "hover:bg-accent/50"
              }`}
            >
              <span className="font-medium">{DAY_NAMES_SHORT[i]}</span>
              <span>{format(day, "dd")}</span>
            </button>
          ))}
        </div>
        <DayCalendar
          date={selectedDay}
          appointments={appointments}
          onSelectAppointment={onSelectAppointment}
          onSelectSlot={onSelectSlot}
        />
      </div>

      {/* Desktop: full week grid */}
      <div className="hidden md:block border rounded-lg overflow-hidden">
        {/* Day headers */}
        <div className="grid grid-cols-[56px_repeat(7,1fr)] border-b bg-muted/30">
          <div className="border-r" />
          {days.map((day, i) => (
            <div
              key={i}
              className={`text-center py-2 text-xs font-medium border-r last:border-r-0 ${
                isToday(day) ? "bg-primary/10 text-primary" : ""
              }`}
            >
              <div>{DAY_NAMES_SHORT[i]}</div>
              <div className={`text-lg font-bold ${isToday(day) ? "text-primary" : ""}`}>
                {format(day, "dd")}
              </div>
            </div>
          ))}
        </div>

        {/* Time grid */}
        <div className="max-h-[600px] overflow-y-auto">
          {timeSlots.map((slot) => (
            <div
              key={slot}
              className="grid grid-cols-[56px_repeat(7,1fr)] border-b last:border-b-0 min-h-[36px]"
            >
              <div className="flex items-start justify-end pr-2 pt-0.5 text-[10px] text-muted-foreground border-r bg-muted/30">
                {slot}
              </div>
              {days.map((day, i) => {
                const slotAppointments = getAppointmentsForDaySlot(appointments, day, slot)
                return (
                  <button
                    key={i}
                    type="button"
                    className={`border-r last:border-r-0 p-0.5 text-left hover:bg-accent/30 transition-colors cursor-pointer ${
                      isToday(day) ? "bg-primary/5" : ""
                    }`}
                    onClick={() => {
                      if (slotAppointments.length === 0) {
                        handleSlotClick(day, slot)
                      }
                    }}
                  >
                    {slotAppointments.map((apt) => (
                      <AppointmentCard
                        key={apt.id}
                        appointment={apt}
                        onSelect={onSelectAppointment}
                      />
                    ))}
                  </button>
                )
              })}
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}
