package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"nutrometra/api/internal/patient/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository provides patient data access.
type Repository struct {
	pool *pgxpool.Pool
}

// New creates a patient Repository.
func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// Create inserts a new patient.
func (r *Repository) Create(ctx context.Context, p *domain.Patient) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO patients
			(id, tenant_id, professional_id, full_name, email, phone, cpf,
			 date_of_birth, gender, notes, active, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)`,
		p.ID, p.TenantID, p.ProfessionalID, p.FullName,
		nilIfEmpty(p.Email), nilIfEmpty(p.Phone), nilIfEmpty(p.CPF),
		p.DateOfBirth, nilIfEmpty(p.Gender), nilIfEmpty(p.Notes),
		p.Active, p.CreatedAt, p.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			if strings.Contains(pgErr.ConstraintName, "cpf") {
				return domain.ErrDuplicateCPF
			}
		}
		return fmt.Errorf("patient: create: %w", err)
	}
	return nil
}

// GetByID returns a patient by ID within a tenant.
func (r *Repository) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*domain.Patient, error) {
	p := &domain.Patient{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, tenant_id, professional_id, full_name,
			COALESCE(email,''), COALESCE(phone,''), COALESCE(cpf,''),
			date_of_birth, COALESCE(gender,''), COALESCE(notes,''),
			active, created_at, updated_at
		 FROM patients WHERE id = $1 AND tenant_id = $2`,
		id, tenantID,
	).Scan(&p.ID, &p.TenantID, &p.ProfessionalID, &p.FullName,
		&p.Email, &p.Phone, &p.CPF, &p.DateOfBirth, &p.Gender, &p.Notes,
		&p.Active, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("patient: get_by_id: %w", err)
	}
	return p, nil
}

// List returns patients for a tenant.
func (r *Repository) List(ctx context.Context, tenantID uuid.UUID) ([]domain.Patient, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, tenant_id, professional_id, full_name,
			COALESCE(email,''), COALESCE(phone,''), COALESCE(cpf,''),
			date_of_birth, COALESCE(gender,''), COALESCE(notes,''),
			active, created_at, updated_at
		 FROM patients WHERE tenant_id = $1 ORDER BY full_name`,
		tenantID,
	)
	if err != nil {
		return nil, fmt.Errorf("patient: list: %w", err)
	}
	defer rows.Close()

	var result []domain.Patient
	for rows.Next() {
		var p domain.Patient
		if err := rows.Scan(&p.ID, &p.TenantID, &p.ProfessionalID, &p.FullName,
			&p.Email, &p.Phone, &p.CPF, &p.DateOfBirth, &p.Gender, &p.Notes,
			&p.Active, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("patient: scan: %w", err)
		}
		result = append(result, p)
	}
	return result, rows.Err()
}

// Update updates a patient.
func (r *Repository) Update(ctx context.Context, p *domain.Patient) error {
	tag, err := r.pool.Exec(ctx,
		`UPDATE patients SET
			full_name = $3, email = $4, phone = $5, cpf = $6,
			date_of_birth = $7, gender = $8, notes = $9, active = $10, updated_at = NOW()
		 WHERE id = $1 AND tenant_id = $2`,
		p.ID, p.TenantID, p.FullName,
		nilIfEmpty(p.Email), nilIfEmpty(p.Phone), nilIfEmpty(p.CPF),
		p.DateOfBirth, nilIfEmpty(p.Gender), nilIfEmpty(p.Notes), p.Active,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.ErrDuplicateCPF
		}
		return fmt.Errorf("patient: update: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// CountByTenant returns the number of active patients in a tenant.
func (r *Repository) CountByTenant(ctx context.Context, tenantID uuid.UUID) (int64, error) {
	var count int64
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM patients WHERE tenant_id = $1 AND active = TRUE`,
		tenantID,
	).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("patient: count: %w", err)
	}
	return count, nil
}

// --- Invites ---

// CreateInvite inserts a patient invite.
func (r *Repository) CreateInvite(ctx context.Context, inv *domain.PatientInvite) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO patient_invites
			(id, tenant_id, patient_id, invited_by, code_hash, max_uses, used_count, expires_at, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
		inv.ID, inv.TenantID, inv.PatientID, inv.InvitedBy,
		inv.CodeHash, inv.MaxUses, inv.UsedCount, inv.ExpiresAt, inv.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("patient: create_invite: %w", err)
	}
	return nil
}

