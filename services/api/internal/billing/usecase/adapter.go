package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// EntitlementAdapter wraps EntitlementService to satisfy simpler interfaces
// used by other modules (professional, patient).
type EntitlementAdapter struct {
	svc *EntitlementService
}

// NewEntitlementAdapter creates an adapter around EntitlementService.
func NewEntitlementAdapter(svc *EntitlementService) *EntitlementAdapter {
	return &EntitlementAdapter{svc: svc}
}

// CheckEntitlement returns (enabled, limit, error) for a feature.
func (a *EntitlementAdapter) CheckEntitlement(ctx context.Context, tenantID uuid.UUID, featureKey string) (bool, *int64, error) {
	ent, err := a.svc.CheckEntitlement(ctx, tenantID, featureKey)
	if err != nil {
		return false, nil, fmt.Errorf("entitlement_adapter: %w", err)
	}
	return ent.Enabled, ent.Limit, nil
}
