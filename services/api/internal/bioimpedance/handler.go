package bioimpedance

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"nutrometra/api/internal/bioimpedance/domain"
	"nutrometra/api/internal/bioimpedance/usecase"

	identitydomain "nutrometra/api/internal/identity/domain"
	"nutrometra/api/internal/platform/audit"
	"nutrometra/api/internal/platform/server"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Handler exposes HTTP endpoints for bioimpedance operations.
type Handler struct {
	uc       *usecase.Usecase
	pool     *pgxpool.Pool
	auditSvc *audit.Service
}

// NewHandler creates a bioimpedance Handler.
func NewHandler(uc *usecase.Usecase, pool *pgxpool.Pool, auditSvc *audit.Service) *Handler {
	return &Handler{uc: uc, pool: pool, auditSvc: auditSvc}
}

type measurementRequest struct {
	PatientID      string   `json:"patient_id"`
	ProfessionalID string   `json:"professional_id"`
	MeasuredAt     *string  `json:"measured_at"` // RFC3339, optional
	WeightKg       *float64 `json:"weight_kg"`
	HeightCm       *float64 `json:"height_cm"`
	BodyFatPct     *float64 `json:"body_fat_pct"`
	LeanMassKg     *float64 `json:"lean_mass_kg"`
	FatMassKg      *float64 `json:"fat_mass_kg"`
	MuscleMassKg   *float64 `json:"muscle_mass_kg"`
	BoneMassKg     *float64 `json:"bone_mass_kg"`
	WaterPct       *float64 `json:"water_pct"`
	VisceralFat    *float64 `json:"visceral_fat"`
	BasalMetabolicRate *int `json:"basal_metabolic_rate"`
	WaistCm        *float64 `json:"waist_cm"`
	HipCm          *float64 `json:"hip_cm"`
	ChestCm        *float64 `json:"chest_cm"`
	RightArmCm     *float64 `json:"right_arm_cm"`
	LeftArmCm      *float64 `json:"left_arm_cm"`
	RightThighCm   *float64 `json:"right_thigh_cm"`
	LeftThighCm    *float64 `json:"left_thigh_cm"`
	RightCalfCm    *float64 `json:"right_calf_cm"`
	LeftCalfCm     *float64 `json:"left_calf_cm"`
	NeckCm         *float64 `json:"neck_cm"`
	AbdomenCm      *float64 `json:"abdomen_cm"`
	TricepsSfMm    *float64 `json:"triceps_sf_mm"`
	BicepsSfMm     *float64 `json:"biceps_sf_mm"`
	SubscapularSfMm *float64 `json:"subscapular_sf_mm"`
	SuprailiacSfMm *float64 `json:"suprailiac_sf_mm"`
	AbdominalSfMm  *float64 `json:"abdominal_sf_mm"`
	ThighSfMm      *float64 `json:"thigh_sf_mm"`
	CalfSfMm       *float64 `json:"calf_sf_mm"`
	Source         string   `json:"source"`
	DeviceModel    string   `json:"device_model"`
	Notes          string   `json:"notes"`
}

type measurementResponse struct {
	ID             string   `json:"id"`
	PatientID      string   `json:"patient_id"`
	ProfessionalID string   `json:"professional_id"`
	MeasuredAt     string   `json:"measured_at"`
	WeightKg       *float64 `json:"weight_kg,omitempty"`
	HeightCm       *float64 `json:"height_cm,omitempty"`
	BMI            *float64 `json:"bmi,omitempty"`
	BodyFatPct     *float64 `json:"body_fat_pct,omitempty"`
	LeanMassKg     *float64 `json:"lean_mass_kg,omitempty"`
	FatMassKg      *float64 `json:"fat_mass_kg,omitempty"`
	MuscleMassKg   *float64 `json:"muscle_mass_kg,omitempty"`
	BoneMassKg     *float64 `json:"bone_mass_kg,omitempty"`
	WaterPct       *float64 `json:"water_pct,omitempty"`
	VisceralFat    *float64 `json:"visceral_fat,omitempty"`
	BasalMetabolicRate *int `json:"basal_metabolic_rate,omitempty"`
	WaistCm        *float64 `json:"waist_cm,omitempty"`
	HipCm          *float64 `json:"hip_cm,omitempty"`
	Source         string   `json:"source"`
	DeviceModel    string   `json:"device_model,omitempty"`
	Notes          string   `json:"notes,omitempty"`
	CreatedAt      string   `json:"created_at"`
}

