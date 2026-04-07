package domain

import (
	"time"

	"github.com/google/uuid"
)

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
