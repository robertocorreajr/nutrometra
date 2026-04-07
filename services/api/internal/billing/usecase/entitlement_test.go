package usecase_test

import (
	"testing"

	"nutrometra/api/internal/billing/domain"
	"nutrometra/api/internal/billing/usecase"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolveEntitlement_OverridePrecedence(t *testing.T) {
	planFeature := &domain.PlanFeature{
		FeatureKey: "ai:assist",
		Enabled:    false,
		LimitValue: nil,
	}

	enabled := true
	limit := int64(100)
	override := &domain.FeatureOverride{
		FeatureKey: "ai:assist",
		Enabled:    &enabled,
		LimitValue: &limit,
	}

	result := usecase.ResolveEntitlement("ai:assist", planFeature, override)

	assert.True(t, result.Enabled, "override should enable the feature")
	require.NotNil(t, result.Limit)
	assert.Equal(t, int64(100), *result.Limit)
	assert.Equal(t, "override", result.Source)
}

func TestResolveEntitlement_OverrideDisables(t *testing.T) {
	planFeature := &domain.PlanFeature{
		FeatureKey: "ai:assist",
		Enabled:    true,
	}

	disabled := false
	override := &domain.FeatureOverride{
		FeatureKey: "ai:assist",
		Enabled:    &disabled,
	}

	result := usecase.ResolveEntitlement("ai:assist", planFeature, override)

	assert.False(t, result.Enabled, "override should disable the feature")
	assert.Equal(t, "override", result.Source)
}

func TestResolveEntitlement_PlanFallback(t *testing.T) {
	limit := int64(50)
	planFeature := &domain.PlanFeature{
		FeatureKey: "pdf:export",
		Enabled:    true,
		LimitValue: &limit,
	}

	result := usecase.ResolveEntitlement("pdf:export", planFeature, nil)

	assert.True(t, result.Enabled)
	require.NotNil(t, result.Limit)
	assert.Equal(t, int64(50), *result.Limit)
	assert.Equal(t, "plan", result.Source)
}

func TestResolveEntitlement_PlanNoLimit(t *testing.T) {
	planFeature := &domain.PlanFeature{
		FeatureKey: "schedule:manage",
		Enabled:    true,
		LimitValue: nil,
	}

	result := usecase.ResolveEntitlement("schedule:manage", planFeature, nil)

	assert.True(t, result.Enabled)
	assert.Nil(t, result.Limit, "nil limit = unlimited")
	assert.Equal(t, "plan", result.Source)
}

func TestResolveEntitlement_DisabledByPlan(t *testing.T) {
	planFeature := &domain.PlanFeature{
		FeatureKey: "ai:assist",
		Enabled:    false,
	}

	result := usecase.ResolveEntitlement("ai:assist", planFeature, nil)

	assert.False(t, result.Enabled)
	assert.Equal(t, "plan", result.Source)
}

func TestResolveEntitlement_NoPlanFeature_DenyByDefault(t *testing.T) {
	result := usecase.ResolveEntitlement("unknown:feature", nil, nil)

	assert.False(t, result.Enabled)
	assert.Nil(t, result.Limit)
	assert.Equal(t, "default", result.Source)
}

func TestResolveEntitlement_OverrideWithNilEnabled(t *testing.T) {
	// Override exists but Enabled is nil → defaults to true (override means "enable")
	limit := int64(200)
	override := &domain.FeatureOverride{
		FeatureKey: "ai:assist",
		Enabled:    nil,
		LimitValue: &limit,
	}

	result := usecase.ResolveEntitlement("ai:assist", nil, override)

	assert.True(t, result.Enabled, "nil Enabled in override defaults to true")
	assert.Equal(t, "override", result.Source)
}

func TestSubscription_IsActive(t *testing.T) {
	tests := []struct {
		status   domain.SubscriptionStatus
		expected bool
	}{
		{domain.SubscriptionActive, true},
		{domain.SubscriptionTrialing, true},
		{domain.SubscriptionPastDue, false},
		{domain.SubscriptionCancelled, false},
		{domain.SubscriptionExpired, false},
	}
	for _, tc := range tests {
		sub := domain.Subscription{Status: tc.status}
		assert.Equal(t, tc.expected, sub.IsActive(), "status: %s", tc.status)
	}
}
