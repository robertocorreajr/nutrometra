package audit

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

// Scope define o escopo do ator que gerou o evento.
type Scope string

const (
	ScopeTenant     Scope = "tenant"
	ScopeBackoffice Scope = "backoffice"
	ScopeSystem     Scope = "system"
)

// Executor é satisfeito por *pgxpool.Pool e por pgx.Tx.
// Permite que Write seja chamado dentro de uma transação existente.
type Executor interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
}

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
	IPAddress    string // armazenado como string pós-parse para simplificar callers
	UserAgent    *string
	CreatedAt    time.Time
}

// Option é uma função de configuração de Entry.
type Option func(*Entry)

// NewEntry cria um Entry com ID e timestamp gerados automaticamente.
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

// WithIPAddress aceita endereços com ou sem porta (ex: "1.2.3.4" ou "1.2.3.4:5678").
// Input inválido é ignorado silenciosamente.
func WithIPAddress(ipOrAddr string) Option {
	return func(e *Entry) {
		host := ipOrAddr
		if h, _, err := net.SplitHostPort(ipOrAddr); err == nil {
			host = h
		}
		if parsed := net.ParseIP(host); parsed != nil {
			e.IPAddress = parsed.String()
		}
	}
}

func WithReason(reason string) Option {
	return func(e *Entry) { e.Reason = &reason }
}

func WithUserAgent(ua string) Option {
	return func(e *Entry) { e.UserAgent = &ua }
}

func WithMetadata(v any) Option {
	return func(e *Entry) {
		if b, err := json.Marshal(v); err == nil {
			e.MetadataJSON = b
		}
	}
}

// Service grava eventos de auditoria no banco.
// É stateless — o executor (pool ou tx) é passado em cada chamada de Write
// para permitir participação na mesma transação da operação auditada.
type Service struct{}

func NewService() *Service {
	return &Service{}
}

// Write grava um Entry em audit_logs.
// db pode ser *pgxpool.Pool ou pgx.Tx — use pgx.Tx para garantir atomicidade
// com a operação que está sendo auditada.
func (s *Service) Write(ctx context.Context, db Executor, e Entry) error {
	if db == nil {
		return fmt.Errorf("audit: executor is nil")
	}
	if e.ActorScope == "" {
		return fmt.Errorf("audit: actor_scope is required")
	}

	var ipArg *string
	if e.IPAddress != "" {
		ipArg = &e.IPAddress
	}

	var metaArg *[]byte
	if len(e.MetadataJSON) > 0 {
		metaArg = &e.MetadataJSON
	}

	_, err := db.Exec(ctx,
		`INSERT INTO audit_logs
			(id, tenant_id, actor_user_id, actor_scope,
			 entity_type, entity_id, action, reason,
			 metadata_json, ip_address, user_agent, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10::inet,$11,$12)`,
		e.ID, e.TenantID, e.ActorUserID, string(e.ActorScope),
		e.EntityType, e.EntityID, e.Action, e.Reason,
		metaArg, ipArg, e.UserAgent, e.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("audit: write: %w", err)
	}
	return nil
}
