package domain

import (
	"time"

	"github.com/google/uuid"
)

type InvoiceStatus string

const (
	InvoiceStatusDraft         InvoiceStatus = "draft"
	InvoiceStatusOpen          InvoiceStatus = "open"
	InvoiceStatusPaid          InvoiceStatus = "paid"
	InvoiceStatusVoid          InvoiceStatus = "void"
	InvoiceStatusUncollectible InvoiceStatus = "uncollectible"
)

type Invoice struct {
	ID                uuid.UUID
	TenantID          uuid.UUID
	SubscriptionID    uuid.UUID
	ProviderInvoiceID *string
	Status            InvoiceStatus
	AmountCents       int64
	Currency          string
	HostedURL         *string
	PeriodStart       *time.Time
	PeriodEnd         *time.Time
	DueDate           *time.Time
	PaidAt            *time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusSucceeded PaymentStatus = "succeeded"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusRefunded  PaymentStatus = "refunded"
)

type Payment struct {
	ID                uuid.UUID
	TenantID          uuid.UUID
	InvoiceID         uuid.UUID
	ProviderPaymentID *string
	Status            PaymentStatus
	AmountCents       int64
	Currency          string
	FailureReason     *string
	PaidAt            *time.Time
	CreatedAt         time.Time
}
