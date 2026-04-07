package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Gender constants.
type Gender string

const (
	GenderMale           Gender = "male"
	GenderFemale         Gender = "female"
	GenderOther          Gender = "other"
	GenderPreferNotToSay Gender = "prefer_not_to_say"
)

// Sentinel errors.
var (
	ErrNotFound            = errors.New("patient: not found")
	ErrDuplicateCPF        = errors.New("patient: CPF already registered in this tenant")
	ErrInviteNotFound      = errors.New("patient: invite not found or expired")
	ErrInviteMaxUsed       = errors.New("patient: invite has reached maximum uses")
	ErrInviteExpired       = errors.New("patient: invite has expired")
	ErrAccessLinkExists    = errors.New("patient: user already linked to this patient")
	ErrEntitlementExceeded = errors.New("patient: entitlement limit exceeded")
	ErrProfileNotFound     = errors.New("patient: profile not found")
)

// Patient represents a patient within a tenant.
type Patient struct {
	ID             uuid.UUID
	TenantID       uuid.UUID
	ProfessionalID uuid.UUID
	FullName       string
	Email          string
	Phone          string
	CPF            string
	DateOfBirth    *time.Time
	Gender         string
	Notes          string
	Active         bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// Validate checks business rules for a Patient.
func (p *Patient) Validate() error {
	if p.FullName == "" {
		return errors.New("patient: full_name is required")
	}
	if p.ProfessionalID == uuid.Nil {
		return errors.New("patient: professional_id is required")
	}
	return nil
}

// PatientInvite represents an invitation for a patient to access the portal.
type PatientInvite struct {
	ID        uuid.UUID
	TenantID  uuid.UUID
	PatientID uuid.UUID
	InvitedBy uuid.UUID
	CodeHash  string
	MaxUses   int
	UsedCount int
	ExpiresAt time.Time
	CreatedAt time.Time
}

// PatientAccessLink links a user account to a patient record.
type PatientAccessLink struct {
	ID        uuid.UUID
	TenantID  uuid.UUID
	PatientID uuid.UUID
	UserID    uuid.UUID
	InviteID  uuid.UUID
	Active    bool
	CreatedAt time.Time
}

// PatientProfile contains extended patient information.
type PatientProfile struct {
	ID                    uuid.UUID
	TenantID              uuid.UUID
	PatientID             uuid.UUID
	Occupation            string
	MaritalStatus         string
	Ethnicity             string
	BloodType             string
	Allergies             []string
	ChronicConditions     []string
	Medications           []string
	EmergencyContactName  string
	EmergencyContactPhone string
	CreatedAt             time.Time
	UpdatedAt             time.Time
}
