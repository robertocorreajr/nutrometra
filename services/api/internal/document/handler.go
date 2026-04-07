package document

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"nutrometra/api/internal/document/domain"
	"nutrometra/api/internal/document/usecase"
	identitydomain "nutrometra/api/internal/identity/domain"
	"nutrometra/api/internal/platform/audit"
	"nutrometra/api/internal/platform/server"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Handler exposes HTTP endpoints for clinical document operations.
type Handler struct {
	uc       *usecase.Usecase
	pool     *pgxpool.Pool
	auditSvc *audit.Service
}

// NewHandler creates a document Handler.
func NewHandler(uc *usecase.Usecase, pool *pgxpool.Pool, auditSvc *audit.Service) *Handler {
	return &Handler{uc: uc, pool: pool, auditSvc: auditSvc}
}

// --- Request / Response types ---

type createDocumentRequest struct {
	ProfessionalID string           `json:"professional_id"`
	AppointmentID  *string          `json:"appointment_id"`
	DocumentType   string           `json:"document_type"`
	Title          string           `json:"title"`
	ContentJSON    json.RawMessage  `json:"content_json"`
}

type updateDocumentRequest struct {
	Title         string          `json:"title"`
	ContentJSON   json.RawMessage `json:"content_json"`
	AppointmentID *string         `json:"appointment_id"`
}

type documentResponse struct {
	ID                string          `json:"id"`
	PatientID         string          `json:"patient_id"`
	ProfessionalID    string          `json:"professional_id"`
	AppointmentID     *string         `json:"appointment_id,omitempty"`
	DocumentType      string          `json:"document_type"`
	Title             string          `json:"title"`
	Status            string          `json:"status"`
	ContentJSON       json.RawMessage `json:"content_json"`
	VersionNumber     int             `json:"version_number"`
	PreviousVersionID *string         `json:"previous_version_id,omitempty"`
	CreatedAt         string          `json:"created_at"`
	UpdatedAt         string          `json:"updated_at"`
}

// --- Converters ---

func toDocumentResponse(d *domain.ClinicalDocument) documentResponse {
	resp := documentResponse{
		ID:             d.ID.String(),
		PatientID:      d.PatientID.String(),
		ProfessionalID: d.ProfessionalID.String(),
		DocumentType:   string(d.DocumentType),
		Title:          d.Title,
		Status:         string(d.Status),
		ContentJSON:    d.ContentJSON,
		VersionNumber:  d.VersionNumber,
		CreatedAt:      d.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      d.UpdatedAt.Format(time.RFC3339),
	}
	if d.AppointmentID != nil {
		s := d.AppointmentID.String()
		resp.AppointmentID = &s
	}
	if d.PreviousVersionID != nil {
		s := d.PreviousVersionID.String()
		resp.PreviousVersionID = &s
	}
	return resp
}

func toDocumentListResponse(d *domain.ClinicalDocument) documentResponse {
	resp := documentResponse{
		ID:             d.ID.String(),
		PatientID:      d.PatientID.String(),
		ProfessionalID: d.ProfessionalID.String(),
		DocumentType:   string(d.DocumentType),
		Title:          d.Title,
		Status:         string(d.Status),
		ContentJSON:    d.ContentJSON,
		VersionNumber:  d.VersionNumber,
		CreatedAt:      d.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      d.UpdatedAt.Format(time.RFC3339),
	}
	if d.AppointmentID != nil {
		s := d.AppointmentID.String()
		resp.AppointmentID = &s
	}
	if d.PreviousVersionID != nil {
		s := d.PreviousVersionID.String()
		resp.PreviousVersionID = &s
	}
	return resp
}

// --- Handlers ---

