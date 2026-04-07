package domain_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"nutrometra/api/internal/tenancy/domain"
)

func TestTenant_IsActive(t *testing.T) {
	active := &domain.Tenant{Status: domain.TenantStatusActive}
	suspended := &domain.Tenant{Status: domain.TenantStatusSuspended}
	cancelled := &domain.Tenant{Status: domain.TenantStatusCancelled}

	assert.True(t, active.IsActive())
	assert.False(t, suspended.IsActive())
	assert.False(t, cancelled.IsActive())
}

func TestTenant_IsSuspended(t *testing.T) {
	active := &domain.Tenant{Status: domain.TenantStatusActive}
	suspended := &domain.Tenant{Status: domain.TenantStatusSuspended}
	cancelled := &domain.Tenant{Status: domain.TenantStatusCancelled}

	assert.False(t, active.IsSuspended())
	assert.True(t, suspended.IsSuspended())
	assert.False(t, cancelled.IsSuspended())
}

func TestTenant_IsCancelled(t *testing.T) {
	active := &domain.Tenant{Status: domain.TenantStatusActive}
	suspended := &domain.Tenant{Status: domain.TenantStatusSuspended}
	cancelled := &domain.Tenant{Status: domain.TenantStatusCancelled}

	assert.False(t, active.IsCancelled())
	assert.False(t, suspended.IsCancelled())
	assert.True(t, cancelled.IsCancelled())
}

func TestTenant_IsTrialing(t *testing.T) {
	base := &domain.Tenant{
		ID:          uuid.New(),
		Type:        domain.TenantTypeSoloProfessional,
		LegalName:   "Test Ltda",
		DisplayName: "Test",
		Slug:        "test",
		Status:      domain.TenantStatusActive,
		Timezone:    "America/Sao_Paulo",
		Locale:      "pt-BR",
	}

	// nil TrialEndsAt → not trialing
	base.TrialEndsAt = nil
	assert.False(t, base.IsTrialing())

	// future date → trialing
	future := time.Now().UTC().Add(24 * time.Hour)
	base.TrialEndsAt = &future
	assert.True(t, base.IsTrialing())

	// past date → not trialing
	past := time.Now().UTC().Add(-24 * time.Hour)
	base.TrialEndsAt = &past
	assert.False(t, base.IsTrialing())
}
