package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"nutrometra/api/internal/billing/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Executor is satisfied by *pgxpool.Pool and pgx.Tx, allowing repo methods
// to participate in an existing transaction.
type Executor interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type Repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) ListActivePlans(ctx context.Context) ([]domain.Plan, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, code, name, active, billing_cycle, currency, price_cents,
		        provider_price_id, provider_product_id
		 FROM subscription_plans WHERE active = TRUE ORDER BY price_cents`,
	)
	if err != nil {
		return nil, fmt.Errorf("billing: list_plans: %w", err)
	}
	defer rows.Close()

	var plans []domain.Plan
	for rows.Next() {
		var p domain.Plan
		if err := rows.Scan(&p.ID, &p.Code, &p.Name, &p.Active, &p.BillingCycle, &p.Currency, &p.PriceCents,
			&p.ProviderPriceID, &p.ProviderProductID); err != nil {
			return nil, fmt.Errorf("billing: scan_plan: %w", err)
		}
		plans = append(plans, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("billing: list_plans: %w", err)
	}
	return plans, nil
}

func (r *Repository) GetActiveSubscription(ctx context.Context, tenantID uuid.UUID) (*domain.Subscription, error) {
	var s domain.Subscription
	err := r.pool.QueryRow(ctx,
		`SELECT id, tenant_id, plan_id, status, started_at, trial_ends_at, renews_at, canceled_at,
		        provider_customer_id, provider_subscription_id, previous_plan_id, plan_changed_at
		 FROM tenant_subscriptions
		 WHERE tenant_id = $1 AND status IN ('active','trialing')
		   AND (status != 'trialing' OR trial_ends_at IS NULL OR trial_ends_at > NOW())
		 ORDER BY started_at DESC LIMIT 1`,
		tenantID,
	).Scan(&s.ID, &s.TenantID, &s.PlanID, &s.Status, &s.StartedAt, &s.TrialEndsAt, &s.RenewsAt, &s.CanceledAt,
		&s.ProviderCustomerID, &s.ProviderSubscriptionID, &s.PreviousPlanID, &s.PlanChangedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // no active subscription
		}
		return nil, fmt.Errorf("billing: get_subscription: %w", err)
	}
	return &s, nil
}

func (r *Repository) GetPlanFeature(ctx context.Context, planID uuid.UUID, featureKey string) (*domain.PlanFeature, error) {
	var f domain.PlanFeature
	err := r.pool.QueryRow(ctx,
		`SELECT id, plan_id, feature_key, enabled, limit_value, trial_enabled, trial_days
		 FROM plan_features WHERE plan_id = $1 AND feature_key = $2`,
		planID, featureKey,
	).Scan(&f.ID, &f.PlanID, &f.FeatureKey, &f.Enabled, &f.LimitValue, &f.TrialEnabled, &f.TrialDays)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("billing: get_plan_feature: %w", err)
	}
	return &f, nil
}

func (r *Repository) GetAllPlanFeatures(ctx context.Context, planID uuid.UUID) ([]domain.PlanFeature, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, plan_id, feature_key, enabled, limit_value, trial_enabled, trial_days
		 FROM plan_features WHERE plan_id = $1`,
		planID,
	)
	if err != nil {
		return nil, fmt.Errorf("billing: get_all_plan_features: %w", err)
	}
	defer rows.Close()

	var features []domain.PlanFeature
	for rows.Next() {
		var f domain.PlanFeature
		if err := rows.Scan(&f.ID, &f.PlanID, &f.FeatureKey, &f.Enabled, &f.LimitValue, &f.TrialEnabled, &f.TrialDays); err != nil {
			return nil, fmt.Errorf("billing: scan_plan_feature: %w", err)
		}
		features = append(features, f)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("billing: get_all_plan_features: %w", err)
	}
	return features, nil
}

func (r *Repository) GetActiveOverride(ctx context.Context, tenantID uuid.UUID, featureKey string, now time.Time) (*domain.FeatureOverride, error) {
	var o domain.FeatureOverride
	err := r.pool.QueryRow(ctx,
		`SELECT tenant_id, feature_key, enabled, limit_value, starts_at, ends_at
		 FROM tenant_feature_overrides
		 WHERE tenant_id = $1 AND feature_key = $2
		   AND (starts_at IS NULL OR starts_at <= $3)
		   AND (ends_at IS NULL OR ends_at > $3)`,
		tenantID, featureKey, now,
	).Scan(&o.TenantID, &o.FeatureKey, &o.Enabled, &o.LimitValue, &o.StartsAt, &o.EndsAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("billing: get_override: %w", err)
	}
	return &o, nil
}

func (r *Repository) CreateSubscription(ctx context.Context, s domain.Subscription) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO tenant_subscriptions
		 (id, tenant_id, plan_id, status, started_at, trial_ends_at, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,NOW(),NOW())`,
		s.ID, s.TenantID, s.PlanID, string(s.Status), s.StartedAt, s.TrialEndsAt,
	)
	if err != nil {
		return fmt.Errorf("billing: create_subscription: %w", err)
	}
	return nil
}
