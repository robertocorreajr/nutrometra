package google

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	identitydomain "nutrometra/api/internal/identity/domain"
	"nutrometra/api/internal/platform/audit"
	"nutrometra/api/internal/platform/config"
	"nutrometra/api/internal/platform/queue"
	"nutrometra/api/internal/platform/server"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/oauth2"
	googleoauth "golang.org/x/oauth2/google"
)

const (
	calendarEventsScope = "https://www.googleapis.com/auth/calendar.events"
	providerGoogle      = "google"
)

// Handler exposes HTTP endpoints for Google Calendar OAuth2 integration.
type Handler struct {
	pool          *pgxpool.Pool
	repo          *Repository
	provider      CalendarProvider
	auditSvc      *audit.Service
	oauthConfig   *oauth2.Config
	encryptionKey string
	jobQueue      queue.Queue
}

// NewHandler creates a new Google Calendar integration Handler.
func NewHandler(
	pool *pgxpool.Pool,
	repo *Repository,
	calProvider CalendarProvider,
	auditSvc *audit.Service,
	googleCfg config.GoogleConfig,
	jobQueue queue.Queue,
) *Handler {
	oauthCfg := &oauth2.Config{
		ClientID:     googleCfg.ClientID,
		ClientSecret: googleCfg.ClientSecret,
		RedirectURL:  googleCfg.RedirectURL,
		Scopes:       []string{calendarEventsScope},
		Endpoint:     googleoauth.Endpoint,
	}
	return &Handler{
		pool:          pool,
		repo:          repo,
		provider:      calProvider,
		auditSvc:      auditSvc,
		oauthConfig:   oauthCfg,
		encryptionKey: googleCfg.EncryptionKey,
		jobQueue:      jobQueue,
	}
}

// generateState creates an HMAC-signed state parameter containing tenantID and userID.
// Format: "tenantID:userID:hmac"
func (h *Handler) generateState(tenantID, userID uuid.UUID) string {
	data := tenantID.String() + ":" + userID.String()
	mac := hmac.New(sha256.New, []byte(h.encryptionKey))
	mac.Write([]byte(data))
	sig := hex.EncodeToString(mac.Sum(nil))
	return data + ":" + sig
}

// verifyState validates the HMAC-signed state and returns tenantID and userID.
func (h *Handler) verifyState(state string) (uuid.UUID, uuid.UUID, error) {
	parts := strings.SplitN(state, ":", 3)
	if len(parts) != 3 {
		return uuid.Nil, uuid.Nil, fmt.Errorf("invalid state format")
	}
	data := parts[0] + ":" + parts[1]
	mac := hmac.New(sha256.New, []byte(h.encryptionKey))
	mac.Write([]byte(data))
	expectedSig := hex.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(parts[2]), []byte(expectedSig)) {
		return uuid.Nil, uuid.Nil, fmt.Errorf("state signature mismatch")
	}
	tenantID, err := uuid.Parse(parts[0])
	if err != nil {
		return uuid.Nil, uuid.Nil, fmt.Errorf("invalid tenant_id in state: %w", err)
	}
	userID, err := uuid.Parse(parts[1])
	if err != nil {
		return uuid.Nil, uuid.Nil, fmt.Errorf("invalid user_id in state: %w", err)
	}
	return tenantID, userID, nil
}

// Authorize generates a Google OAuth2 consent URL and redirects the user.
// GET /integrations/google/authorize
func (h *Handler) Authorize(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := identitydomain.TenantIDFromContext(ctx)
	if !ok {
		server.RenderError(w, r, http.StatusUnauthorized, "missing_tenant", "Tenant ID required")
		return
	}
	userID, ok := identitydomain.UserIDFromContext(ctx)
	if !ok {
		server.RenderError(w, r, http.StatusUnauthorized, "missing_user", "User ID required")
		return
	}

	state := h.generateState(tenantID, userID)
	url := h.oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.SetAuthURLParam("prompt", "consent"))

	server.RenderJSON(w, http.StatusOK, map[string]string{
		"authorize_url": url,
	})
}

