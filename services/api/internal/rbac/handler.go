package rbac

import (
	"encoding/json"
	"errors"
	"net/http"

	identitydomain "nutrometra/api/internal/identity/domain"
	"nutrometra/api/internal/platform/audit"
	"nutrometra/api/internal/platform/server"
	"nutrometra/api/internal/rbac/domain"
	"nutrometra/api/internal/rbac/repository"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Handler exposes RBAC HTTP endpoints.
type Handler struct {
	repo     *repository.Repository
	pool     *pgxpool.Pool // needed as Executor for audit writes
	auditSvc *audit.Service
}

// NewHandler creates an RBAC handler.
func NewHandler(repo *repository.Repository, pool *pgxpool.Pool, auditSvc *audit.Service) *Handler {
	return &Handler{repo: repo, pool: pool, auditSvc: auditSvc}
}

type roleResponse struct {
	ID               string `json:"id"`
	Code             string `json:"code"`
	ApplicationScope string `json:"application_scope"`
	Name             string `json:"name"`
	Description      string `json:"description"`
}

func toRoleResponse(r domain.Role) roleResponse {
	return roleResponse{
		ID:               r.ID.String(),
		Code:             r.Code,
		ApplicationScope: r.ApplicationScope,
		Name:             r.Name,
		Description:      r.Description,
	}
}

// ListRoles returns available roles filtered by application_scope query param.
// GET /roles?scope=tenant (default: "tenant")
func (h *Handler) ListRoles(w http.ResponseWriter, r *http.Request) {
	scope := r.URL.Query().Get("scope")
	if scope == "" {
		scope = "tenant"
	}
	// Validate scope
	if scope != "tenant" && scope != "backoffice" && scope != "system" {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_scope", "scope must be tenant, backoffice, or system")
		return
	}

	roles, err := h.repo.ListRoles(r.Context(), scope)
	if err != nil {
		server.RenderError(w, r, http.StatusInternalServerError, "list_roles_failed", "Failed to list roles")
		return
	}

	resp := make([]roleResponse, len(roles))
	for i, role := range roles {
		resp[i] = toRoleResponse(role)
	}
	server.RenderJSON(w, http.StatusOK, resp)
}

type assignRoleRequest struct {
	RoleCode string `json:"role_code"`
}

// AssignRole assigns a role to a tenant member.
// POST /members/{member_id}/roles
// Requires tenant in context (from TenantMiddleware).
// member_id is the tenant_user_id from tenant_users table.
func (h *Handler) AssignRole(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := identitydomain.TenantIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusBadRequest, "missing_tenant", "No tenant in context")
		return
	}

	memberID, err := uuid.Parse(chi.URLParam(r, "member_id"))
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_member_id", "Invalid member ID")
		return
	}

	var req assignRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.RoleCode == "" {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_body", "role_code is required")
		return
	}

	if err := h.repo.AssignRole(r.Context(), tenantID, memberID, req.RoleCode); err != nil {
		if errors.Is(err, repository.ErrRoleNotFound) {
			server.RenderError(w, r, http.StatusUnprocessableEntity, "role_not_found", "Role not found")
			return
		}
		server.RenderError(w, r, http.StatusInternalServerError, "assign_role_failed", "Failed to assign role")
		return
	}

	// Audit log
	actorID, _ := identitydomain.UserIDFromContext(r.Context())
	entry := audit.NewEntry(
		audit.WithTenantID(tenantID),
		audit.WithActor(actorID, audit.ScopeTenant),
		audit.WithEntity("tenant_user_role", memberID),
		audit.WithAction("role_assigned"),
		audit.WithIPAddress(r.RemoteAddr),
	)
	_ = h.auditSvc.Write(r.Context(), h.pool, entry)

	w.WriteHeader(http.StatusNoContent)
}

// RevokeRole removes a role from a tenant member.
// DELETE /members/{member_id}/roles/{role_id}
// Requires tenant in context (from TenantMiddleware).
func (h *Handler) RevokeRole(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := identitydomain.TenantIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusBadRequest, "missing_tenant", "No tenant in context")
		return
	}

	memberID, err := uuid.Parse(chi.URLParam(r, "member_id"))
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_member_id", "Invalid member ID")
		return
	}

	roleID, err := uuid.Parse(chi.URLParam(r, "role_id"))
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_role_id", "Invalid role ID")
		return
	}

	if err := h.repo.RevokeRole(r.Context(), tenantID, memberID, roleID); err != nil {
		server.RenderError(w, r, http.StatusInternalServerError, "revoke_role_failed", "Failed to revoke role")
		return
	}

	// Audit log
	actorID, _ := identitydomain.UserIDFromContext(r.Context())
	entry := audit.NewEntry(
		audit.WithTenantID(tenantID),
		audit.WithActor(actorID, audit.ScopeTenant),
		audit.WithEntity("tenant_user_role", memberID),
		audit.WithAction("role_revoked"),
		audit.WithIPAddress(r.RemoteAddr),
	)
	_ = h.auditSvc.Write(r.Context(), h.pool, entry)

	w.WriteHeader(http.StatusNoContent)
}
