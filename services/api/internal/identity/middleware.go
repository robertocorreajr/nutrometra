package identity

import (
	"net/http"
	"strings"

	"nutrometra/api/internal/identity/domain"
	"nutrometra/api/internal/identity/oidc"
	"nutrometra/api/internal/platform/server"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// AuthMiddleware validates the JWT and injects the Zitadel subject and email
// into the request context. Does NOT resolve the internal user UUID — that is
// done by UserResolverMiddleware.
// If validator is nil (e.g. in tests), any token triggers 401.
func AuthMiddleware(validator *oidc.Validator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rawToken := ExtractBearerToken(r)
			if rawToken == "" {
				server.RenderError(w, r, http.StatusUnauthorized, "missing_token", "Authorization header required")
				return
			}

			if validator == nil {
				server.RenderError(w, r, http.StatusUnauthorized, "validator_unavailable", "Auth validator not configured")
				return
			}

			claims, err := validator.Validate(r.Context(), rawToken)
			if err != nil {
				server.RenderError(w, r, http.StatusUnauthorized, "invalid_token", "Invalid or expired token")
				return
			}

			// Store Zitadel subject and email; UserResolverMiddleware resolves the UUID.
			ctx := r.Context()
			ctx = domain.SetExternalAuthInContext(ctx, claims.Subject, claims.Email)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// ExtractBearerToken extracts the token from Authorization: Bearer <token>.
// Returns empty string if the header is missing or not Bearer format.
func ExtractBearerToken(r *http.Request) string {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return ""
	}
	parts := strings.SplitN(auth, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
		return ""
	}
	return strings.TrimSpace(parts[1])
}

// UserResolverMiddleware upserts the user in the DB based on the Zitadel subject
// and injects the internal UUID into the context. Must run after AuthMiddleware.
func UserResolverMiddleware(pool *pgxpool.Pool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			subject, email, ok := domain.ExternalAuthFromContext(r.Context())
			if !ok {
				next.ServeHTTP(w, r)
				return
			}

			var userID uuid.UUID
			err := pool.QueryRow(r.Context(),
				`INSERT INTO users (id, email, external_auth_id, status, created_at, updated_at)
				 VALUES (gen_random_uuid(), $1, $2, 'active', NOW(), NOW())
				 ON CONFLICT (external_auth_id) DO UPDATE SET
				   email = EXCLUDED.email,
				   last_login_at = NOW(),
				   updated_at = NOW()
				 RETURNING id`,
				email, subject,
			).Scan(&userID)
			if err != nil {
				server.RenderError(w, r, http.StatusInternalServerError, "user_resolve_failed", "Failed to resolve user")
				return
			}

			ctx := domain.SetUserInContext(r.Context(), userID, email)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
