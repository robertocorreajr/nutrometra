package usecase

import (
	"context"

	"github.com/google/uuid"
)

// PatientContext holds clinical data about a patient for AI prompt construction.
type PatientContext struct {
	PatientID   uuid.UUID              `json:"patient_id"`
	PatientName string                 `json:"patient_name,omitempty"`
	Age         *int                   `json:"age,omitempty"`
	Profile     map[string]interface{} `json:"profile,omitempty"`
	Note        string                 `json:"note"`
}

// PatientContextProvider fetches clinical context for AI suggestions.
// Will be implemented by patient/clinical modules in future phases.
type PatientContextProvider interface {
	GetPatientContext(ctx context.Context, tenantID, patientID uuid.UUID) (*PatientContext, error)
}

// StubPatientContextProvider returns empty context -- placeholder until clinical modules are wired.
type StubPatientContextProvider struct{}

// GetPatientContext returns a stub patient context with a note indicating modules are not yet available.
func (s *StubPatientContextProvider) GetPatientContext(_ context.Context, _, patientID uuid.UUID) (*PatientContext, error) {
	return &PatientContext{
		PatientID: patientID,
		Note:      "patient data modules not yet available",
	}, nil
}
