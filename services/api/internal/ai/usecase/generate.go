package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	aidomain "nutrometra/api/internal/ai/domain"
	"nutrometra/api/internal/ai/provider"
	airepo "nutrometra/api/internal/ai/repository"
	"nutrometra/api/internal/platform/queue"

	"github.com/google/uuid"
)

// GenerateHandler processes ai_generate jobs from the queue.
type GenerateHandler struct {
	repo *airepo.Repository
	llm  provider.LLMProvider
}

// NewGenerateHandler creates a new GenerateHandler.
func NewGenerateHandler(repo *airepo.Repository, llm provider.LLMProvider) *GenerateHandler {
	return &GenerateHandler{repo: repo, llm: llm}
}

// Handle processes a single ai_generate job.
func (h *GenerateHandler) Handle(ctx context.Context, job *queue.JobRecord) error {
	var payload struct {
		SuggestionID string `json:"suggestion_id"`
	}
	if err := json.Unmarshal(job.PayloadJSON, &payload); err != nil {
		return fmt.Errorf("ai_generate: unmarshal payload: %w", err)
	}

	suggestionID, err := uuid.Parse(payload.SuggestionID)
	if err != nil {
		return fmt.Errorf("ai_generate: parse suggestion_id: %w", err)
	}

	// Get suggestion - use job.TenantID if available
	var tenantID uuid.UUID
	if job.TenantID != nil {
		tenantID = *job.TenantID
	}

	sug, err := h.repo.GetByID(ctx, tenantID, suggestionID)
	if err != nil {
		return fmt.Errorf("ai_generate: get suggestion: %w", err)
	}

	// Update status to generating
	if err := h.repo.UpdateStatus(ctx, sug.ID, aidomain.StatusGenerating); err != nil {
		return fmt.Errorf("ai_generate: update status: %w", err)
	}

	// Build prompt
	systemPrompt := buildSystemPrompt(sug.SuggestionType)
	userPrompt := buildUserPrompt(sug.SuggestionType, sug.InputContextJSON)

	slog.Info("ai_generate: calling LLM", "suggestion_id", sug.ID, "type", sug.SuggestionType)

	resp, err := h.llm.Complete(ctx, provider.CompletionRequest{
		SystemPrompt: systemPrompt,
		UserPrompt:   userPrompt,
	})
	if err != nil {
		_ = h.repo.UpdateStatus(ctx, sug.ID, aidomain.StatusFailed)
		return fmt.Errorf("ai_generate: llm complete: %w", err)
	}

	// Save completion
	if err := h.repo.UpdateCompletion(ctx, sug.ID, resp.Text, resp.Model, resp.InputTokens, resp.OutputTokens); err != nil {
		return fmt.Errorf("ai_generate: update completion: %w", err)
	}

	// Save prompt used for audit trail
	promptText := fmt.Sprintf("System: %s\n\nUser: %s", systemPrompt, userPrompt)
	if err := h.repo.UpdatePromptUsed(ctx, sug.ID, promptText); err != nil {
		slog.Warn("ai_generate: failed to save prompt_used", "error", err)
	}

	slog.Info("ai_generate: completed", "suggestion_id", sug.ID, "tokens_in", resp.InputTokens, "tokens_out", resp.OutputTokens)
	return nil
}

func buildSystemPrompt(sugType aidomain.SuggestionType) string {
	base := `Você é um assistente de nutrição. Suas sugestões são ASSISTIVAS e requerem revisão por um nutricionista habilitado.
IMPORTANTE: Sempre indicar que:
- Esta é uma sugestão gerada por IA
- Requer revisão e aprovação de profissional habilitado
- O contexto utilizado pode estar incompleto
- Não substitui o julgamento clínico profissional`

	switch sugType {
	case aidomain.SuggestionDietDraft:
		return base + "\n\nGere um rascunho de plano alimentar com refeições, horários e alimentos sugeridos."
	case aidomain.SuggestionMealStructure:
		return base + "\n\nSugira uma estrutura de refeições (número, horários, distribuição calórica)."
	case aidomain.SuggestionSubstitutions:
		return base + "\n\nSugira substituições alimentares equivalentes em valor nutricional."
	case aidomain.SuggestionClinicalSummary:
		return base + "\n\nGere um resumo clínico com base nos dados disponíveis do paciente."
	case aidomain.SuggestionReviewChecklist:
		return base + "\n\nGere um checklist de revisão para o plano alimentar."
	default:
		return base
	}
}

func buildUserPrompt(sugType aidomain.SuggestionType, contextJSON json.RawMessage) string {
	return fmt.Sprintf("Tipo de sugestão: %s\n\nContexto disponível:\n%s", sugType, string(contextJSON))
}
