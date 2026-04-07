package repository

import (
	"context"
	"fmt"

	"nutrometra/api/internal/rbac/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Enforcer checks whether a user has a specific permission within a tenant.
type Enforcer interface {
	HasPermission(ctx context.Context, userID, tenantID uuid.UUID, permissionCode string) (bool, error)
}

// Repository provides RBAC data access.
type Repository struct {
	pool *pgxpool.Pool
}

// New creates a Repository.
func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// HasPermission checks if a user has a permission in a tenant via the role chain:
// tenant_users -> tenant_user_roles -> roles -> role_permissions -> permissions.
// Only active tenant_users are considered.
func (r *Repository) HasPermission(ctx context.Context, userID, tenantID uuid.UUID, permissionCode string) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx, `
		SELECT EXISTS(
			SELECT 1
			FROM tenant_user_roles tur
			JOIN tenant_users tu ON tur.tenant_user_id = tu.id
			JOIN role_permissions rp ON tur.role_id = rp.role_id
			JOIN permissions p ON rp.permission_id = p.id
			WHERE tu.user_id = $1
			  AND tu.tenant_id = $2
			  AND tu.status = 'active'
			  AND p.code = $3
		)
	`, userID, tenantID, permissionCode).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("rbac: has_permission: %w", err)
	}
	return exists, nil
}

// AssignRole assigns a role (by code) to a tenant_user.
// Idempotent via ON CONFLICT DO NOTHING.
func (r *Repository) AssignRole(ctx context.Context, tenantID, tenantUserID uuid.UUID, roleCode string) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO tenant_user_roles (id, tenant_id, tenant_user_id, role_id)
		SELECT gen_random_uuid(), $1, $2, r.id
		FROM roles r
		WHERE r.code = $3
		ON CONFLICT (tenant_user_id, role_id) DO NOTHING
	`, tenantID, tenantUserID, roleCode)
	if err != nil {
		return fmt.Errorf("rbac: assign_role: %w", err)
	}
	return nil
}

// RevokeRole removes a role assignment from a tenant_user.
func (r *Repository) RevokeRole(ctx context.Context, tenantUserID, roleID uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`DELETE FROM tenant_user_roles WHERE tenant_user_id = $1 AND role_id = $2`,
		tenantUserID, roleID,
	)
	if err != nil {
		return fmt.Errorf("rbac: revoke_role: %w", err)
	}
	return nil
}

// ListRoles returns all available roles.
func (r *Repository) ListRoles(ctx context.Context) ([]domain.Role, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, code, application_scope, name, description FROM roles ORDER BY code`,
	)
	if err != nil {
		return nil, fmt.Errorf("rbac: list_roles: %w", err)
	}
	defer rows.Close()

	var roles []domain.Role
	for rows.Next() {
		var role domain.Role
		if err := rows.Scan(&role.ID, &role.Code, &role.ApplicationScope, &role.Name, &role.Description); err != nil {
			return nil, fmt.Errorf("rbac: scan_role: %w", err)
		}
		roles = append(roles, role)
	}
	return roles, rows.Err()
}
