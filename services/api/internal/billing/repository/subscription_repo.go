package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"nutrometra/api/internal/billing/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// SubscriptionRepository handles subscription-specific mutations beyond
// what the base Repository provides.
type SubscriptionRepository struct {
	pool *pgxpool.Pool
}

func NewSubscriptionRepository(pool *pgxpool.Pool) *SubscriptionRepository {
	return &SubscriptionRepository{pool: pool}
}

// SetProviderIDs sets the Stripe customer and subscription IDs on a subscription.
func (r *SubscriptionRepository) SetProviderIDs(ctx context.Context, exec Executor, subID uuid.UUID, customerID, subscriptionID string) error {
	_, err := exec.Exec(ctx,
		`UPDATE tenant_subscriptions
		 SET provider_customer_id = $1, provider_subscription_id = $2, updated_at = NOW()
		 WHERE id = $3`,
		customerID, subscriptionID, subID,
	)
	if err != nil {
		return fmt.Errorf("billing: set_provider_ids: %w", err)
	}
	return nil
}

// UpdateStatus updates the subscription status.
func (r *SubscriptionRepository) UpdateStatus(ctx context.Context, exec Executor, subID uuid.UUID, status domain.SubscriptionStatus) error {
	_, err := exec.Exec(ctx,
		`UPDATE tenant_subscriptions SET status = $1, updated_at = NOW() WHERE id = $2`,
		string(status), subID,
	)
	if err != nil {
		return fmt.Errorf("billing: update_sub_status: %w", err)
	}
	return nil
}

// UpdatePlan changes the subscription plan, recording the previous plan.
func (r *SubscriptionRepository) UpdatePlan(ctx context.Context, exec Executor, subID, newPlanID uuid.UUID) error {
	now := time.Now().UTC()
	_, err := exec.Exec(ctx,
		`UPDATE tenant_subscriptions
		 SET previous_plan_id = plan_id, plan_id = $1, plan_changed_at = $2, updated_at = NOW()
		 WHERE id = $3`,
		newPlanID, now, subID,
	)
	if err != nil {
		return fmt.Errorf("billing: update_sub_plan: %w", err)
	}
	return nil
}

// Cancel marks a subscription as cancelled.
func (r *SubscriptionRepository) Cancel(ctx context.Context, exec Executor, subID uuid.UUID) error {
	now := time.Now().UTC()
	_, err := exec.Exec(ctx,
		`UPDATE tenant_subscriptions SET status = 'cancelled', canceled_at = $1, updated_at = NOW() WHERE id = $2`,
		now, subID,
	)
	if err != nil {
		return fmt.Errorf("billing: cancel_subscription: %w", err)
	}
	return nil
}

// Reactivate sets a cancelled subscription back to active.
func (r *SubscriptionRepository) Reactivate(ctx context.Context, exec Executor, subID uuid.UUID) error {
	_, err := exec.Exec(ctx,
		`UPDATE tenant_subscriptions SET status = 'active', canceled_at = NULL, updated_at = NOW() WHERE id = $1`,
		subID,
	)
	if err != nil {
		return fmt.Errorf("billing: reactivate_subscription: %w", err)
	}
	return nil
}

// GetByProviderSubscriptionID finds a subscription by its Stripe subscription ID.
func (r *SubscriptionRepository) GetByProviderSubscriptionID(ctx context.Context, providerSubID string) (*domain.Subscription, error) {
	var s domain.Subscription
	err := r.pool.QueryRow(ctx,
		`SELECT id, tenant_id, plan_id, status, started_at, trial_ends_at, renews_at, canceled_at,
		        provider_customer_id, provider_subscription_id, previous_plan_id, plan_changed_at
		 FROM tenant_subscriptions
		 WHERE provider_subscription_id = $1
		 ORDER BY started_at DESC LIMIT 1`,
		providerSubID,
	).Scan(&s.ID, &s.TenantID, &s.PlanID, &s.Status, &s.StartedAt, &s.TrialEndsAt, &s.RenewsAt, &s.CanceledAt,
		&s.ProviderCustomerID, &s.ProviderSubscriptionID, &s.PreviousPlanID, &s.PlanChangedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("billing: get_sub_by_provider_id: %w", err)
	}
	return &s, nil
}

// GetPlanByProviderPriceID finds a plan by its Stripe price ID.
func (r *SubscriptionRepository) GetPlanByProviderPriceID(ctx context.Context, priceID string) (*domain.Plan, error) {
	var p domain.Plan
	err := r.pool.QueryRow(ctx,
		`SELECT id, code, name, active, billing_cycle, currency, price_cents,
		        provider_price_id, provider_product_id
		 FROM subscription_plans
		 WHERE provider_price_id = $1`,
		priceID,
	).Scan(&p.ID, &p.Code, &p.Name, &p.Active, &p.BillingCycle, &p.Currency, &p.PriceCents,
		&p.ProviderPriceID, &p.ProviderProductID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("billing: get_plan_by_price_id: %w", err)
	}
	return &p, nil
}

// GetAnySubscription returns the most recent subscription for a tenant regardless of status.
func (r *SubscriptionRepository) GetAnySubscription(ctx context.Context, tenantID uuid.UUID) (*domain.Subscription, error) {
	var s domain.Subscription
	err := r.pool.QueryRow(ctx,
		`SELECT id, tenant_id, plan_id, status, started_at, trial_ends_at, renews_at, canceled_at,
		        provider_customer_id, provider_subscription_id, previous_plan_id, plan_changed_at
		 FROM tenant_subscriptions
		 WHERE tenant_id = $1
		 ORDER BY started_at DESC LIMIT 1`,
		tenantID,
	).Scan(&s.ID, &s.TenantID, &s.PlanID, &s.Status, &s.StartedAt, &s.TrialEndsAt, &s.RenewsAt, &s.CanceledAt,
		&s.ProviderCustomerID, &s.ProviderSubscriptionID, &s.PreviousPlanID, &s.PlanChangedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("billing: get_any_subscription: %w", err)
	}
	return &s, nil
}

// GetPlanByID returns a plan by its UUID.
func (r *SubscriptionRepository) GetPlanByID(ctx context.Context, planID uuid.UUID) (*domain.Plan, error) {
	var p domain.Plan
	err := r.pool.QueryRow(ctx,
		`SELECT id, code, name, active, billing_cycle, currency, price_cents,
		        provider_price_id, provider_product_id
		 FROM subscription_plans
		 WHERE id = $1`,
		planID,
	).Scan(&p.ID, &p.Code, &p.Name, &p.Active, &p.BillingCycle, &p.Currency, &p.PriceCents,
		&p.ProviderPriceID, &p.ProviderProductID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("billing: get_plan_by_id: %w", err)
	}
	return &p, nil
}