func toMeasurementResponse(m *domain.BodyMeasurement) measurementResponse {
	return measurementResponse{
		ID:             m.ID.String(),
		PatientID:      m.PatientID.String(),
		ProfessionalID: m.ProfessionalID.String(),
		MeasuredAt:     m.MeasuredAt.Format(time.RFC3339),
		WeightKg:       m.WeightKg,
		HeightCm:       m.HeightCm,
		BMI:            m.BMI,
		BodyFatPct:     m.BodyFatPct,
		LeanMassKg:     m.LeanMassKg,
		FatMassKg:      m.FatMassKg,
		MuscleMassKg:   m.MuscleMassKg,
		BoneMassKg:     m.BoneMassKg,
		WaterPct:       m.WaterPct,
		VisceralFat:    m.VisceralFat,
		BasalMetabolicRate: m.BasalMetabolicRate,
		WaistCm:        m.WaistCm,
		HipCm:          m.HipCm,
		Source:         string(m.Source),
		DeviceModel:    m.DeviceModel,
		Notes:          m.Notes,
		CreatedAt:      m.CreatedAt.Format(time.RFC3339),
	}
}

// Create handles POST /measurements
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := identitydomain.TenantIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusBadRequest, "missing_tenant", "No tenant in context")
		return
	}

	var req measurementRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_body", "Invalid request body")
		return
	}

	patientID, err := uuid.Parse(req.PatientID)
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_patient_id", "Invalid patient_id")
		return
	}
	profID, err := uuid.Parse(req.ProfessionalID)
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_professional_id", "Invalid professional_id")
		return
	}

	m := &domain.BodyMeasurement{
		TenantID:       tenantID,
		PatientID:      patientID,
		ProfessionalID: profID,
		WeightKg:       req.WeightKg,
		HeightCm:       req.HeightCm,
		BodyFatPct:     req.BodyFatPct,
		LeanMassKg:     req.LeanMassKg,
		FatMassKg:      req.FatMassKg,
		MuscleMassKg:   req.MuscleMassKg,
		BoneMassKg:     req.BoneMassKg,
		WaterPct:       req.WaterPct,
		VisceralFat:    req.VisceralFat,
		BasalMetabolicRate: req.BasalMetabolicRate,
		WaistCm:        req.WaistCm,
		HipCm:          req.HipCm,
		ChestCm:        req.ChestCm,
		RightArmCm:     req.RightArmCm,
		LeftArmCm:      req.LeftArmCm,
		RightThighCm:   req.RightThighCm,
		LeftThighCm:    req.LeftThighCm,
		RightCalfCm:    req.RightCalfCm,
		LeftCalfCm:     req.LeftCalfCm,
		NeckCm:         req.NeckCm,
		AbdomenCm:      req.AbdomenCm,
		TricepsSfMm:    req.TricepsSfMm,
		BicepsSfMm:     req.BicepsSfMm,
		SubscapularSfMm: req.SubscapularSfMm,
		SuprailiacSfMm: req.SuprailiacSfMm,
		AbdominalSfMm:  req.AbdominalSfMm,
		ThighSfMm:      req.ThighSfMm,
		CalfSfMm:       req.CalfSfMm,
		Source:         domain.MeasurementSource(req.Source),
		DeviceModel:    req.DeviceModel,
		Notes:          req.Notes,
	}

	if req.MeasuredAt != nil {
		t, err := time.Parse(time.RFC3339, *req.MeasuredAt)
		if err != nil {
			server.RenderError(w, r, http.StatusBadRequest, "invalid_measured_at", "measured_at must be RFC3339")
			return
		}
		m.MeasuredAt = t
	}

	if err := h.uc.Create(r.Context(), m); err != nil {
		if errors.Is(err, domain.ErrInvalidWeight) || errors.Is(err, domain.ErrInvalidHeight) {
			server.RenderError(w, r, http.StatusBadRequest, "validation_failed", err.Error())
			return
		}
		server.RenderError(w, r, http.StatusInternalServerError, "create_failed", "Failed to create measurement")
		return
	}

	actorID, _ := identitydomain.UserIDFromContext(r.Context())
	entry := audit.NewEntry(
		audit.WithTenantID(tenantID),
		audit.WithActor(actorID, audit.ScopeTenant),
		audit.WithEntity("body_measurement", m.ID),
		audit.WithAction("measurement_created"),
		audit.WithIPAddress(r.RemoteAddr),
	)
	_ = h.auditSvc.Write(r.Context(), h.pool, entry)

	server.RenderJSON(w, http.StatusCreated, toMeasurementResponse(m))
}