// Create handles POST /patients/{patient_id}/documents
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
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

	var req createDocumentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_body", "Invalid request body")
		return
	}

	profID, err := uuid.Parse(req.ProfessionalID)
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_professional_id", "Invalid professional_id")
		return
	}

	d := &domain.ClinicalDocument{
		TenantID:       tenantID,
		PatientID:      patientID,
		ProfessionalID: profID,
		DocumentType:   domain.DocumentType(req.DocumentType),
		Title:          req.Title,
		ContentJSON:    req.ContentJSON,
	}

	if req.AppointmentID != nil {
		apptID, err := uuid.Parse(*req.AppointmentID)
		if err != nil {
			server.RenderError(w, r, http.StatusBadRequest, "invalid_appointment_id", "Invalid appointment_id")
			return
		}
		d.AppointmentID = &apptID
	}

	if err := h.uc.Create(r.Context(), d); err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "create_document_failed", err.Error())
		return
	}

	actorID, _ := identitydomain.UserIDFromContext(r.Context())
	entry := audit.NewEntry(
		audit.WithTenantID(tenantID),
		audit.WithActor(actorID, audit.ScopeTenant),
		audit.WithEntity("clinical_document", d.ID),
		audit.WithAction("created"),
		audit.WithIPAddress(r.RemoteAddr),
	)
	_ = h.auditSvc.Write(r.Context(), h.pool, entry)

	server.RenderJSON(w, http.StatusCreated, toDocumentResponse(d))
}

// ListByPatient handles GET /patients/{patient_id}/documents
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

	docs, err := h.uc.ListByPatient(r.Context(), tenantID, patientID)
	if err != nil {
		server.RenderError(w, r, http.StatusInternalServerError, "list_documents_failed", "Failed to list documents")
		return
	}

	resp := make([]documentResponse, len(docs))
	for i := range docs {
		resp[i] = toDocumentListResponse(&docs[i])
	}
	server.RenderJSON(w, http.StatusOK, resp)
}

// GetByID handles GET /documents/{id}
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := identitydomain.TenantIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusBadRequest, "missing_tenant", "No tenant in context")
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_id", "Invalid document ID")
		return
	}

	d, err := h.uc.GetByID(r.Context(), tenantID, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			server.RenderError(w, r, http.StatusNotFound, "not_found", "Document not found")
			return
		}
		server.RenderError(w, r, http.StatusInternalServerError, "get_document_failed", "Failed to get document")
		return
	}
	server.RenderJSON(w, http.StatusOK, toDocumentResponse(d))
}

// Update handles PUT /documents/{id}
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := identitydomain.TenantIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusBadRequest, "missing_tenant", "No tenant in context")
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_id", "Invalid document ID")
		return
	}

	var req updateDocumentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_body", "Invalid request body")
		return
	}

	d := &domain.ClinicalDocument{
		ID:          id,
		TenantID:    tenantID,
		Title:       req.Title,
		ContentJSON: req.ContentJSON,
	}

	if req.AppointmentID != nil {
		apptID, err := uuid.Parse(*req.AppointmentID)
		if err != nil {
			server.RenderError(w, r, http.StatusBadRequest, "invalid_appointment_id", "Invalid appointment_id")
			return
		}
		d.AppointmentID = &apptID
	}

	if err := h.uc.Update(r.Context(), d); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			server.RenderError(w, r, http.StatusNotFound, "not_found", "Document not found")
			return
		}
		if errors.Is(err, domain.ErrNotDraft) {
			server.RenderError(w, r, http.StatusConflict, "not_draft", "Only draft documents can be updated")
			return
		}
		server.RenderError(w, r, http.StatusBadRequest, "update_document_failed", err.Error())
		return
	}

	actorID, _ := identitydomain.UserIDFromContext(r.Context())
	entry := audit.NewEntry(
		audit.WithTenantID(tenantID),
		audit.WithActor(actorID, audit.ScopeTenant),
		audit.WithEntity("clinical_document", id),
		audit.WithAction("updated"),
		audit.WithIPAddress(r.RemoteAddr),
	)
	_ = h.auditSvc.Write(r.Context(), h.pool, entry)

	server.RenderJSON(w, http.StatusOK, toDocumentResponse(d))
}

