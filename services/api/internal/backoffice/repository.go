package backoffice

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository provides backoffice data access.
type Repository struct {
	pool *pgxpool.Pool
}

// NewRepository creates a backoffice repository.
func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// IsBackofficeUser checks if a user has an active backoffice account.
func (r *Repository) IsBackofficeUser(ctx context.Context, userID uuid.UUID) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx,
		`SELECT EXISTS(
			SELECT 1 FROM backoffice_users
			WHERE user_id = $1 AND is_active = TRUE
		)`, userID,
	).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("backoffice: is_user: %w", err)
	}
	return exists, nil
}

// HasBackofficePermission checks if a user has a specific backoffice permission
// via: backoffice_users → backoffice_user_roles → role_permissions → permissions.
func (r *Repository) HasBackofficePermission(ctx context.Context, userID uuid.UUID, permissionCode string) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx,
		`SELECT EXISTS(
			SELECT 1
			FROM backoffice_users bu
			JOIN backoffice_user_roles bur ON bu.id = bur.backoffice_user_id
			JOIN role_permissions rp ON bur.role_id = rp.role_id
			JOIN permissions p ON rp.permission_id = p.id
			WHERE bu.user_id = $1
			  AND bu.is_active = TRUE
			  AND p.code = $2
		)`, userID, permissionCode,
	).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("backoffice: has_permission: %w", err)
	}
	return exists, nil
}

// TenantSummary is used for listing tenants in the backoffice.
type TenantSummary struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// ListTenants returns a paginated list of tenants.
func (r *Repository) ListTenants(ctx context.Context, limit, offset int) ([]TenantSummary, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, name, slug, status, created_at
		 FROM tenants
		 ORDER BY created_at DESC
		 LIMIT $1 OFFSET $2`,
		limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("backoffice: list_tenants: %w", err)
	}
	defer rows.Close()

	var tenants []TenantSummary
	for rows.Next() {
		var t TenantSummary
		if err := rows.Scan(&t.ID, &t.Name, &t.Slug, &t.Status, &t.CreatedAt); err != nil {
			return nil, fmt.Errorf("backoffice: scan_tenant: %w", err)
		}
		tenants = append(tenants, t)
	}
	return tenants, rows.Err()
}

// TenantDetail includes tenant info with subscription details.
type TenantDetail struct {
	ID                     uuid.UUID  `json:"id"`
	Name                   string     `json:"name"`
	Slug                   string     `json:"slug"`
	Status                 string     `json:"status"`
	CreatedAt              time.Time  `json:"created_at"`
	SubscriptionID         *uuid.UUID `json:"subscription_id,omitempty"`
	SubscriptionStatus     *string    `json:"subscription_status,omitempty"`
	PlanCode               *string    `json:"plan_code,omitempty"`
	PlanName               *string    `json:"plan_name,omitempty"`
	ProviderCustomerID     *string    `json:"provider_customer_id,omitempty"`
	ProviderSubscriptionID *string    `json:"provider_subscription_id,omitempty"`
}

// GetTenantDetail returns tenant info joined with subscription and plan.
func (r *Repository) GetTenantDetail(ctx context.Context, tenantID uuid.UUID) (*TenantDetail, error) {
	var d TenantDetail
	err := r.pool.QueryRow(ctx,
		`SELECT t.id, t.name, t.slug, t.status, t.created_at,
		        ts.id, ts.status, sp.code, sp.name,
		        ts.provider_customer_id, ts.provider_subscription_id
		 FROM tenants t
		 LEFT JOIN tenant_subscriptions ts ON ts.tenant_id = t.id
		   AND ts.started_at = (SELECT MAX(started_at) FROM tenant_subscriptions WHERE tenant_id = t.id)
		 LEFT JOIN subscription_plans sp ON ts.plan_id = sp.id
		 WHERE t.id = $1`,
		tenantID,
	).Scan(
		&d.ID, &d.Name, &d.Slug, &d.Status, &d.CreatedAt,
		&d.SubscriptionID, &d.SubscriptionStatus, &d.PlanCode, &d.PlanName,
		&d.ProviderCustomerID, &d.ProviderSubscriptionID,
	)
	if err != nil {
		return nil, fmt.Errorf("backoffice: get_tenant_detail: %w", err)
	}
	return &d, nil
}

// AuditEntry is a simplified audit log entry for backoffice display.
type AuditEntry struct {
	ID          uuid.UUID  `json:"id"`
	ActorUserID *uuid.UUID `json:"actor_user_id,omitempty"`
	ActorScope  string     `json:"actor_scope"`
	EntityType  string     `json:"entity_type"`
	EntityID    *uuid.UUID `json:"entity_id,omitempty"`
	Action      string     `json:"action"`
	Reason      *string    `json:"reason,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

// ListAuditLogs returns paginated audit logs for a tenant.
func (r *Repository) ListAuditLogs(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]AuditEntry, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, actor_user_id, actor_scope, entity_type, entity_id, action, reason, created_at
		 FROM audit_logs
		 WHERE tenant_id = $1
		 ORDER BY created_at DESC
		 LIMIT $2 OFFSET $3`,
		tenantID, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("backoffice: list_audit: %w", err)
	}
	defer rows.Close()

	var entries []AuditEntry
	for rows.Next() {
		var e AuditEntry
		if err := rows.Scan(&e.ID, &e.ActorUserID, &e.ActorScope, &e.EntityType, &e.EntityID,
			&e.Action, &e.Reason, &e.CreatedAt); err != nil {
			return nil, fmt.Errorf("backoffice: scan_audit: %w", err)
		}
		entries = append(entries, e)
	}
	return entries, rows.Err()
}
