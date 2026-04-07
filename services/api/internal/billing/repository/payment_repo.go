package repository

import (
	"context"
	"fmt"

	"nutrometra/api/internal/billing/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PaymentRepository handles payment persistence.
type PaymentRepository struct {
	pool *pgxpool.Pool
}

func NewPaymentRepository(pool *pgxpool.Pool) *PaymentRepository {
	return &PaymentRepository{pool: pool}
}

// Create inserts a new payment record.
func (r *PaymentRepository) Create(ctx context.Context, exec Executor, p domain.Payment) error {
	_, err := exec.Exec(ctx,
		`INSERT INTO billing_payments
		 (id, tenant_id, invoice_id, provider_payment_id, status,
		  amount_cents, currency, failure_reason, paid_at, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,NOW())`,
		p.ID, p.TenantID, p.InvoiceID, p.ProviderPaymentID,
		string(p.Status), p.AmountCents, p.Currency, p.FailureReason, p.PaidAt,
	)
	if err != nil {
		return fmt.Errorf("billing: create_payment: %w", err)
	}
	return nil
}

// GetByTenantID returns paginated payments for a tenant.
func (r *PaymentRepository) GetByTenantID(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]domain.Payment, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, tenant_id, invoice_id, provider_payment_id, status,
		        amount_cents, currency, failure_reason, paid_at, created_at
		 FROM billing_payments
		 WHERE tenant_id = $1
		 ORDER BY created_at DESC
		 LIMIT $2 OFFSET $3`,
		tenantID, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("billing: list_payments: %w", err)
	}
	defer rows.Close()

	var payments []domain.Payment
	for rows.Next() {
		var p domain.Payment
		if err := rows.Scan(
			&p.ID, &p.TenantID, &p.InvoiceID, &p.ProviderPaymentID,
			&p.Status, &p.AmountCents, &p.Currency, &p.FailureReason,
			&p.PaidAt, &p.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("billing: scan_payment: %w", err)
		}
		payments = append(payments, p)
	}
	return payments, rows.Err()
}
