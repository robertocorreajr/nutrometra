package professional

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"nutrometra/api/internal/professional/domain"
	"nutrometra/api/internal/professional/usecase"

	identitydomain "nutrometra/api/internal/identity/domain"
	"nutrometra/api/internal/platform/audit"
	"nutrometra/api/internal/platform/server"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Handler exposes HTTP endpoints for professional operations.
type Handler struct {
	uc       *usecase.Usecase
	pool     *pgxpool.Pool
	auditSvc *audit.Service
}

// NewHandler creates a professional Handler.
func NewHandler(uc *usecase.Usecase, pool *pgxpool.Pool, auditSvc *audit.Service) *Handler {
	return &Handler{uc: uc, pool: pool, auditSvc: auditSvc}
}

// --- Request/Response types ---

type createProfessionalRequest struct {
	UserID             string `json:"user_id"`
	FullName           string `json:"full_name"`
	RegistrationType   string `json:"registration_type"`
	RegistrationNumber string `json:"registration_number"`
	RegistrationState  string `json:"registration_state"`
	Specialty          string `json:"specialty"`
	Bio                string `json:"bio"`
	Phone              string `json:"phone"`
	AvatarURL          string `json:"avatar_url"`
}

type professionalResponse struct {
	ID                 string `json:"id"`
	TenantID           string `json:"tenant_id"`
	UserID             string `json:"user_id"`
	FullName           string `json:"full_name"`
	RegistrationType   string `json:"registration_type"`
	RegistrationNumber string `json:"registration_number"`
	RegistrationState  string `json:"registration_state,omitempty"`
	Specialty          string `json:"specialty,omitempty"`
	Bio                string `json:"bio,omitempty"`
	Phone              string `json:"phone,omitempty"`
	AvatarURL          string `json:"avatar_url,omitempty"`
	Active             bool   `json:"active"`
	CreatedAt          string `json:"created_at"`
	UpdatedAt          string `json:"updated_at"`
}

func toProfResponse(p *domain.Professional) professionalResponse {
	return professionalResponse{
		ID:                 p.ID.String(),
		TenantID:           p.TenantID.String(),
		UserID:             p.UserID.String(),
		FullName:           p.FullName,
		RegistrationType:   string(p.RegistrationType),
		RegistrationNumber: p.RegistrationNumber,
		RegistrationState:  p.RegistrationState,
		Specialty:          p.Specialty,
		Bio:                p.Bio,
		Phone:              p.Phone,
		AvatarURL:          p.AvatarURL,
		Active:             p.Active,
		CreatedAt:          p.CreatedAt.Format(time.RFC3339),
		UpdatedAt:          p.UpdatedAt.Format(time.RFC3339),
	}
}

// GetMe handles GET /professionals/me — returns the professional for the authenticated user.
func (h *Handler) GetMe(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := identitydomain.TenantIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusBadRequest, "missing_tenant", "No tenant in context")
		return
	}

	userID, ok := identitydomain.UserIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusUnauthorized, "unauthenticated", "Not authenticated")
		return
	}

	p, err := h.uc.GetByUserID(r.Context(), tenantID, userID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			server.RenderError(w, r, http.StatusNotFound, "not_found", "No professional record for this user")
			return
		}
		server.RenderError(w, r, http.StatusInternalServerError, "get_failed", "Failed to get professional")
		return
	}
	server.RenderJSON(w, http.StatusOK, toProfResponse(p))
}

// Create handles POST /professionals
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := identitydomain.TenantIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusBadRequest, "missing_tenant", "No tenant in context")
		return
	}

	var req createProfessionalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_body", "Invalid request body")
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_user_id", "Invalid user_id")
		return
	}

	p := &domain.Professional{
		TenantID:           tenantID,
		UserID:             userID,
		FullName:           req.FullName,
		RegistrationType:   domain.RegistrationType(req.RegistrationType),
		RegistrationNumber: req.RegistrationNumber,
		RegistrationState:  req.RegistrationState,
		Specialty:          req.Specialty,
		Bio:                req.Bio,
		Phone:              req.Phone,
		AvatarURL:          req.AvatarURL,
	}

	if err := h.uc.Create(r.Context(), p); err != nil {
		switch {
		case errors.Is(err, domain.ErrDuplicateUser):
			server.RenderError(w, r, http.StatusConflict, "duplicate_user", "User already registered as professional in this tenant")
		case errors.Is(err, domain.ErrDuplicateRegistration):
			server.RenderError(w, r, http.StatusConflict, "duplicate_registration", "Registration already exists in this tenant")
		case errors.Is(err, domain.ErrEntitlementExceeded):
			server.RenderError(w, r, http.StatusForbidden, "entitlement_exceeded", "Professional limit exceeded for your plan")
		default:
			server.RenderError(w, r, http.StatusInternalServerError, "create_failed", "Failed to create professional")
		}
		return
	}

	actorID, _ := identitydomain.UserIDFromContext(r.Context())
	entry := audit.NewEntry(
		audit.WithTenantID(tenantID),
		audit.WithActor(actorID, audit.ScopeTenant),
		audit.WithEntity("professional", p.ID),
		audit.WithAction("professional_created"),
		audit.WithIPAddress(r.RemoteAddr),
	)
	_ = h.auditSvc.Write(r.Context(), h.pool, entry)

	server.RenderJSON(w, http.StatusCreated, toProfResponse(p))
}

