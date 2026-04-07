package provider

import (
	"context"
	"encoding/json"
)

// WebhookEvent represents a parsed webhook event from a billing provider.
type WebhookEvent struct {
	ID   string
	Type string
	Data json.RawMessage
}

// BillingProvider abstracts the interaction with an external billing/payment provider.
type BillingProvider interface {
	CreateCustomer(ctx context.Context, email, name string) (customerID string, err error)
	CreateSubscription(ctx context.Context, customerID, priceID string) (subscriptionID string, err error)
	UpdateSubscription(ctx context.Context, subscriptionID, newPriceID string) error
	CancelSubscription(ctx context.Context, subscriptionID string) error
	ConstructWebhookEvent(payload []byte, sigHeader string) (*WebhookEvent, error)
}
