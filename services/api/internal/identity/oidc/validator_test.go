package oidc_test

import (
	"context"
	"testing"
	"time"

	"nutrometra/api/internal/identity/oidc"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidator_New(t *testing.T) {
	v := oidc.NewValidator("http://localhost:8080", "test-client")
	assert.NotNil(t, v)
}

func TestValidator_ValidateRaw_MalformedToken(t *testing.T) {
	v := oidc.NewValidator("http://localhost:8080", "test-client")
	_, err := v.ValidateRaw(context.Background(), "not.a.jwt")
	assert.Error(t, err)
}

func TestValidator_ValidateRaw_ExpiredToken(t *testing.T) {
	v := oidc.NewValidator("http://localhost:8080", "test-client")

	key := []byte("test-secret")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "user-123",
		"exp": time.Now().Add(-1 * time.Hour).Unix(),
		"iss": "http://localhost:8080",
		"aud": jwt.ClaimStrings{"test-client"},
	})
	signed, err := token.SignedString(key)
	require.NoError(t, err)

	_, err = v.ValidateRaw(context.Background(), signed)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "token")
}

func TestValidator_ValidateRaw_ValidStructure(t *testing.T) {
	v := oidc.NewValidator("http://localhost:8080", "test-client")

	key := []byte("test-secret")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   "user-abc",
		"email": "user@example.com",
		"exp":   time.Now().Add(1 * time.Hour).Unix(),
		"iss":   "http://localhost:8080",
		"aud":   jwt.ClaimStrings{"test-client"},
	})
	signed, err := token.SignedString(key)
	require.NoError(t, err)

	claims, err := v.ValidateRaw(context.Background(), signed)
	require.NoError(t, err)
	assert.Equal(t, "user-abc", claims.Subject)
	assert.Equal(t, "user@example.com", claims.Email)
}

func TestValidator_Validate_WithoutInit(t *testing.T) {
	v := oidc.NewValidator("http://localhost:8080", "test-client")

	key := []byte("test-secret")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "user-123",
		"exp": time.Now().Add(1 * time.Hour).Unix(),
		"iss": "http://localhost:8080",
		"aud": jwt.ClaimStrings{"test-client"},
	})
	signed, err := token.SignedString(key)
	require.NoError(t, err)

	_, err = v.Validate(context.Background(), signed)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not initialized")
}

func TestValidator_ValidateRaw_MissingSub(t *testing.T) {
	v := oidc.NewValidator("http://localhost:8080", "test-client")

	key := []byte("test-secret")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": "user@example.com",
		"exp":   time.Now().Add(1 * time.Hour).Unix(),
		"iss":   "http://localhost:8080",
		"aud":   jwt.ClaimStrings{"test-client"},
	})
	signed, err := token.SignedString(key)
	require.NoError(t, err)

	_, err = v.ValidateRaw(context.Background(), signed)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "subject")
}