// List handles GET /professionals
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := identitydomain.TenantIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusBadRequest, "missing_tenant", "No tenant in context")
		return
	}

	profs, err := h.uc.List(r.Context(), tenantID)
	if err != nil {
		server.RenderError(w, r, http.StatusInternalServerError, "list_failed", "Failed to list professionals")
		return
	}

	resp := make([]professionalResponse, len(profs))
	for i := range profs {
		resp[i] = toProfResponse(&profs[i])
	}
	server.RenderJSON(w, http.StatusOK, resp)
}

// GetByID handles GET /professionals/{id}
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := identitydomain.TenantIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusBadRequest, "missing_tenant", "No tenant in context")
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_id", "Invalid professional ID")
		return
	}

	p, err := h.uc.GetByID(r.Context(), tenantID, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			server.RenderError(w, r, http.StatusNotFound, "not_found", "Professional not found")
			return
		}
		server.RenderError(w, r, http.StatusInternalServerError, "get_failed", "Failed to get professional")
		return
	}
	server.RenderJSON(w, http.StatusOK, toProfResponse(p))
}

// Update handles PUT /professionals/{id}
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := identitydomain.TenantIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusBadRequest, "missing_tenant", "No tenant in context")
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_id", "Invalid professional ID")
		return
	}

	var req createProfessionalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_body", "Invalid request body")
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_user_id", "Invalid user_id")
		return
	}

	p := &domain.Professional{
		ID:                 id,
		TenantID:           tenantID,
		UserID:             userID,
		FullName:           req.FullName,
		RegistrationType:   domain.RegistrationType(req.RegistrationType),
		RegistrationNumber: req.RegistrationNumber,
		RegistrationState:  req.RegistrationState,
		Specialty:          req.Specialty,
		Bio:                req.Bio,
		Phone:              req.Phone,
		AvatarURL:          req.AvatarURL,
	}

	if err := h.uc.Update(r.Context(), p); err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			server.RenderError(w, r, http.StatusNotFound, "not_found", "Professional not found")
		case errors.Is(err, domain.ErrDuplicateRegistration):
			server.RenderError(w, r, http.StatusConflict, "duplicate_registration", "Registration already exists")
		default:
			server.RenderError(w, r, http.StatusInternalServerError, "update_failed", "Failed to update professional")
		}
		return
	}

	actorID, _ := identitydomain.UserIDFromContext(r.Context())
	entry := audit.NewEntry(
		audit.WithTenantID(tenantID),
		audit.WithActor(actorID, audit.ScopeTenant),
		audit.WithEntity("professional", id),
		audit.WithAction("professional_updated"),
		audit.WithIPAddress(r.RemoteAddr),
	)
	_ = h.auditSvc.Write(r.Context(), h.pool, entry)

	server.RenderJSON(w, http.StatusOK, toProfResponse(p))
}

// --- Address endpoints ---

type createAddressRequest struct {
	Label        string   `json:"label"`
	Street       string   `json:"street"`
	Number       string   `json:"number"`
	Complement   string   `json:"complement"`
	Neighborhood string   `json:"neighborhood"`
	City         string   `json:"city"`
	State        string   `json:"state"`
	ZipCode      string   `json:"zip_code"`
	Country      string   `json:"country"`
	Latitude     *float64 `json:"latitude"`
	Longitude    *float64 `json:"longitude"`
	Phone        string   `json:"phone"`
	Notes        string   `json:"notes"`
}

