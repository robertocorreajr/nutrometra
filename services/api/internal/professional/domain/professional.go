package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// RegistrationType defines the professional registration body.
type RegistrationType string

const (
	RegistrationCRN   RegistrationType = "CRN"
	RegistrationCRM   RegistrationType = "CRM"
	RegistrationOther RegistrationType = "other"
)

// ServiceMode defines how a professional delivers care.
type ServiceMode string

const (
	ServiceModeOnsite    ServiceMode = "onsite"
	ServiceModeOnline    ServiceMode = "online"
	ServiceModeHomeVisit ServiceMode = "home_visit"
)

// Sentinel errors.
var (
	ErrNotFound              = errors.New("professional: not found")
	ErrDuplicateRegistration = errors.New("professional: duplicate registration")
	ErrDuplicateUser         = errors.New("professional: user already registered in tenant")
	ErrAddressNotFound       = errors.New("professional: address not found")
	ErrServiceModeNotFound   = errors.New("professional: service mode not found")
	ErrEntitlementExceeded   = errors.New("professional: entitlement limit exceeded")
)

// Professional represents a healthcare professional within a tenant.
type Professional struct {
	ID                 uuid.UUID
	TenantID           uuid.UUID
	UserID             uuid.UUID
	FullName           string
	RegistrationType   RegistrationType
	RegistrationNumber string
	RegistrationState  string
	Specialty          string
	Bio                string
	Phone              string
	AvatarURL          string
	Active             bool
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

// Validate checks business rules for a Professional.
func (p *Professional) Validate() error {
	if p.FullName == "" {
		return errors.New("professional: full_name is required")
	}
	if p.RegistrationType == "" {
		return errors.New("professional: registration_type is required")
	}
	if p.RegistrationNumber == "" {
		return errors.New("professional: registration_number is required")
	}
	if (p.RegistrationType == RegistrationCRN || p.RegistrationType == RegistrationCRM) && p.RegistrationState == "" {
		return errors.New("professional: registration_state is required for CRN/CRM")
	}
	return nil
}

// Address represents a physical location where a professional provides care.
type Address struct {
	ID             uuid.UUID
	TenantID       uuid.UUID
	ProfessionalID uuid.UUID
	Label          string
	Street         string
	Number         string
	Complement     string
	Neighborhood   string
	City           string
	State          string
	ZipCode        string
	Country        string
	Latitude       *float64
	Longitude      *float64
	Phone          string
	Notes          string
	Active         bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// Validate checks business rules for an Address.
func (a *Address) Validate() error {
	if a.Label == "" {
		return errors.New("address: label is required")
	}
	if a.Street == "" {
		return errors.New("address: street is required")
	}
	if a.City == "" {
		return errors.New("address: city is required")
	}
	if a.State == "" {
		return errors.New("address: state is required")
	}
	if a.ZipCode == "" {
		return errors.New("address: zip_code is required")
	}
	return nil
}

// ProfessionalServiceMode represents a service mode for a professional.
type ProfessionalServiceMode struct {
	ID             uuid.UUID
	TenantID       uuid.UUID
	ProfessionalID uuid.UUID
	Mode           ServiceMode
	AddressID      *uuid.UUID
	DurationMin    int
	Active         bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
