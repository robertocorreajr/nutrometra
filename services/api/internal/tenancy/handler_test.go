package tenancy_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	identitydomain "nutrometra/api/internal/identity/domain"
	"nutrometra/api/internal/tenancy"
	"nutrometra/api/internal/tenancy/domain"
	"nutrometra/api/internal/tenancy/usecase"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandler_GetCurrent_OK(t *testing.T) {
	tenant := activeTenant()
	repo := &stubTenantRepo{
		tenants: map[uuid.UUID]*domain.Tenant{tenant.ID: tenant},
		members: make(map[string]bool),
	}
	uc := usecase.NewTenantUsecase(repo)
	h := tenancy.NewHandler(uc)

	req := httptest.NewRequest(http.MethodGet, "/tenants/current", nil)
	ctx := identitydomain.SetTenantInContext(req.Context(), tenant.ID)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	h.GetCurrent(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	var body map[string]interface{}
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &body))
	assert.Equal(t, tenant.ID.String(), body["id"])
	assert.Equal(t, string(tenant.Type), body["type"])
	assert.Equal(t, tenant.DisplayName, body["display_name"])
}

func TestHandler_GetCurrent_NoTenantInContext(t *testing.T) {
	repo := &stubTenantRepo{
		tenants: make(map[uuid.UUID]*domain.Tenant),
		members: make(map[string]bool),
	}
	uc := usecase.NewTenantUsecase(repo)
	h := tenancy.NewHandler(uc)

	req := httptest.NewRequest(http.MethodGet, "/tenants/current", nil)
	rr := httptest.NewRecorder()

	h.GetCurrent(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}
