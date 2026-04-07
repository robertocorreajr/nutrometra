package backoffice

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	billingdomain "nutrometra/api/internal/billing/domain"
	billingrepo "nutrometra/api/internal/billing/repository"
	"nutrometra/api/internal/billing/usecase"
	identitydomain "nutrometra/api/internal/identity/domain"
	"nutrometra/api/internal/platform/audit"
	"nutrometra/api/internal/platform/server"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// Handler exposes backoffice HTTP endpoints.
type Handler struct {
	repo        *Repository
	subUC       *usecase.SubscriptionUsecase
	overrideUC  *usecase.OverrideUsecase
	invoiceRepo *billingrepo.InvoiceRepository
	paymentRepo *billingrepo.PaymentRepository
}

// NewHandler creates a backoffice handler.
func NewHandler(
	repo *Repository,
	subUC *usecase.SubscriptionUsecase,
	overrideUC *usecase.OverrideUsecase,
	invoiceRepo *billingrepo.InvoiceRepository,
	paymentRepo *billingrepo.PaymentRepository,
) *Handler {
	return &Handler{
		repo:        repo,
		subUC:       subUC,
		overrideUC:  overrideUC,
		invoiceRepo: invoiceRepo,
		paymentRepo: paymentRepo,
	}
}

// ListTenants returns a paginated list of tenants.
// GET /backoffice/tenants
func (h *Handler) ListTenants(w http.ResponseWriter, r *http.Request) {
	limit, offset := parsePagination(r)
	tenants, err := h.repo.ListTenants(r.Context(), limit, offset)
	if err != nil {
		server.RenderError(w, r, http.StatusInternalServerError, "list_tenants_failed", "Failed to list tenants")
		return
	}
	server.RenderJSON(w, http.StatusOK, tenants)
}

// GetTenantDetail returns detailed tenant information.
// GET /backoffice/tenants/{id}
func (h *Handler) GetTenantDetail(w http.ResponseWriter, r *http.Request) {
	tenantID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_id", "Invalid tenant ID")
		return
	}

	detail, err := h.repo.GetTenantDetail(r.Context(), tenantID)
	if err != nil {
		server.RenderError(w, r, http.StatusInternalServerError, "get_tenant_failed", "Failed to get tenant detail")
		return
	}
	server.RenderJSON(w, http.StatusOK, detail)
}

type boChangePlanRequest struct {
	PlanID string `json:"plan_id"`
	Reason string `json:"reason"`
}

// ChangePlan forces a plan change from backoffice.
// PATCH /backoffice/tenants/{id}/subscription/plan
func (h *Handler) ChangePlan(w http.ResponseWriter, r *http.Request) {
	tenantID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_id", "Invalid tenant ID")
		return
	}
	actorID, _ := identitydomain.UserIDFromContext(r.Context())

	var req boChangePlanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_body", "Invalid request body")
		return
	}
	if req.Reason == "" {
		server.RenderError(w, r, http.StatusBadRequest, "missing_reason", "reason is required")
		return
	}

	planID, err := uuid.Parse(req.PlanID)
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_plan_id", "plan_id must be a valid UUID")
		return
	}

	if err := h.subUC.ChangePlan(r.Context(), tenantID, planID, actorID, audit.ScopeBackoffice, req.Reason); err != nil {
		switch {
		case errors.Is(err, usecase.ErrNoSubscription):
			server.RenderError(w, r, http.StatusNotFound, "no_subscription", err.Error())
		case errors.Is(err, usecase.ErrPlanNotFound):
			server.RenderError(w, r, http.StatusNotFound, "plan_not_found", err.Error())
		case errors.Is(err, usecase.ErrSamePlan):
			server.RenderError(w, r, http.StatusConflict, "same_plan", err.Error())
		default:
			server.RenderError(w, r, http.StatusInternalServerError, "change_plan_failed", "Failed to change plan")
		}
		return
	}
	server.RenderJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

type boReasonRequest struct {
	Reason string `json:"reason"`
}

// CancelSubscription forces a subscription cancellation from backoffice.
// POST /backoffice/tenants/{id}/subscription/cancel
func (h *Handler) CancelSubscription(w http.ResponseWriter, r *http.Request) {
	tenantID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_id", "Invalid tenant ID")
		return
	}
	actorID, _ := identitydomain.UserIDFromContext(r.Context())

	var req boReasonRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_body", "Invalid request body")
		return
	}
	if req.Reason == "" {
		server.RenderError(w, r, http.StatusBadRequest, "missing_reason", "reason is required")
		return
	}

	if err := h.subUC.Cancel(r.Context(), tenantID, actorID, audit.ScopeBackoffice, req.Reason); err != nil {
		if errors.Is(err, usecase.ErrNoSubscription) {
			server.RenderError(w, r, http.StatusNotFound, "no_subscription", err.Error())
			return
		}
		server.RenderError(w, r, http.StatusInternalServerError, "cancel_failed", "Failed to cancel subscription")
		return
	}
	server.RenderJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// ReactivateSubscription reactivates a cancelled subscription from backoffice.
