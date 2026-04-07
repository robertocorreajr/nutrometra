package tenancy

import (
	"errors"
	"net/http"

	identitydomain "nutrometra/api/internal/identity/domain"
	"nutrometra/api/internal/platform/server"
	"nutrometra/api/internal/tenancy/domain"
	"nutrometra/api/internal/tenancy/usecase"

	"github.com/google/uuid"
)

// TenantMiddleware resolves and validates the tenant from the X-Tenant-ID header.
// Guarantees:
//  1. Header present and valid UUID.
//  2. Tenant exists in DB.
//  3. Tenant is active (not suspended or cancelled).
//  4. Authenticated user is an active member of the tenant (when user in context).
//
// Panics if uc is nil — fail fast at wiring time, not at request time.
func TenantMiddleware(uc usecase.TenantUsecase) func(http.Handler) http.Handler {
	if uc == nil {
		panic("TenantMiddleware requires a non-nil TenantUsecase")
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tenantIDStr := r.Header.Get("X-Tenant-ID")
			if tenantIDStr == "" {
				server.RenderError(w, r, http.StatusBadRequest, "missing_tenant_id", "X-Tenant-ID header is required")
				return
			}

			tenantID, err := uuid.Parse(tenantIDStr)
			if err != nil {
				server.RenderError(w, r, http.StatusBadRequest, "invalid_tenant_id", "X-Tenant-ID must be a valid UUID")
				return
			}

			tenant, err := uc.GetByID(r.Context(), tenantID)
			if err != nil {
				if errors.Is(err, domain.ErrTenantNotFound) {
					server.RenderError(w, r, http.StatusForbidden, "tenant_not_found", "Tenant not found")
					return
				}
				server.RenderError(w, r, http.StatusInternalServerError, "tenant_lookup_failed", "Failed to resolve tenant")
				return
			}
			if tenant.IsSuspended() {
				server.RenderError(w, r, http.StatusForbidden, "tenant_suspended", "Tenant is suspended")
				return
			}
			if !tenant.IsActive() {
				server.RenderError(w, r, http.StatusForbidden, "tenant_not_active", "Tenant is not active")
				return
			}

			// Membership check: if user is authenticated, verify they belong to this tenant.
			userID, ok := identitydomain.UserIDFromContext(r.Context())
			if ok {
				isMember, err := uc.IsMember(r.Context(), tenantID, userID)
				if err != nil {
					server.RenderError(w, r, http.StatusInternalServerError, "membership_check_failed", "Failed to check membership")
					return
				}
				if !isMember {
					server.RenderError(w, r, http.StatusForbidden, "not_a_member", "User is not a member of this tenant")
					return
				}
			}

			ctx := identitydomain.SetTenantInContext(r.Context(), tenantID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
