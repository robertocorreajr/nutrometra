package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// ExportStatus represents the lifecycle state of an export.
type ExportStatus string

const (
	StatusPending    ExportStatus = "pending"
	StatusProcessing ExportStatus = "processing"
	StatusCompleted  ExportStatus = "completed"
	StatusFailed     ExportStatus = "failed"
)

// ExportType represents the kind of export.
type ExportType string

const (
	TypePDF      ExportType = "pdf"
	TypePrintJob ExportType = "print_job"
)

var (
	ErrNotFound = errors.New("export: not found")
)

// ExportedFile represents a requested or completed file export.
type ExportedFile struct {
	ID                uuid.UUID
	TenantID          uuid.UUID
	RelatedEntityType string
	RelatedEntityID   uuid.UUID
	ExportType        ExportType
	FileKey           string
	Status            ExportStatus
	RequestedByUserID uuid.UUID
	CreatedAt         time.Time
	CompletedAt       *time.Time
	FailureReason     string
}
