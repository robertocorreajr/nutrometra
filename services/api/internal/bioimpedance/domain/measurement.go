package domain

import (
	"errors"
	"math"
	"time"

	"github.com/google/uuid"
)

// MeasurementSource defines how the measurement was captured.
type MeasurementSource string

const (
	SourceManual MeasurementSource = "manual"
	SourceDevice MeasurementSource = "device"
	SourceImport MeasurementSource = "import"
)

// Sentinel errors.
var (
	ErrNotFound         = errors.New("bioimpedance: not found")
	ErrAlreadyPublished = errors.New("bioimpedance: measurement already published")
	ErrInvalidWeight    = errors.New("bioimpedance: weight must be between 1 and 500 kg")
	ErrInvalidHeight    = errors.New("bioimpedance: height must be between 30 and 300 cm")
)

// BodyMeasurement represents a bioimpedance/anthropometric assessment.
type BodyMeasurement struct {
	ID             uuid.UUID
	TenantID       uuid.UUID
	PatientID      uuid.UUID
	ProfessionalID uuid.UUID
	MeasuredAt     time.Time

	// Anthropometric
	WeightKg *float64
	HeightCm *float64
	BMI      *float64

	// Body composition
	BodyFatPct       *float64
	LeanMassKg       *float64
	FatMassKg        *float64
	MuscleMassKg     *float64
	BoneMassKg       *float64
	WaterPct         *float64
	VisceralFat      *float64
	BasalMetabolicRate *int

	// Circumferences
	WaistCm     *float64
	HipCm       *float64
	ChestCm     *float64
	RightArmCm  *float64
	LeftArmCm   *float64
	RightThighCm *float64
	LeftThighCm  *float64
	RightCalfCm  *float64
	LeftCalfCm   *float64
	NeckCm       *float64
	AbdomenCm    *float64

	// Skinfolds
	TricepsSfMm      *float64
	BicepsSfMm       *float64
	SubscapularSfMm  *float64
	SuprailiacSfMm   *float64
	AbdominalSfMm    *float64
	ThighSfMm        *float64
	CalfSfMm         *float64

	Source      MeasurementSource
	DeviceModel string
	Notes       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Validate checks measurement constraints and auto-calculates BMI.
func (m *BodyMeasurement) Validate() error {
	if m.WeightKg != nil {
		if *m.WeightKg < 1 || *m.WeightKg > 500 {
			return ErrInvalidWeight
		}
	}
	if m.HeightCm != nil {
		if *m.HeightCm < 30 || *m.HeightCm > 300 {
			return ErrInvalidHeight
		}
	}
	return nil
}

// CalculateBMI computes BMI from weight and height if both are present.
func (m *BodyMeasurement) CalculateBMI() {
	if m.WeightKg != nil && m.HeightCm != nil && *m.HeightCm > 0 {
		heightM := *m.HeightCm / 100.0
		bmi := *m.WeightKg / (heightM * heightM)
		bmi = math.Round(bmi*100) / 100
		m.BMI = &bmi
	}
}

// MeasurementPublication records when a measurement is shared with the patient.
type MeasurementPublication struct {
	ID            uuid.UUID
	TenantID      uuid.UUID
	MeasurementID uuid.UUID
	PatientID     uuid.UUID
	PublishedBy   uuid.UUID
	PublishedAt   time.Time
	Message       string
}
