package rbac_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	identitydomain "nutrometra/api/internal/identity/domain"
	"nutrometra/api/internal/rbac"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockEnforcer implements rbac.Enforcer for tests.
type mockEnforcer struct {
	result bool
	err    error
}

func (m *mockEnforcer) HasPermission(_ context.Context, _, _ uuid.UUID, _ string) (bool, error) {
	return m.result, m.err
}

func TestRequirePermission_Allowed(t *testing.T) {
	mw := rbac.RequirePermission("tenant:manage", &mockEnforcer{result: true})
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx := identitydomain.SetUserInContext(req.Context(), uuid.New(), "user@test.com")
	ctx = identitydomain.SetTenantInContext(ctx, uuid.New())
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	mw(next).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestRequirePermission_Denied(t *testing.T) {
	mw := rbac.RequirePermission("tenant:manage", &mockEnforcer{result: false})
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx := identitydomain.SetUserInContext(req.Context(), uuid.New(), "user@test.com")
	ctx = identitydomain.SetTenantInContext(ctx, uuid.New())
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	mw(next).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusForbidden, rr.Code)
	var body map[string]string
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &body))
	assert.Equal(t, "permission_denied", body["code"])
}

func TestRequirePermission_NoUserInContext(t *testing.T) {
	mw := rbac.RequirePermission("tenant:manage", &mockEnforcer{result: true})
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// No user, no tenant
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	mw(next).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusForbidden, rr.Code)
}

func TestRequirePermission_NoTenantInContext(t *testing.T) {
	mw := rbac.RequirePermission("tenant:manage", &mockEnforcer{result: true})
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// User but no tenant
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx := identitydomain.SetUserInContext(req.Context(), uuid.New(), "user@test.com")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()
	mw(next).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusForbidden, rr.Code)
}

func TestRequirePermission_EnforcerError(t *testing.T) {
	mw := rbac.RequirePermission("tenant:manage", &mockEnforcer{result: false, err: fmt.Errorf("db down")})
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx := identitydomain.SetUserInContext(req.Context(), uuid.New(), "user@test.com")
	ctx = identitydomain.SetTenantInContext(ctx, uuid.New())
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	mw(next).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestRequirePermission_NilEnforcer_Panics(t *testing.T) {
	assert.Panics(t, func() {
		rbac.RequirePermission("tenant:manage", nil)
	})
}
