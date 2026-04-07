package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// SuggestionType enumerates the kinds of AI suggestions the system can produce.
type SuggestionType string

const (
	SuggestionDietDraft       SuggestionType = "diet_draft"
	SuggestionMealStructure   SuggestionType = "meal_structure"
	SuggestionSubstitutions   SuggestionType = "substitutions"
	SuggestionClinicalSummary SuggestionType = "clinical_summary"
	SuggestionReviewChecklist SuggestionType = "review_checklist"
)

// SuggestionStatus tracks the lifecycle of an AI suggestion.
type SuggestionStatus string

const (
	StatusPending    SuggestionStatus = "pending"
	StatusGenerating SuggestionStatus = "generating"
	StatusCompleted  SuggestionStatus = "completed"
	StatusFailed     SuggestionStatus = "failed"
	StatusAccepted   SuggestionStatus = "accepted"
	StatusRejected   SuggestionStatus = "rejected"
)

// AISuggestion represents a single AI-generated suggestion record.
type AISuggestion struct {
	ID               uuid.UUID
	TenantID         uuid.UUID
	UserID           uuid.UUID
	SuggestionType   SuggestionType
	Status           SuggestionStatus
	InputContextJSON json.RawMessage
	PromptUsed       *string
	ResponseText     *string
	ModelID          *string
	InputTokens      *int
	OutputTokens     *int
	CreatedAt        time.Time
	CompletedAt      *time.Time
	ReviewedAt       *time.Time
	ReviewedBy       *uuid.UUID
}

// ValidSuggestionType returns true if the given string is a valid SuggestionType.
func ValidSuggestionType(t string) bool {
	switch SuggestionType(t) {
	case SuggestionDietDraft, SuggestionMealStructure, SuggestionSubstitutions,
		SuggestionClinicalSummary, SuggestionReviewChecklist:
		return true
	}
	return false
}
