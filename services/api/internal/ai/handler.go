package ai

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	aidomain "nutrometra/api/internal/ai/domain"
	"nutrometra/api/internal/ai/usecase"
	identitydomain "nutrometra/api/internal/identity/domain"
	"nutrometra/api/internal/platform/audit"
	"nutrometra/api/internal/platform/server"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Handler exposes HTTP endpoints for AI suggestion management.
type Handler struct {
	service  *usecase.SuggestionService
	pool     *pgxpool.Pool
	auditSvc *audit.Service
}

// NewHandler creates a new AI suggestions Handler.
func NewHandler(service *usecase.SuggestionService, pool *pgxpool.Pool, auditSvc *audit.Service) *Handler {
	return &Handler{service: service, pool: pool, auditSvc: auditSvc}
}

type createRequest struct {
	SuggestionType string                 `json:"suggestion_type"`
	PatientID      *uuid.UUID             `json:"patient_id,omitempty"`
	ExtraContext   map[string]interface{} `json:"extra_context,omitempty"`
}

type suggestionResponse struct {
	ID             uuid.UUID               `json:"id"`
	TenantID       uuid.UUID               `json:"tenant_id"`
	UserID         uuid.UUID               `json:"user_id"`
	SuggestionType aidomain.SuggestionType `json:"suggestion_type"`
	Status         aidomain.SuggestionStatus `json:"status"`
	ResponseText   *string                 `json:"response_text,omitempty"`
	ModelID        *string                 `json:"model_id,omitempty"`
	InputTokens    *int                    `json:"input_tokens,omitempty"`
	OutputTokens   *int                    `json:"output_tokens,omitempty"`
	CreatedAt      string                  `json:"created_at"`
	CompletedAt    *string                 `json:"completed_at,omitempty"`
	ReviewedAt     *string                 `json:"reviewed_at,omitempty"`
	ReviewedBy     *uuid.UUID              `json:"reviewed_by,omitempty"`
}

func toResponse(s *aidomain.AISuggestion) suggestionResponse {
	r := suggestionResponse{
		ID:             s.ID,
		TenantID:       s.TenantID,
		UserID:         s.UserID,
		SuggestionType: s.SuggestionType,
		Status:         s.Status,
		ResponseText:   s.ResponseText,
		ModelID:        s.ModelID,
		InputTokens:    s.InputTokens,
		OutputTokens:   s.OutputTokens,
		CreatedAt:      s.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		ReviewedBy:     s.ReviewedBy,
	}
	if s.CompletedAt != nil {
		t := s.CompletedAt.Format("2006-01-02T15:04:05Z07:00")
		r.CompletedAt = &t
	}
	if s.ReviewedAt != nil {
		t := s.ReviewedAt.Format("2006-01-02T15:04:05Z07:00")
		r.ReviewedAt = &t
	}
	return r
}

// Create handles POST /ai/suggestions
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := identitydomain.TenantIDFromContext(ctx)
	if !ok {
		server.RenderError(w, r, http.StatusUnauthorized, "missing_tenant", "Tenant ID required")
		return
	}
	userID, ok := identitydomain.UserIDFromContext(ctx)
	if !ok {
		server.RenderError(w, r, http.StatusUnauthorized, "missing_user", "User ID required")
		return
	}

	var req createRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_body", "Invalid request body")
		return
	}

	if !aidomain.ValidSuggestionType(req.SuggestionType) {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_type", "Invalid suggestion type")
		return
	}

	sug, err := h.service.Create(ctx, usecase.CreateSuggestionInput{
		TenantID:       tenantID,
		UserID:         userID,
		SuggestionType: req.SuggestionType,
		PatientID:      req.PatientID,
		ExtraContext:   req.ExtraContext,
	})
	if err != nil {
		slog.Error("ai: create suggestion failed", "error", err)
		server.RenderError(w, r, http.StatusInternalServerError, "create_failed", "Failed to create suggestion")
		return
	}

	server.RenderJSON(w, http.StatusCreated, toResponse(sug))
}

