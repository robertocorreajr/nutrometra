package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"nutrometra/api/internal/professional/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository provides professional data access.
type Repository struct {
	pool *pgxpool.Pool
}

// New creates a Repository.
func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// Create inserts a new professional.
func (r *Repository) Create(ctx context.Context, p *domain.Professional) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO professionals
			(id, tenant_id, user_id, full_name, registration_type, registration_number,
			 registration_state, specialty, bio, phone, avatar_url, active, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)`,
		p.ID, p.TenantID, p.UserID, p.FullName, string(p.RegistrationType), p.RegistrationNumber,
		nilIfEmpty(p.RegistrationState), nilIfEmpty(p.Specialty), nilIfEmpty(p.Bio),
		nilIfEmpty(p.Phone), nilIfEmpty(p.AvatarURL), p.Active, p.CreatedAt, p.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			if strings.Contains(pgErr.ConstraintName, "user_id") {
				return domain.ErrDuplicateUser
			}
			return domain.ErrDuplicateRegistration
		}
		return fmt.Errorf("professional: create: %w", err)
	}
	return nil
}

// GetByID returns a professional by ID within a tenant.
func (r *Repository) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*domain.Professional, error) {
	p := &domain.Professional{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, tenant_id, user_id, full_name, registration_type, registration_number,
			COALESCE(registration_state,''), COALESCE(specialty,''), COALESCE(bio,''),
			COALESCE(phone,''), COALESCE(avatar_url,''), active, created_at, updated_at
		 FROM professionals WHERE id = $1 AND tenant_id = $2`,
		id, tenantID,
	).Scan(&p.ID, &p.TenantID, &p.UserID, &p.FullName, &p.RegistrationType, &p.RegistrationNumber,
		&p.RegistrationState, &p.Specialty, &p.Bio, &p.Phone, &p.AvatarURL, &p.Active,
		&p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("professional: get_by_id: %w", err)
	}
	return p, nil
}

// GetByUserID returns the professional linked to a user within a tenant.
func (r *Repository) GetByUserID(ctx context.Context, tenantID, userID uuid.UUID) (*domain.Professional, error) {
	p := &domain.Professional{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, tenant_id, user_id, full_name, registration_type, registration_number,
			COALESCE(registration_state,''), COALESCE(specialty,''), COALESCE(bio,''),
			COALESCE(phone,''), COALESCE(avatar_url,''), active, created_at, updated_at
		 FROM professionals WHERE tenant_id = $1 AND user_id = $2`,
		tenantID, userID,
	).Scan(&p.ID, &p.TenantID, &p.UserID, &p.FullName, &p.RegistrationType, &p.RegistrationNumber,
		&p.RegistrationState, &p.Specialty, &p.Bio, &p.Phone, &p.AvatarURL, &p.Active,
		&p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("professional: get_by_user_id: %w", err)
	}
	return p, nil
}

// List returns all professionals for a tenant.
func (r *Repository) List(ctx context.Context, tenantID uuid.UUID) ([]domain.Professional, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, tenant_id, user_id, full_name, registration_type, registration_number,
			COALESCE(registration_state,''), COALESCE(specialty,''), COALESCE(bio,''),
			COALESCE(phone,''), COALESCE(avatar_url,''), active, created_at, updated_at
		 FROM professionals WHERE tenant_id = $1 ORDER BY full_name`,
		tenantID,
	)
	if err != nil {
		return nil, fmt.Errorf("professional: list: %w", err)
	}
	defer rows.Close()

	var result []domain.Professional
	for rows.Next() {
		var p domain.Professional
		if err := rows.Scan(&p.ID, &p.TenantID, &p.UserID, &p.FullName, &p.RegistrationType,
			&p.RegistrationNumber, &p.RegistrationState, &p.Specialty, &p.Bio, &p.Phone,
			&p.AvatarURL, &p.Active, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("professional: scan: %w", err)
		}
		result = append(result, p)
	}
	return result, rows.Err()
}