// Finalize handles POST /documents/{id}/finalize
func (h *Handler) Finalize(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := identitydomain.TenantIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusBadRequest, "missing_tenant", "No tenant in context")
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_id", "Invalid document ID")
		return
	}

	if err := h.uc.Finalize(r.Context(), tenantID, id); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			server.RenderError(w, r, http.StatusNotFound, "not_found", "Document not found")
			return
		}
		if errors.Is(err, domain.ErrAlreadyFinalized) {
			server.RenderError(w, r, http.StatusConflict, "already_finalized", "Document is already finalized")
			return
		}
		if errors.Is(err, domain.ErrNotDraft) {
			server.RenderError(w, r, http.StatusConflict, "not_draft", "Only draft documents can be finalized")
			return
		}
		server.RenderError(w, r, http.StatusInternalServerError, "finalize_failed", "Failed to finalize document")
		return
	}

	actorID, _ := identitydomain.UserIDFromContext(r.Context())
	entry := audit.NewEntry(
		audit.WithTenantID(tenantID),
		audit.WithActor(actorID, audit.ScopeTenant),
		audit.WithEntity("clinical_document", id),
		audit.WithAction("finalized"),
		audit.WithIPAddress(r.RemoteAddr),
	)
	_ = h.auditSvc.Write(r.Context(), h.pool, entry)

	w.WriteHeader(http.StatusNoContent)
}

// Publish handles POST /documents/{id}/publish
func (h *Handler) Publish(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := identitydomain.TenantIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusBadRequest, "missing_tenant", "No tenant in context")
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_id", "Invalid document ID")
		return
	}

	actorID, _ := identitydomain.UserIDFromContext(r.Context())

	if err := h.uc.Publish(r.Context(), tenantID, id, actorID); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			server.RenderError(w, r, http.StatusNotFound, "not_found", "Document not found")
			return
		}
		if errors.Is(err, domain.ErrNotFinalized) {
			server.RenderError(w, r, http.StatusConflict, "not_finalized", "Only finalized documents can be published")
			return
		}
		server.RenderError(w, r, http.StatusInternalServerError, "publish_failed", "Failed to publish document")
		return
	}

	entry := audit.NewEntry(
		audit.WithTenantID(tenantID),
		audit.WithActor(actorID, audit.ScopeTenant),
		audit.WithEntity("clinical_document", id),
		audit.WithAction("published"),
		audit.WithIPAddress(r.RemoteAddr),
	)
	_ = h.auditSvc.Write(r.Context(), h.pool, entry)

	w.WriteHeader(http.StatusNoContent)
}

// NewVersion handles POST /documents/{id}/new-version
func (h *Handler) NewVersion(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := identitydomain.TenantIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusBadRequest, "missing_tenant", "No tenant in context")
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_id", "Invalid document ID")
		return
	}

	newDoc, err := h.uc.CreateNewVersion(r.Context(), tenantID, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			server.RenderError(w, r, http.StatusNotFound, "not_found", "Document not found")
			return
		}
		server.RenderError(w, r, http.StatusInternalServerError, "new_version_failed", "Failed to create new version")
		return
	}

	actorID, _ := identitydomain.UserIDFromContext(r.Context())
	entry := audit.NewEntry(
		audit.WithTenantID(tenantID),
		audit.WithActor(actorID, audit.ScopeTenant),
		audit.WithEntity("clinical_document", newDoc.ID),
		audit.WithAction("version_created"),
		audit.WithIPAddress(r.RemoteAddr),
	)
	_ = h.auditSvc.Write(r.Context(), h.pool, entry)

	server.RenderJSON(w, http.StatusCreated, toDocumentResponse(newDoc))
}

// ListVersions handles GET /documents/{id}/versions
func (h *Handler) ListVersions(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := identitydomain.TenantIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusBadRequest, "missing_tenant", "No tenant in context")
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_id", "Invalid document ID")
		return
	}

	docs, err := h.uc.ListVersions(r.Context(), tenantID, id)
	if err != nil {
		server.RenderError(w, r, http.StatusInternalServerError, "list_versions_failed", "Failed to list versions")
		return
	}

	resp := make([]documentResponse, len(docs))
	for i := range docs {
		resp[i] = toDocumentListResponse(&docs[i])
	}
	server.RenderJSON(w, http.StatusOK, resp)
}
