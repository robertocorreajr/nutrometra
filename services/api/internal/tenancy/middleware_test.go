package tenancy_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	identitydomain "nutrometra/api/internal/identity/domain"
	"nutrometra/api/internal/tenancy"
	"nutrometra/api/internal/tenancy/domain"
	"nutrometra/api/internal/tenancy/usecase"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// stubTenantRepo implements repository.TenantRepository for tests.
type stubTenantRepo struct {
	tenants map[uuid.UUID]*domain.Tenant
	members map[string]bool
}

func (s *stubTenantRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Tenant, error) {
	t, ok := s.tenants[id]
	if !ok {
		return nil, domain.ErrTenantNotFound
	}
	return t, nil
}

func (s *stubTenantRepo) GetBySlug(ctx context.Context, slug string) (*domain.Tenant, error) {
	for _, t := range s.tenants {
		if t.Slug == slug {
			return t, nil
		}
	}
	return nil, domain.ErrTenantNotFound
}

func (s *stubTenantRepo) IsMember(ctx context.Context, tenantID, userID uuid.UUID) (bool, error) {
	return s.members[tenantID.String()+":"+userID.String()], nil
}

func emptyUsecase() usecase.TenantUsecase {
	return usecase.NewTenantUsecase(&stubTenantRepo{
		tenants: make(map[uuid.UUID]*domain.Tenant),
		members: make(map[string]bool),
	})
}

func activeTenant() *domain.Tenant {
	return &domain.Tenant{
		ID:          uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"),
		Type:        domain.TenantTypeClinic,
		LegalName:   "Test Clinic Ltda",
		DisplayName: "Test Clinic",
		Slug:        "test-clinic",
		Status:      domain.TenantStatusActive,
		Timezone:    "America/Sao_Paulo",
		Locale:      "pt-BR",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func TestTenantMiddleware_NilUsecase_Panics(t *testing.T) {
	assert.Panics(t, func() {
		tenancy.TenantMiddleware(nil)
	})
}

func TestTenantMiddleware_MissingHeader(t *testing.T) {
	mw := tenancy.TenantMiddleware(emptyUsecase())
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	mw(next).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	var body map[string]string
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &body))
	assert.Equal(t, "missing_tenant_id", body["code"])
}

func TestTenantMiddleware_InvalidUUID(t *testing.T) {
	mw := tenancy.TenantMiddleware(emptyUsecase())
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Tenant-ID", "not-a-uuid")
	rr := httptest.NewRecorder()
	mw(next).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	var body map[string]string
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &body))
	assert.Equal(t, "invalid_tenant_id", body["code"])
}

func TestTenantMiddleware_ValidUUID_NotInDB(t *testing.T) {
	mw := tenancy.TenantMiddleware(emptyUsecase())
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Tenant-ID", uuid.New().String())
	rr := httptest.NewRecorder()
	mw(next).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusForbidden, rr.Code)
	var body map[string]string
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &body))
	assert.Equal(t, "tenant_not_found", body["code"])
}

func TestTenantMiddleware_TenantNotFound(t *testing.T) {
	repo := &stubTenantRepo{tenants: make(map[uuid.UUID]*domain.Tenant), members: make(map[string]bool)}
	uc := usecase.NewTenantUsecase(repo)
	mw := tenancy.TenantMiddleware(uc)
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Tenant-ID", uuid.New().String())
	rr := httptest.NewRecorder()
	mw(next).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusForbidden, rr.Code)
}

func TestTenantMiddleware_TenantSuspended(t *testing.T) {
	tenant := activeTenant()
	tenant.Status = domain.TenantStatusSuspended
	repo := &stubTenantRepo{
		tenants: map[uuid.UUID]*domain.Tenant{tenant.ID: tenant},
		members: make(map[string]bool),
	}
	uc := usecase.NewTenantUsecase(repo)
	mw := tenancy.TenantMiddleware(uc)
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Tenant-ID", tenant.ID.String())
	rr := httptest.NewRecorder()
	mw(next).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusForbidden, rr.Code)
	var body map[string]string
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &body))
	assert.Equal(t, "tenant_suspended", body["code"])
}

func TestTenantMiddleware_TenantCancelled(t *testing.T) {
	tenant := activeTenant()
	tenant.Status = domain.TenantStatusCancelled
	repo := &stubTenantRepo{
		tenants: map[uuid.UUID]*domain.Tenant{tenant.ID: tenant},
		members: make(map[string]bool),
	}
	uc := usecase.NewTenantUsecase(repo)
	mw := tenancy.TenantMiddleware(uc)
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Tenant-ID", tenant.ID.String())
	rr := httptest.NewRecorder()
	mw(next).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusForbidden, rr.Code)
	var body map[string]string
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &body))
	assert.Equal(t, "tenant_not_active", body["code"])
}

func TestTenantMiddleware_ActiveTenant_NoUser_OK(t *testing.T) {
	// No user in context → membership check skipped, request passes
	tenant := activeTenant()
	repo := &stubTenantRepo{
		tenants: map[uuid.UUID]*domain.Tenant{tenant.ID: tenant},
		members: make(map[string]bool),
	}
	uc := usecase.NewTenantUsecase(repo)
	mw := tenancy.TenantMiddleware(uc)
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Tenant-ID", tenant.ID.String())
	rr := httptest.NewRecorder()
	mw(next).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestTenantMiddleware_ActiveTenant_UserNotMember(t *testing.T) {
	tenant := activeTenant()
	userID := uuid.New()
	repo := &stubTenantRepo{
		tenants: map[uuid.UUID]*domain.Tenant{tenant.ID: tenant},
		members: make(map[string]bool), // no membership
	}
	uc := usecase.NewTenantUsecase(repo)
	mw := tenancy.TenantMiddleware(uc)
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Tenant-ID", tenant.ID.String())
	// Inject authenticated user into context
	ctx := identitydomain.SetUserInContext(req.Context(), userID, "user@test.com")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()
	mw(next).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusForbidden, rr.Code)
}

func TestTenantMiddleware_ActiveTenant_UserIsMember_OK(t *testing.T) {
	tenant := activeTenant()
	userID := uuid.New()
	repo := &stubTenantRepo{
		tenants: map[uuid.UUID]*domain.Tenant{tenant.ID: tenant},
		members: map[string]bool{
			tenant.ID.String() + ":" + userID.String(): true,
		},
	}
	uc := usecase.NewTenantUsecase(repo)
	mw := tenancy.TenantMiddleware(uc)

	var capturedTenantID uuid.UUID
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, ok := identitydomain.TenantIDFromContext(r.Context())
		assert.True(t, ok)
		capturedTenantID = id
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Tenant-ID", tenant.ID.String())
	ctx := identitydomain.SetUserInContext(req.Context(), userID, "user@test.com")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()
	mw(next).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, tenant.ID, capturedTenantID)
}
