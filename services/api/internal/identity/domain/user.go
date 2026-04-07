package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// User representa um usuario autenticado no sistema.
type User struct {
	ID             uuid.UUID
	Email          string
	ExternalAuthID string // sub do JWT (identificador no Zitadel)
	Status         string
	LastLoginAt    *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// Claims sao os campos extraidos do JWT do Zitadel.
type Claims struct {
	Subject string // JWT "sub" = ExternalAuthID
	Email   string
	Name    string
}

// contextKey evita colisoes de chave no context.
type contextKey string

const (
	ContextKeyUserID         contextKey = "user_id"
	ContextKeyUserEmail      contextKey = "user_email"
	ContextKeyTenantID       contextKey = "tenant_id"
	ContextKeyActorRole      contextKey = "actor_role"
	ContextKeyActorScope     contextKey = "actor_scope"
	ContextKeyExternalAuthID contextKey = "external_auth_id"
)

// SetUserInContext adiciona user_id e user_email ao contexto.
func SetUserInContext(ctx context.Context, userID uuid.UUID, email string) context.Context {
	ctx = context.WithValue(ctx, ContextKeyUserID, userID)
	ctx = context.WithValue(ctx, ContextKeyUserEmail, email)
	return ctx
}

// UserIDFromContext extrai o user_id do contexto.
func UserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(ContextKeyUserID).(uuid.UUID)
	return id, ok
}

// UserEmailFromContext extrai o user_email do contexto.
func UserEmailFromContext(ctx context.Context) (string, bool) {
	email, ok := ctx.Value(ContextKeyUserEmail).(string)
	return email, ok
}

// SetTenantInContext adiciona tenant_id ao contexto.
func SetTenantInContext(ctx context.Context, tenantID uuid.UUID) context.Context {
	return context.WithValue(ctx, ContextKeyTenantID, tenantID)
}

// TenantIDFromContext extrai o tenant_id do contexto.
func TenantIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(ContextKeyTenantID).(uuid.UUID)
	return id, ok
}

// SetExternalAuthInContext stores the Zitadel subject and email before UUID resolution.
// Used by AuthMiddleware before UserResolverMiddleware runs.
func SetExternalAuthInContext(ctx context.Context, subject, email string) context.Context {
	ctx = context.WithValue(ctx, ContextKeyExternalAuthID, subject)
	ctx = context.WithValue(ctx, ContextKeyUserEmail, email)
	return ctx
}

// ExternalAuthFromContext retrieves the Zitadel subject and email stored by AuthMiddleware.
func ExternalAuthFromContext(ctx context.Context) (subject, email string, ok bool) {
	subject, ok = ctx.Value(ContextKeyExternalAuthID).(string)
	if !ok || subject == "" {
		return "", "", false
	}
	email, _ = ctx.Value(ContextKeyUserEmail).(string)
	return subject, email, true
}
