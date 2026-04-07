package provider

import "context"

// CompletionRequest holds the input for an LLM completion call.
type CompletionRequest struct {
	SystemPrompt string
	UserPrompt   string
	Model        string
	MaxTokens    int
}

// CompletionResponse holds the output from an LLM completion call.
type CompletionResponse struct {
	Text         string
	Model        string
	InputTokens  int
	OutputTokens int
}

// LLMProvider abstracts an LLM backend for AI suggestions.
type LLMProvider interface {
	Complete(ctx context.Context, req CompletionRequest) (*CompletionResponse, error)
}