// GetInviteByHash returns an invite by its code hash.
func (r *Repository) GetInviteByHash(ctx context.Context, codeHash string) (*domain.PatientInvite, error) {
	inv := &domain.PatientInvite{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, tenant_id, patient_id, invited_by, code_hash, max_uses, used_count, expires_at, created_at
		 FROM patient_invites WHERE code_hash = $1`,
		codeHash,
	).Scan(&inv.ID, &inv.TenantID, &inv.PatientID, &inv.InvitedBy,
		&inv.CodeHash, &inv.MaxUses, &inv.UsedCount, &inv.ExpiresAt, &inv.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrInviteNotFound
		}
		return nil, fmt.Errorf("patient: get_invite_by_hash: %w", err)
	}
	return inv, nil
}

// IncrementInviteUsedCount increments the used_count of an invite.
func (r *Repository) IncrementInviteUsedCount(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE patient_invites SET used_count = used_count + 1 WHERE id = $1`,
		id,
	)
	if err != nil {
		return fmt.Errorf("patient: increment_invite: %w", err)
	}
	return nil
}

// --- Access Links ---

// CreateAccessLink inserts a patient access link.
func (r *Repository) CreateAccessLink(ctx context.Context, link *domain.PatientAccessLink) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO patient_access_links
			(id, tenant_id, patient_id, user_id, invite_id, active, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7)`,
		link.ID, link.TenantID, link.PatientID, link.UserID,
		link.InviteID, link.Active, link.CreatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.ErrAccessLinkExists
		}
		return fmt.Errorf("patient: create_access_link: %w", err)
	}
	return nil
}

// --- Tenant Users ---

// CreateTenantUser inserts a tenant_users membership entry for the patient.
// Uses ON CONFLICT DO NOTHING for idempotency.
func (r *Repository) CreateTenantUser(ctx context.Context, tenantID, userID, invitedBy uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO tenant_users (id, tenant_id, user_id, status, invited_by_user_id)
		 VALUES ($1, $2, $3, 'active', $4)
		 ON CONFLICT (tenant_id, user_id) DO NOTHING`,
		uuid.New(), tenantID, userID, invitedBy,
	)
	if err != nil {
		return fmt.Errorf("patient: create_tenant_user: %w", err)
	}
	return nil
}

// --- Profiles ---

// UpsertProfile creates or updates a patient profile.
func (r *Repository) UpsertProfile(ctx context.Context, p *domain.PatientProfile) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO patient_profiles
			(id, tenant_id, patient_id, occupation, marital_status, ethnicity, blood_type,
			 allergies, chronic_conditions, medications,
			 emergency_contact_name, emergency_contact_phone, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)
		 ON CONFLICT (tenant_id, patient_id)
		 DO UPDATE SET
			occupation = EXCLUDED.occupation,
			marital_status = EXCLUDED.marital_status,
			ethnicity = EXCLUDED.ethnicity,
			blood_type = EXCLUDED.blood_type,
			allergies = EXCLUDED.allergies,
			chronic_conditions = EXCLUDED.chronic_conditions,
			medications = EXCLUDED.medications,
			emergency_contact_name = EXCLUDED.emergency_contact_name,
			emergency_contact_phone = EXCLUDED.emergency_contact_phone,
			updated_at = NOW()`,
		p.ID, p.TenantID, p.PatientID,
		nilIfEmpty(p.Occupation), nilIfEmpty(p.MaritalStatus),
		nilIfEmpty(p.Ethnicity), nilIfEmpty(p.BloodType),
		p.Allergies, p.ChronicConditions, p.Medications,
		nilIfEmpty(p.EmergencyContactName), nilIfEmpty(p.EmergencyContactPhone),
		p.CreatedAt, p.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("patient: upsert_profile: %w", err)
	}
	return nil
}

// GetProfile returns a patient profile.
func (r *Repository) GetProfile(ctx context.Context, tenantID, patientID uuid.UUID) (*domain.PatientProfile, error) {
	p := &domain.PatientProfile{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, tenant_id, patient_id,
			COALESCE(occupation,''), COALESCE(marital_status,''),
			COALESCE(ethnicity,''), COALESCE(blood_type,''),
			COALESCE(allergies,'{}'), COALESCE(chronic_conditions,'{}'), COALESCE(medications,'{}'),
			COALESCE(emergency_contact_name,''), COALESCE(emergency_contact_phone,''),
			created_at, updated_at
		 FROM patient_profiles WHERE tenant_id = $1 AND patient_id = $2`,
		tenantID, patientID,
	).Scan(&p.ID, &p.TenantID, &p.PatientID,
		&p.Occupation, &p.MaritalStatus, &p.Ethnicity, &p.BloodType,
		&p.Allergies, &p.ChronicConditions, &p.Medications,
		&p.EmergencyContactName, &p.EmergencyContactPhone,
		&p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrProfileNotFound
		}
		return nil, fmt.Errorf("patient: get_profile: %w", err)
	}
	return p, nil
}

func nilIfEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
