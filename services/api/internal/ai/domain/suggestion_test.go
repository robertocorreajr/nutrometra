package domain

import "testing"

func TestValidSuggestionType(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"diet_draft is valid", "diet_draft", true},
		{"meal_structure is valid", "meal_structure", true},
		{"substitutions is valid", "substitutions", true},
		{"clinical_summary is valid", "clinical_summary", true},
		{"review_checklist is valid", "review_checklist", true},
		{"empty is invalid", "", false},
		{"unknown type is invalid", "unknown_type", false},
		{"partial match is invalid", "diet", false},
		{"uppercase is invalid", "DIET_DRAFT", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidSuggestionType(tt.input)
			if got != tt.want {
				t.Errorf("ValidSuggestionType(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestSuggestionTypeConstants(t *testing.T) {
	// Verify constants match expected string values
	if string(SuggestionDietDraft) != "diet_draft" {
		t.Errorf("SuggestionDietDraft = %q, want %q", SuggestionDietDraft, "diet_draft")
	}
	if string(SuggestionMealStructure) != "meal_structure" {
		t.Errorf("SuggestionMealStructure = %q, want %q", SuggestionMealStructure, "meal_structure")
	}
	if string(SuggestionSubstitutions) != "substitutions" {
		t.Errorf("SuggestionSubstitutions = %q, want %q", SuggestionSubstitutions, "substitutions")
	}
	if string(SuggestionClinicalSummary) != "clinical_summary" {
		t.Errorf("SuggestionClinicalSummary = %q, want %q", SuggestionClinicalSummary, "clinical_summary")
	}
	if string(SuggestionReviewChecklist) != "review_checklist" {
		t.Errorf("SuggestionReviewChecklist = %q, want %q", SuggestionReviewChecklist, "review_checklist")
	}
}

func TestSuggestionStatusConstants(t *testing.T) {
	if string(StatusPending) != "pending" {
		t.Errorf("StatusPending = %q, want %q", StatusPending, "pending")
	}
	if string(StatusGenerating) != "generating" {
		t.Errorf("StatusGenerating = %q, want %q", StatusGenerating, "generating")
	}
	if string(StatusCompleted) != "completed" {
		t.Errorf("StatusCompleted = %q, want %q", StatusCompleted, "completed")
	}
	if string(StatusFailed) != "failed" {
		t.Errorf("StatusFailed = %q, want %q", StatusFailed, "failed")
	}
	if string(StatusAccepted) != "accepted" {
		t.Errorf("StatusAccepted = %q, want %q", StatusAccepted, "accepted")
	}
	if string(StatusRejected) != "rejected" {
		t.Errorf("StatusRejected = %q, want %q", StatusRejected, "rejected")
	}
}
