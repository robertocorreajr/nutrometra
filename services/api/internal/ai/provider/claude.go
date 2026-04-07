package provider

import (
	"context"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

// ClaudeProvider implements LLMProvider using the Anthropic Messages API.
type ClaudeProvider struct {
	client          anthropic.Client
	defaultModel    string
	defaultMaxTokens int
}

// NewClaudeProvider creates a new ClaudeProvider with the given API key and defaults.
func NewClaudeProvider(apiKey, model string, maxTokens int) *ClaudeProvider {
	if model == "" {
		model = "claude-sonnet-4-20250514"
	}
	if maxTokens <= 0 {
		maxTokens = 4096
	}
	client := anthropic.NewClient(option.WithAPIKey(apiKey))
	return &ClaudeProvider{
		client:          client,
		defaultModel:    model,
		defaultMaxTokens: maxTokens,
	}
}

// Complete sends a completion request to the Anthropic Messages API and returns the response.
func (p *ClaudeProvider) Complete(ctx context.Context, req CompletionRequest) (*CompletionResponse, error) {
	model := req.Model
	if model == "" {
		model = p.defaultModel
	}
	maxTokens := req.MaxTokens
	if maxTokens <= 0 {
		maxTokens = p.defaultMaxTokens
	}

	params := anthropic.MessageNewParams{
		Model:     anthropic.Model(model),
		MaxTokens: int64(maxTokens),
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(req.UserPrompt)),
		},
	}
	if req.SystemPrompt != "" {
		params.System = []anthropic.TextBlockParam{
			{Text: req.SystemPrompt},
		}
	}

	resp, err := p.client.Messages.New(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("claude: complete: %w", err)
	}

	text := ""
	for _, block := range resp.Content {
		if block.Type == "text" {
			text += block.Text
		}
	}

	return &CompletionResponse{
		Text:         text,
		Model:        string(resp.Model),
		InputTokens:  int(resp.Usage.InputTokens),
		OutputTokens: int(resp.Usage.OutputTokens),
	}, nil
}
