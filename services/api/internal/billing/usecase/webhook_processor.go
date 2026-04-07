package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"nutrometra/api/internal/billing/domain"
	"nutrometra/api/internal/billing/repository"
	"nutrometra/api/internal/platform/audit"
	"nutrometra/api/internal/platform/db"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// WebhookProcessor handles idempotent processing of billing provider webhook events.
type WebhookProcessor struct {
	pool        *db.Pool
	webhookRepo *repository.WebhookRepository
	subRepo     *repository.SubscriptionRepository
	invoiceRepo *repository.InvoiceRepository
	paymentRepo *repository.PaymentRepository
	auditSvc    *audit.Service
}

func NewWebhookProcessor(
	pool *db.Pool,
	webhookRepo *repository.WebhookRepository,
	subRepo *repository.SubscriptionRepository,
	invoiceRepo *repository.InvoiceRepository,
	paymentRepo *repository.PaymentRepository,
	auditSvc *audit.Service,
) *WebhookProcessor {
	return &WebhookProcessor{
		pool:        pool,
		webhookRepo: webhookRepo,
		subRepo:     subRepo,
		invoiceRepo: invoiceRepo,
		paymentRepo: paymentRepo,
		auditSvc:    auditSvc,
	}
}

// ProcessEvent processes a webhook event idempotently.
func (wp *WebhookProcessor) ProcessEvent(ctx context.Context, providerEventID, eventType string, data json.RawMessage) error {
	evt := domain.WebhookEvent{
		ID:              uuid.New(),
		Provider:        "stripe",
		ProviderEventID: providerEventID,
		EventType:       eventType,
		PayloadJSON:     data,
		Status:          domain.WebhookStatusPending,
	}

	inserted, err := wp.webhookRepo.InsertIdempotent(ctx, evt)
	if err != nil {
		return fmt.Errorf("webhook: insert: %w", err)
	}
	if !inserted {
		slog.InfoContext(ctx, "webhook event already processed", "event_id", providerEventID)
		return nil
	}

	var processErr error
	switch eventType {
	case "invoice.paid":
		processErr = wp.handleInvoicePaid(ctx, data)
	case "invoice.payment_failed":
		processErr = wp.handleInvoicePaymentFailed(ctx, data)
	case "customer.subscription.updated":
		processErr = wp.handleSubscriptionUpdated(ctx, data)
	case "customer.subscription.deleted":
		processErr = wp.handleSubscriptionDeleted(ctx, data)
	default:
		slog.InfoContext(ctx, "unhandled webhook event type", "type", eventType)
		_ = wp.webhookRepo.MarkProcessed(ctx, evt.ID)
		return nil
	}

	if processErr != nil {
		slog.ErrorContext(ctx, "webhook processing failed", "event_id", providerEventID, "type", eventType, "error", processErr)
		_ = wp.webhookRepo.MarkFailed(ctx, evt.ID, processErr.Error())
		return processErr
	}

	_ = wp.webhookRepo.MarkProcessed(ctx, evt.ID)
	return nil
}

type stripeInvoiceData struct {
	ID           string `json:"id"`
	Customer     string `json:"customer"`
	Subscription string `json:"subscription"`
	AmountPaid   int64  `json:"amount_paid"`
	Currency     string `json:"currency"`
	HostedURL    string `json:"hosted_invoice_url"`
	Status       string `json:"status"`
}

