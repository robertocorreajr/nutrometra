package identity_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"nutrometra/api/internal/identity"
	"nutrometra/api/internal/identity/domain"

	"github.com/stretchr/testify/assert"
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
