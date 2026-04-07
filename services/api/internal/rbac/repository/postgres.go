package repository

import (
	"context"
	"errors"
	"fmt"

	"nutrometra/api/internal/rbac/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ErrRoleNotFound is returned when a role code does not match any existing role.
var ErrRoleNotFound = errors.New("rbac: role not found")

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
// The SQL verifies that tenant_user_id actually belongs to tenant_id (defense-in-depth).
// Idempotent via ON CONFLICT DO NOTHING.
// Returns ErrRoleNotFound if the roleCode does not match any role or the tenant_user
// does not belong to the specified tenant.
func (r *Repository) AssignRole(ctx context.Context, tenantID, tenantUserID uuid.UUID, roleCode string) error {
	tag, err := r.pool.Exec(ctx, `
		INSERT INTO tenant_user_roles (id, tenant_id, tenant_user_id, role_id)
		SELECT gen_random_uuid(), tu.tenant_id, tu.id, r.id
		FROM tenant_users tu
		CROSS JOIN roles r
		WHERE tu.id = $1
		  AND tu.tenant_id = $2
		  AND r.code = $3
		ON CONFLICT (tenant_user_id, role_id) DO NOTHING
	`, tenantUserID, tenantID, roleCode)
	if err != nil {
		return fmt.Errorf("rbac: assign_role: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrRoleNotFound
	}
	return nil
}

// RevokeRole removes a role assignment from a tenant_user.
// Includes tenant_id in the WHERE clause for defense-in-depth tenant isolation.
func (r *Repository) RevokeRole(ctx context.Context, tenantID, tenantUserID, roleID uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`DELETE FROM tenant_user_roles WHERE tenant_id = $1 AND tenant_user_id = $2 AND role_id = $3`,
		tenantID, tenantUserID, roleID,
	)
	if err != nil {
		return fmt.Errorf("rbac: revoke_role: %w", err)
	}
	return nil
}

// ListRoles returns roles filtered by application scope.
func (r *Repository) ListRoles(ctx context.Context, applicationScope string) ([]domain.Role, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, code, application_scope, name, description FROM roles WHERE application_scope = $1 ORDER BY code`,
		applicationScope,
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
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rbac: list_roles: %w", err)
	}
	return roles, nil
}
