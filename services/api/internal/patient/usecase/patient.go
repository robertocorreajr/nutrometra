package usecase

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"nutrometra/api/internal/patient/domain"

	"github.com/google/uuid"
)

// PatientRepository defines the data access contract.
type PatientRepository interface {
	Create(ctx context.Context, p *domain.Patient) error
	GetByID(ctx context.Context, tenantID, id uuid.UUID) (*domain.Patient, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Patient, error)
	List(ctx context.Context, tenantID uuid.UUID) ([]domain.Patient, error)
	Update(ctx context.Context, p *domain.Patient) error
	CountByTenant(ctx context.Context, tenantID uuid.UUID) (int64, error)

	CreateInvite(ctx context.Context, inv *domain.PatientInvite) error
	GetInviteByHash(ctx context.Context, codeHash string) (*domain.PatientInvite, error)
	IncrementInviteUsedCount(ctx context.Context, id uuid.UUID) error

	CreateAccessLink(ctx context.Context, link *domain.PatientAccessLink) error

	CreateTenantUser(ctx context.Context, tenantID, userID, invitedBy uuid.UUID) error

	UpsertProfile(ctx context.Context, p *domain.PatientProfile) error
	GetProfile(ctx context.Context, tenantID, patientID uuid.UUID) (*domain.PatientProfile, error)
}

// EntitlementChecker resolves feature limits for a tenant.
type EntitlementChecker interface {
	CheckEntitlement(ctx context.Context, tenantID uuid.UUID, featureKey string) (enabled bool, limit *int64, err error)
}

// Usecase contains business logic for patients.
type Usecase struct {
	repo         PatientRepository
	entitlements EntitlementChecker
}

// New creates a patient Usecase.
func New(repo PatientRepository, entitlements EntitlementChecker) *Usecase {
	return &Usecase{repo: repo, entitlements: entitlements}
}

// Create creates a new patient, checking entitlement limits.
func (uc *Usecase) Create(ctx context.Context, p *domain.Patient) error {
	if err := p.Validate(); err != nil {
		return err
	}

	enabled, limit, err := uc.entitlements.CheckEntitlement(ctx, p.TenantID, "patients:create")
	if err != nil {
		return fmt.Errorf("patient: check_entitlement: %w", err)
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

// GetByID returns a patient by ID.
func (uc *Usecase) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*domain.Patient, error) {
	return uc.repo.GetByID(ctx, tenantID, id)
}

// GetByUserID returns the patient linked to a user via access links.
func (uc *Usecase) GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Patient, error) {
	return uc.repo.GetByUserID(ctx, userID)
}

// List returns all patients for a tenant.
func (uc *Usecase) List(ctx context.Context, tenantID uuid.UUID) ([]domain.Patient, error) {
	return uc.repo.List(ctx, tenantID)
}

// Update updates a patient.
func (uc *Usecase) Update(ctx context.Context, p *domain.Patient) error {
	if err := p.Validate(); err != nil {
		return err
	}
	return uc.repo.Update(ctx, p)
}

// InviteResult contains the plaintext code and invite metadata.
type InviteResult struct {
	InviteID  uuid.UUID
	Code      string // plaintext — return to caller once, never stored
	ExpiresAt time.Time
}

// GenerateInvite creates an invite code for a patient.
func (uc *Usecase) GenerateInvite(ctx context.Context, tenantID, patientID, invitedBy uuid.UUID) (*InviteResult, error) {
	// Generate random code
	codeBytes := make([]byte, 16)
	if _, err := rand.Read(codeBytes); err != nil {
		return nil, fmt.Errorf("patient: generate_code: %w", err)
	}
	code := hex.EncodeToString(codeBytes)

	// Hash the code for storage
	hash := sha256.Sum256([]byte(code))
	codeHash := hex.EncodeToString(hash[:])

	now := time.Now().UTC()
	expiresAt := now.Add(7 * 24 * time.Hour) // 7 days

	inv := &domain.PatientInvite{
		ID:        uuid.New(),
		TenantID:  tenantID,
		PatientID: patientID,
		InvitedBy: invitedBy,
		CodeHash:  codeHash,
		MaxUses:   1,
		UsedCount: 0,
		ExpiresAt: expiresAt,
		CreatedAt: now,
	}

	if err := uc.repo.CreateInvite(ctx, inv); err != nil {
		return nil, err
	}

	return &InviteResult{
		InviteID:  inv.ID,
		Code:      code,
		ExpiresAt: expiresAt,
	}, nil
}

// ActivatePortalAccess validates an invite code and creates an access link.
func (uc *Usecase) ActivatePortalAccess(ctx context.Context, code string, userID uuid.UUID) (*domain.PatientAccessLink, error) {
	// Hash the provided code
	hash := sha256.Sum256([]byte(code))
	codeHash := hex.EncodeToString(hash[:])

	inv, err := uc.repo.GetInviteByHash(ctx, codeHash)
	if err != nil {
		return nil, err
	}

	// Validate invite
	if time.Now().UTC().After(inv.ExpiresAt) {
		return nil, domain.ErrInviteExpired
	}
	if inv.UsedCount >= inv.MaxUses {
		return nil, domain.ErrInviteMaxUsed
	}

	// Create access link
	link := &domain.PatientAccessLink{
		ID:        uuid.New(),
		TenantID:  inv.TenantID,
		PatientID: inv.PatientID,
		UserID:    userID,
		InviteID:  inv.ID,
		Active:    true,
		CreatedAt: time.Now().UTC(),
	}

	if err := uc.repo.CreateAccessLink(ctx, link); err != nil {
		return nil, err
	}

	// Create tenant membership so patient can access tenant-scoped routes
	if err := uc.repo.CreateTenantUser(ctx, inv.TenantID, userID, inv.InvitedBy); err != nil {
		return nil, fmt.Errorf("patient: create tenant membership: %w", err)
	}

	// Increment used count
	if err := uc.repo.IncrementInviteUsedCount(ctx, inv.ID); err != nil {
		return nil, err
	}

	return link, nil
}

// UpsertProfile creates or updates a patient profile.
func (uc *Usecase) UpsertProfile(ctx context.Context, p *domain.PatientProfile) error {
	p.ID = uuid.New()
	now := time.Now().UTC()
	p.CreatedAt = now
	p.UpdatedAt = now
	return uc.repo.UpsertProfile(ctx, p)
}

// GetProfile returns a patient profile.
func (uc *Usecase) GetProfile(ctx context.Context, tenantID, patientID uuid.UUID) (*domain.PatientProfile, error) {
	return uc.repo.GetProfile(ctx, tenantID, patientID)
}
