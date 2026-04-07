package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// AppointmentStatus represents the lifecycle state of an appointment.
type AppointmentStatus string

const (
	StatusScheduled AppointmentStatus = "scheduled"
	StatusConfirmed AppointmentStatus = "confirmed"
	StatusCompleted AppointmentStatus = "completed"
	StatusCancelled AppointmentStatus = "cancelled"
	StatusNoShow    AppointmentStatus = "no_show"
)

// ServiceMode defines how care is delivered.
type ServiceMode string

const (
	ModeOnsite    ServiceMode = "onsite"
	ModeOnline    ServiceMode = "online"
	ModeHomeVisit ServiceMode = "home_visit"
)

// AppointmentSource defines who created the appointment.
type AppointmentSource string

const (
	SourceProfessional AppointmentSource = "professional"
	SourcePatient      AppointmentSource = "patient"
	SourceSystem       AppointmentSource = "system"
)

// Sentinel errors.
var (
	ErrNotFound           = errors.New("scheduling: not found")
	ErrConflict           = errors.New("scheduling: time conflict with existing appointment")
	ErrInvalidTransition  = errors.New("scheduling: invalid status transition")
	ErrSlotUnavailable    = errors.New("scheduling: requested slot is not available")
	ErrEntitlementExceeded = errors.New("scheduling: entitlement limit exceeded")
)

// ValidTransitions defines allowed appointment status transitions.
var ValidTransitions = map[AppointmentStatus][]AppointmentStatus{
	StatusScheduled: {StatusConfirmed, StatusCancelled, StatusNoShow},
	StatusConfirmed: {StatusCompleted, StatusCancelled, StatusNoShow},
}

// IsValidTransition checks if a status transition is allowed.
func IsValidTransition(from, to AppointmentStatus) bool {
	allowed, ok := ValidTransitions[from]
	if !ok {
		return false
	}
	for _, s := range allowed {
		if s == to {
			return true
		}
	}
	return false
}

// AvailabilityRule defines a recurring time slot when a professional is available.
type AvailabilityRule struct {
	ID             uuid.UUID
	TenantID       uuid.UUID
	ProfessionalID uuid.UUID
	DayOfWeek      int       // 0=Sunday
	StartTime      time.Time // only time part used
	EndTime        time.Time // only time part used
	ServiceMode    ServiceMode
	AddressID      *uuid.UUID
	Active         bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// ScheduleBlock represents a time period when a professional is unavailable.
type ScheduleBlock struct {
	ID             uuid.UUID
	TenantID       uuid.UUID
	ProfessionalID uuid.UUID
	StartAt        time.Time
	EndAt          time.Time
	Reason         string
	AllDay         bool
	CreatedAt      time.Time
}

// Appointment represents a scheduled consultation.
type Appointment struct {
	ID                 uuid.UUID
	TenantID           uuid.UUID
	ProfessionalID     uuid.UUID
	PatientID          *uuid.UUID
	StartAt            time.Time
	EndAt              time.Time
	ServiceMode        ServiceMode
	AddressID          *uuid.UUID
	Status             AppointmentStatus
	Source             AppointmentSource
	Notes              string
	CancellationReason string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

// AppointmentAuditEvent records changes to an appointment.
type AppointmentAuditEvent struct {
	ID            uuid.UUID
	TenantID      uuid.UUID
	AppointmentID uuid.UUID
	ActorUserID   *uuid.UUID
	EventType     string
	OldStatus     string
	NewStatus     string
	OldStartAt    *time.Time
	NewStartAt    *time.Time
	Notes         string
	CreatedAt     time.Time
}

// TimeSlot represents an available time slot for booking.
type TimeSlot struct {
	Start       time.Time   `json:"start"`
	End         time.Time   `json:"end"`
	ServiceMode ServiceMode `json:"service_mode"`
	AddressID   *uuid.UUID  `json:"address_id,omitempty"`
}
