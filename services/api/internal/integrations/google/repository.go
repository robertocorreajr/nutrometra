package google

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// CalendarConnection represents a stored OAuth2 connection to a calendar provider.
type CalendarConnection struct {
	ID                    uuid.UUID
	TenantID              uuid.UUID
	UserID                uuid.UUID
	Provider              string
	Status                string
	ExternalAccountID     *string
	EncryptedAccessToken  string
	EncryptedRefreshToken string
	TokenExpiresAt        *time.Time
	Scopes                string
	SyncedAt              *time.Time
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

// CalendarSyncLog represents a sync operation record.
type CalendarSyncLog struct {
	ID                   uuid.UUID
	TenantID             uuid.UUID
	CalendarConnectionID uuid.UUID
	EntityType           string
	EntityID             uuid.UUID
	ProviderEventID      *string
	Operation            string
	Status               string
	ErrorMessage         *string
	CreatedAt            time.Time
}

// Repository handles persistence for calendar connections and sync logs.
type Repository struct {
	pool          *pgxpool.Pool
	encryptionKey string
}

// NewRepository creates a new calendar Repository.
func NewRepository(pool *pgxpool.Pool, encryptionKey string) *Repository {
	return &Repository{pool: pool, encryptionKey: encryptionKey}
}

// CreateConnection inserts a new calendar connection with encrypted tokens.
func (r *Repository) CreateConnection(ctx context.Context, conn CalendarConnection) error {
	encAccess, err := EncryptToken(conn.EncryptedAccessToken, r.encryptionKey)
	if err != nil {
		return fmt.Errorf("calendar repo: encrypt access token: %w", err)
	}
	encRefresh, err := EncryptToken(conn.EncryptedRefreshToken, r.encryptionKey)
	if err != nil {
		return fmt.Errorf("calendar repo: encrypt refresh token: %w", err)
	}

	_, err = r.pool.Exec(ctx,
		`INSERT INTO calendar_connections
			(id, tenant_id, user_id, provider, status, external_account_id,
			 encrypted_access_token, encrypted_refresh_token, token_expires_at,
			 scopes, synced_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		ON CONFLICT (tenant_id, user_id, provider)
		DO UPDATE SET
			status = EXCLUDED.status,
			external_account_id = EXCLUDED.external_account_id,
			encrypted_access_token = EXCLUDED.encrypted_access_token,
			encrypted_refresh_token = EXCLUDED.encrypted_refresh_token,
			token_expires_at = EXCLUDED.token_expires_at,
			scopes = EXCLUDED.scopes,
			updated_at = EXCLUDED.updated_at`,
		conn.ID, conn.TenantID, conn.UserID, conn.Provider, conn.Status,
		conn.ExternalAccountID, encAccess, encRefresh, conn.TokenExpiresAt,
		conn.Scopes, conn.SyncedAt, conn.CreatedAt, conn.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("calendar repo: create connection: %w", err)
	}
	return nil
}

// GetConnectionByUserAndProvider returns the calendar connection for a given user and provider.
// Tokens are returned encrypted — caller must decrypt as needed.
func (r *Repository) GetConnectionByUserAndProvider(ctx context.Context, tenantID, userID uuid.UUID, provider string) (*CalendarConnection, error) {
	var conn CalendarConnection
	err := r.pool.QueryRow(ctx,
		`SELECT id, tenant_id, user_id, provider, status, external_account_id,
			encrypted_access_token, encrypted_refresh_token, token_expires_at,
			scopes, synced_at, created_at, updated_at
		FROM calendar_connections
		WHERE tenant_id = $1 AND user_id = $2 AND provider = $3`,
		tenantID, userID, provider,
	).Scan(
		&conn.ID, &conn.TenantID, &conn.UserID, &conn.Provider, &conn.Status,
		&conn.ExternalAccountID, &conn.EncryptedAccessToken, &conn.EncryptedRefreshToken,
		&conn.TokenExpiresAt, &conn.Scopes, &conn.SyncedAt, &conn.CreatedAt, &conn.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("calendar repo: get connection: %w", err)
	}
	return &conn, nil
}

// UpdateConnectionTokens updates the encrypted tokens and expiry for a connection.
func (r *Repository) UpdateConnectionTokens(ctx context.Context, connID uuid.UUID, accessToken, refreshToken string, expiresAt time.Time) error {
	encAccess, err := EncryptToken(accessToken, r.encryptionKey)
	if err != nil {
		return fmt.Errorf("calendar repo: encrypt access token: %w", err)
	}
	encRefresh, err := EncryptToken(refreshToken, r.encryptionKey)
	if err != nil {
		return fmt.Errorf("calendar repo: encrypt refresh token: %w", err)
	}

	_, err = r.pool.Exec(ctx,
		`UPDATE calendar_connections
		SET encrypted_access_token = $1, encrypted_refresh_token = $2,
			token_expires_at = $3, status = 'active', updated_at = NOW()
		WHERE id = $4`,
		encAccess, encRefresh, expiresAt, connID,
	)
	if err != nil {
		return fmt.Errorf("calendar repo: update tokens: %w", err)
	}
	return nil
}

// RevokeConnection sets the connection status to 'revoked'.
func (r *Repository) RevokeConnection(ctx context.Context, connID uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE calendar_connections SET status = 'revoked', updated_at = NOW() WHERE id = $1`,
		connID,
	)
	if err != nil {
		return fmt.Errorf("calendar repo: revoke connection: %w", err)
	}
	return nil
}

// DecryptConnectionTokens decrypts access and refresh tokens from an encrypted connection.
func (r *Repository) DecryptConnectionTokens(conn *CalendarConnection) (accessToken, refreshToken string, err error) {
	accessToken, err = DecryptToken(conn.EncryptedAccessToken, r.encryptionKey)
	if err != nil {
		return "", "", fmt.Errorf("calendar repo: decrypt access token: %w", err)
	}
	refreshToken, err = DecryptToken(conn.EncryptedRefreshToken, r.encryptionKey)
	if err != nil {
		return "", "", fmt.Errorf("calendar repo: decrypt refresh token: %w", err)
	}
	return accessToken, refreshToken, nil
}

// CreateSyncLog inserts a new sync log entry.
func (r *Repository) CreateSyncLog(ctx context.Context, log CalendarSyncLog) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO calendar_sync_log
			(id, tenant_id, calendar_connection_id, entity_type, entity_id,
			 provider_event_id, operation, status, error_message, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		log.ID, log.TenantID, log.CalendarConnectionID, log.EntityType, log.EntityID,
		log.ProviderEventID, log.Operation, log.Status, log.ErrorMessage, log.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("calendar repo: create sync log: %w", err)
	}
	return nil
}

// UpdateSyncLogStatus updates the status, error message, and provider event ID of a sync log.
func (r *Repository) UpdateSyncLogStatus(ctx context.Context, logID uuid.UUID, status string, errorMsg *string, providerEventID *string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE calendar_sync_log
		SET status = $1, error_message = $2, provider_event_id = $3
		WHERE id = $4`,
		status, errorMsg, providerEventID, logID,
	)
	if err != nil {
		return fmt.Errorf("calendar repo: update sync log status: %w", err)
	}
	return nil
}
