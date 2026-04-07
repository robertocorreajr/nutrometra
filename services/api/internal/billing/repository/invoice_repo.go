package repository

import (
	"context"
	"errors"
	"fmt"

	"nutrometra/api/internal/billing/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// InvoiceRepository handles invoice persistence.
type InvoiceRepository struct {
	pool *pgxpool.Pool
}

func NewInvoiceRepository(pool *pgxpool.Pool) *InvoiceRepository {
	return &InvoiceRepository{pool: pool}
}

// Upsert inserts or updates an invoice by provider_invoice_id.
func (r *InvoiceRepository) Upsert(ctx context.Context, exec Executor, inv domain.Invoice) error {
	_, err := exec.Exec(ctx,
		`INSERT INTO billing_invoices
		 (id, tenant_id, subscription_id, provider_invoice_id, status,
		  amount_cents, currency, hosted_url, period_start, period_end,
		  due_date, paid_at, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,NOW(),NOW())
		 ON CONFLICT (provider_invoice_id) DO UPDATE SET
		   status = EXCLUDED.status,
		   amount_cents = EXCLUDED.amount_cents,
		   hosted_url = EXCLUDED.hosted_url,
		   paid_at = EXCLUDED.paid_at,
		   updated_at = NOW()`,
		inv.ID, inv.TenantID, inv.SubscriptionID, inv.ProviderInvoiceID,
		string(inv.Status), inv.AmountCents, inv.Currency, inv.HostedURL,
		inv.PeriodStart, inv.PeriodEnd, inv.DueDate, inv.PaidAt,
	)
	if err != nil {
		return fmt.Errorf("billing: upsert_invoice: %w", err)
	}
	return nil
}

// GetByTenantID returns paginated invoices for a tenant.
func (r *InvoiceRepository) GetByTenantID(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]domain.Invoice, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, tenant_id, subscription_id, provider_invoice_id, status,
		        amount_cents, currency, hosted_url, period_start, period_end,
		        due_date, paid_at, created_at, updated_at
		 FROM billing_invoices
		 WHERE tenant_id = $1
		 ORDER BY created_at DESC
		 LIMIT $2 OFFSET $3`,
		tenantID, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("billing: list_invoices: %w", err)
	}
	defer rows.Close()

	var invoices []domain.Invoice
	for rows.Next() {
		var inv domain.Invoice
		if err := rows.Scan(
			&inv.ID, &inv.TenantID, &inv.SubscriptionID, &inv.ProviderInvoiceID,
			&inv.Status, &inv.AmountCents, &inv.Currency, &inv.HostedURL,
			&inv.PeriodStart, &inv.PeriodEnd, &inv.DueDate, &inv.PaidAt,
			&inv.CreatedAt, &inv.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("billing: scan_invoice: %w", err)
		}
		invoices = append(invoices, inv)
	}
	return invoices, rows.Err()
}

// GetByID returns a single invoice.
func (r *InvoiceRepository) GetByID(ctx context.Context, tenantID, invoiceID uuid.UUID) (*domain.Invoice, error) {
	var inv domain.Invoice
	err := r.pool.QueryRow(ctx,
		`SELECT id, tenant_id, subscription_id, provider_invoice_id, status,
		        amount_cents, currency, hosted_url, period_start, period_end,
		        due_date, paid_at, created_at, updated_at
		 FROM billing_invoices
		 WHERE id = $1 AND tenant_id = $2`,
		invoiceID, tenantID,
	).Scan(
		&inv.ID, &inv.TenantID, &inv.SubscriptionID, &inv.ProviderInvoiceID,
		&inv.Status, &inv.AmountCents, &inv.Currency, &inv.HostedURL,
		&inv.PeriodStart, &inv.PeriodEnd, &inv.DueDate, &inv.PaidAt,
		&inv.CreatedAt, &inv.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("billing: get_invoice: %w", err)
	}
	return &inv, nil
}
