package provider

import "testing"

func TestNewClaudeProvider_Defaults(t *testing.T) {
	p := NewClaudeProvider("test-key", "", 0)

	if p.defaultModel != "claude-sonnet-4-20250514" {
		t.Errorf("default model = %q, want %q", p.defaultModel, "claude-sonnet-4-20250514")
	}
	if p.defaultMaxTokens != 4096 {
		t.Errorf("default max tokens = %d, want %d", p.defaultMaxTokens, 4096)
	}
}

func TestNewClaudeProvider_CustomValues(t *testing.T) {
	p := NewClaudeProvider("test-key", "claude-opus-4-20250514", 8192)

	if p.defaultModel != "claude-opus-4-20250514" {
		t.Errorf("model = %q, want %q", p.defaultModel, "claude-opus-4-20250514")
	}
	if p.defaultMaxTokens != 8192 {
		t.Errorf("max tokens = %d, want %d", p.defaultMaxTokens, 8192)
	}
}

func TestClaudeProviderImplementsLLMProvider(t *testing.T) {
	var _ LLMProvider = (*ClaudeProvider)(nil)
}
