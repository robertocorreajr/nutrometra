package google

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"nutrometra/api/internal/platform/queue"

	"github.com/google/uuid"
	"golang.org/x/oauth2"
	googleoauth "golang.org/x/oauth2/google"
)

// SyncPayload is the job payload for calendar sync operations.
type SyncPayload struct {
	TenantID   uuid.UUID `json:"tenant_id"`
	UserID     uuid.UUID `json:"user_id"`
	EntityType string    `json:"entity_type"`
	EntityID   uuid.UUID `json:"entity_id"`
	Operation  string    `json:"operation"` // create, update, delete

	// Event details for create/update operations.
	Summary     string    `json:"summary,omitempty"`
	Description string    `json:"description,omitempty"`
	StartTime   time.Time `json:"start_time,omitempty"`
	EndTime     time.Time `json:"end_time,omitempty"`
	Location    string    `json:"location,omitempty"`
	TimeZone    string    `json:"time_zone,omitempty"`
}

// SyncHandler processes calendar_sync jobs from the background queue.
type SyncHandler struct {
	repo          *Repository
	provider      CalendarProvider
	encryptionKey string
}

// NewSyncHandler creates a new SyncHandler.
func NewSyncHandler(repo *Repository, provider CalendarProvider, encryptionKey string) queue.HandlerFunc {
	h := &SyncHandler{
		repo:          repo,
		provider:      provider,
		encryptionKey: encryptionKey,
	}
	return h.Handle
}

// Handle processes a single calendar_sync job.
func (h *SyncHandler) Handle(ctx context.Context, job *queue.JobRecord) error {
	var payload SyncPayload
	if err := json.Unmarshal(job.PayloadJSON, &payload); err != nil {
		return fmt.Errorf("calendar sync: unmarshal payload: %w", err)
	}

	log := slog.With(
		"job_id", job.ID,
		"tenant_id", payload.TenantID,
		"user_id", payload.UserID,
		"entity_type", payload.EntityType,
		"entity_id", payload.EntityID,
		"operation", payload.Operation,
	)

	// Look up the calendar connection.
	conn, err := h.repo.GetConnectionByUserAndProvider(ctx, payload.TenantID, payload.UserID, providerGoogle)
	if err != nil {
		return fmt.Errorf("calendar sync: get connection: %w", err)
	}
	if conn == nil || conn.Status != "active" {
		log.Warn("calendar sync: no active connection, skipping")
		return nil
	}

	// Create sync log entry.
	syncLog := CalendarSyncLog{
		ID:                   uuid.New(),
		TenantID:             payload.TenantID,
		CalendarConnectionID: conn.ID,
		EntityType:           payload.EntityType,
		EntityID:             payload.EntityID,
		Operation:            payload.Operation,
		Status:               "pending",
		CreatedAt:            time.Now().UTC(),
	}
	if err := h.repo.CreateSyncLog(ctx, syncLog); err != nil {
		return fmt.Errorf("calendar sync: create sync log: %w", err)
	}

	// Decrypt tokens.
	accessToken, refreshToken, err := h.repo.DecryptConnectionTokens(conn)
	if err != nil {
		errMsg := "failed to decrypt tokens"
		_ = h.repo.UpdateSyncLogStatus(ctx, syncLog.ID, "failed", &errMsg, nil)
		return fmt.Errorf("calendar sync: decrypt tokens: %w", err)
	}

	// If token expired, attempt refresh.
	if conn.TokenExpiresAt != nil && conn.TokenExpiresAt.Before(time.Now()) {
		log.Info("calendar sync: token expired, attempting refresh")
		newToken, refreshErr := h.refreshToken(ctx, refreshToken)
		if refreshErr != nil {
			errMsg := fmt.Sprintf("token refresh failed: %v", refreshErr)
			_ = h.repo.UpdateSyncLogStatus(ctx, syncLog.ID, "failed", &errMsg, nil)
			return fmt.Errorf("calendar sync: refresh token: %w", refreshErr)
		}
		accessToken = newToken.AccessToken
		if newToken.RefreshToken != "" {
			refreshToken = newToken.RefreshToken
		}
		// Update stored tokens.
		if err := h.repo.UpdateConnectionTokens(ctx, conn.ID, accessToken, refreshToken, newToken.Expiry); err != nil {
			log.Error("calendar sync: failed to update refreshed tokens", "error", err)
		}
	}

	// Execute the calendar operation.
	var providerEventID *string
	switch payload.Operation {
	case "create":
		event := CalendarEvent{
			Summary:     payload.Summary,
			Description: payload.Description,
			StartTime:   payload.StartTime,
			EndTime:     payload.EndTime,
			Location:    payload.Location,
			TimeZone:    payload.TimeZone,
		}
		extID, createErr := h.provider.CreateEvent(ctx, accessToken, event)
		if createErr != nil {
			errMsg := fmt.Sprintf("create event failed: %v", createErr)
			_ = h.repo.UpdateSyncLogStatus(ctx, syncLog.ID, "failed", &errMsg, nil)
			return fmt.Errorf("calendar sync: create event: %w", createErr)
		}
		providerEventID = &extID
		log.Info("calendar sync: event created", "provider_event_id", extID)

	case "update":
		// Look up the last synced event ID for this entity.
		extEventID := h.findProviderEventID(ctx, conn.ID, payload.EntityType, payload.EntityID)
		if extEventID == "" {
			errMsg := "no existing provider event ID found for update"
			_ = h.repo.UpdateSyncLogStatus(ctx, syncLog.ID, "failed", &errMsg, nil)
			return fmt.Errorf("calendar sync: no provider event ID for entity %s/%s", payload.EntityType, payload.EntityID)
		}
		event := CalendarEvent{
			Summary:     payload.Summary,
			Description: payload.Description,
			StartTime:   payload.StartTime,
			EndTime:     payload.EndTime,
			Location:    payload.Location,
			TimeZone:    payload.TimeZone,
		}
		if updateErr := h.provider.UpdateEvent(ctx, accessToken, extEventID, event); updateErr != nil {
			errMsg := fmt.Sprintf("update event failed: %v", updateErr)
			_ = h.repo.UpdateSyncLogStatus(ctx, syncLog.ID, "failed", &errMsg, nil)
			return fmt.Errorf("calendar sync: update event: %w", updateErr)
		}
		providerEventID = &extEventID
		log.Info("calendar sync: event updated", "provider_event_id", extEventID)

	case "delete":
		extEventID := h.findProviderEventID(ctx, conn.ID, payload.EntityType, payload.EntityID)
		if extEventID == "" {
			log.Warn("calendar sync: no provider event ID found for delete, skipping")
			status := "synced"
			_ = h.repo.UpdateSyncLogStatus(ctx, syncLog.ID, status, nil, nil)
			return nil
		}
		if deleteErr := h.provider.DeleteEvent(ctx, accessToken, extEventID); deleteErr != nil {
			errMsg := fmt.Sprintf("delete event failed: %v", deleteErr)
			_ = h.repo.UpdateSyncLogStatus(ctx, syncLog.ID, "failed", &errMsg, nil)
			return fmt.Errorf("calendar sync: delete event: %w", deleteErr)
		}
		providerEventID = &extEventID
		log.Info("calendar sync: event deleted", "provider_event_id", extEventID)

	default:
		errMsg := fmt.Sprintf("unknown operation: %s", payload.Operation)
		_ = h.repo.UpdateSyncLogStatus(ctx, syncLog.ID, "failed", &errMsg, nil)
		return fmt.Errorf("calendar sync: unknown operation: %s", payload.Operation)
	}

	// Mark as synced.
	_ = h.repo.UpdateSyncLogStatus(ctx, syncLog.ID, "synced", nil, providerEventID)
	return nil
}

