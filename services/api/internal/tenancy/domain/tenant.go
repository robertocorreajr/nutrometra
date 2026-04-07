package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Sentinel errors.
var ErrTenantNotFound = errors.New("tenant not found")

type TenantStatus string

const (
	TenantStatusActive    TenantStatus = "active"
	TenantStatusSuspended TenantStatus = "suspended"
	TenantStatusCancelled TenantStatus = "cancelled"
)

type Tenant struct {
	ID          uuid.UUID
	Name        string
	Slug        string
	PlanID      *uuid.UUID
	Status      TenantStatus
	TrialEndsAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (t *Tenant) IsActive() bool {
	return t.Status == TenantStatusActive
}

func (t *Tenant) IsSuspended() bool {
	return t.Status == TenantStatusSuspended
}

func (t *Tenant) IsCancelled() bool {
	return t.Status == TenantStatusCancelled
}

func (t *Tenant) IsTrialing() bool {
	return t.TrialEndsAt != nil && time.Now().UTC().Before(t.TrialEndsAt.UTC())
}
