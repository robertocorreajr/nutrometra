package rbac

import (
	"net/http"

	identitydomain "nutrometra/api/internal/identity/domain"
	"nutrometra/api/internal/platform/server"
	"nutrometra/api/internal/rbac/repository"
)

// RequirePermission returns a middleware that denies access if the authenticated
// user does not have the specified permission in the current tenant.
// Deny-by-default: missing user or tenant in context = 403.
func RequirePermission(permission string, enforcer repository.Enforcer) func(http.Handler) http.Handler {
	if enforcer == nil {
		panic("RequirePermission requires a non-nil Enforcer")
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := identitydomain.UserIDFromContext(r.Context())
			if !ok {
				server.RenderError(w, r, http.StatusForbidden, "permission_denied", "Access denied")
				return
			}

			tenantID, ok := identitydomain.TenantIDFromContext(r.Context())
			if !ok {
				server.RenderError(w, r, http.StatusForbidden, "permission_denied", "Access denied")
				return
			}

			allowed, err := enforcer.HasPermission(r.Context(), userID, tenantID, permission)
			if err != nil {
				server.RenderError(w, r, http.StatusInternalServerError, "permission_check_failed", "Failed to verify permission")
				return
			}
			if !allowed {
				server.RenderError(w, r, http.StatusForbidden, "permission_denied", "You do not have permission to perform this action")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