// Update updates a professional's mutable fields.
func (r *Repository) Update(ctx context.Context, p *domain.Professional) error {
	tag, err := r.pool.Exec(ctx,
		`UPDATE professionals SET
			full_name = $3, registration_type = $4, registration_number = $5,
			registration_state = $6, specialty = $7, bio = $8, phone = $9,
			avatar_url = $10, active = $11, updated_at = NOW()
		 WHERE id = $1 AND tenant_id = $2`,
		p.ID, p.TenantID, p.FullName, string(p.RegistrationType), p.RegistrationNumber,
		nilIfEmpty(p.RegistrationState), nilIfEmpty(p.Specialty), nilIfEmpty(p.Bio),
		nilIfEmpty(p.Phone), nilIfEmpty(p.AvatarURL), p.Active,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.ErrDuplicateRegistration
		}
		return fmt.Errorf("professional: update: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// CountByTenant returns the number of active professionals in a tenant.
func (r *Repository) CountByTenant(ctx context.Context, tenantID uuid.UUID) (int64, error) {
	var count int64
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM professionals WHERE tenant_id = $1 AND active = TRUE`,
		tenantID,
	).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("professional: count: %w", err)
	}
	return count, nil
}

// --- Addresses ---

// CreateAddress inserts a professional address.
func (r *Repository) CreateAddress(ctx context.Context, a *domain.Address) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO professional_addresses
			(id, tenant_id, professional_id, label, street, number, complement,
			 neighborhood, city, state, zip_code, country, latitude, longitude,
			 phone, notes, active, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19)`,
		a.ID, a.TenantID, a.ProfessionalID, a.Label, a.Street,
		nilIfEmpty(a.Number), nilIfEmpty(a.Complement), nilIfEmpty(a.Neighborhood),
		a.City, a.State, a.ZipCode, a.Country, a.Latitude, a.Longitude,
		nilIfEmpty(a.Phone), nilIfEmpty(a.Notes), a.Active, a.CreatedAt, a.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("professional: create_address: %w", err)
	}
	return nil
}

// ListAddresses returns addresses for a professional within a tenant.
func (r *Repository) ListAddresses(ctx context.Context, tenantID, professionalID uuid.UUID) ([]domain.Address, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, tenant_id, professional_id, label, street,
			COALESCE(number,''), COALESCE(complement,''), COALESCE(neighborhood,''),
			city, state, zip_code, country, latitude, longitude,
			COALESCE(phone,''), COALESCE(notes,''), active, created_at, updated_at
		 FROM professional_addresses
		 WHERE tenant_id = $1 AND professional_id = $2
		 ORDER BY label`,
		tenantID, professionalID,
	)
	if err != nil {
		return nil, fmt.Errorf("professional: list_addresses: %w", err)
	}
	defer rows.Close()

	var result []domain.Address
	for rows.Next() {
		var a domain.Address
		if err := rows.Scan(&a.ID, &a.TenantID, &a.ProfessionalID, &a.Label, &a.Street,
			&a.Number, &a.Complement, &a.Neighborhood, &a.City, &a.State, &a.ZipCode,
			&a.Country, &a.Latitude, &a.Longitude, &a.Phone, &a.Notes, &a.Active,
			&a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, fmt.Errorf("professional: scan_address: %w", err)
		}
		result = append(result, a)
	}
	return result, rows.Err()
}

// UpdateAddress updates a professional address.
func (r *Repository) UpdateAddress(ctx context.Context, a *domain.Address) error {
	tag, err := r.pool.Exec(ctx,
		`UPDATE professional_addresses SET
			label = $4, street = $5, number = $6, complement = $7, neighborhood = $8,
			city = $9, state = $10, zip_code = $11, country = $12, latitude = $13,
			longitude = $14, phone = $15, notes = $16, active = $17, updated_at = NOW()
		 WHERE id = $1 AND tenant_id = $2 AND professional_id = $3`,
		a.ID, a.TenantID, a.ProfessionalID, a.Label, a.Street,
		nilIfEmpty(a.Number), nilIfEmpty(a.Complement), nilIfEmpty(a.Neighborhood),
		a.City, a.State, a.ZipCode, a.Country, a.Latitude, a.Longitude,
		nilIfEmpty(a.Phone), nilIfEmpty(a.Notes), a.Active,
	)
	if err != nil {
		return fmt.Errorf("professional: update_address: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrAddressNotFound
	}
	return nil
}

// CountAddresses returns the number of active addresses for a professional.
func (r *Repository) CountAddresses(ctx context.Context, tenantID, professionalID uuid.UUID) (int64, error) {
	var count int64
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM professional_addresses WHERE tenant_id = $1 AND professional_id = $2 AND active = TRUE`,
		tenantID, professionalID,
	).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("professional: count_addresses: %w", err)
	}
	return count, nil
}

// --- Service Modes ---

// SetServiceMode upserts a service mode for a professional.
func (r *Repository) SetServiceMode(ctx context.Context, sm *domain.ProfessionalServiceMode) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO professional_service_modes
			(id, tenant_id, professional_id, mode, address_id, duration_min, active, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		 ON CONFLICT (professional_id, mode, address_id)
		 DO UPDATE SET duration_min = EXCLUDED.duration_min, active = EXCLUDED.active, updated_at = NOW()`,
		sm.ID, sm.TenantID, sm.ProfessionalID, string(sm.Mode), sm.AddressID,
		sm.DurationMin, sm.Active, sm.CreatedAt, sm.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("professional: set_service_mode: %w", err)
	}
	return nil
}

// ListServiceModes returns service modes for a professional.
func (r *Repository) ListServiceModes(ctx context.Context, tenantID, professionalID uuid.UUID) ([]domain.ProfessionalServiceMode, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, tenant_id, professional_id, mode, address_id, duration_min, active, created_at, updated_at
		 FROM professional_service_modes
		 WHERE tenant_id = $1 AND professional_id = $2
		 ORDER BY mode`,
		tenantID, professionalID,
	)
	if err != nil {
		return nil, fmt.Errorf("professional: list_service_modes: %w", err)
	}
	defer rows.Close()

	var result []domain.ProfessionalServiceMode
	for rows.Next() {
		var sm domain.ProfessionalServiceMode
		if err := rows.Scan(&sm.ID, &sm.TenantID, &sm.ProfessionalID, &sm.Mode, &sm.AddressID,
			&sm.DurationMin, &sm.Active, &sm.CreatedAt, &sm.UpdatedAt); err != nil {
			return nil, fmt.Errorf("professional: scan_service_mode: %w", err)
		}
		result = append(result, sm)
	}
	return result, rows.Err()
}

func nilIfEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
