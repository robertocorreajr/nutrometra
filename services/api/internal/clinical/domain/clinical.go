package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// AnamnesisStatus represents the lifecycle state of an anamnesis.
type AnamnesisStatus string

const (
	AnamnesisStatusDraft     AnamnesisStatus = "draft"
	AnamnesisStatusFinalized AnamnesisStatus = "finalized"
)

// AttachmentCategory defines the type of clinical attachment.
type AttachmentCategory string

const (
	CategoryGeneral      AttachmentCategory = "general"
	CategoryExam         AttachmentCategory = "exam"
	CategoryLabResult    AttachmentCategory = "lab_result"
	CategoryPrescription AttachmentCategory = "prescription"
	CategoryPhoto        AttachmentCategory = "photo"
	CategoryOther        AttachmentCategory = "other"
)

// Sentinel errors.
var (
	ErrNotFound         = errors.New("clinical: not found")
	ErrAlreadyFinalized = errors.New("clinical: anamnesis already finalized")
)

// Anamnesis represents a clinical intake assessment.
type Anamnesis struct {
	ID                     uuid.UUID
	TenantID               uuid.UUID
	PatientID              uuid.UUID
	ProfessionalID         uuid.UUID
	Status                 AnamnesisStatus
	ChiefComplaint         string
	HistoryPresentIllness  string
	PastMedicalHistory     string
	FamilyHistory          string
	SocialHistory          string
	DietaryHistory         string
	PhysicalActivity       string
	SleepPattern           string
	BowelHabits            string
	WaterIntake            string
	Supplements            string
	Observations           string
	FinalizedAt            *time.Time
	CreatedAt              time.Time
	UpdatedAt              time.Time
}

// ProgressNote represents a clinical progress note.
type ProgressNote struct {
	ID               uuid.UUID
	TenantID         uuid.UUID
	PatientID        uuid.UUID
	ProfessionalID   uuid.UUID
	AppointmentID    *uuid.UUID
	Title            string
	Content          string
	VisibleToPatient bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// Validate checks business rules for a ProgressNote.
func (n *ProgressNote) Validate() error {
	if n.Title == "" {
		return errors.New("clinical: title is required")
	}
	if n.Content == "" {
		return errors.New("clinical: content is required")
	}
	return nil
}

// ClinicalAttachment represents a file attached to a patient's record.
type ClinicalAttachment struct {
	ID             uuid.UUID
	TenantID       uuid.UUID
	PatientID      uuid.UUID
	ProfessionalID uuid.UUID
	FileName       string
	FileType       string
	FileSizeBytes  int64
	StorageKey     string
	Category       AttachmentCategory
	Description    string
	CreatedAt      time.Time
}
