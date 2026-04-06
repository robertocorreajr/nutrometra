package config_test

import (
	"os"
	"testing"

	"nutrometra/api/internal/platform/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad_AllDefaults(t *testing.T) {
	os.Clearenv()
	os.Setenv("ZITADEL_ISSUER", "http://localhost:8080")
	os.Setenv("ZITADEL_CLIENT_ID", "test-client")

	cfg, err := config.Load()
	require.NoError(t, err)

	assert.Equal(t, "localhost", cfg.Postgres.Host)
	assert.Equal(t, 5432, cfg.Postgres.Port)
	assert.Equal(t, 8081, cfg.API.Port)
	assert.Equal(t, "development", cfg.API.Env)
}

func TestLoad_MissingZitadelIssuer(t *testing.T) {
	os.Clearenv()
	_, err := config.Load()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ZITADEL_ISSUER")
}

func TestLoad_FromEnv(t *testing.T) {
	os.Clearenv()
	os.Setenv("POSTGRES_HOST", "db.example.com")
	os.Setenv("POSTGRES_PORT", "5433")
	os.Setenv("ZITADEL_ISSUER", "https://auth.example.com")
	os.Setenv("ZITADEL_CLIENT_ID", "my-client")
	os.Setenv("API_PORT", "9090")

	cfg, err := config.Load()
	require.NoError(t, err)

	assert.Equal(t, "db.example.com", cfg.Postgres.Host)
	assert.Equal(t, 5433, cfg.Postgres.Port)
	assert.Equal(t, 9090, cfg.API.Port)
}
