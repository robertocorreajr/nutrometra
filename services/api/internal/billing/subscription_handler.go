package billing

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"nutrometra/api/internal/billing/repository"
	"nutrometra/api/internal/billing/usecase"
	identitydomain "nutrometra/api/internal/identity/domain"
	"nutrometra/api/internal/platform/audit"
	"nutrometra/api/internal/platform/server"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// SubscriptionHandler handles subscription self-service endpoints.
type SubscriptionHandler struct {
	subUC       *usecase.SubscriptionUsecase
	invoiceRepo *repository.InvoiceRepository
	paymentRepo *repository.PaymentRepository
}

// NewSubscriptionHandler creates a subscription handler.
func NewSubscriptionHandler(
	subUC *usecase.SubscriptionUsecase,
	invoiceRepo *repository.InvoiceRepository,
	paymentRepo *repository.PaymentRepository,
) *SubscriptionHandler {
	return &SubscriptionHandler{
		subUC:       subUC,
		invoiceRepo: invoiceRepo,
		paymentRepo: paymentRepo,
	}
}

type checkoutRequest struct {
	PlanID string `json:"plan_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
}

// Checkout creates a new subscription via Stripe.
// POST /subscription/checkout
func (h *SubscriptionHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := identitydomain.TenantIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusBadRequest, "missing_tenant", "No tenant in context")
		return
	}
	actorID, _ := identitydomain.UserIDFromContext(r.Context())

	var req checkoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_body", "Invalid request body")
		return
	}

	planID, err := uuid.Parse(req.PlanID)
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_plan_id", "plan_id must be a valid UUID")
		return
	}
	if req.Email == "" {
		server.RenderError(w, r, http.StatusBadRequest, "missing_email", "email is required")
		return
	}

	if err := h.subUC.Checkout(r.Context(), tenantID, planID, req.Email, req.Name, actorID, audit.ScopeTenant); err != nil {
		switch {
		case errors.Is(err, usecase.ErrSubscriptionExists):
			server.RenderError(w, r, http.StatusConflict, "subscription_exists", err.Error())
		case errors.Is(err, usecase.ErrPlanNotFound):
			server.RenderError(w, r, http.StatusNotFound, "plan_not_found", err.Error())
		case errors.Is(err, usecase.ErrPlanHasNoProviderID):
			server.RenderError(w, r, http.StatusUnprocessableEntity, "plan_not_configured", err.Error())
		case errors.Is(err, usecase.ErrNoProviderConfigured):
			server.RenderError(w, r, http.StatusServiceUnavailable, "provider_unavailable", "Billing provider not configured")
		default:
			server.RenderError(w, r, http.StatusInternalServerError, "checkout_failed", "Failed to create subscription")
		}
		return
	}

	server.RenderJSON(w, http.StatusCreated, map[string]string{"status": "ok"})
}

type changePlanRequest struct {
	PlanID string `json:"plan_id"`
}

// ChangePlan upgrades or downgrades the subscription.
// PATCH /subscription/plan
func (h *SubscriptionHandler) ChangePlan(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := identitydomain.TenantIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusBadRequest, "missing_tenant", "No tenant in context")
		return
	}
	actorID, _ := identitydomain.UserIDFromContext(r.Context())

	var req changePlanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_body", "Invalid request body")
		return
	}

	planID, err := uuid.Parse(req.PlanID)
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_plan_id", "plan_id must be a valid UUID")
		return
	}

	if err := h.subUC.ChangePlan(r.Context(), tenantID, planID, actorID, audit.ScopeTenant, "self-service plan change"); err != nil {
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

// CancelSubscription cancels the current subscription.
// POST /subscription/cancel
func (h *SubscriptionHandler) CancelSubscription(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := identitydomain.TenantIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusBadRequest, "missing_tenant", "No tenant in context")
		return
	}
	actorID, _ := identitydomain.UserIDFromContext(r.Context())

	if err := h.subUC.Cancel(r.Context(), tenantID, actorID, audit.ScopeTenant, "self-service cancellation"); err != nil {
		if errors.Is(err, usecase.ErrNoSubscription) {
			server.RenderError(w, r, http.StatusNotFound, "no_subscription", err.Error())
			return
		}
		server.RenderError(w, r, http.StatusInternalServerError, "cancel_failed", "Failed to cancel subscription")
		return
	}

	server.RenderJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// ListInvoices returns invoices for the current tenant.
// GET /invoices
func (h *SubscriptionHandler) ListInvoices(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := identitydomain.TenantIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusBadRequest, "missing_tenant", "No tenant in context")
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

// GetInvoice returns a single invoice.
// GET /invoices/{id}
func (h *SubscriptionHandler) GetInvoice(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := identitydomain.TenantIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusBadRequest, "missing_tenant", "No tenant in context")
		return
	}

	invoiceID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_id", "Invalid invoice ID")
		return
	}

	invoice, err := h.invoiceRepo.GetByID(r.Context(), tenantID, invoiceID)
	if err != nil {
		server.RenderError(w, r, http.StatusInternalServerError, "get_invoice_failed", "Failed to get invoice")
		return
	}
	if invoice == nil {
		server.RenderError(w, r, http.StatusNotFound, "not_found", "Invoice not found")
		return
	}
	server.RenderJSON(w, http.StatusOK, invoice)
}

// ListPayments returns payments for the current tenant.
// GET /payments
func (h *SubscriptionHandler) ListPayments(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := identitydomain.TenantIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusBadRequest, "missing_tenant", "No tenant in context")
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
