package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/webhook"
)

// StripeProvider implements BillingProvider using Stripe.
type StripeProvider struct {
	client        *stripe.Client
	webhookSecret string
}

// NewStripeProvider creates a new Stripe billing provider.
func NewStripeProvider(secretKey, webhookSecret string) *StripeProvider {
	sc := stripe.NewClient(secretKey)
	return &StripeProvider{
		client:        sc,
		webhookSecret: webhookSecret,
	}
}

func (s *StripeProvider) CreateCustomer(ctx context.Context, email, name string) (string, error) {
	params := &stripe.CustomerCreateParams{
		Email: stripe.String(email),
		Name:  stripe.String(name),
	}
	cust, err := s.client.V1Customers.Create(ctx, params)
	if err != nil {
		return "", fmt.Errorf("stripe: create_customer: %w", err)
	}
	return cust.ID, nil
}

func (s *StripeProvider) CreateSubscription(ctx context.Context, customerID, priceID string) (string, error) {
	params := &stripe.SubscriptionCreateParams{
		Customer: stripe.String(customerID),
		Items: []*stripe.SubscriptionCreateItemParams{
			{Price: stripe.String(priceID)},
		},
	}
	sub, err := s.client.V1Subscriptions.Create(ctx, params)
	if err != nil {
		return "", fmt.Errorf("stripe: create_subscription: %w", err)
	}
	return sub.ID, nil
}

func (s *StripeProvider) UpdateSubscription(ctx context.Context, subscriptionID, newPriceID string) error {
	sub, err := s.client.V1Subscriptions.Retrieve(ctx, subscriptionID, nil)
	if err != nil {
		return fmt.Errorf("stripe: get_subscription: %w", err)
	}
	if len(sub.Items.Data) == 0 {
		return fmt.Errorf("stripe: subscription has no items")
	}
	itemID := sub.Items.Data[0].ID

	params := &stripe.SubscriptionUpdateParams{
		Items: []*stripe.SubscriptionUpdateItemParams{
			{
				ID:    stripe.String(itemID),
				Price: stripe.String(newPriceID),
			},
		},
		ProrationBehavior: stripe.String("create_prorations"),
	}
	_, err = s.client.V1Subscriptions.Update(ctx, subscriptionID, params)
	if err != nil {
		return fmt.Errorf("stripe: update_subscription: %w", err)
	}
	return nil
}

func (s *StripeProvider) CancelSubscription(ctx context.Context, subscriptionID string) error {
	_, err := s.client.V1Subscriptions.Cancel(ctx, subscriptionID, nil)
	if err != nil {
		return fmt.Errorf("stripe: cancel_subscription: %w", err)
	}
	return nil
}

func (s *StripeProvider) ConstructWebhookEvent(payload []byte, sigHeader string) (*WebhookEvent, error) {
	event, err := webhook.ConstructEvent(payload, sigHeader, s.webhookSecret)
	if err != nil {
		return nil, fmt.Errorf("stripe: invalid_webhook_signature: %w", err)
	}
	data, err := json.Marshal(event.Data.Object)
	if err != nil {
		return nil, fmt.Errorf("stripe: marshal_event_data: %w", err)
	}
	return &WebhookEvent{
		ID:   event.ID,
		Type: string(event.Type),
		Data: data,
	}, nil
}
