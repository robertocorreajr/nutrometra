package export

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"nutrometra/api/internal/export/domain"
	"nutrometra/api/internal/export/usecase"
	identitydomain "nutrometra/api/internal/identity/domain"
	"nutrometra/api/internal/platform/audit"
	"nutrometra/api/internal/platform/server"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Handler exposes HTTP endpoints for export operations.
type Handler struct {
	uc       *usecase.Usecase
	pool     *pgxpool.Pool
	auditSvc *audit.Service
}

// NewHandler creates an export Handler.
func NewHandler(uc *usecase.Usecase, pool *pgxpool.Pool, auditSvc *audit.Service) *Handler {
	return &Handler{uc: uc, pool: pool, auditSvc: auditSvc}
}

// --- Request / Response types ---

type requestExportRequest struct {
	EntityType string `json:"entity_type"`
	EntityID   string `json:"entity_id"`
}

type exportResponse struct {
	ID                string  `json:"id"`
	TenantID          string  `json:"tenant_id"`
	RelatedEntityType string  `json:"related_entity_type"`
	RelatedEntityID   string  `json:"related_entity_id"`
	ExportType        string  `json:"export_type"`
	FileKey           string  `json:"file_key,omitempty"`
	Status            string  `json:"status"`
	RequestedByUserID string  `json:"requested_by_user_id"`
	CreatedAt         string  `json:"created_at"`
	CompletedAt       *string `json:"completed_at,omitempty"`
	FailureReason     string  `json:"failure_reason,omitempty"`
}

// --- Converter ---

func toExportResponse(ef *domain.ExportedFile) exportResponse {
	resp := exportResponse{
		ID:                ef.ID.String(),
		TenantID:          ef.TenantID.String(),
		RelatedEntityType: ef.RelatedEntityType,
		RelatedEntityID:   ef.RelatedEntityID.String(),
		ExportType:        string(ef.ExportType),
		FileKey:           ef.FileKey,
		Status:            string(ef.Status),
		RequestedByUserID: ef.RequestedByUserID.String(),
		CreatedAt:         ef.CreatedAt.Format(time.RFC3339),
		FailureReason:     ef.FailureReason,
	}
	if ef.CompletedAt != nil {
		s := ef.CompletedAt.Format(time.RFC3339)
		resp.CompletedAt = &s
	}
	return resp
}

// --- Handlers ---

// RequestExport handles POST /exports
func (h *Handler) RequestExport(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := identitydomain.TenantIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusBadRequest, "missing_tenant", "No tenant in context")
		return
	}

	actorID, ok := identitydomain.UserIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusBadRequest, "missing_user", "No user in context")
		return
	}

	var req requestExportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_body", "Invalid request body")
		return
	}

	if req.EntityType == "" {
		server.RenderError(w, r, http.StatusBadRequest, "missing_entity_type", "entity_type is required")
		return
	}

	entityID, err := uuid.Parse(req.EntityID)
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_entity_id", "Invalid entity_id")
		return
	}

	ef, err := h.uc.RequestExport(r.Context(), tenantID, actorID, req.EntityType, entityID)
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "request_export_failed", err.Error())
		return
	}

	entry := audit.NewEntry(
		audit.WithTenantID(tenantID),
		audit.WithActor(actorID, audit.ScopeTenant),
		audit.WithEntity("export", ef.ID),
		audit.WithAction("export_requested"),
		audit.WithIPAddress(r.RemoteAddr),
	)
	_ = h.auditSvc.Write(r.Context(), h.pool, entry)

	server.RenderJSON(w, http.StatusAccepted, toExportResponse(ef))
}

// List handles GET /exports
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := identitydomain.TenantIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusBadRequest, "missing_tenant", "No tenant in context")
		return
	}

	exports, err := h.uc.ListByTenant(r.Context(), tenantID)
	if err != nil {
		server.RenderError(w, r, http.StatusInternalServerError, "list_exports_failed", "Failed to list exports")
		return
	}

	resp := make([]exportResponse, len(exports))
	for i := range exports {
		resp[i] = toExportResponse(&exports[i])
	}
	server.RenderJSON(w, http.StatusOK, resp)
}

// GetByID handles GET /exports/{id}
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := identitydomain.TenantIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusBadRequest, "missing_tenant", "No tenant in context")
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_id", "Invalid export ID")
		return
	}

	ef, err := h.uc.GetByID(r.Context(), tenantID, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			server.RenderError(w, r, http.StatusNotFound, "not_found", "Export not found")
			return
		}
		server.RenderError(w, r, http.StatusInternalServerError, "get_export_failed", "Failed to get export")
		return
	}
	server.RenderJSON(w, http.StatusOK, toExportResponse(ef))
}

// Download handles GET /exports/{id}/download
func (h *Handler) Download(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := identitydomain.TenantIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusBadRequest, "missing_tenant", "No tenant in context")
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_id", "Invalid export ID")
		return
	}

	ef, err := h.uc.GetByID(r.Context(), tenantID, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			server.RenderError(w, r, http.StatusNotFound, "not_found", "Export not found")
			return
		}
		server.RenderError(w, r, http.StatusInternalServerError, "get_export_failed", "Failed to get export")
		return
	}

	if ef.Status != domain.StatusCompleted {
		server.RenderError(w, r, http.StatusConflict, "export_not_ready",
			fmt.Sprintf("Export is %s, not ready for download", string(ef.Status)))
		return
	}

	if ef.FileKey == "" {
		server.RenderError(w, r, http.StatusNotFound, "file_not_found", "Export file not found")
		return
	}

	data, err := os.ReadFile(ef.FileKey)
	if err != nil {
		server.RenderError(w, r, http.StatusInternalServerError, "read_file_failed", "Failed to read export file")
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", `attachment; filename="export.pdf"`)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}
