package clinical

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"nutrometra/api/internal/clinical/domain"
	"nutrometra/api/internal/clinical/usecase"

	identitydomain "nutrometra/api/internal/identity/domain"
	"nutrometra/api/internal/platform/audit"
	"nutrometra/api/internal/platform/server"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Handler exposes HTTP endpoints for clinical operations.
type Handler struct {
	uc       *usecase.Usecase
	pool     *pgxpool.Pool
	auditSvc *audit.Service
}

// NewHandler creates a clinical Handler.
func NewHandler(uc *usecase.Usecase, pool *pgxpool.Pool, auditSvc *audit.Service) *Handler {
	return &Handler{uc: uc, pool: pool, auditSvc: auditSvc}
}

// --- Anamnesis ---

type anamnesisRequest struct {
	PatientID              string `json:"patient_id"`
	ProfessionalID         string `json:"professional_id"`
	ChiefComplaint         string `json:"chief_complaint"`
	HistoryPresentIllness  string `json:"history_present_illness"`
	PastMedicalHistory     string `json:"past_medical_history"`
	FamilyHistory          string `json:"family_history"`
	SocialHistory          string `json:"social_history"`
	DietaryHistory         string `json:"dietary_history"`
	PhysicalActivity       string `json:"physical_activity"`
	SleepPattern           string `json:"sleep_pattern"`
	BowelHabits            string `json:"bowel_habits"`
	WaterIntake            string `json:"water_intake"`
	Supplements            string `json:"supplements"`
	Observations           string `json:"observations"`
}

type anamnesisResponse struct {
	ID                     string  `json:"id"`
	PatientID              string  `json:"patient_id"`
	ProfessionalID         string  `json:"professional_id"`
	Status                 string  `json:"status"`
	ChiefComplaint         string  `json:"chief_complaint,omitempty"`
	HistoryPresentIllness  string  `json:"history_present_illness,omitempty"`
	PastMedicalHistory     string  `json:"past_medical_history,omitempty"`
	FamilyHistory          string  `json:"family_history,omitempty"`
	SocialHistory          string  `json:"social_history,omitempty"`
	DietaryHistory         string  `json:"dietary_history,omitempty"`
	PhysicalActivity       string  `json:"physical_activity,omitempty"`
	SleepPattern           string  `json:"sleep_pattern,omitempty"`
	BowelHabits            string  `json:"bowel_habits,omitempty"`
	WaterIntake            string  `json:"water_intake,omitempty"`
	Supplements            string  `json:"supplements,omitempty"`
	Observations           string  `json:"observations,omitempty"`
	FinalizedAt            *string `json:"finalized_at,omitempty"`
	CreatedAt              string  `json:"created_at"`
	UpdatedAt              string  `json:"updated_at"`
}

func toAnamnesisResponse(a *domain.Anamnesis) anamnesisResponse {
	resp := anamnesisResponse{
		ID:                    a.ID.String(),
		PatientID:             a.PatientID.String(),
		ProfessionalID:        a.ProfessionalID.String(),
		Status:                string(a.Status),
		ChiefComplaint:        a.ChiefComplaint,
		HistoryPresentIllness: a.HistoryPresentIllness,
		PastMedicalHistory:    a.PastMedicalHistory,
		FamilyHistory:         a.FamilyHistory,
		SocialHistory:         a.SocialHistory,
		DietaryHistory:        a.DietaryHistory,
		PhysicalActivity:      a.PhysicalActivity,
		SleepPattern:          a.SleepPattern,
		BowelHabits:           a.BowelHabits,
		WaterIntake:           a.WaterIntake,
		Supplements:           a.Supplements,
		Observations:          a.Observations,
		CreatedAt:             a.CreatedAt.Format(time.RFC3339),
		UpdatedAt:             a.UpdatedAt.Format(time.RFC3339),
	}
	if a.FinalizedAt != nil {
		s := a.FinalizedAt.Format(time.RFC3339)
		resp.FinalizedAt = &s
	}
	return resp
}