// POST /backoffice/tenants/{id}/subscription/reactivate
func (h *Handler) ReactivateSubscription(w http.ResponseWriter, r *http.Request) {
	tenantID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_id", "Invalid tenant ID")
		return
	}
	actorID, _ := identitydomain.UserIDFromContext(r.Context())

	var req boReasonRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_body", "Invalid request body")
		return
	}
	if req.Reason == "" {
		server.RenderError(w, r, http.StatusBadRequest, "missing_reason", "reason is required")
		return
	}

	if err := h.subUC.Reactivate(r.Context(), tenantID, actorID, audit.ScopeBackoffice, req.Reason); err != nil {
		if errors.Is(err, usecase.ErrNoSubscription) {
			server.RenderError(w, r, http.StatusNotFound, "no_subscription", err.Error())
			return
		}
		server.RenderError(w, r, http.StatusInternalServerError, "reactivate_failed", "Failed to reactivate subscription")
		return
	}
	server.RenderJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// GetSubscription returns subscription details for a tenant.
// GET /backoffice/tenants/{id}/subscription
func (h *Handler) GetSubscription(w http.ResponseWriter, r *http.Request) {
	tenantID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_id", "Invalid tenant ID")
		return
	}

	detail, err := h.repo.GetTenantDetail(r.Context(), tenantID)
	if err != nil {
		server.RenderError(w, r, http.StatusInternalServerError, "get_subscription_failed", "Failed to get subscription")
		return
	}
	server.RenderJSON(w, http.StatusOK, detail)
}

// ListInvoices returns invoices for a specific tenant.
// GET /backoffice/tenants/{id}/invoices
func (h *Handler) ListInvoices(w http.ResponseWriter, r *http.Request) {
	tenantID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_id", "Invalid tenant ID")
		return
	}

	limit, offset := parsePagination(r)
	invoices, err := h.invoiceRepo.GetByTenantID(r.Context(), tenantID, limit, offset)
	if err != nil {
		server.RenderError(w, r, http.StatusInternalServerError, "list_invoices_failed", "Failed to list invoices")
		return
	}
	server.RenderJSON(w, http.StatusOK, invoices)
}

// ListPayments returns payments for a specific tenant.
// GET /backoffice/tenants/{id}/payments
func (h *Handler) ListPayments(w http.ResponseWriter, r *http.Request) {
	tenantID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_id", "Invalid tenant ID")
		return
	}

	limit, offset := parsePagination(r)
	payments, err := h.paymentRepo.GetByTenantID(r.Context(), tenantID, limit, offset)
	if err != nil {
		server.RenderError(w, r, http.StatusInternalServerError, "list_payments_failed", "Failed to list payments")
		return
	}
	server.RenderJSON(w, http.StatusOK, payments)
}

// ListOverrides returns feature overrides for a tenant.
// GET /backoffice/tenants/{id}/overrides
func (h *Handler) ListOverrides(w http.ResponseWriter, r *http.Request) {
	tenantID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_id", "Invalid tenant ID")
		return
	}

	overrides, err := h.overrideUC.ListByTenantID(r.Context(), tenantID)
	if err != nil {
		server.RenderError(w, r, http.StatusInternalServerError, "list_overrides_failed", "Failed to list overrides")
		return
	}
	server.RenderJSON(w, http.StatusOK, overrides)
}

type createOverrideRequest struct {
	FeatureKey string  `json:"feature_key"`
	Enabled    *bool   `json:"enabled"`
	LimitValue *int64  `json:"limit_value"`
	StartsAt   *string `json:"starts_at"`
	EndsAt     *string `json:"ends_at"`
	Reason     string  `json:"reason"`
}

// CreateOverride creates a feature override for a tenant.
// POST /backoffice/tenants/{id}/overrides
func (h *Handler) CreateOverride(w http.ResponseWriter, r *http.Request) {
	tenantID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_id", "Invalid tenant ID")
		return
	}
	actorID, _ := identitydomain.UserIDFromContext(r.Context())

	var req createOverrideRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_body", "Invalid request body")
		return
	}
	if req.FeatureKey == "" {
		server.RenderError(w, r, http.StatusBadRequest, "missing_feature_key", "feature_key is required")
		return
	}

	o := billingdomain.FeatureOverride{
		TenantID:   tenantID,
		FeatureKey: req.FeatureKey,
		Enabled:    req.Enabled,
		LimitValue: req.LimitValue,
		Reason:     req.Reason,
	}
	if req.StartsAt != nil {
		t, err := time.Parse(time.RFC3339, *req.StartsAt)
		if err != nil {
			server.RenderError(w, r, http.StatusBadRequest, "invalid_starts_at", "starts_at must be RFC3339")
			return
		}
		o.StartsAt = &t
	}
	if req.EndsAt != nil {
		t, err := time.Parse(time.RFC3339, *req.EndsAt)
		if err != nil {
			server.RenderError(w, r, http.StatusBadRequest, "invalid_ends_at", "ends_at must be RFC3339")
			return
		}
		o.EndsAt = &t
	}

	if err := h.overrideUC.Create(r.Context(), o, actorID, audit.ScopeBackoffice); err != nil {
		switch {
		case errors.Is(err, usecase.ErrReasonRequired):
			server.RenderError(w, r, http.StatusBadRequest, "missing_reason", err.Error())
		case errors.Is(err, usecase.ErrOverrideExists):
			server.RenderError(w, r, http.StatusConflict, "override_exists", err.Error())
		default:
			server.RenderError(w, r, http.StatusInternalServerError, "create_override_failed", "Failed to create override")
		}
		return
	}
	server.RenderJSON(w, http.StatusCreated, map[string]string{"status": "ok"})
}