type addressResponse struct {
	ID             string   `json:"id"`
	TenantID       string   `json:"tenant_id"`
	ProfessionalID string   `json:"professional_id"`
	Label          string   `json:"label"`
	Street         string   `json:"street"`
	Number         string   `json:"number,omitempty"`
	Complement     string   `json:"complement,omitempty"`
	Neighborhood   string   `json:"neighborhood,omitempty"`
	City           string   `json:"city"`
	State          string   `json:"state"`
	ZipCode        string   `json:"zip_code"`
	Country        string   `json:"country"`
	Latitude       *float64 `json:"latitude,omitempty"`
	Longitude      *float64 `json:"longitude,omitempty"`
	Phone          string   `json:"phone,omitempty"`
	Notes          string   `json:"notes,omitempty"`
	Active         bool     `json:"active"`
	CreatedAt      string   `json:"created_at"`
	UpdatedAt      string   `json:"updated_at"`
}

func toAddrResponse(a *domain.Address) addressResponse {
	return addressResponse{
		ID:             a.ID.String(),
		TenantID:       a.TenantID.String(),
		ProfessionalID: a.ProfessionalID.String(),
		Label:          a.Label,
		Street:         a.Street,
		Number:         a.Number,
		Complement:     a.Complement,
		Neighborhood:   a.Neighborhood,
		City:           a.City,
		State:          a.State,
		ZipCode:        a.ZipCode,
		Country:        a.Country,
		Latitude:       a.Latitude,
		Longitude:      a.Longitude,
		Phone:          a.Phone,
		Notes:          a.Notes,
		Active:         a.Active,
		CreatedAt:      a.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      a.UpdatedAt.Format(time.RFC3339),
	}
}

// CreateAddress handles POST /professionals/{id}/addresses
func (h *Handler) CreateAddress(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := identitydomain.TenantIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusBadRequest, "missing_tenant", "No tenant in context")
		return
	}

	profID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_id", "Invalid professional ID")
		return
	}

	var req createAddressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_body", "Invalid request body")
		return
	}

	a := &domain.Address{
		TenantID:       tenantID,
		ProfessionalID: profID,
		Label:          req.Label,
		Street:         req.Street,
		Number:         req.Number,
		Complement:     req.Complement,
		Neighborhood:   req.Neighborhood,
		City:           req.City,
		State:          req.State,
		ZipCode:        req.ZipCode,
		Country:        req.Country,
		Latitude:       req.Latitude,
		Longitude:      req.Longitude,
		Phone:          req.Phone,
		Notes:          req.Notes,
	}

	if err := h.uc.CreateAddress(r.Context(), a); err != nil {
		if errors.Is(err, domain.ErrEntitlementExceeded) {
			server.RenderError(w, r, http.StatusForbidden, "entitlement_exceeded", "Address limit exceeded for your plan")
			return
		}
		server.RenderError(w, r, http.StatusInternalServerError, "create_address_failed", "Failed to create address")
		return
	}

	actorID, _ := identitydomain.UserIDFromContext(r.Context())
	entry := audit.NewEntry(
		audit.WithTenantID(tenantID),
		audit.WithActor(actorID, audit.ScopeTenant),
		audit.WithEntity("professional_address", a.ID),
		audit.WithAction("address_created"),
		audit.WithIPAddress(r.RemoteAddr),
	)
	_ = h.auditSvc.Write(r.Context(), h.pool, entry)

	server.RenderJSON(w, http.StatusCreated, toAddrResponse(a))
}

// ListAddresses handles GET /professionals/{id}/addresses
func (h *Handler) ListAddresses(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := identitydomain.TenantIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusBadRequest, "missing_tenant", "No tenant in context")
		return
	}

	profID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_id", "Invalid professional ID")
		return
	}

	addrs, err := h.uc.ListAddresses(r.Context(), tenantID, profID)
	if err != nil {
		server.RenderError(w, r, http.StatusInternalServerError, "list_addresses_failed", "Failed to list addresses")
		return
	}

	resp := make([]addressResponse, len(addrs))
	for i := range addrs {
		resp[i] = toAddrResponse(&addrs[i])
	}
	server.RenderJSON(w, http.StatusOK, resp)
}

