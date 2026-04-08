package identity

import (
	"net/http"

	"nutrometra/api/internal/identity/domain"
	"nutrometra/api/internal/platform/server"
	tenancyrepo "nutrometra/api/internal/tenancy/repository"
)

type meResponse struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
}

// MeHandler returns the authenticated user's profile.
// Requires AuthMiddleware + UserResolverMiddleware to have run.
func MeHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := domain.UserIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusUnauthorized, "unauthenticated", "Not authenticated")
		return
	}
	email, _ := domain.UserEmailFromContext(r.Context())

	server.RenderJSON(w, http.StatusOK, meResponse{
		UserID: userID.String(),
		Email:  email,
	})
}

type tenantItem struct {
	TenantID    string `json:"tenant_id"`
	DisplayName string `json:"display_name"`
	Slug        string `json:"slug"`
	Type        string `json:"type"`
}

// MeTenantsHandler returns the tenants the authenticated user belongs to.
func MeTenantsHandler(tenantRepo tenancyrepo.TenantRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := domain.UserIDFromContext(r.Context())
		if !ok {
			server.RenderError(w, r, http.StatusUnauthorized, "unauthenticated", "Not authenticated")
			return
		}

		tenants, err := tenantRepo.ListByUser(r.Context(), userID)
		if err != nil {
			server.RenderError(w, r, http.StatusInternalServerError, "internal", "Failed to list tenants")
			return
		}

		items := make([]tenantItem, 0, len(tenants))
		for _, t := range tenants {
			items = append(items, tenantItem{
				TenantID:    t.ID.String(),
				DisplayName: t.DisplayName,
				Slug:        t.Slug,
				Type:        string(t.Type),
			})
		}

		server.RenderJSON(w, http.StatusOK, items)
	}
}