func (wp *WebhookProcessor) handleInvoicePaid(ctx context.Context, data json.RawMessage) error {
	var inv stripeInvoiceData
	if err := json.Unmarshal(data, &inv); err != nil {
		return fmt.Errorf("webhook: unmarshal invoice: %w", err)
	}

	sub, err := wp.subRepo.GetByProviderSubscriptionID(ctx, inv.Subscription)
	if err != nil {
		return fmt.Errorf("webhook: find subscription: %w", err)
	}
	if sub == nil {
		slog.WarnContext(ctx, "webhook: subscription not found for invoice", "provider_sub_id", inv.Subscription)
		return nil
	}

	now := time.Now().UTC()
	return db.RunInTx(ctx, wp.pool, func(ctx context.Context, tx pgx.Tx) error {
		invoice := domain.Invoice{
			ID:                uuid.New(),
			TenantID:          sub.TenantID,
			SubscriptionID:    sub.ID,
			ProviderInvoiceID: &inv.ID,
			Status:            domain.InvoiceStatusPaid,
			AmountCents:       inv.AmountPaid,
			Currency:          inv.Currency,
			HostedURL:         strPtr(inv.HostedURL),
			PaidAt:            &now,
		}
		if err := wp.invoiceRepo.Upsert(ctx, tx, invoice); err != nil {
			return err
		}

		payment := domain.Payment{
			ID:                uuid.New(),
			TenantID:          sub.TenantID,
			InvoiceID:         invoice.ID,
			ProviderPaymentID: &inv.ID,
			Status:            domain.PaymentStatusSucceeded,
			AmountCents:       inv.AmountPaid,
			Currency:          inv.Currency,
			PaidAt:            &now,
		}
		if err := wp.paymentRepo.Create(ctx, tx, payment); err != nil {
			return err
		}

		if err := wp.subRepo.UpdateStatus(ctx, tx, sub.ID, domain.SubscriptionActive); err != nil {
			return err
		}

		entry := audit.NewEntry(
			audit.WithTenantID(sub.TenantID),
			audit.WithActor(uuid.Nil, audit.ScopeSystem),
			audit.WithEntity("tenant_subscription", sub.ID),
			audit.WithAction("invoice_paid"),
			audit.WithMetadata(map[string]string{"provider_invoice_id": inv.ID}),
		)
		return wp.auditSvc.Write(ctx, tx, entry)
	})
}

func (wp *WebhookProcessor) handleInvoicePaymentFailed(ctx context.Context, data json.RawMessage) error {
	var inv stripeInvoiceData
	if err := json.Unmarshal(data, &inv); err != nil {
		return fmt.Errorf("webhook: unmarshal invoice: %w", err)
	}

	sub, err := wp.subRepo.GetByProviderSubscriptionID(ctx, inv.Subscription)
	if err != nil {
		return fmt.Errorf("webhook: find subscription: %w", err)
	}
	if sub == nil {
		return nil
	}

	return db.RunInTx(ctx, wp.pool, func(ctx context.Context, tx pgx.Tx) error {
		invoice := domain.Invoice{
			ID:                uuid.New(),
			TenantID:          sub.TenantID,
			SubscriptionID:    sub.ID,
			ProviderInvoiceID: &inv.ID,
			Status:            domain.InvoiceStatusOpen,
			AmountCents:       inv.AmountPaid,
			Currency:          inv.Currency,
			HostedURL:         strPtr(inv.HostedURL),
		}
		if err := wp.invoiceRepo.Upsert(ctx, tx, invoice); err != nil {
			return err
		}

		payment := domain.Payment{
			ID:                uuid.New(),
			TenantID:          sub.TenantID,
			InvoiceID:         invoice.ID,
			ProviderPaymentID: &inv.ID,
			Status:            domain.PaymentStatusFailed,
			AmountCents:       inv.AmountPaid,
			Currency:          inv.Currency,
		}
		if err := wp.paymentRepo.Create(ctx, tx, payment); err != nil {
			return err
		}

		if err := wp.subRepo.UpdateStatus(ctx, tx, sub.ID, domain.SubscriptionPastDue); err != nil {
			return err
		}

		entry := audit.NewEntry(
			audit.WithTenantID(sub.TenantID),
			audit.WithActor(uuid.Nil, audit.ScopeSystem),
			audit.WithEntity("tenant_subscription", sub.ID),
			audit.WithAction("invoice_payment_failed"),
			audit.WithMetadata(map[string]string{"provider_invoice_id": inv.ID}),
		)
		return wp.auditSvc.Write(ctx, tx, entry)
	})
}

