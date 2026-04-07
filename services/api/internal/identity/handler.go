package identity

import (
	"net/http"

	"nutrometra/api/internal/identity/domain"
	"nutrometra/api/internal/platform/server"
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
