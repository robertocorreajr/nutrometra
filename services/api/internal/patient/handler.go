package patient

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"nutrometra/api/internal/patient/domain"
	"nutrometra/api/internal/patient/usecase"

	identitydomain "nutrometra/api/internal/identity/domain"
	"nutrometra/api/internal/platform/audit"
	"nutrometra/api/internal/platform/server"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Handler exposes HTTP endpoints for patient operations.
type Handler struct {
	uc       *usecase.Usecase
	pool     *pgxpool.Pool
	auditSvc *audit.Service
}

// NewHandler creates a patient Handler.
func NewHandler(uc *usecase.Usecase, pool *pgxpool.Pool, auditSvc *audit.Service) *Handler {
	return &Handler{uc: uc, pool: pool, auditSvc: auditSvc}
}

// --- Request/Response types ---

type createPatientRequest struct {
	ProfessionalID string  `json:"professional_id"`
	FullName       string  `json:"full_name"`
	Email          string  `json:"email"`
	Phone          string  `json:"phone"`
	CPF            string  `json:"cpf"`
	DateOfBirth    *string `json:"date_of_birth"` // YYYY-MM-DD
	Gender         string  `json:"gender"`
	Notes          string  `json:"notes"`
}

type patientResponse struct {
	ID             string  `json:"id"`
	TenantID       string  `json:"tenant_id"`
	ProfessionalID string  `json:"professional_id"`
	FullName       string  `json:"full_name"`
	Email          string  `json:"email,omitempty"`
	Phone          string  `json:"phone,omitempty"`
	CPF            string  `json:"cpf,omitempty"`
	DateOfBirth    *string `json:"date_of_birth,omitempty"`
	Gender         string  `json:"gender,omitempty"`
	Notes          string  `json:"notes,omitempty"`
	Active         bool    `json:"active"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
}

func toPatientResponse(p *domain.Patient) patientResponse {
	resp := patientResponse{
		ID:             p.ID.String(),
		TenantID:       p.TenantID.String(),
		ProfessionalID: p.ProfessionalID.String(),
		FullName:       p.FullName,
		Email:          p.Email,
		Phone:          p.Phone,
		CPF:            p.CPF,
		Gender:         p.Gender,
		Notes:          p.Notes,
		Active:         p.Active,
		CreatedAt:      p.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      p.UpdatedAt.Format(time.RFC3339),
	}
	if p.DateOfBirth != nil {
		d := p.DateOfBirth.Format("2006-01-02")
		resp.DateOfBirth = &d
	}
	return resp
}

// Create handles POST /patients
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := identitydomain.TenantIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusBadRequest, "missing_tenant", "No tenant in context")
		return
	}

	var req createPatientRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_body", "Invalid request body")
		return
	}

	profID, err := uuid.Parse(req.ProfessionalID)
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_professional_id", "Invalid professional_id")
		return
	}

	p := &domain.Patient{
		TenantID:       tenantID,
		ProfessionalID: profID,
		FullName:       req.FullName,
		Email:          req.Email,
		Phone:          req.Phone,
		CPF:            req.CPF,
		Gender:         req.Gender,
		Notes:          req.Notes,
	}

	if req.DateOfBirth != nil {
		dob, err := time.Parse("2006-01-02", *req.DateOfBirth)
		if err != nil {
			server.RenderError(w, r, http.StatusBadRequest, "invalid_date_of_birth", "date_of_birth must be YYYY-MM-DD")
			return
		}
		p.DateOfBirth = &dob
	}

	if err := h.uc.Create(r.Context(), p); err != nil {
		switch {
		case errors.Is(err, domain.ErrDuplicateCPF):
			server.RenderError(w, r, http.StatusConflict, "duplicate_cpf", "CPF already registered in this tenant")
		case errors.Is(err, domain.ErrEntitlementExceeded):
			server.RenderError(w, r, http.StatusForbidden, "entitlement_exceeded", "Patient limit exceeded for your plan")
		default:
			server.RenderError(w, r, http.StatusInternalServerError, "create_failed", "Failed to create patient")
		}
		return
	}

	actorID, _ := identitydomain.UserIDFromContext(r.Context())
	entry := audit.NewEntry(
		audit.WithTenantID(tenantID),
		audit.WithActor(actorID, audit.ScopeTenant),
		audit.WithEntity("patient", p.ID),
		audit.WithAction("patient_created"),
		audit.WithIPAddress(r.RemoteAddr),
	)
	_ = h.auditSvc.Write(r.Context(), h.pool, entry)

	server.RenderJSON(w, http.StatusCreated, toPatientResponse(p))
}

// List handles GET /patients
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := identitydomain.TenantIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusBadRequest, "missing_tenant", "No tenant in context")
		return
	}

	patients, err := h.uc.List(r.Context(), tenantID)
	if err != nil {
		server.RenderError(w, r, http.StatusInternalServerError, "list_failed", "Failed to list patients")
		return
	}

	resp := make([]patientResponse, len(patients))
	for i := range patients {
		resp[i] = toPatientResponse(&patients[i])
	}
	server.RenderJSON(w, http.StatusOK, resp)
}

// GetByID handles GET /patients/{id}
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := identitydomain.TenantIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusBadRequest, "missing_tenant", "No tenant in context")
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_id", "Invalid patient ID")
		return
	}

	p, err := h.uc.GetByID(r.Context(), tenantID, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			server.RenderError(w, r, http.StatusNotFound, "not_found", "Patient not found")
			return
		}
		server.RenderError(w, r, http.StatusInternalServerError, "get_failed", "Failed to get patient")
		return
	}
	server.RenderJSON(w, http.StatusOK, toPatientResponse(p))
}

// Update handles PUT /patients/{id}
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := identitydomain.TenantIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusBadRequest, "missing_tenant", "No tenant in context")
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_id", "Invalid patient ID")
		return
	}

	var req createPatientRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_body", "Invalid request body")
		return
	}

	profID, err := uuid.Parse(req.ProfessionalID)
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_professional_id", "Invalid professional_id")
		return
	}

	p := &domain.Patient{
		ID:             id,
		TenantID:       tenantID,
		ProfessionalID: profID,
		FullName:       req.FullName,
		Email:          req.Email,
		Phone:          req.Phone,
		CPF:            req.CPF,
		Gender:         req.Gender,
		Notes:          req.Notes,
	}

	if req.DateOfBirth != nil {
		dob, err := time.Parse("2006-01-02", *req.DateOfBirth)
		if err != nil {
			server.RenderError(w, r, http.StatusBadRequest, "invalid_date_of_birth", "date_of_birth must be YYYY-MM-DD")
			return
		}
		p.DateOfBirth = &dob
	}

	if err := h.uc.Update(r.Context(), p); err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			server.RenderError(w, r, http.StatusNotFound, "not_found", "Patient not found")
		case errors.Is(err, domain.ErrDuplicateCPF):
			server.RenderError(w, r, http.StatusConflict, "duplicate_cpf", "CPF already registered")
		default:
			server.RenderError(w, r, http.StatusInternalServerError, "update_failed", "Failed to update patient")
		}
		return
	}

	actorID, _ := identitydomain.UserIDFromContext(r.Context())
	entry := audit.NewEntry(
		audit.WithTenantID(tenantID),
		audit.WithActor(actorID, audit.ScopeTenant),
		audit.WithEntity("patient", id),
		audit.WithAction("patient_updated"),
		audit.WithIPAddress(r.RemoteAddr),
	)
	_ = h.auditSvc.Write(r.Context(), h.pool, entry)

	server.RenderJSON(w, http.StatusOK, toPatientResponse(p))
}

// --- Invite endpoints ---

type inviteResponse struct {
	InviteID  string `json:"invite_id"`
	Code      string `json:"code"`
	ExpiresAt string `json:"expires_at"`
}

// GenerateInvite handles POST /patients/{id}/invites
func (h *Handler) GenerateInvite(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := identitydomain.TenantIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusBadRequest, "missing_tenant", "No tenant in context")
		return
	}

	patientID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_id", "Invalid patient ID")
		return
	}

	actorID, _ := identitydomain.UserIDFromContext(r.Context())

	result, err := h.uc.GenerateInvite(r.Context(), tenantID, patientID, actorID)
	if err != nil {
		server.RenderError(w, r, http.StatusInternalServerError, "invite_failed", "Failed to generate invite")
		return
	}

	entry := audit.NewEntry(
		audit.WithTenantID(tenantID),
		audit.WithActor(actorID, audit.ScopeTenant),
		audit.WithEntity("patient_invite", result.InviteID),
		audit.WithAction("invite_generated"),
		audit.WithIPAddress(r.RemoteAddr),
	)
	_ = h.auditSvc.Write(r.Context(), h.pool, entry)

	server.RenderJSON(w, http.StatusCreated, inviteResponse{
		InviteID:  result.InviteID.String(),
		Code:      result.Code,
		ExpiresAt: result.ExpiresAt.Format(time.RFC3339),
	})
}

// ActivatePortalAccess handles POST /invites/activate
type activateRequest struct {
	Code string `json:"code"`
}

func (h *Handler) ActivatePortalAccess(w http.ResponseWriter, r *http.Request) {
	userID, ok := identitydomain.UserIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusUnauthorized, "unauthorized", "Authentication required")
		return
	}

	var req activateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Code == "" {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_body", "code is required")
		return
	}

	link, err := h.uc.ActivatePortalAccess(r.Context(), req.Code, userID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInviteNotFound):
			server.RenderError(w, r, http.StatusNotFound, "invite_not_found", "Invite not found")
		case errors.Is(err, domain.ErrInviteExpired):
			server.RenderError(w, r, http.StatusGone, "invite_expired", "Invite has expired")
		case errors.Is(err, domain.ErrInviteMaxUsed):
			server.RenderError(w, r, http.StatusConflict, "invite_max_used", "Invite has reached maximum uses")
		case errors.Is(err, domain.ErrAccessLinkExists):
			server.RenderError(w, r, http.StatusConflict, "already_linked", "User already linked to this patient")
		default:
			server.RenderError(w, r, http.StatusInternalServerError, "activate_failed", "Failed to activate portal access")
		}
		return
	}

	server.RenderJSON(w, http.StatusCreated, map[string]string{
		"access_link_id": link.ID.String(),
		"patient_id":     link.PatientID.String(),
		"tenant_id":      link.TenantID.String(),
	})
}

// --- Profile endpoints ---

type profileRequest struct {
	Occupation            string   `json:"occupation"`
	MaritalStatus         string   `json:"marital_status"`
	Ethnicity             string   `json:"ethnicity"`
	BloodType             string   `json:"blood_type"`
	Allergies             []string `json:"allergies"`
	ChronicConditions     []string `json:"chronic_conditions"`
	Medications           []string `json:"medications"`
	EmergencyContactName  string   `json:"emergency_contact_name"`
	EmergencyContactPhone string   `json:"emergency_contact_phone"`
}

type profileResponse struct {
	ID                    string   `json:"id"`
	PatientID             string   `json:"patient_id"`
	Occupation            string   `json:"occupation,omitempty"`
	MaritalStatus         string   `json:"marital_status,omitempty"`
	Ethnicity             string   `json:"ethnicity,omitempty"`
	BloodType             string   `json:"blood_type,omitempty"`
	Allergies             []string `json:"allergies"`
	ChronicConditions     []string `json:"chronic_conditions"`
	Medications           []string `json:"medications"`
	EmergencyContactName  string   `json:"emergency_contact_name,omitempty"`
	EmergencyContactPhone string   `json:"emergency_contact_phone,omitempty"`
	CreatedAt             string   `json:"created_at"`
	UpdatedAt             string   `json:"updated_at"`
}

func toProfileResponse(p *domain.PatientProfile) profileResponse {
	return profileResponse{
		ID:                    p.ID.String(),
		PatientID:             p.PatientID.String(),
		Occupation:            p.Occupation,
		MaritalStatus:         p.MaritalStatus,
		Ethnicity:             p.Ethnicity,
		BloodType:             p.BloodType,
		Allergies:             p.Allergies,
		ChronicConditions:     p.ChronicConditions,
		Medications:           p.Medications,
		EmergencyContactName:  p.EmergencyContactName,
		EmergencyContactPhone: p.EmergencyContactPhone,
		CreatedAt:             p.CreatedAt.Format(time.RFC3339),
		UpdatedAt:             p.UpdatedAt.Format(time.RFC3339),
	}
}

// UpsertProfile handles POST /patients/{id}/profiles
func (h *Handler) UpsertProfile(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := identitydomain.TenantIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusBadRequest, "missing_tenant", "No tenant in context")
		return
	}

	patientID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_id", "Invalid patient ID")
		return
	}

	var req profileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_body", "Invalid request body")
		return
	}

	p := &domain.PatientProfile{
		TenantID:              tenantID,
		PatientID:             patientID,
		Occupation:            req.Occupation,
		MaritalStatus:         req.MaritalStatus,
		Ethnicity:             req.Ethnicity,
		BloodType:             req.BloodType,
		Allergies:             req.Allergies,
		ChronicConditions:     req.ChronicConditions,
		Medications:           req.Medications,
		EmergencyContactName:  req.EmergencyContactName,
		EmergencyContactPhone: req.EmergencyContactPhone,
	}

	if err := h.uc.UpsertProfile(r.Context(), p); err != nil {
		server.RenderError(w, r, http.StatusInternalServerError, "profile_failed", "Failed to save profile")
		return
	}

	server.RenderJSON(w, http.StatusOK, toProfileResponse(p))
}

// GetProfile handles GET /patients/{id}/profiles
func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := identitydomain.TenantIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusBadRequest, "missing_tenant", "No tenant in context")
		return
	}

	patientID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_id", "Invalid patient ID")
		return
	}

	p, err := h.uc.GetProfile(r.Context(), tenantID, patientID)
	if err != nil {
		if errors.Is(err, domain.ErrProfileNotFound) {
			server.RenderError(w, r, http.StatusNotFound, "not_found", "Profile not found")
			return
		}
		server.RenderError(w, r, http.StatusInternalServerError, "get_profile_failed", "Failed to get profile")
		return
	}
	server.RenderJSON(w, http.StatusOK, toProfileResponse(p))
}