// CreateAnamnesis handles POST /patients/{patient_id}/anamneses
func (h *Handler) CreateAnamnesis(w http.ResponseWriter, r *http.Request) {
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

	var req anamnesisRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_body", "Invalid request body")
		return
	}

	profID, err := uuid.Parse(req.ProfessionalID)
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_professional_id", "Invalid professional_id")
		return
	}

	a := &domain.Anamnesis{
		TenantID:              tenantID,
		PatientID:             patientID,
		ProfessionalID:        profID,
		ChiefComplaint:        req.ChiefComplaint,
		HistoryPresentIllness: req.HistoryPresentIllness,
		PastMedicalHistory:    req.PastMedicalHistory,
		FamilyHistory:         req.FamilyHistory,
		SocialHistory:         req.SocialHistory,
		DietaryHistory:        req.DietaryHistory,
		PhysicalActivity:      req.PhysicalActivity,
		SleepPattern:          req.SleepPattern,
		BowelHabits:           req.BowelHabits,
		WaterIntake:           req.WaterIntake,
		Supplements:           req.Supplements,
		Observations:          req.Observations,
	}

	if err := h.uc.CreateAnamnesis(r.Context(), a); err != nil {
		server.RenderError(w, r, http.StatusInternalServerError, "create_failed", "Failed to create anamnesis")
		return
	}

	actorID, _ := identitydomain.UserIDFromContext(r.Context())
	entry := audit.NewEntry(
		audit.WithTenantID(tenantID),
		audit.WithActor(actorID, audit.ScopeTenant),
		audit.WithEntity("anamnesis", a.ID),
		audit.WithAction("anamnesis_created"),
		audit.WithIPAddress(r.RemoteAddr),
	)
	_ = h.auditSvc.Write(r.Context(), h.pool, entry)

	server.RenderJSON(w, http.StatusCreated, toAnamnesisResponse(a))
}

// GetAnamnesis handles GET /anamneses/{id}
func (h *Handler) GetAnamnesis(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := identitydomain.TenantIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusBadRequest, "missing_tenant", "No tenant in context")
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_id", "Invalid anamnesis ID")
		return
	}

	a, err := h.uc.GetAnamnesisByID(r.Context(), tenantID, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			server.RenderError(w, r, http.StatusNotFound, "not_found", "Anamnesis not found")
			return
		}
		server.RenderError(w, r, http.StatusInternalServerError, "get_failed", "Failed to get anamnesis")
		return
	}
	server.RenderJSON(w, http.StatusOK, toAnamnesisResponse(a))
}

// ListAnamneses handles GET /patients/{patient_id}/anamneses
func (h *Handler) ListAnamneses(w http.ResponseWriter, r *http.Request) {
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

	list, err := h.uc.ListAnamnesesByPatient(r.Context(), tenantID, patientID)
	if err != nil {
		server.RenderError(w, r, http.StatusInternalServerError, "list_failed", "Failed to list anamneses")
		return
	}

	resp := make([]anamnesisResponse, len(list))
	for i := range list {
		resp[i] = toAnamnesisResponse(&list[i])
	}
	server.RenderJSON(w, http.StatusOK, resp)
}

// UpdateAnamnesis handles PUT /anamneses/{id}
func (h *Handler) UpdateAnamnesis(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := identitydomain.TenantIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusBadRequest, "missing_tenant", "No tenant in context")
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_id", "Invalid anamnesis ID")
		return
	}

	var req anamnesisRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_body", "Invalid request body")
		return
	}

	profID, err := uuid.Parse(req.ProfessionalID)
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_professional_id", "Invalid professional_id")
		return
	}

	a := &domain.Anamnesis{
		ID:                    id,
		TenantID:              tenantID,
		ProfessionalID:        profID,
		ChiefComplaint:        req.ChiefComplaint,
		HistoryPresentIllness: req.HistoryPresentIllness,
		PastMedicalHistory:    req.PastMedicalHistory,
		FamilyHistory:         req.FamilyHistory,
		SocialHistory:         req.SocialHistory,
		DietaryHistory:        req.DietaryHistory,
		PhysicalActivity:      req.PhysicalActivity,
		SleepPattern:          req.SleepPattern,
		BowelHabits:           req.BowelHabits,
		WaterIntake:           req.WaterIntake,
		Supplements:           req.Supplements,
		Observations:          req.Observations,
		Status:                domain.AnamnesisStatusDraft,
	}

	if err := h.uc.UpdateAnamnesis(r.Context(), a); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			server.RenderError(w, r, http.StatusNotFound, "not_found", "Anamnesis not found")
			return
		}
		if errors.Is(err, domain.ErrAlreadyFinalized) {
			server.RenderError(w, r, http.StatusConflict, "already_finalized", "Anamnesis has been finalized and cannot be modified")
			return
		}
		server.RenderError(w, r, http.StatusInternalServerError, "update_failed", "Failed to update anamnesis")
		return
	}

	server.RenderJSON(w, http.StatusOK, toAnamnesisResponse(a))
}