// refreshToken uses the OAuth2 refresh token to obtain a new access token.
func (h *SyncHandler) refreshToken(ctx context.Context, refreshToken string) (*oauth2.Token, error) {
	// Build a minimal oauth2.Config just for token refresh.
	// The client ID/secret are needed for refresh, but they are stored in the handler config
	// which is not directly available here. Instead we use the token source approach.
	ts := (&oauth2.Config{
		Endpoint: googleoauth.Endpoint,
	}).TokenSource(ctx, &oauth2.Token{
		RefreshToken: refreshToken,
	})
	newToken, err := ts.Token()
	if err != nil {
		return nil, fmt.Errorf("refresh token: %w", err)
	}
	return newToken, nil
}

// findProviderEventID looks up the most recent successfully synced provider event ID
// for the given entity.
func (h *SyncHandler) findProviderEventID(ctx context.Context, connectionID uuid.UUID, entityType string, entityID uuid.UUID) string {
	var providerEventID *string
	err := h.repo.pool.QueryRow(ctx,
		`SELECT provider_event_id FROM calendar_sync_log
		WHERE calendar_connection_id = $1 AND entity_type = $2 AND entity_id = $3
			AND status = 'synced' AND provider_event_id IS NOT NULL
		ORDER BY created_at DESC LIMIT 1`,
		connectionID, entityType, entityID,
	).Scan(&providerEventID)
	if err != nil || providerEventID == nil {
		return ""
	}
	return *providerEventID
}
