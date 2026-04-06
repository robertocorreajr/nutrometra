package identity

import (
	"log/slog"
	"net/http"
	"strings"

	"nutrometra/api/internal/identity/domain"
	"nutrometra/api/internal/identity/oidc"
	"nutrometra/api/internal/platform/server"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// AuthMiddleware validates the JWT and stores the Zitadel subject and email in
// the request context. Does NOT resolve the internal user UUID — that is done
// by UserResolverMiddleware.
// If the token is absent, returns 401 immediately.
// If validator is nil and a token is present, returns 401 (validator not configured).
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

			// Store Zitadel subject and email; UserResolverMiddleware converts to internal UUID.
			ctx := domain.SetExternalAuthInContext(r.Context(), claims.Subject, claims.Email)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// ExtractBearerToken extracts the token from Authorization: Bearer <token>.
// The "bearer" scheme prefix is matched case-insensitively per RFC 7235.
// Returns empty string if the header is absent, scheme is not bearer, or token is empty.
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

// UserResolverMiddleware upserts the user in the DB using the Zitadel subject as
// the external_auth_id, then injects the internal UUID into the context.
// Must run after AuthMiddleware.
//
// Pass-through behaviour: if no external auth is present in the context (e.g.
// on a public route that does not go through AuthMiddleware), the request is
// forwarded without modification. Downstream handlers on protected routes must
// call domain.UserIDFromContext and handle the missing-UUID case explicitly.
func UserResolverMiddleware(pool *pgxpool.Pool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			subject, email, ok := domain.ExternalAuthFromContext(r.Context())
			if !ok {
				// No authenticated identity in context — pass through to handler.
				next.ServeHTTP(w, r)
				return
			}

			var userID uuid.UUID
			// Argument order: $1 = email, $2 = external_auth_id (Zitadel subject).
			// ON CONFLICT targets the unique index on external_auth_id.
			err := pool.QueryRow(r.Context(),
				`INSERT INTO users (id, email, external_auth_id, status, last_login_at, created_at, updated_at)
				 VALUES (gen_random_uuid(), $1, $2, 'active', NOW(), NOW(), NOW())
				 ON CONFLICT (external_auth_id) DO UPDATE SET
				   email          = EXCLUDED.email,
				   last_login_at  = NOW(),
				   updated_at     = NOW()
				 RETURNING id`,
				email, subject, // $1=email, $2=external_auth_id
			).Scan(&userID)
			if err != nil {
				slog.ErrorContext(r.Context(), "failed to resolve user",
					"external_auth_id", subject,
					"error", err,
				)
				server.RenderError(w, r, http.StatusInternalServerError, "user_resolve_failed", "Failed to resolve user")
				return
			}

			ctx := domain.SetUserInContext(r.Context(), userID, email)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
