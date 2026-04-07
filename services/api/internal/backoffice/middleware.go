package backoffice

import (
	"context"
	"net/http"

	identitydomain "nutrometra/api/internal/identity/domain"
	"nutrometra/api/internal/platform/server"
)

type contextKey string

const contextKeyActorScope contextKey = "backoffice_actor_scope"

// ActorScopeFromContext returns the actor scope set by BackofficeMiddleware.
func ActorScopeFromContext(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(contextKeyActorScope).(string)
	return v, ok
}

// BackofficeMiddleware verifies that the authenticated user is an active backoffice user.
// Does NOT require X-Tenant-ID header (unlike TenantMiddleware).
func BackofficeMiddleware(repo *Repository) func(http.Handler) http.Handler {
	if repo == nil {
		panic("BackofficeMiddleware requires a non-nil Repository")
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := identitydomain.UserIDFromContext(r.Context())
			if !ok {
				server.RenderError(w, r, http.StatusForbidden, "access_denied", "Authentication required")
				return
			}

			isBO, err := repo.IsBackofficeUser(r.Context(), userID)
			if err != nil {
				server.RenderError(w, r, http.StatusInternalServerError, "backoffice_check_failed", "Failed to verify backoffice access")
				return
			}
			if !isBO {
				server.RenderError(w, r, http.StatusForbidden, "not_backoffice_user", "Backoffice access denied")
				return
			}

			ctx := context.WithValue(r.Context(), contextKeyActorScope, "backoffice")
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireBackofficePermission returns a middleware that checks if the user
// has a specific backoffice permission.
func RequireBackofficePermission(permission string, repo *Repository) func(http.Handler) http.Handler {
	if repo == nil {
		panic("RequireBackofficePermission requires a non-nil Repository")
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := identitydomain.UserIDFromContext(r.Context())
			if !ok {
				server.RenderError(w, r, http.StatusForbidden, "permission_denied", "Access denied")
				return
			}

			allowed, err := repo.HasBackofficePermission(r.Context(), userID, permission)
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