// Callback exchanges the authorization code for tokens and stores them encrypted.
// GET /integrations/google/callback
func (h *Handler) Callback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	stateParam := r.URL.Query().Get("state")
	if stateParam == "" {
		server.RenderError(w, r, http.StatusBadRequest, "missing_state", "State parameter required")
		return
	}

	tenantID, userID, err := h.verifyState(stateParam)
	if err != nil {
		slog.Warn("google callback: invalid state", "error", err)
		server.RenderError(w, r, http.StatusBadRequest, "invalid_state", "Invalid or tampered state parameter")
		return
	}

	code := r.URL.Query().Get("code")
	if code == "" {
		server.RenderError(w, r, http.StatusBadRequest, "missing_code", "Authorization code required")
		return
	}

	token, err := h.oauthConfig.Exchange(ctx, code)
	if err != nil {
		slog.Error("google callback: token exchange failed", "error", err)
		server.RenderError(w, r, http.StatusBadGateway, "token_exchange_failed", "Failed to exchange authorization code")
		return
	}

	now := time.Now().UTC()
	conn := CalendarConnection{
		ID:                    uuid.New(),
		TenantID:              tenantID,
		UserID:                userID,
		Provider:              providerGoogle,
		Status:                "active",
		EncryptedAccessToken:  token.AccessToken,
		EncryptedRefreshToken: token.RefreshToken,
		TokenExpiresAt:        &token.Expiry,
		Scopes:                calendarEventsScope,
		CreatedAt:             now,
		UpdatedAt:             now,
	}

	if err := h.repo.CreateConnection(ctx, conn); err != nil {
		slog.Error("google callback: store connection failed", "error", err)
		server.RenderError(w, r, http.StatusInternalServerError, "store_connection_failed", "Failed to store calendar connection")
		return
	}

	// Audit
	entry := audit.NewEntry(
		audit.WithTenantID(tenantID),
		audit.WithActor(userID, audit.ScopeTenant),
		audit.WithEntity("calendar_connection", conn.ID),
		audit.WithAction("google_calendar_connected"),
		audit.WithIPAddress(r.RemoteAddr),
	)
	if err := h.auditSvc.Write(ctx, h.pool, entry); err != nil {
		slog.Error("google callback: audit write failed", "error", err)
	}

	server.RenderJSON(w, http.StatusOK, map[string]string{
		"status":  "connected",
		"message": "Google Calendar connected successfully",
	})
}

// Disconnect revokes the Google Calendar connection for the current user.
// POST /integrations/google/disconnect
func (h *Handler) Disconnect(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := identitydomain.TenantIDFromContext(ctx)
	if !ok {
		server.RenderError(w, r, http.StatusUnauthorized, "missing_tenant", "Tenant ID required")
		return
	}
	userID, ok := identitydomain.UserIDFromContext(ctx)
	if !ok {
		server.RenderError(w, r, http.StatusUnauthorized, "missing_user", "User ID required")
		return
	}

	conn, err := h.repo.GetConnectionByUserAndProvider(ctx, tenantID, userID, providerGoogle)
	if err != nil {
		slog.Error("google disconnect: lookup failed", "error", err)
		server.RenderError(w, r, http.StatusInternalServerError, "lookup_failed", "Failed to look up connection")
		return
	}
	if conn == nil {
		server.RenderError(w, r, http.StatusNotFound, "not_connected", "No Google Calendar connection found")
		return
	}

	if err := h.repo.RevokeConnection(ctx, conn.ID); err != nil {
		slog.Error("google disconnect: revoke failed", "error", err)
		server.RenderError(w, r, http.StatusInternalServerError, "revoke_failed", "Failed to revoke connection")
		return
	}

	// Audit
	entry := audit.NewEntry(
		audit.WithTenantID(tenantID),
		audit.WithActor(userID, audit.ScopeTenant),
		audit.WithEntity("calendar_connection", conn.ID),
		audit.WithAction("google_calendar_disconnected"),
		audit.WithIPAddress(r.RemoteAddr),
	)
	if err := h.auditSvc.Write(ctx, h.pool, entry); err != nil {
		slog.Error("google disconnect: audit write failed", "error", err)
	}

	server.RenderJSON(w, http.StatusOK, map[string]string{
		"status":  "disconnected",
		"message": "Google Calendar disconnected successfully",
	})
}

type statusResponse struct {
	Connected bool    `json:"connected"`
	Provider  string  `json:"provider"`
	Status    string  `json:"status"`
	SyncedAt  *string `json:"synced_at,omitempty"`
}

// Status returns the Google Calendar connection status for the current user.
// GET /integrations/google/status
func (h *Handler) Status(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := identitydomain.TenantIDFromContext(ctx)
	if !ok {
		server.RenderError(w, r, http.StatusUnauthorized, "missing_tenant", "Tenant ID required")
		return
	}
	userID, ok := identitydomain.UserIDFromContext(ctx)
	if !ok {
		server.RenderError(w, r, http.StatusUnauthorized, "missing_user", "User ID required")
		return
	}

	conn, err := h.repo.GetConnectionByUserAndProvider(ctx, tenantID, userID, providerGoogle)
	if err != nil {
		slog.Error("google status: lookup failed", "error", err)
		server.RenderError(w, r, http.StatusInternalServerError, "lookup_failed", "Failed to look up connection")
		return
	}

	if conn == nil {
		server.RenderJSON(w, http.StatusOK, statusResponse{
			Connected: false,
			Provider:  providerGoogle,
			Status:    "not_connected",
		})
		return
	}

	resp := statusResponse{
		Connected: conn.Status == "active",
		Provider:  conn.Provider,
		Status:    conn.Status,
	}
	if conn.SyncedAt != nil {
		formatted := conn.SyncedAt.Format(time.RFC3339)
		resp.SyncedAt = &formatted
	}
	server.RenderJSON(w, http.StatusOK, resp)
}
