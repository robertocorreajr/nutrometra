package usecase

import (
	"context"
	"fmt"
	"time"

	"nutrometra/api/internal/professional/domain"

	"github.com/google/uuid"
)

// ProfessionalRepository defines the data access contract.
type ProfessionalRepository interface {
	Create(ctx context.Context, p *domain.Professional) error
	GetByID(ctx context.Context, tenantID, id uuid.UUID) (*domain.Professional, error)
	GetByUserID(ctx context.Context, tenantID, userID uuid.UUID) (*domain.Professional, error)
	List(ctx context.Context, tenantID uuid.UUID) ([]domain.Professional, error)
	Update(ctx context.Context, p *domain.Professional) error
	CountByTenant(ctx context.Context, tenantID uuid.UUID) (int64, error)

	CreateAddress(ctx context.Context, a *domain.Address) error
	ListAddresses(ctx context.Context, tenantID, professionalID uuid.UUID) ([]domain.Address, error)
	UpdateAddress(ctx context.Context, a *domain.Address) error
	CountAddresses(ctx context.Context, tenantID, professionalID uuid.UUID) (int64, error)

	SetServiceMode(ctx context.Context, sm *domain.ProfessionalServiceMode) error
	ListServiceModes(ctx context.Context, tenantID, professionalID uuid.UUID) ([]domain.ProfessionalServiceMode, error)
}

// EntitlementChecker resolves feature limits for a tenant.
type EntitlementChecker interface {
	CheckEntitlement(ctx context.Context, tenantID uuid.UUID, featureKey string) (enabled bool, limit *int64, err error)
}

// Usecase contains business logic for professionals.
type Usecase struct {
	repo         ProfessionalRepository
	entitlements EntitlementChecker
}

// New creates a professional Usecase.
func New(repo ProfessionalRepository, entitlements EntitlementChecker) *Usecase {
	return &Usecase{repo: repo, entitlements: entitlements}
}

// Create creates a new professional, checking entitlement limits.
func (uc *Usecase) Create(ctx context.Context, p *domain.Professional) error {
	if err := p.Validate(); err != nil {
		return err
	}

	// Check entitlement: professionals:create
	enabled, limit, err := uc.entitlements.CheckEntitlement(ctx, p.TenantID, "professionals:create")
	if err != nil {
		return fmt.Errorf("professional: check_entitlement: %w", err)
	}
	if !enabled {
		return domain.ErrEntitlementExceeded
	}
	if limit != nil {
		count, err := uc.repo.CountByTenant(ctx, p.TenantID)
		if err != nil {
			return err
		}
		if count >= *limit {
			return domain.ErrEntitlementExceeded
		}
	}

	p.ID = uuid.New()
	now := time.Now().UTC()
	p.CreatedAt = now
	p.UpdatedAt = now
	p.Active = true

	return uc.repo.Create(ctx, p)
}

// GetByID returns a professional by ID.
func (uc *Usecase) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*domain.Professional, error) {
	return uc.repo.GetByID(ctx, tenantID, id)
}

// GetByUserID returns the professional linked to a user within a tenant.
func (uc *Usecase) GetByUserID(ctx context.Context, tenantID, userID uuid.UUID) (*domain.Professional, error) {
	return uc.repo.GetByUserID(ctx, tenantID, userID)
}

// List returns all professionals for a tenant.
func (uc *Usecase) List(ctx context.Context, tenantID uuid.UUID) ([]domain.Professional, error) {
	return uc.repo.List(ctx, tenantID)
}

// Update updates a professional.
func (uc *Usecase) Update(ctx context.Context, p *domain.Professional) error {
	if err := p.Validate(); err != nil {
		return err
	}
	return uc.repo.Update(ctx, p)
}

// CreateAddress creates a new address for a professional, checking entitlement limits.
func (uc *Usecase) CreateAddress(ctx context.Context, a *domain.Address) error {
	if err := a.Validate(); err != nil {
		return err
	}

	// Check entitlement: addresses:create
	enabled, limit, err := uc.entitlements.CheckEntitlement(ctx, a.TenantID, "addresses:create")
	if err != nil {
		return fmt.Errorf("professional: check_entitlement: %w", err)
	}
	if !enabled {
		return domain.ErrEntitlementExceeded
	}
	if limit != nil {
		count, err := uc.repo.CountAddresses(ctx, a.TenantID, a.ProfessionalID)
		if err != nil {
			return err
		}
		if count >= *limit {
			return domain.ErrEntitlementExceeded
		}
	}

	a.ID = uuid.New()
	now := time.Now().UTC()
	a.CreatedAt = now
	a.UpdatedAt = now
	a.Active = true
	if a.Country == "" {
		a.Country = "BR"
	}

	return uc.repo.CreateAddress(ctx, a)
}

// CreateSelf creates the professional record for the logged-in user.
// Does not check entitlements since self-provisioning is always allowed.
func (uc *Usecase) CreateSelf(ctx context.Context, p *domain.Professional) error {
	if err := p.Validate(); err != nil {
		return err
	}
	p.ID = uuid.New()
	now := time.Now().UTC()
	p.CreatedAt = now
	p.UpdatedAt = now
	p.Active = true
	return uc.repo.Create(ctx, p)
}

// ListAddresses returns addresses for a professional.
func (uc *Usecase) ListAddresses(ctx context.Context, tenantID, professionalID uuid.UUID) ([]domain.Address, error) {
	return uc.repo.ListAddresses(ctx, tenantID, professionalID)
}

// UpdateAddress updates a professional address.
func (uc *Usecase) UpdateAddress(ctx context.Context, a *domain.Address) error {
	if err := a.Validate(); err != nil {
		return err
	}
	return uc.repo.UpdateAddress(ctx, a)
}

// SetServiceMode upserts a service mode for a professional.
func (uc *Usecase) SetServiceMode(ctx context.Context, sm *domain.ProfessionalServiceMode) error {
	if sm.Mode == domain.ServiceModeOnsite && sm.AddressID == nil {
		return fmt.Errorf("professional: onsite mode requires address_id")
	}
	if sm.DurationMin <= 0 {
		sm.DurationMin = 50
	}

	sm.ID = uuid.New()
	now := time.Now().UTC()
	sm.CreatedAt = now
	sm.UpdatedAt = now
	sm.Active = true

	return uc.repo.SetServiceMode(ctx, sm)
}

// ListServiceModes returns service modes for a professional.
func (uc *Usecase) ListServiceModes(ctx context.Context, tenantID, professionalID uuid.UUID) ([]domain.ProfessionalServiceMode, error) {
	return uc.repo.ListServiceModes(ctx, tenantID, professionalID)
}
