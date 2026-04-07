package billing

import (
	"net/http"
	"time"

	"nutrometra/api/internal/billing/domain"
	"nutrometra/api/internal/billing/repository"
	"nutrometra/api/internal/billing/usecase"
	identitydomain "nutrometra/api/internal/identity/domain"
	"nutrometra/api/internal/platform/audit"
	"nutrometra/api/internal/platform/server"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Handler exposes HTTP endpoints for billing operations.
type Handler struct {
	repo              *repository.Repository
	entitlement       *usecase.EntitlementService
	cachedEntitlement *usecase.CachedEntitlementService
	pool              *pgxpool.Pool
	auditSvc          *audit.Service
}

// NewHandler creates a billing Handler.
func NewHandler(repo *repository.Repository, entitlement *usecase.EntitlementService, cachedEntitlement *usecase.CachedEntitlementService, pool *pgxpool.Pool, auditSvc *audit.Service) *Handler {
	return &Handler{repo: repo, entitlement: entitlement, cachedEntitlement: cachedEntitlement, pool: pool, auditSvc: auditSvc}
}

type planResponse struct {
	ID           string `json:"id"`
	Code         string `json:"code"`
	Name         string `json:"name"`
	BillingCycle string `json:"billing_cycle"`
	Currency     string `json:"currency"`
	PriceCents   int64  `json:"price_cents"`
}

// ListPlans returns all active subscription plans. Public endpoint.
// GET /plans
func (h *Handler) ListPlans(w http.ResponseWriter, r *http.Request) {
	plans, err := h.repo.ListActivePlans(r.Context())
	if err != nil {
		server.RenderError(w, r, http.StatusInternalServerError, "list_plans_failed", "Failed to list plans")
		return
	}
	resp := make([]planResponse, len(plans))
	for i, p := range plans {
		resp[i] = planResponse{
			ID:           p.ID.String(),
			Code:         p.Code,
			Name:         p.Name,
			BillingCycle: p.BillingCycle,
			Currency:     p.Currency,
			PriceCents:   p.PriceCents,
		}
	}
	server.RenderJSON(w, http.StatusOK, resp)
}

type subscriptionResponse struct {
	ID          string  `json:"id"`
	TenantID    string  `json:"tenant_id"`
	PlanID      string  `json:"plan_id"`
	Status      string  `json:"status"`
	StartedAt   string  `json:"started_at"`
	TrialEndsAt *string `json:"trial_ends_at,omitempty"`
	RenewsAt    *string `json:"renews_at,omitempty"`
}

func toSubResponse(s *domain.Subscription) subscriptionResponse {
	r := subscriptionResponse{
		ID:        s.ID.String(),
		TenantID:  s.TenantID.String(),
		PlanID:    s.PlanID.String(),
		Status:    string(s.Status),
		StartedAt: s.StartedAt.Format(time.RFC3339),
	}
	if s.TrialEndsAt != nil {
		v := s.TrialEndsAt.Format(time.RFC3339)
		r.TrialEndsAt = &v
	}
	if s.RenewsAt != nil {
		v := s.RenewsAt.Format(time.RFC3339)
		r.RenewsAt = &v
	}
	return r
}

// GetSubscription returns the active subscription for the current tenant.
// GET /subscription (tenant from context)
func (h *Handler) GetSubscription(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := identitydomain.TenantIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusBadRequest, "missing_tenant", "No tenant in context")
		return
	}

	sub, err := h.repo.GetActiveSubscription(r.Context(), tenantID)
	if err != nil {
		server.RenderError(w, r, http.StatusInternalServerError, "subscription_failed", "Failed to get subscription")
		return
	}
	if sub == nil {
		server.RenderError(w, r, http.StatusNotFound, "no_subscription", "No active subscription found")
		return
	}
	server.RenderJSON(w, http.StatusOK, toSubResponse(sub))
}

// ActivateTrial starts a 14-day Pro trial for the current tenant.
// POST /subscription/trial (tenant from context)
func (h *Handler) ActivateTrial(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := identitydomain.TenantIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusBadRequest, "missing_tenant", "No tenant in context")
		return
	}

	existing, _ := h.repo.GetActiveSubscription(r.Context(), tenantID)
	if existing != nil {
		server.RenderError(w, r, http.StatusConflict, "subscription_exists", "Tenant already has an active subscription")
		return
	}

	// Pro plan UUID from seed migration 000015
	proID := uuid.MustParse("20000000-0000-0000-0000-000000000003")
	now := time.Now().UTC()
	trialEnd := now.Add(14 * 24 * time.Hour)
	sub := domain.Subscription{
		ID:          uuid.New(),
		TenantID:    tenantID,
		PlanID:      proID,
		Status:      domain.SubscriptionTrialing,
		StartedAt:   now,
		TrialEndsAt: &trialEnd,
	}

	if err := h.repo.CreateSubscription(r.Context(), sub); err != nil {
		server.RenderError(w, r, http.StatusInternalServerError, "trial_failed", "Failed to activate trial")
		return
	}

	actorID, _ := identitydomain.UserIDFromContext(r.Context())
	entry := audit.NewEntry(
		audit.WithTenantID(tenantID),
		audit.WithActor(actorID, audit.ScopeTenant),
		audit.WithEntity("tenant_subscription", sub.ID),
		audit.WithAction("trial_activated"),
		audit.WithIPAddress(r.RemoteAddr),
	)
	_ = h.auditSvc.Write(r.Context(), h.pool, entry)

	if h.cachedEntitlement != nil {
		h.cachedEntitlement.InvalidateEntitlements(r.Context(), tenantID)
	}

	server.RenderJSON(w, http.StatusCreated, toSubResponse(&sub))
}

// GetEntitlements returns all entitlements for the current tenant.
// GET /entitlements (tenant from context)
func (h *Handler) GetEntitlements(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := identitydomain.TenantIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusBadRequest, "missing_tenant", "No tenant in context")
		return
	}

	var entitlements map[string]*domain.Entitlement
	var err error
	if h.cachedEntitlement != nil {
		entitlements, err = h.cachedEntitlement.GetAllEntitlements(r.Context(), tenantID)
	} else {
		entitlements, err = h.entitlement.GetAllEntitlements(r.Context(), tenantID)
	}
	if err != nil {
		server.RenderError(w, r, http.StatusInternalServerError, "entitlements_failed", "Failed to get entitlements")
		return
	}
	server.RenderJSON(w, http.StatusOK, entitlements)
}