// FinalizeAnamnesis handles POST /anamneses/{id}/finalize
func (h *Handler) FinalizeAnamnesis(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := identitydomain.TenantIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusBadRequest, "missing_tenant", "No tenant in context")
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_id", "Invalid anamnesis ID")
		return
	}

	if err := h.uc.FinalizeAnamnesis(r.Context(), tenantID, id); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			server.RenderError(w, r, http.StatusNotFound, "not_found", "Anamnesis not found")
			return
		}
		if errors.Is(err, domain.ErrAlreadyFinalized) {
			server.RenderError(w, r, http.StatusConflict, "already_finalized", "Anamnesis already finalized")
			return
		}
		server.RenderError(w, r, http.StatusInternalServerError, "finalize_failed", "Failed to finalize anamnesis")
		return
	}

	actorID, _ := identitydomain.UserIDFromContext(r.Context())
	entry := audit.NewEntry(
		audit.WithTenantID(tenantID),
		audit.WithActor(actorID, audit.ScopeTenant),
		audit.WithEntity("anamnesis", id),
		audit.WithAction("anamnesis_finalized"),
		audit.WithIPAddress(r.RemoteAddr),
	)
	_ = h.auditSvc.Write(r.Context(), h.pool, entry)

	w.WriteHeader(http.StatusNoContent)
}

// --- Progress Notes ---

type noteRequest struct {
	ProfessionalID   string  `json:"professional_id"`
	AppointmentID    *string `json:"appointment_id"`
	Title            string  `json:"title"`
	Content          string  `json:"content"`
	VisibleToPatient bool    `json:"visible_to_patient"`
}

type noteResponse struct {
	ID               string  `json:"id"`
	PatientID        string  `json:"patient_id"`
	ProfessionalID   string  `json:"professional_id"`
	AppointmentID    *string `json:"appointment_id,omitempty"`
	Title            string  `json:"title"`
	Content          string  `json:"content"`
	VisibleToPatient bool    `json:"visible_to_patient"`
	CreatedAt        string  `json:"created_at"`
	UpdatedAt        string  `json:"updated_at"`
}

func toNoteResponse(n *domain.ProgressNote) noteResponse {
	resp := noteResponse{
		ID:               n.ID.String(),
		PatientID:        n.PatientID.String(),
		ProfessionalID:   n.ProfessionalID.String(),
		Title:            n.Title,
		Content:          n.Content,
		VisibleToPatient: n.VisibleToPatient,
		CreatedAt:        n.CreatedAt.Format(time.RFC3339),
		UpdatedAt:        n.UpdatedAt.Format(time.RFC3339),
	}
	if n.AppointmentID != nil {
		s := n.AppointmentID.String()
		resp.AppointmentID = &s
	}
	return resp
}

// CreateProgressNote handles POST /patients/{patient_id}/notes
func (h *Handler) CreateProgressNote(w http.ResponseWriter, r *http.Request) {
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

	var req noteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_body", "Invalid request body")
		return
	}

	profID, err := uuid.Parse(req.ProfessionalID)
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_professional_id", "Invalid professional_id")
		return
	}

	n := &domain.ProgressNote{
		TenantID:         tenantID,
		PatientID:        patientID,
		ProfessionalID:   profID,
		Title:            req.Title,
		Content:          req.Content,
		VisibleToPatient: req.VisibleToPatient,
	}

	if req.AppointmentID != nil {
		aid, err := uuid.Parse(*req.AppointmentID)
		if err != nil {
			server.RenderError(w, r, http.StatusBadRequest, "invalid_appointment_id", "Invalid appointment_id")
			return
		}
		n.AppointmentID = &aid
	}

	if err := h.uc.CreateProgressNote(r.Context(), n); err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "create_note_failed", err.Error())
		return
	}

	actorID, _ := identitydomain.UserIDFromContext(r.Context())
	entry := audit.NewEntry(
		audit.WithTenantID(tenantID),
		audit.WithActor(actorID, audit.ScopeTenant),
		audit.WithEntity("progress_note", n.ID),
		audit.WithAction("progress_note_created"),
		audit.WithIPAddress(r.RemoteAddr),
	)
	_ = h.auditSvc.Write(r.Context(), h.pool, entry)

	server.RenderJSON(w, http.StatusCreated, toNoteResponse(n))
}