// GetByID handles GET /measurements/{id}
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := identitydomain.TenantIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusBadRequest, "missing_tenant", "No tenant in context")
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_id", "Invalid measurement ID")
		return
	}

	m, err := h.uc.GetByID(r.Context(), tenantID, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			server.RenderError(w, r, http.StatusNotFound, "not_found", "Measurement not found")
			return
		}
		server.RenderError(w, r, http.StatusInternalServerError, "get_failed", "Failed to get measurement")
		return
	}
	server.RenderJSON(w, http.StatusOK, toMeasurementResponse(m))
}

// ListByPatient handles GET /patients/{patient_id}/measurements
func (h *Handler) ListByPatient(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := identitydomain.TenantIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusBadRequest, "missing_tenant", "No tenant in context")
		return
	}

	patientID, err := uuid.Parse(chi.URLParam(r, "patient_id"))
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_patient_id", "Invalid patient ID")
		return
	}

	measurements, err := h.uc.ListByPatient(r.Context(), tenantID, patientID)
	if err != nil {
		server.RenderError(w, r, http.StatusInternalServerError, "list_failed", "Failed to list measurements")
		return
	}

	resp := make([]measurementResponse, len(measurements))
	for i := range measurements {
		resp[i] = toMeasurementResponse(&measurements[i])
	}
	server.RenderJSON(w, http.StatusOK, resp)
}

// Publish handles POST /measurements/{id}/publish
type publishRequest struct {
	Message string `json:"message"`
}

func (h *Handler) Publish(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := identitydomain.TenantIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusBadRequest, "missing_tenant", "No tenant in context")
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_id", "Invalid measurement ID")
		return
	}

	var req publishRequest
	_ = json.NewDecoder(r.Body).Decode(&req)

	actorID, _ := identitydomain.UserIDFromContext(r.Context())

	if err := h.uc.PublishToPatient(r.Context(), tenantID, id, actorID, req.Message); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			server.RenderError(w, r, http.StatusNotFound, "not_found", "Measurement not found")
			return
		}
		if errors.Is(err, domain.ErrAlreadyPublished) {
			server.RenderError(w, r, http.StatusConflict, "already_published", "Measurement already published")
			return
		}
		server.RenderError(w, r, http.StatusInternalServerError, "publish_failed", "Failed to publish measurement")
		return
	}

	entry := audit.NewEntry(
		audit.WithTenantID(tenantID),
		audit.WithActor(actorID, audit.ScopeTenant),
		audit.WithEntity("body_measurement", id),
		audit.WithAction("measurement_published"),
		audit.WithIPAddress(r.RemoteAddr),
	)
	_ = h.auditSvc.Write(r.Context(), h.pool, entry)

	w.WriteHeader(http.StatusNoContent)
}
