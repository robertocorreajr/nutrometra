package domain

import (
	"time"

	"github.com/google/uuid"
)

type Plan struct {
	ID                uuid.UUID
	Code              string
	Name              string
	Active            bool
	BillingCycle      string
	Currency          string
	PriceCents        int64
	ProviderPriceID   *string
	ProviderProductID *string
}

type PlanFeature struct {
	ID           uuid.UUID
	PlanID       uuid.UUID
	FeatureKey   string
	Enabled      bool
	LimitValue   *int64
	TrialEnabled bool
	TrialDays    *int
}

type SubscriptionStatus string

const (
	SubscriptionTrialing  SubscriptionStatus = "trialing"
	SubscriptionActive    SubscriptionStatus = "active"
	SubscriptionPastDue   SubscriptionStatus = "past_due"
	SubscriptionCancelled SubscriptionStatus = "cancelled"
	SubscriptionExpired   SubscriptionStatus = "expired"
)

type Subscription struct {
	ID                     uuid.UUID
	TenantID               uuid.UUID
	PlanID                 uuid.UUID
	Status                 SubscriptionStatus
	StartedAt              time.Time
	TrialEndsAt            *time.Time
	RenewsAt               *time.Time
	CanceledAt             *time.Time
	ProviderCustomerID     *string
	ProviderSubscriptionID *string
	PreviousPlanID         *uuid.UUID
	PlanChangedAt          *time.Time
}

func (s Subscription) IsActive() bool {
	return s.Status == SubscriptionActive || s.Status == SubscriptionTrialing
}

type FeatureOverride struct {
	ID              uuid.UUID
	TenantID        uuid.UUID
	FeatureKey      string
	Enabled         *bool
	LimitValue      *int64
	StartsAt        *time.Time
	EndsAt          *time.Time
	Reason          string
	CreatedByUserID *uuid.UUID
}

// Entitlement is the result of resolving a feature for a tenant.
type Entitlement struct {
	FeatureKey string
	Enabled    bool
	Limit      *int64 // nil = unlimited
	Source     string // "override" | "plan" | "default"
}
