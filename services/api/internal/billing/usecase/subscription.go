package usecase

import (
	"context"
	"errors"
	"fmt"

	"nutrometra/api/internal/billing/domain"
	"nutrometra/api/internal/billing/provider"
	"nutrometra/api/internal/billing/repository"
	"nutrometra/api/internal/platform/audit"
	"nutrometra/api/internal/platform/db"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var (
	ErrSubscriptionExists   = errors.New("billing: tenant already has an active subscription")
	ErrNoSubscription       = errors.New("billing: no active subscription found")
	ErrPlanNotFound         = errors.New("billing: plan not found")
	ErrPlanHasNoProviderID  = errors.New("billing: plan has no provider price ID configured")
	ErrNoProviderConfigured = errors.New("billing: no billing provider configured")
	ErrSamePlan             = errors.New("billing: already on this plan")
)

// SubscriptionUsecase manages subscription lifecycle operations.
type SubscriptionUsecase struct {
	pool        *db.Pool
	bp          provider.BillingProvider
	billingRepo *repository.Repository
	subRepo     *repository.SubscriptionRepository
	auditSvc    *audit.Service
}

func NewSubscriptionUsecase(
	pool *db.Pool,
	bp provider.BillingProvider,
	billingRepo *repository.Repository,
	subRepo *repository.SubscriptionRepository,
	auditSvc *audit.Service,
) *SubscriptionUsecase {
	return &SubscriptionUsecase{
		pool:        pool,
		bp:          bp,
		billingRepo: billingRepo,
		subRepo:     subRepo,
		auditSvc:    auditSvc,
	}
}

// Checkout creates a Stripe customer and subscription, then records it locally.
func (uc *SubscriptionUsecase) Checkout(ctx context.Context, tenantID uuid.UUID, planID uuid.UUID, email, name string, actorID uuid.UUID, actorScope audit.Scope) error {
	if uc.bp == nil {
		return ErrNoProviderConfigured
	}

	existing, err := uc.billingRepo.GetActiveSubscription(ctx, tenantID)
	if err != nil {
		return fmt.Errorf("checkout: %w", err)
	}
	if existing != nil {
		return ErrSubscriptionExists
	}

	plan, err := uc.subRepo.GetPlanByID(ctx, planID)
	if err != nil {
		return fmt.Errorf("checkout: %w", err)
	}
	if plan == nil {
		return ErrPlanNotFound
	}
	if plan.ProviderPriceID == nil || *plan.ProviderPriceID == "" {
		return ErrPlanHasNoProviderID
	}

	customerID, err := uc.bp.CreateCustomer(ctx, email, name)
	if err != nil {
		return fmt.Errorf("checkout: %w", err)
	}

	providerSubID, err := uc.bp.CreateSubscription(ctx, customerID, *plan.ProviderPriceID)
	if err != nil {
		return fmt.Errorf("checkout: %w", err)
	}

	sub := domain.Subscription{
		ID:       uuid.New(),
		TenantID: tenantID,
		PlanID:   planID,
		Status:   domain.SubscriptionActive,
	}

	return db.RunInTx(ctx, uc.pool, func(ctx context.Context, tx pgx.Tx) error {
		if err := uc.billingRepo.CreateSubscription(ctx, sub); err != nil {
			return err
		}
		if err := uc.subRepo.SetProviderIDs(ctx, tx, sub.ID, customerID, providerSubID); err != nil {
			return err
		}

		entry := audit.NewEntry(
			audit.WithTenantID(tenantID),
			audit.WithActor(actorID, actorScope),
			audit.WithEntity("tenant_subscription", sub.ID),
			audit.WithAction("subscription_checkout"),
			audit.WithMetadata(map[string]string{"plan": plan.Code}),
		)
		return uc.auditSvc.Write(ctx, tx, entry)
	})
}

// ChangePlan upgrades or downgrades the subscription plan.
func (uc *SubscriptionUsecase) ChangePlan(ctx context.Context, tenantID, newPlanID, actorID uuid.UUID, actorScope audit.Scope, reason string) error {
	if uc.bp == nil {
		return ErrNoProviderConfigured
	}

	sub, err := uc.billingRepo.GetActiveSubscription(ctx, tenantID)
	if err != nil {
		return fmt.Errorf("change_plan: %w", err)
	}
	if sub == nil {
		return ErrNoSubscription
	}
	if sub.PlanID == newPlanID {
		return ErrSamePlan
	}

	newPlan, err := uc.subRepo.GetPlanByID(ctx, newPlanID)
	if err != nil {
		return fmt.Errorf("change_plan: %w", err)
	}
	if newPlan == nil {
		return ErrPlanNotFound
	}
	if newPlan.ProviderPriceID == nil || *newPlan.ProviderPriceID == "" {
		return ErrPlanHasNoProviderID
	}

	if sub.ProviderSubscriptionID != nil && *sub.ProviderSubscriptionID != "" {
		if err := uc.bp.UpdateSubscription(ctx, *sub.ProviderSubscriptionID, *newPlan.ProviderPriceID); err != nil {
			return fmt.Errorf("change_plan: %w", err)
		}
	}

	return db.RunInTx(ctx, uc.pool, func(ctx context.Context, tx pgx.Tx) error {
		if err := uc.subRepo.UpdatePlan(ctx, tx, sub.ID, newPlanID); err != nil {
			return err
		}

		entry := audit.NewEntry(
			audit.WithTenantID(tenantID),
			audit.WithActor(actorID, actorScope),
			audit.WithEntity("tenant_subscription", sub.ID),
			audit.WithAction("plan_changed"),
			audit.WithReason(reason),
			audit.WithMetadata(map[string]string{
				"from_plan": sub.PlanID.String(),
				"to_plan":   newPlanID.String(),
			}),
		)
		return uc.auditSvc.Write(ctx, tx, entry)
	})
}

// Cancel cancels the active subscription.
func (uc *SubscriptionUsecase) Cancel(ctx context.Context, tenantID, actorID uuid.UUID, actorScope audit.Scope, reason string) error {
	sub, err := uc.billingRepo.GetActiveSubscription(ctx, tenantID)
	if err != nil {
		return fmt.Errorf("cancel: %w", err)
	}
	if sub == nil {
		return ErrNoSubscription
	}

	if uc.bp != nil && sub.ProviderSubscriptionID != nil && *sub.ProviderSubscriptionID != "" {
		if err := uc.bp.CancelSubscription(ctx, *sub.ProviderSubscriptionID); err != nil {
			return fmt.Errorf("cancel: %w", err)
		}
	}

	return db.RunInTx(ctx, uc.pool, func(ctx context.Context, tx pgx.Tx) error {
		if err := uc.subRepo.Cancel(ctx, tx, sub.ID); err != nil {
			return err
		}

		entry := audit.NewEntry(
			audit.WithTenantID(tenantID),
			audit.WithActor(actorID, actorScope),
			audit.WithEntity("tenant_subscription", sub.ID),
			audit.WithAction("subscription_cancelled"),
			audit.WithReason(reason),
		)
		return uc.auditSvc.Write(ctx, tx, entry)
	})
}

// Reactivate reactivates a cancelled subscription.
func (uc *SubscriptionUsecase) Reactivate(ctx context.Context, tenantID, actorID uuid.UUID, actorScope audit.Scope, reason string) error {
	sub, err := uc.subRepo.GetAnySubscription(ctx, tenantID)
	if err != nil {
		return fmt.Errorf("reactivate: %w", err)
	}
	if sub == nil {
		return ErrNoSubscription
	}
	if sub.Status != domain.SubscriptionCancelled {
		return fmt.Errorf("reactivate: subscription is not cancelled (status: %s)", sub.Status)
	}

	return db.RunInTx(ctx, uc.pool, func(ctx context.Context, tx pgx.Tx) error {
		if err := uc.subRepo.Reactivate(ctx, tx, sub.ID); err != nil {
			return err
		}

		entry := audit.NewEntry(
			audit.WithTenantID(tenantID),
			audit.WithActor(actorID, actorScope),
			audit.WithEntity("tenant_subscription", sub.ID),
			audit.WithAction("subscription_reactivated"),
			audit.WithReason(reason),
		)
		return uc.auditSvc.Write(ctx, tx, entry)
	})
}