// ListProgressNotes handles GET /patients/{patient_id}/notes
func (h *Handler) ListProgressNotes(w http.ResponseWriter, r *http.Request) {
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

	notes, err := h.uc.ListProgressNotesByPatient(r.Context(), tenantID, patientID)
	if err != nil {
		server.RenderError(w, r, http.StatusInternalServerError, "list_notes_failed", "Failed to list notes")
		return
	}

	resp := make([]noteResponse, len(notes))
	for i := range notes {
		resp[i] = toNoteResponse(&notes[i])
	}
	server.RenderJSON(w, http.StatusOK, resp)
}

// --- Attachments ---

type attachmentRequest struct {
	ProfessionalID string `json:"professional_id"`
	FileName       string `json:"file_name"`
	FileType       string `json:"file_type"`
	FileSizeBytes  int64  `json:"file_size_bytes"`
	StorageKey     string `json:"storage_key"`
	Category       string `json:"category"`
	Description    string `json:"description"`
}

type attachmentResponse struct {
	ID             string `json:"id"`
	PatientID      string `json:"patient_id"`
	ProfessionalID string `json:"professional_id"`
	FileName       string `json:"file_name"`
	FileType       string `json:"file_type"`
	FileSizeBytes  int64  `json:"file_size_bytes"`
	StorageKey     string `json:"storage_key"`
	Category       string `json:"category"`
	Description    string `json:"description,omitempty"`
	CreatedAt      string `json:"created_at"`
}

func toAttachmentResponse(a *domain.ClinicalAttachment) attachmentResponse {
	return attachmentResponse{
		ID:             a.ID.String(),
		PatientID:      a.PatientID.String(),
		ProfessionalID: a.ProfessionalID.String(),
		FileName:       a.FileName,
		FileType:       a.FileType,
		FileSizeBytes:  a.FileSizeBytes,
		StorageKey:     a.StorageKey,
		Category:       string(a.Category),
		Description:    a.Description,
		CreatedAt:      a.CreatedAt.Format(time.RFC3339),
	}
}

// CreateAttachment handles POST /patients/{patient_id}/attachments
func (h *Handler) CreateAttachment(w http.ResponseWriter, r *http.Request) {
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

	var req attachmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_body", "Invalid request body")
		return
	}

	profID, err := uuid.Parse(req.ProfessionalID)
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_professional_id", "Invalid professional_id")
		return
	}

	a := &domain.ClinicalAttachment{
		TenantID:       tenantID,
		PatientID:      patientID,
		ProfessionalID: profID,
		FileName:       req.FileName,
		FileType:       req.FileType,
		FileSizeBytes:  req.FileSizeBytes,
		StorageKey:     req.StorageKey,
		Category:       domain.AttachmentCategory(req.Category),
		Description:    req.Description,
	}

	if err := h.uc.CreateAttachment(r.Context(), a); err != nil {
		server.RenderError(w, r, http.StatusInternalServerError, "create_attachment_failed", "Failed to create attachment")
		return
	}

	actorID, _ := identitydomain.UserIDFromContext(r.Context())
	entry := audit.NewEntry(
		audit.WithTenantID(tenantID),
		audit.WithActor(actorID, audit.ScopeTenant),
		audit.WithEntity("clinical_attachment", a.ID),
		audit.WithAction("attachment_created"),
		audit.WithIPAddress(r.RemoteAddr),
	)
	_ = h.auditSvc.Write(r.Context(), h.pool, entry)

	server.RenderJSON(w, http.StatusCreated, toAttachmentResponse(a))
}

// ListAttachments handles GET /patients/{patient_id}/attachments
func (h *Handler) ListAttachments(w http.ResponseWriter, r *http.Request) {
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

	attachments, err := h.uc.ListAttachmentsByPatient(r.Context(), tenantID, patientID)
	if err != nil {
		server.RenderError(w, r, http.StatusInternalServerError, "list_attachments_failed", "Failed to list attachments")
		return
	}

	resp := make([]attachmentResponse, len(attachments))
	for i := range attachments {
		resp[i] = toAttachmentResponse(&attachments[i])
	}
	server.RenderJSON(w, http.StatusOK, resp)
}
