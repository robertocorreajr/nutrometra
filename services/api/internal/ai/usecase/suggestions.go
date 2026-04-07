package usecase

import (
	"context"
	"encoding/json"
	"fmt"

	aidomain "nutrometra/api/internal/ai/domain"
	"nutrometra/api/internal/ai/provider"
	airepo "nutrometra/api/internal/ai/repository"
	"nutrometra/api/internal/platform/queue"

	"github.com/google/uuid"
)

// SuggestionService orchestrates AI suggestion creation, retrieval, and review.
type SuggestionService struct {
	repo            *airepo.Repository
	llm             provider.LLMProvider
	contextProvider PatientContextProvider
	queue           queue.Queue
}

// NewSuggestionService creates a new SuggestionService.
func NewSuggestionService(repo *airepo.Repository, llm provider.LLMProvider, ctxProvider PatientContextProvider, q queue.Queue) *SuggestionService {
	return &SuggestionService{repo: repo, llm: llm, contextProvider: ctxProvider, queue: q}
}

// CreateSuggestionInput holds the input for creating a new AI suggestion.
type CreateSuggestionInput struct {
	TenantID       uuid.UUID
	UserID         uuid.UUID
	SuggestionType string
	PatientID      *uuid.UUID
	ExtraContext   map[string]interface{}
}

// Create creates a new suggestion record and enqueues a generation job.
func (s *SuggestionService) Create(ctx context.Context, input CreateSuggestionInput) (*aidomain.AISuggestion, error) {
	// Build input context
	inputCtx := map[string]interface{}{
		"suggestion_type": input.SuggestionType,
		"extra":           input.ExtraContext,
	}
	if input.PatientID != nil {
		patientCtx, err := s.contextProvider.GetPatientContext(ctx, input.TenantID, *input.PatientID)
		if err != nil {
			return nil, fmt.Errorf("ai: get patient context: %w", err)
		}
		inputCtx["patient"] = patientCtx
	}

	contextJSON, err := json.Marshal(inputCtx)
	if err != nil {
		return nil, fmt.Errorf("ai: marshal context: %w", err)
	}

	suggestion := &aidomain.AISuggestion{
		ID:               uuid.New(),
		TenantID:         input.TenantID,
		UserID:           input.UserID,
		SuggestionType:   aidomain.SuggestionType(input.SuggestionType),
		Status:           aidomain.StatusPending,
		InputContextJSON: contextJSON,
	}

	if err := s.repo.Create(ctx, suggestion); err != nil {
		return nil, fmt.Errorf("ai: create suggestion: %w", err)
	}

	// Enqueue generation job
	payload, _ := json.Marshal(map[string]string{
		"suggestion_id": suggestion.ID.String(),
	})
	if err := s.queue.Enqueue(ctx, queue.Job{
		TenantID: &input.TenantID,
		JobType:  "ai_generate",
		Payload:  payload,
	}); err != nil {
		return nil, fmt.Errorf("ai: enqueue generation: %w", err)
	}

	return suggestion, nil
}

// GetByID retrieves a suggestion by ID, scoped to tenant.
func (s *SuggestionService) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*aidomain.AISuggestion, error) {
	return s.repo.GetByID(ctx, tenantID, id)
}

// List returns suggestions for a user.
func (s *SuggestionService) List(ctx context.Context, tenantID, userID uuid.UUID, limit, offset int) ([]aidomain.AISuggestion, error) {
	return s.repo.ListByUser(ctx, tenantID, userID, limit, offset)
}

// Accept marks a suggestion as accepted.
func (s *SuggestionService) Accept(ctx context.Context, tenantID, id, reviewerID uuid.UUID) error {
	sug, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return err
	}
	if sug.Status != aidomain.StatusCompleted {
		return fmt.Errorf("ai: can only accept completed suggestions")
	}
	return s.repo.UpdateReview(ctx, id, reviewerID, aidomain.StatusAccepted)
}

// Reject marks a suggestion as rejected.
func (s *SuggestionService) Reject(ctx context.Context, tenantID, id, reviewerID uuid.UUID) error {
	sug, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return err
	}
	if sug.Status != aidomain.StatusCompleted {
		return fmt.Errorf("ai: can only reject completed suggestions")
	}
	return s.repo.UpdateReview(ctx, id, reviewerID, aidomain.StatusRejected)
}