// UpdateAddress handles PUT /professionals/{id}/addresses/{addr_id}
func (h *Handler) UpdateAddress(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := identitydomain.TenantIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusBadRequest, "missing_tenant", "No tenant in context")
		return
	}

	profID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_id", "Invalid professional ID")
		return
	}

	addrID, err := uuid.Parse(chi.URLParam(r, "addr_id"))
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_addr_id", "Invalid address ID")
		return
	}

	var req createAddressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_body", "Invalid request body")
		return
	}

	a := &domain.Address{
		ID:             addrID,
		TenantID:       tenantID,
		ProfessionalID: profID,
		Label:          req.Label,
		Street:         req.Street,
		Number:         req.Number,
		Complement:     req.Complement,
		Neighborhood:   req.Neighborhood,
		City:           req.City,
		State:          req.State,
		ZipCode:        req.ZipCode,
		Country:        req.Country,
		Latitude:       req.Latitude,
		Longitude:      req.Longitude,
		Phone:          req.Phone,
		Notes:          req.Notes,
	}

	if err := h.uc.UpdateAddress(r.Context(), a); err != nil {
		if errors.Is(err, domain.ErrAddressNotFound) {
			server.RenderError(w, r, http.StatusNotFound, "not_found", "Address not found")
			return
		}
		server.RenderError(w, r, http.StatusInternalServerError, "update_address_failed", "Failed to update address")
		return
	}

	actorID, _ := identitydomain.UserIDFromContext(r.Context())
	entry := audit.NewEntry(
		audit.WithTenantID(tenantID),
		audit.WithActor(actorID, audit.ScopeTenant),
		audit.WithEntity("professional_address", addrID),
		audit.WithAction("address_updated"),
		audit.WithIPAddress(r.RemoteAddr),
	)
	_ = h.auditSvc.Write(r.Context(), h.pool, entry)

	server.RenderJSON(w, http.StatusOK, toAddrResponse(a))
}

// --- Service Mode endpoints ---

type setServiceModeRequest struct {
	Mode        string  `json:"mode"`
	AddressID   *string `json:"address_id"`
	DurationMin int     `json:"duration_min"`
}

type serviceModeResponse struct {
	ID             string  `json:"id"`
	TenantID       string  `json:"tenant_id"`
	ProfessionalID string  `json:"professional_id"`
	Mode           string  `json:"mode"`
	AddressID      *string `json:"address_id,omitempty"`
	DurationMin    int     `json:"duration_min"`
	Active         bool    `json:"active"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
}

func toSMResponse(sm *domain.ProfessionalServiceMode) serviceModeResponse {
	r := serviceModeResponse{
		ID:             sm.ID.String(),
		TenantID:       sm.TenantID.String(),
		ProfessionalID: sm.ProfessionalID.String(),
		Mode:           string(sm.Mode),
		DurationMin:    sm.DurationMin,
		Active:         sm.Active,
		CreatedAt:      sm.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      sm.UpdatedAt.Format(time.RFC3339),
	}
	if sm.AddressID != nil {
		s := sm.AddressID.String()
		r.AddressID = &s
	}
	return r
}

// SetServiceMode handles PUT /professionals/{id}/service-modes
func (h *Handler) SetServiceMode(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := identitydomain.TenantIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusBadRequest, "missing_tenant", "No tenant in context")
		return
	}

	profID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_id", "Invalid professional ID")
		return
	}

	var req setServiceModeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_body", "Invalid request body")
		return
	}

	sm := &domain.ProfessionalServiceMode{
		TenantID:       tenantID,
		ProfessionalID: profID,
		Mode:           domain.ServiceMode(req.Mode),
		DurationMin:    req.DurationMin,
	}
	if req.AddressID != nil {
		aid, err := uuid.Parse(*req.AddressID)
		if err != nil {
			server.RenderError(w, r, http.StatusBadRequest, "invalid_address_id", "Invalid address_id")
			return
		}
		sm.AddressID = &aid
	}

	if err := h.uc.SetServiceMode(r.Context(), sm); err != nil {
		server.RenderError(w, r, http.StatusInternalServerError, "set_service_mode_failed", "Failed to set service mode")
		return
	}

	actorID, _ := identitydomain.UserIDFromContext(r.Context())
	entry := audit.NewEntry(
		audit.WithTenantID(tenantID),
		audit.WithActor(actorID, audit.ScopeTenant),
		audit.WithEntity("professional_service_mode", sm.ID),
		audit.WithAction("service_mode_set"),
		audit.WithIPAddress(r.RemoteAddr),
	)
	_ = h.auditSvc.Write(r.Context(), h.pool, entry)

	server.RenderJSON(w, http.StatusOK, toSMResponse(sm))
}

// ListServiceModes handles GET /professionals/{id}/service-modes
func (h *Handler) ListServiceModes(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := identitydomain.TenantIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusBadRequest, "missing_tenant", "No tenant in context")
		return
	}

	profID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_id", "Invalid professional ID")
		return
	}

	modes, err := h.uc.ListServiceModes(r.Context(), tenantID, profID)
	if err != nil {
		server.RenderError(w, r, http.StatusInternalServerError, "list_service_modes_failed", "Failed to list service modes")
		return
	}

	resp := make([]serviceModeResponse, len(modes))
	for i := range modes {
		resp[i] = toSMResponse(&modes[i])
	}
	server.RenderJSON(w, http.StatusOK, resp)
}