type updateOverrideRequest struct {
	Enabled    *bool   `json:"enabled"`
	LimitValue *int64  `json:"limit_value"`
	StartsAt   *string `json:"starts_at"`
	EndsAt     *string `json:"ends_at"`
	Reason     string  `json:"reason"`
}

// UpdateOverride updates a feature override.
// PUT /backoffice/tenants/{id}/overrides/{feature_key}
func (h *Handler) UpdateOverride(w http.ResponseWriter, r *http.Request) {
	tenantID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_id", "Invalid tenant ID")
		return
	}
	featureKey := chi.URLParam(r, "feature_key")
	actorID, _ := identitydomain.UserIDFromContext(r.Context())

	var req updateOverrideRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_body", "Invalid request body")
		return
	}

	o := billingdomain.FeatureOverride{
		TenantID:   tenantID,
		FeatureKey: featureKey,
		Enabled:    req.Enabled,
		LimitValue: req.LimitValue,
		Reason:     req.Reason,
	}
	if req.StartsAt != nil {
		t, err := time.Parse(time.RFC3339, *req.StartsAt)
		if err != nil {
			server.RenderError(w, r, http.StatusBadRequest, "invalid_starts_at", "starts_at must be RFC3339")
			return
		}
		o.StartsAt = &t
	}
	if req.EndsAt != nil {
		t, err := time.Parse(time.RFC3339, *req.EndsAt)
		if err != nil {
			server.RenderError(w, r, http.StatusBadRequest, "invalid_ends_at", "ends_at must be RFC3339")
			return
		}
		o.EndsAt = &t
	}

	if err := h.overrideUC.Update(r.Context(), o, actorID, audit.ScopeBackoffice); err != nil {
		switch {
		case errors.Is(err, usecase.ErrReasonRequired):
			server.RenderError(w, r, http.StatusBadRequest, "missing_reason", err.Error())
		case errors.Is(err, usecase.ErrOverrideNotFound):
			server.RenderError(w, r, http.StatusNotFound, "override_not_found", err.Error())
		default:
			server.RenderError(w, r, http.StatusInternalServerError, "update_override_failed", "Failed to update override")
		}
		return
	}
	server.RenderJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// DeleteOverride removes a feature override.
// DELETE /backoffice/tenants/{id}/overrides/{feature_key}
func (h *Handler) DeleteOverride(w http.ResponseWriter, r *http.Request) {
	tenantID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_id", "Invalid tenant ID")
		return
	}
	featureKey := chi.URLParam(r, "feature_key")
	actorID, _ := identitydomain.UserIDFromContext(r.Context())

	var req boReasonRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_body", "Invalid request body")
		return
	}

	if err := h.overrideUC.Delete(r.Context(), tenantID, featureKey, actorID, audit.ScopeBackoffice, req.Reason); err != nil {
		switch {
		case errors.Is(err, usecase.ErrReasonRequired):
			server.RenderError(w, r, http.StatusBadRequest, "missing_reason", err.Error())
		case errors.Is(err, usecase.ErrOverrideNotFound):
			server.RenderError(w, r, http.StatusNotFound, "override_not_found", err.Error())
		default:
			server.RenderError(w, r, http.StatusInternalServerError, "delete_override_failed", "Failed to delete override")
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ListAuditLogs returns audit trail for a tenant.
// GET /backoffice/tenants/{id}/audit
func (h *Handler) ListAuditLogs(w http.ResponseWriter, r *http.Request) {
	tenantID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_id", "Invalid tenant ID")
		return
	}

	limit, offset := parsePagination(r)
	entries, err := h.repo.ListAuditLogs(r.Context(), tenantID, limit, offset)
	if err != nil {
		server.RenderError(w, r, http.StatusInternalServerError, "list_audit_failed", "Failed to list audit logs")
		return
	}
	server.RenderJSON(w, http.StatusOK, entries)
}

func parsePagination(r *http.Request) (limit, offset int) {
	limit = 50
	offset = 0
	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 100 {
			limit = n
		}
	}
	if v := r.URL.Query().Get("offset"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			offset = n
		}
	}
	return
}
