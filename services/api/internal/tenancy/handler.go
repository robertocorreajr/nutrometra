package tenancy

import (
	"errors"
	"net/http"

	identitydomain "nutrometra/api/internal/identity/domain"
	"nutrometra/api/internal/platform/server"
	"nutrometra/api/internal/tenancy/domain"
	"nutrometra/api/internal/tenancy/usecase"
)

// Handler exposes HTTP endpoints for tenant operations.
type Handler struct {
	uc usecase.TenantUsecase
}

// NewHandler creates a tenant Handler.
func NewHandler(uc usecase.TenantUsecase) *Handler {
	return &Handler{uc: uc}
}

type tenantResponse struct {
	ID          string  `json:"id"`
	Type        string  `json:"type"`
	LegalName   string  `json:"legal_name"`
	DisplayName string  `json:"display_name"`
	Slug        string  `json:"slug"`
	Status      string  `json:"status"`
	Timezone    string  `json:"timezone"`
	Locale      string  `json:"locale"`
	TrialEndsAt *string `json:"trial_ends_at,omitempty"`
}

func toTenantResponse(t *domain.Tenant) tenantResponse {
	resp := tenantResponse{
		ID:          t.ID.String(),
		Type:        string(t.Type),
		LegalName:   t.LegalName,
		DisplayName: t.DisplayName,
		Slug:        t.Slug,
		Status:      string(t.Status),
		Timezone:    t.Timezone,
		Locale:      t.Locale,
	}
	if t.TrialEndsAt != nil {
		s := t.TrialEndsAt.Format("2006-01-02T15:04:05Z07:00")
		resp.TrialEndsAt = &s
	}
	return resp
}

// GetCurrent returns the tenant from context (set by TenantMiddleware).
// GET /tenants/current
func (h *Handler) GetCurrent(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := identitydomain.TenantIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusBadRequest, "missing_tenant", "No tenant in context")
		return
	}

	tenant, err := h.uc.GetByID(r.Context(), tenantID)
	if err != nil {
		if errors.Is(err, domain.ErrTenantNotFound) {
			server.RenderError(w, r, http.StatusNotFound, "tenant_not_found", "Tenant not found")
			return
		}
		server.RenderError(w, r, http.StatusInternalServerError, "internal_error", "Failed to fetch tenant")
		return
	}

	server.RenderJSON(w, http.StatusOK, toTenantResponse(tenant))
}
