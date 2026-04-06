package identity_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"nutrometra/api/internal/identity"
	"nutrometra/api/internal/identity/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthMiddleware_MissingToken(t *testing.T) {
	mw := identity.AuthMiddleware(nil)
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	rr := httptest.NewRecorder()

	mw(next).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestAuthMiddleware_InvalidBearerFormat(t *testing.T) {
	mw := identity.AuthMiddleware(nil)
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "NotBearer token")
	rr := httptest.NewRecorder()

	mw(next).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestExtractBearerToken(t *testing.T) {
	tests := []struct {
		header   string
		expected string
	}{
		{"Bearer my-token", "my-token"},
		{"bearer my-token", "my-token"},
		{"Token my-token", ""},
		{"", ""},
		{"Bearer", ""},
	}
	for _, tc := range tests {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		if tc.header != "" {
			req.Header.Set("Authorization", tc.header)
		}
		got := identity.ExtractBearerToken(req)
		assert.Equal(t, tc.expected, got, "header: %q", tc.header)
	}
}

func TestUserIDFromContext_NotSet(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	_, ok := domain.UserIDFromContext(req.Context())
	assert.False(t, ok)
}

func TestMeHandler_Authenticated(t *testing.T) {
	userID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	email := "user@example.com"

	req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	ctx := domain.SetUserInContext(req.Context(), userID, email)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	identity.MeHandler(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	var body map[string]string
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &body))
	assert.Equal(t, userID.String(), body["user_id"])
	assert.Equal(t, email, body["email"])
}

func TestMeHandler_Unauthenticated(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	rr := httptest.NewRecorder()

	identity.MeHandler(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestUserResolverMiddleware_PassThrough_NoExternalAuth(t *testing.T) {
	// When no external auth is in context (e.g. public route skipping AuthMiddleware),
	// UserResolverMiddleware should forward the request unchanged.
	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		// Internal user UUID must NOT have been injected.
		_, ok := domain.UserIDFromContext(r.Context())
		assert.False(t, ok, "user_id should not be set on pass-through")
		w.WriteHeader(http.StatusOK)
	})

	// Pass nil pool — if the middleware tries to use the pool it will panic,
	// which would fail the test; a nil pool is safe only on the pass-through path.
	mw := identity.UserResolverMiddleware(nil)
	req := httptest.NewRequest(http.MethodGet, "/public", nil)
	rr := httptest.NewRecorder()

	mw(next).ServeHTTP(rr, req)

	assert.True(t, called, "next handler should have been called")
	assert.Equal(t, http.StatusOK, rr.Code)
}