type stripeSubscriptionData struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Plan   struct {
		ID string `json:"id"`
	} `json:"plan"`
	Items struct {
		Data []struct {
			Price struct {
				ID string `json:"id"`
			} `json:"price"`
		} `json:"data"`
	} `json:"items"`
}

func (wp *WebhookProcessor) handleSubscriptionUpdated(ctx context.Context, data json.RawMessage) error {
	var subData stripeSubscriptionData
	if err := json.Unmarshal(data, &subData); err != nil {
		return fmt.Errorf("webhook: unmarshal subscription: %w", err)
	}

	sub, err := wp.subRepo.GetByProviderSubscriptionID(ctx, subData.ID)
	if err != nil {
		return fmt.Errorf("webhook: find subscription: %w", err)
	}
	if sub == nil {
		return nil
	}

	return db.RunInTx(ctx, wp.pool, func(ctx context.Context, tx pgx.Tx) error {
		newStatus := mapStripeStatus(subData.Status)
		if newStatus != "" && newStatus != sub.Status {
			if err := wp.subRepo.UpdateStatus(ctx, tx, sub.ID, newStatus); err != nil {
				return err
			}
		}

		// Sync plan if price changed
		if len(subData.Items.Data) > 0 {
			priceID := subData.Items.Data[0].Price.ID
			if priceID != "" {
				plan, err := wp.subRepo.GetPlanByProviderPriceID(ctx, priceID)
				if err != nil {
					return err
				}
				if plan != nil && plan.ID != sub.PlanID {
					if err := wp.subRepo.UpdatePlan(ctx, tx, sub.ID, plan.ID); err != nil {
						return err
					}
				}
			}
		}

		entry := audit.NewEntry(
			audit.WithTenantID(sub.TenantID),
			audit.WithActor(uuid.Nil, audit.ScopeSystem),
			audit.WithEntity("tenant_subscription", sub.ID),
			audit.WithAction("subscription_updated_via_webhook"),
			audit.WithMetadata(map[string]string{"stripe_status": subData.Status}),
		)
		return wp.auditSvc.Write(ctx, tx, entry)
	})
}

func (wp *WebhookProcessor) handleSubscriptionDeleted(ctx context.Context, data json.RawMessage) error {
	var subData stripeSubscriptionData
	if err := json.Unmarshal(data, &subData); err != nil {
		return fmt.Errorf("webhook: unmarshal subscription: %w", err)
	}

	sub, err := wp.subRepo.GetByProviderSubscriptionID(ctx, subData.ID)
	if err != nil {
		return fmt.Errorf("webhook: find subscription: %w", err)
	}
	if sub == nil {
		return nil
	}

	return db.RunInTx(ctx, wp.pool, func(ctx context.Context, tx pgx.Tx) error {
		if err := wp.subRepo.Cancel(ctx, tx, sub.ID); err != nil {
			return err
		}

		entry := audit.NewEntry(
			audit.WithTenantID(sub.TenantID),
			audit.WithActor(uuid.Nil, audit.ScopeSystem),
			audit.WithEntity("tenant_subscription", sub.ID),
			audit.WithAction("subscription_cancelled_via_webhook"),
		)
		return wp.auditSvc.Write(ctx, tx, entry)
	})
}

func mapStripeStatus(s string) domain.SubscriptionStatus {
	switch s {
	case "active":
		return domain.SubscriptionActive
	case "trialing":
		return domain.SubscriptionTrialing
	case "past_due":
		return domain.SubscriptionPastDue
	case "canceled", "cancelled":
		return domain.SubscriptionCancelled
	case "unpaid":
		return domain.SubscriptionPastDue
	default:
		return ""
	}
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
