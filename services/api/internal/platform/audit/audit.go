package audit

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Scope define o escopo do ator que gerou o evento.
type Scope string

const (
	ScopeTenant     Scope = "tenant"
	ScopeBackoffice Scope = "backoffice"
	ScopeSystem     Scope = "system"
)

// Entry representa um evento de auditoria.
type Entry struct {
	ID           uuid.UUID
	TenantID     *uuid.UUID
	ActorUserID  *uuid.UUID
	ActorScope   Scope
	EntityType   string
	EntityID     *uuid.UUID
	Action       string
	Reason       *string
	MetadataJSON []byte
	IPAddress    *net.IP
	UserAgent    *string
	CreatedAt    time.Time
}

// Option é uma função de configuração de Entry.
type Option func(*Entry)

// NewEntry cria um Entry com ID e timestamp gerados.
func NewEntry(opts ...Option) Entry {
	e := Entry{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
	}
	for _, o := range opts {
		o(&e)
	}
	return e
}

func WithTenantID(id uuid.UUID) Option {
	return func(e *Entry) { e.TenantID = &id }
}

func WithActor(id uuid.UUID, scope Scope) Option {
	return func(e *Entry) {
		e.ActorUserID = &id
		e.ActorScope = scope
	}
}

func WithEntity(entityType string, id uuid.UUID) Option {
	return func(e *Entry) {
		e.EntityType = entityType
		e.EntityID = &id
	}
}

func WithAction(action string) Option {
	return func(e *Entry) { e.Action = action }
}

func WithIPAddress(ip string) Option {
	return func(e *Entry) {
		parsed := net.ParseIP(ip)
		if parsed != nil {
			e.IPAddress = &parsed
		}
	}
}

func WithReason(reason string) Option {
	return func(e *Entry) { e.Reason = &reason }
}

// Service grava eventos de auditoria no banco.
type Service struct {
	pool *pgxpool.Pool
}

func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}

// Write grava um Entry em audit_logs. Deve ser chamado dentro da mesma transação da operação.
func (s *Service) Write(ctx context.Context, e Entry) error {
	if s.pool == nil {
		return fmt.Errorf("audit: db pool is nil")
	}

	ipStr := (*string)(nil)
	if e.IPAddress != nil {
		str := e.IPAddress.String()
		ipStr = &str
	}

	_, err := s.pool.Exec(ctx,
		`INSERT INTO audit_logs
			(id, tenant_id, actor_user_id, actor_scope, entity_type, entity_id, action, reason, ip_address, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9::inet,$10)`,
		e.ID, e.TenantID, e.ActorUserID, string(e.ActorScope),
		e.EntityType, e.EntityID, e.Action, e.Reason,
		ipStr, e.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("audit: write: %w", err)
	}
	return nil
}
