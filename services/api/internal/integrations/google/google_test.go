package google

import (
	"testing"

	"nutrometra/api/internal/platform/audit"
	"nutrometra/api/internal/platform/config"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewHandler(t *testing.T) {
	auditSvc := audit.NewService()
	provider := NewGoogleCalendarProvider()
	testKey := generateTestKey(t)
	cfg := config.GoogleConfig{
		ClientID:      "test-client-id",
		ClientSecret:  "test-client-secret",
		RedirectURL:   "http://localhost:8081/integrations/google/callback",
		EncryptionKey: testKey,
	}

	handler := NewHandler(nil, nil, provider, auditSvc, cfg, nil)

	assert.NotNil(t, handler)
	assert.NotNil(t, handler.oauthConfig)
	assert.Equal(t, cfg.ClientID, handler.oauthConfig.ClientID)
	assert.Equal(t, cfg.ClientSecret, handler.oauthConfig.ClientSecret)
	assert.Equal(t, cfg.RedirectURL, handler.oauthConfig.RedirectURL)
	assert.Len(t, handler.oauthConfig.Scopes, 1)
	assert.Equal(t, calendarEventsScope, handler.oauthConfig.Scopes[0])
}

func TestStateTokenGenerateAndVerify(t *testing.T) {
	testKey := generateTestKey(t)
	handler := &Handler{encryptionKey: testKey}

	tenantID := uuid.New()
	userID := uuid.New()

	state := handler.generateState(tenantID, userID)
	assert.NotEmpty(t, state)

	// Verify should succeed.
	gotTenantID, gotUserID, err := handler.verifyState(state)
	require.NoError(t, err)
	assert.Equal(t, tenantID, gotTenantID)
	assert.Equal(t, userID, gotUserID)
}

func TestStateTokenTampered(t *testing.T) {
	testKey := generateTestKey(t)
	handler := &Handler{encryptionKey: testKey}

	tenantID := uuid.New()
	userID := uuid.New()

	state := handler.generateState(tenantID, userID)

	// Tamper with the state.
	tampered := state + "x"
	_, _, err := handler.verifyState(tampered)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "signature mismatch")
}

func TestStateTokenDifferentKey(t *testing.T) {
	handler1 := &Handler{encryptionKey: generateTestKey(t)}
	handler2 := &Handler{encryptionKey: generateTestKey(t)}

	tenantID := uuid.New()
	userID := uuid.New()

	state := handler1.generateState(tenantID, userID)

	// Should fail with different key.
	_, _, err := handler2.verifyState(state)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "signature mismatch")
}

func TestStateTokenInvalidFormat(t *testing.T) {
	handler := &Handler{encryptionKey: generateTestKey(t)}

	tests := []struct {
		name  string
		state string
	}{
		{"empty", ""},
		{"no colons", "justasinglevalue"},
		{"one colon", "part1:part2"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := handler.verifyState(tt.state)
			assert.Error(t, err)
		})
	}
}

func TestNewGoogleCalendarProvider(t *testing.T) {
	provider := NewGoogleCalendarProvider()
	assert.NotNil(t, provider)
}