// List handles GET /ai/suggestions
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := identitydomain.TenantIDFromContext(ctx)
	if !ok {
		server.RenderError(w, r, http.StatusUnauthorized, "missing_tenant", "Tenant ID required")
		return
	}
	userID, ok := identitydomain.UserIDFromContext(ctx)
	if !ok {
		server.RenderError(w, r, http.StatusUnauthorized, "missing_user", "User ID required")
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	suggestions, err := h.service.List(ctx, tenantID, userID, limit, offset)
	if err != nil {
		slog.Error("ai: list suggestions failed", "error", err)
		server.RenderError(w, r, http.StatusInternalServerError, "list_failed", "Failed to list suggestions")
		return
	}

	result := make([]suggestionResponse, 0, len(suggestions))
	for i := range suggestions {
		result = append(result, toResponse(&suggestions[i]))
	}
	server.RenderJSON(w, http.StatusOK, result)
}

// GetByID handles GET /ai/suggestions/{id}
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := identitydomain.TenantIDFromContext(ctx)
	if !ok {
		server.RenderError(w, r, http.StatusUnauthorized, "missing_tenant", "Tenant ID required")
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_id", "Invalid suggestion ID")
		return
	}

	sug, err := h.service.GetByID(ctx, tenantID, id)
	if err != nil {
		slog.Error("ai: get suggestion failed", "error", err, "id", id)
		server.RenderError(w, r, http.StatusNotFound, "not_found", "Suggestion not found")
		return
	}

	server.RenderJSON(w, http.StatusOK, toResponse(sug))
}

// Accept handles POST /ai/suggestions/{id}/accept
func (h *Handler) Accept(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := identitydomain.TenantIDFromContext(ctx)
	if !ok {
		server.RenderError(w, r, http.StatusUnauthorized, "missing_tenant", "Tenant ID required")
		return
	}
	userID, ok := identitydomain.UserIDFromContext(ctx)
	if !ok {
		server.RenderError(w, r, http.StatusUnauthorized, "missing_user", "User ID required")
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_id", "Invalid suggestion ID")
		return
	}

	if err := h.service.Accept(ctx, tenantID, id, userID); err != nil {
		slog.Error("ai: accept suggestion failed", "error", err, "id", id)
		server.RenderError(w, r, http.StatusBadRequest, "accept_failed", err.Error())
		return
	}

	// Audit
	entry := audit.NewEntry(
		audit.WithTenantID(tenantID),
		audit.WithActor(userID, audit.ScopeTenant),
		audit.WithEntity("ai_suggestion", id),
		audit.WithAction("ai_suggestion_accepted"),
		audit.WithIPAddress(r.RemoteAddr),
	)
	if err := h.auditSvc.Write(ctx, h.pool, entry); err != nil {
		slog.Error("ai: audit write failed", "error", err)
	}

	server.RenderJSON(w, http.StatusOK, map[string]string{"status": "accepted"})
}

// Reject handles POST /ai/suggestions/{id}/reject
func (h *Handler) Reject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := identitydomain.TenantIDFromContext(ctx)
	if !ok {
		server.RenderError(w, r, http.StatusUnauthorized, "missing_tenant", "Tenant ID required")
		return
	}
	userID, ok := identitydomain.UserIDFromContext(ctx)
	if !ok {
		server.RenderError(w, r, http.StatusUnauthorized, "missing_user", "User ID required")
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_id", "Invalid suggestion ID")
		return
	}

	if err := h.service.Reject(ctx, tenantID, id, userID); err != nil {
		slog.Error("ai: reject suggestion failed", "error", err, "id", id)
		server.RenderError(w, r, http.StatusBadRequest, "reject_failed", err.Error())
		return
	}

	// Audit
	entry := audit.NewEntry(
		audit.WithTenantID(tenantID),
		audit.WithActor(userID, audit.ScopeTenant),
		audit.WithEntity("ai_suggestion", id),
		audit.WithAction("ai_suggestion_rejected"),
		audit.WithIPAddress(r.RemoteAddr),
	)
	if err := h.auditSvc.Write(ctx, h.pool, entry); err != nil {
		slog.Error("ai: audit write failed", "error", err)
	}

	server.RenderJSON(w, http.StatusOK, map[string]string{"status": "rejected"})
}
