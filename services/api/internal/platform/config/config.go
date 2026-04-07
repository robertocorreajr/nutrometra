package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	API      APIConfig
	Postgres PostgresConfig
	Redis    RedisConfig
	Zitadel  ZitadelConfig
	Stripe   StripeConfig
	Google   GoogleConfig
}

type StripeConfig struct {
	SecretKey      string
	WebhookSecret  string
	PublishableKey string
}

type GoogleConfig struct {
	ClientID      string
	ClientSecret  string
	RedirectURL   string
	EncryptionKey string // 64 hex chars for AES-256
}

type APIConfig struct {
	Port     int
	Env      string
	LogLevel string
}

type PostgresConfig struct {
	Host     string
	Port     int
	DB       string
	User     string
	Password string
}

func (c PostgresConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",
		c.Host, c.Port, c.DB, c.User, c.Password,
	)
}

type RedisConfig struct {
	Addr     string
	Password string
}

type ZitadelConfig struct {
	Issuer   string
	ClientID string
	Domain   string
}

func Load() (*Config, error) {
	// Carrega .env apenas em desenvolvimento (ignora erro se não existir)
	_ = godotenv.Load()

	zitadelIssuer := os.Getenv("ZITADEL_ISSUER")
	if zitadelIssuer == "" {
		return nil, fmt.Errorf("ZITADEL_ISSUER is required")
	}

	zitadelClientID := os.Getenv("ZITADEL_CLIENT_ID")
	if zitadelClientID == "" {
		env := envStr("API_ENV", "development")
		if env != "development" {
			return nil, fmt.Errorf("ZITADEL_CLIENT_ID is required in non-development environments")
		}
		slog.Warn("ZITADEL_CLIENT_ID not set — OIDC token validation will reject all tokens. Run: make setup-zitadel")
	}

	postgresPassword := os.Getenv("POSTGRES_PASSWORD")
	if postgresPassword == "" {
		return nil, fmt.Errorf("POSTGRES_PASSWORD is required")
	}

	return &Config{
		API: APIConfig{
			Port:     envInt("API_PORT", 8081),
			Env:      envStr("API_ENV", "development"),
			LogLevel: envStr("LOG_LEVEL", "info"),
		},
		Postgres: PostgresConfig{
			Host:     envStr("POSTGRES_HOST", "localhost"),
			Port:     envInt("POSTGRES_PORT", 5432),
			DB:       envStr("POSTGRES_DB", "nutrometra"),
			User:     envStr("POSTGRES_USER", "nutrometra"),
			Password: postgresPassword,
		},
		Redis: RedisConfig{
			Addr:     envStr("REDIS_ADDR", "localhost:6379"),
			Password: envStr("REDIS_PASSWORD", ""),
		},
		Zitadel: ZitadelConfig{
			Issuer:   zitadelIssuer,
			ClientID: zitadelClientID,
			Domain:   envStr("ZITADEL_DOMAIN", "localhost"),
		},
		Stripe: StripeConfig{
			SecretKey:      envStr("STRIPE_SECRET_KEY", ""),
			WebhookSecret:  envStr("STRIPE_WEBHOOK_SECRET", ""),
			PublishableKey: envStr("STRIPE_PUBLISHABLE_KEY", ""),
		},
		Google: GoogleConfig{
			ClientID:      envStr("GOOGLE_CLIENT_ID", ""),
			ClientSecret:  envStr("GOOGLE_CLIENT_SECRET", ""),
			RedirectURL:   envStr("GOOGLE_REDIRECT_URL", "http://localhost:8081/integrations/google/callback"),
			EncryptionKey: envStr("GOOGLE_ENCRYPTION_KEY", ""),
		},
	}, nil
}

func envStr(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

func envInt(key string, defaultVal int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return defaultVal
}
