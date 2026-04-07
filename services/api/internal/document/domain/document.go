package domain

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

// DocumentType represents the kind of clinical document.
type DocumentType string

const (
	TypeExamRequest  DocumentType = "exam_request"
	TypePrescription DocumentType = "prescription"
	TypeLetter       DocumentType = "letter"
	TypeOther        DocumentType = "other"
)

// DocumentStatus represents the lifecycle state of a clinical document.
type DocumentStatus string

const (
	StatusDraft     DocumentStatus = "draft"
	StatusFinalized DocumentStatus = "finalized"
	StatusPublished DocumentStatus = "published"
)

// Sentinel errors.
var (
	ErrNotFound         = errors.New("document: not found")
	ErrNotDraft         = errors.New("document: operation requires draft status")
	ErrNotFinalized     = errors.New("document: operation requires finalized status")
	ErrAlreadyFinalized = errors.New("document: document is already finalized")
)

// validDocumentTypes holds the set of accepted document types.
var validDocumentTypes = map[DocumentType]bool{
	TypeExamRequest:  true,
	TypePrescription: true,
	TypeLetter:       true,
	TypeOther:        true,
}

// ClinicalDocument represents a clinical document (exam request, prescription, letter, etc.).
type ClinicalDocument struct {
	ID                uuid.UUID
	TenantID          uuid.UUID
	PatientID         uuid.UUID
	ProfessionalID    uuid.UUID
	AppointmentID     *uuid.UUID
	DocumentType      DocumentType
	Title             string
	Status            DocumentStatus
	ContentJSON       json.RawMessage
	VersionNumber     int
	PreviousVersionID *uuid.UUID
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// Validate checks business rules for a ClinicalDocument.
func (d *ClinicalDocument) Validate() error {
	if d.Title == "" {
		return errors.New("document: title is required")
	}
	if !validDocumentTypes[d.DocumentType] {
		return errors.New("document: invalid document_type")
	}
	return nil
}

// DocumentPublication represents a publication record linking a published document to a patient.
type DocumentPublication struct {
	ID                 uuid.UUID
	TenantID           uuid.UUID
	ClinicalDocumentID uuid.UUID
	PatientID          uuid.UUID
	PublishedByUserID  uuid.UUID
	PublishedAt        time.Time
}
