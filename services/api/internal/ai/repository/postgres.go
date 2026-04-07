package repository

import (
	"context"
	"fmt"
	"time"

	domain "nutrometra/api/internal/ai/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository provides CRUD operations for ai_suggestions.
type Repository struct {
	pool *pgxpool.Pool
}

// New creates a new Repository backed by the given connection pool.
func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// Create inserts a new ai_suggestion row.
func (r *Repository) Create(ctx context.Context, s *domain.AISuggestion) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO ai_suggestions
			(id, tenant_id, user_id, suggestion_type, status, input_context_json, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		s.ID, s.TenantID, s.UserID, string(s.SuggestionType), string(s.Status),
		s.InputContextJSON, s.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("ai_repo: create: %w", err)
	}
	return nil
}

// GetByID retrieves a suggestion by ID, scoped to tenant for isolation.
func (r *Repository) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*domain.AISuggestion, error) {
	s := &domain.AISuggestion{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, tenant_id, user_id, suggestion_type, status,
			input_context_json, prompt_used, response_text, model_id,
			input_tokens, output_tokens, created_at, completed_at,
			reviewed_at, reviewed_by
		FROM ai_suggestions
		WHERE id = $1 AND tenant_id = $2`,
		id, tenantID,
	).Scan(
		&s.ID, &s.TenantID, &s.UserID, &s.SuggestionType, &s.Status,
		&s.InputContextJSON, &s.PromptUsed, &s.ResponseText, &s.ModelID,
		&s.InputTokens, &s.OutputTokens, &s.CreatedAt, &s.CompletedAt,
		&s.ReviewedAt, &s.ReviewedBy,
	)
	if err != nil {
		return nil, fmt.Errorf("ai_repo: get_by_id: %w", err)
	}
	return s, nil
}

// ListByUser returns suggestions for a user within a tenant, ordered by creation time descending.
func (r *Repository) ListByUser(ctx context.Context, tenantID, userID uuid.UUID, limit, offset int) ([]domain.AISuggestion, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	rows, err := r.pool.Query(ctx,
		`SELECT id, tenant_id, user_id, suggestion_type, status,
			input_context_json, prompt_used, response_text, model_id,
			input_tokens, output_tokens, created_at, completed_at,
			reviewed_at, reviewed_by
		FROM ai_suggestions
		WHERE tenant_id = $1 AND user_id = $2
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4`,
		tenantID, userID, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("ai_repo: list_by_user: %w", err)
	}
	defer rows.Close()

	var suggestions []domain.AISuggestion
	for rows.Next() {
		var s domain.AISuggestion
		if err := rows.Scan(
			&s.ID, &s.TenantID, &s.UserID, &s.SuggestionType, &s.Status,
			&s.InputContextJSON, &s.PromptUsed, &s.ResponseText, &s.ModelID,
			&s.InputTokens, &s.OutputTokens, &s.CreatedAt, &s.CompletedAt,
			&s.ReviewedAt, &s.ReviewedBy,
		); err != nil {
			return nil, fmt.Errorf("ai_repo: list_by_user scan: %w", err)
		}
		suggestions = append(suggestions, s)
	}
	return suggestions, rows.Err()
}

// UpdateStatus sets the status of a suggestion.
func (r *Repository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.SuggestionStatus) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE ai_suggestions SET status = $2 WHERE id = $1`,
		id, string(status),
	)
	if err != nil {
		return fmt.Errorf("ai_repo: update_status: %w", err)
	}
	return nil
}

// UpdateCompletion saves the LLM response and marks the suggestion as completed.
func (r *Repository) UpdateCompletion(ctx context.Context, id uuid.UUID, responseText, modelID string, inputTokens, outputTokens int) error {
	now := time.Now().UTC()
	_, err := r.pool.Exec(ctx,
		`UPDATE ai_suggestions
		SET response_text = $2, model_id = $3, input_tokens = $4, output_tokens = $5,
			status = 'completed', completed_at = $6
		WHERE id = $1`,
		id, responseText, modelID, inputTokens, outputTokens, now,
	)
	if err != nil {
		return fmt.Errorf("ai_repo: update_completion: %w", err)
	}
	return nil
}

// UpdatePromptUsed saves the prompt text that was sent to the LLM for audit trail.
func (r *Repository) UpdatePromptUsed(ctx context.Context, id uuid.UUID, promptText string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE ai_suggestions SET prompt_used = $2 WHERE id = $1`,
		id, promptText,
	)
	if err != nil {
		return fmt.Errorf("ai_repo: update_prompt_used: %w", err)
	}
	return nil
}

// UpdateReview marks a suggestion as accepted or rejected by a reviewer.
func (r *Repository) UpdateReview(ctx context.Context, id, reviewedBy uuid.UUID, status domain.SuggestionStatus) error {
	now := time.Now().UTC()
	_, err := r.pool.Exec(ctx,
		`UPDATE ai_suggestions
		SET status = $2, reviewed_by = $3, reviewed_at = $4
		WHERE id = $1`,
		id, string(status), reviewedBy, now,
	)
	if err != nil {
		return fmt.Errorf("ai_repo: update_review: %w", err)
	}
	return nil
}
