package scheduling

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"nutrometra/api/internal/scheduling/domain"
	"nutrometra/api/internal/scheduling/usecase"

	identitydomain "nutrometra/api/internal/identity/domain"
	"nutrometra/api/internal/platform/audit"
	"nutrometra/api/internal/platform/server"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Handler exposes HTTP endpoints for scheduling operations.
type Handler struct {
	uc       *usecase.Usecase
	pool     *pgxpool.Pool
	auditSvc *audit.Service
}

// NewHandler creates a scheduling Handler.
func NewHandler(uc *usecase.Usecase, pool *pgxpool.Pool, auditSvc *audit.Service) *Handler {
	return &Handler{uc: uc, pool: pool, auditSvc: auditSvc}
}

// --- Availability Rules ---

type availabilityRuleRequest struct {
	DayOfWeek   int    `json:"day_of_week"`
	StartTime   string `json:"start_time"` // HH:MM
	EndTime     string `json:"end_time"`   // HH:MM
	ServiceMode string `json:"service_mode"`
	AddressID   *string `json:"address_id"`
}

type availabilityRuleResponse struct {
	ID             string  `json:"id"`
	ProfessionalID string  `json:"professional_id"`
	DayOfWeek      int     `json:"day_of_week"`
	StartTime      string  `json:"start_time"`
	EndTime        string  `json:"end_time"`
	ServiceMode    string  `json:"service_mode"`
	AddressID      *string `json:"address_id,omitempty"`
	Active         bool    `json:"active"`
}

func toRuleResponse(r *domain.AvailabilityRule) availabilityRuleResponse {
	resp := availabilityRuleResponse{
		ID:             r.ID.String(),
		ProfessionalID: r.ProfessionalID.String(),
		DayOfWeek:      r.DayOfWeek,
		StartTime:      r.StartTime.Format("15:04"),
		EndTime:        r.EndTime.Format("15:04"),
		ServiceMode:    string(r.ServiceMode),
		Active:         r.Active,
	}
	if r.AddressID != nil {
		s := r.AddressID.String()
		resp.AddressID = &s
	}
	return resp
}

// CreateAvailabilityRule handles POST /professionals/{prof_id}/availability
func (h *Handler) CreateAvailabilityRule(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := identitydomain.TenantIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusBadRequest, "missing_tenant", "No tenant in context")
		return
	}

	profID, err := uuid.Parse(chi.URLParam(r, "prof_id"))
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_prof_id", "Invalid professional ID")
		return
	}

	var req availabilityRuleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_body", "Invalid request body")
		return
	}

	startTime, err := time.Parse("15:04", req.StartTime)
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_start_time", "start_time must be HH:MM")
		return
	}
	endTime, err := time.Parse("15:04", req.EndTime)
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_end_time", "end_time must be HH:MM")
		return
	}

	rule := &domain.AvailabilityRule{
		TenantID:       tenantID,
		ProfessionalID: profID,
		DayOfWeek:      req.DayOfWeek,
		StartTime:      startTime,
		EndTime:        endTime,
		ServiceMode:    domain.ServiceMode(req.ServiceMode),
	}

	if req.AddressID != nil {
		aid, err := uuid.Parse(*req.AddressID)
		if err != nil {
			server.RenderError(w, r, http.StatusBadRequest, "invalid_address_id", "Invalid address_id")
			return
		}
		rule.AddressID = &aid
	}

	if err := h.uc.CreateAvailabilityRule(r.Context(), rule); err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "create_rule_failed", err.Error())
		return
	}

	actorID, _ := identitydomain.UserIDFromContext(r.Context())
	entry := audit.NewEntry(
		audit.WithTenantID(tenantID),
		audit.WithActor(actorID, audit.ScopeTenant),
		audit.WithEntity("availability_rule", rule.ID),
		audit.WithAction("availability_rule_created"),
		audit.WithIPAddress(r.RemoteAddr),
	)
	_ = h.auditSvc.Write(r.Context(), h.pool, entry)

	server.RenderJSON(w, http.StatusCreated, toRuleResponse(rule))
}

// ListAvailabilityRules handles GET /professionals/{prof_id}/availability
func (h *Handler) ListAvailabilityRules(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := identitydomain.TenantIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusBadRequest, "missing_tenant", "No tenant in context")
		return
	}

	profID, err := uuid.Parse(chi.URLParam(r, "prof_id"))
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_prof_id", "Invalid professional ID")
		return
	}

	rules, err := h.uc.ListAvailabilityRules(r.Context(), tenantID, profID)
	if err != nil {
		server.RenderError(w, r, http.StatusInternalServerError, "list_rules_failed", "Failed to list availability rules")
		return
	}

	resp := make([]availabilityRuleResponse, len(rules))
	for i := range rules {
		resp[i] = toRuleResponse(&rules[i])
	}
	server.RenderJSON(w, http.StatusOK, resp)
}

// UpdateAvailabilityRule handles PUT /availability/{id}
func (h *Handler) UpdateAvailabilityRule(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := identitydomain.TenantIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusBadRequest, "missing_tenant", "No tenant in context")
		return
	}

	ruleID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_id", "Invalid rule ID")
		return
	}

	var req availabilityRuleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_body", "Invalid request body")
		return
	}

	startTime, err := time.Parse("15:04", req.StartTime)
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_start_time", "start_time must be HH:MM")
		return
	}
	endTime, err := time.Parse("15:04", req.EndTime)
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_end_time", "end_time must be HH:MM")
		return
	}

	rule := &domain.AvailabilityRule{
		ID:          ruleID,
		TenantID:    tenantID,
		DayOfWeek:   req.DayOfWeek,
		StartTime:   startTime,
		EndTime:     endTime,
		ServiceMode: domain.ServiceMode(req.ServiceMode),
	}

	if req.AddressID != nil {
		aid, err := uuid.Parse(*req.AddressID)
		if err != nil {
			server.RenderError(w, r, http.StatusBadRequest, "invalid_address_id", "Invalid address_id")
			return
		}
		rule.AddressID = &aid
	}

	if err := h.uc.UpdateAvailabilityRule(r.Context(), rule); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			server.RenderError(w, r, http.StatusNotFound, "not_found", "Rule not found")
			return
		}
		server.RenderError(w, r, http.StatusInternalServerError, "update_rule_failed", "Failed to update rule")
		return
	}

	server.RenderJSON(w, http.StatusOK, toRuleResponse(rule))
}

// DeleteAvailabilityRule handles DELETE /availability/{id}
func (h *Handler) DeleteAvailabilityRule(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := identitydomain.TenantIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusBadRequest, "missing_tenant", "No tenant in context")
		return
	}

	ruleID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_id", "Invalid rule ID")
		return
	}

	if err := h.uc.DeleteAvailabilityRule(r.Context(), tenantID, ruleID); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			server.RenderError(w, r, http.StatusNotFound, "not_found", "Rule not found")
			return
		}
		server.RenderError(w, r, http.StatusInternalServerError, "delete_rule_failed", "Failed to delete rule")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetAvailableSlots handles GET /professionals/{prof_id}/slots?date=2025-01-15&duration=50
func (h *Handler) GetAvailableSlots(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := identitydomain.TenantIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusBadRequest, "missing_tenant", "No tenant in context")
		return
	}

	profID, err := uuid.Parse(chi.URLParam(r, "prof_id"))
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_prof_id", "Invalid professional ID")
		return
	}

	dateStr := r.URL.Query().Get("date")
	if dateStr == "" {
		server.RenderError(w, r, http.StatusBadRequest, "missing_date", "date query param is required (YYYY-MM-DD)")
		return
	}
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_date", "date must be YYYY-MM-DD")
		return
	}

	duration := 50
	if d := r.URL.Query().Get("duration"); d != "" {
		duration, err = strconv.Atoi(d)
		if err != nil || duration <= 0 {
			server.RenderError(w, r, http.StatusBadRequest, "invalid_duration", "duration must be a positive integer")
			return
		}
	}

	slots, err := h.uc.GetAvailableSlots(r.Context(), tenantID, profID, date, duration)
	if err != nil {
		server.RenderError(w, r, http.StatusInternalServerError, "slots_failed", "Failed to compute available slots")
		return
	}
	server.RenderJSON(w, http.StatusOK, slots)
}

// --- Schedule Blocks ---

type createBlockRequest struct {
	StartAt string `json:"start_at"` // RFC3339
	EndAt   string `json:"end_at"`   // RFC3339
	Reason  string `json:"reason"`
	AllDay  bool   `json:"all_day"`
}

type blockResponse struct {
	ID             string `json:"id"`
	ProfessionalID string `json:"professional_id"`
	StartAt        string `json:"start_at"`
	EndAt          string `json:"end_at"`
	Reason         string `json:"reason,omitempty"`
	AllDay         bool   `json:"all_day"`
}

func toBlockResponse(b *domain.ScheduleBlock) blockResponse {
	return blockResponse{
		ID:             b.ID.String(),
		ProfessionalID: b.ProfessionalID.String(),
		StartAt:        b.StartAt.Format(time.RFC3339),
		EndAt:          b.EndAt.Format(time.RFC3339),
		Reason:         b.Reason,
		AllDay:         b.AllDay,
	}
}

// CreateBlock handles POST /professionals/{prof_id}/blocks
func (h *Handler) CreateBlock(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := identitydomain.TenantIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusBadRequest, "missing_tenant", "No tenant in context")
		return
	}

	profID, err := uuid.Parse(chi.URLParam(r, "prof_id"))
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_prof_id", "Invalid professional ID")
		return
	}

	var req createBlockRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_body", "Invalid request body")
		return
	}

	startAt, err := time.Parse(time.RFC3339, req.StartAt)
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_start_at", "start_at must be RFC3339")
		return
	}
	endAt, err := time.Parse(time.RFC3339, req.EndAt)
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_end_at", "end_at must be RFC3339")
		return
	}

	b := &domain.ScheduleBlock{
		TenantID:       tenantID,
		ProfessionalID: profID,
		StartAt:        startAt,
		EndAt:          endAt,
		Reason:         req.Reason,
		AllDay:         req.AllDay,
	}

	if err := h.uc.CreateBlock(r.Context(), b); err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "create_block_failed", err.Error())
		return
	}

	actorID, _ := identitydomain.UserIDFromContext(r.Context())
	entry := audit.NewEntry(
		audit.WithTenantID(tenantID),
		audit.WithActor(actorID, audit.ScopeTenant),
		audit.WithEntity("schedule_block", b.ID),
		audit.WithAction("schedule_block_created"),
		audit.WithIPAddress(r.RemoteAddr),
	)
	_ = h.auditSvc.Write(r.Context(), h.pool, entry)

	server.RenderJSON(w, http.StatusCreated, toBlockResponse(b))
}

// ListBlocks handles GET /professionals/{prof_id}/blocks?from=...&to=...
func (h *Handler) ListBlocks(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := identitydomain.TenantIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusBadRequest, "missing_tenant", "No tenant in context")
		return
	}

	profID, err := uuid.Parse(chi.URLParam(r, "prof_id"))
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_prof_id", "Invalid professional ID")
		return
	}

	from := time.Now().UTC()
	to := from.Add(30 * 24 * time.Hour) // default 30 days

	if f := r.URL.Query().Get("from"); f != "" {
		from, err = time.Parse(time.RFC3339, f)
		if err != nil {
			server.RenderError(w, r, http.StatusBadRequest, "invalid_from", "from must be RFC3339")
			return
		}
	}
	if t := r.URL.Query().Get("to"); t != "" {
		to, err = time.Parse(time.RFC3339, t)
		if err != nil {
			server.RenderError(w, r, http.StatusBadRequest, "invalid_to", "to must be RFC3339")
			return
		}
	}

	blocks, err := h.uc.ListBlocks(r.Context(), tenantID, profID, from, to)
	if err != nil {
		server.RenderError(w, r, http.StatusInternalServerError, "list_blocks_failed", "Failed to list blocks")
		return
	}

	resp := make([]blockResponse, len(blocks))
	for i := range blocks {
		resp[i] = toBlockResponse(&blocks[i])
	}
	server.RenderJSON(w, http.StatusOK, resp)
}

// DeleteBlock handles DELETE /blocks/{id}
func (h *Handler) DeleteBlock(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := identitydomain.TenantIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusBadRequest, "missing_tenant", "No tenant in context")
		return
	}

	blockID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_id", "Invalid block ID")
		return
	}

	if err := h.uc.DeleteBlock(r.Context(), tenantID, blockID); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			server.RenderError(w, r, http.StatusNotFound, "not_found", "Block not found")
			return
		}
		server.RenderError(w, r, http.StatusInternalServerError, "delete_block_failed", "Failed to delete block")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// --- Appointments ---

type createAppointmentRequest struct {
	ProfessionalID string  `json:"professional_id"`
	PatientID      *string `json:"patient_id"`
	StartAt        string  `json:"start_at"`
	EndAt          string  `json:"end_at"`
	ServiceMode    string  `json:"service_mode"`
	AddressID      *string `json:"address_id"`
	Source         string  `json:"source"`
	Notes          string  `json:"notes"`
}

type appointmentResponse struct {
	ID                 string  `json:"id"`
	TenantID           string  `json:"tenant_id"`
	ProfessionalID     string  `json:"professional_id"`
	PatientID          *string `json:"patient_id,omitempty"`
	StartAt            string  `json:"start_at"`
	EndAt              string  `json:"end_at"`
	ServiceMode        string  `json:"service_mode"`
	AddressID          *string `json:"address_id,omitempty"`
	Status             string  `json:"status"`
	Source             string  `json:"source"`
	Notes              string  `json:"notes,omitempty"`
	CancellationReason string  `json:"cancellation_reason,omitempty"`
	CreatedAt          string  `json:"created_at"`
	UpdatedAt          string  `json:"updated_at"`
}

func toApptResponse(a *domain.Appointment) appointmentResponse {
	resp := appointmentResponse{
		ID:                 a.ID.String(),
		TenantID:           a.TenantID.String(),
		ProfessionalID:     a.ProfessionalID.String(),
		StartAt:            a.StartAt.Format(time.RFC3339),
		EndAt:              a.EndAt.Format(time.RFC3339),
		ServiceMode:        string(a.ServiceMode),
		Status:             string(a.Status),
		Source:             string(a.Source),
		Notes:              a.Notes,
		CancellationReason: a.CancellationReason,
		CreatedAt:          a.CreatedAt.Format(time.RFC3339),
		UpdatedAt:          a.UpdatedAt.Format(time.RFC3339),
	}
	if a.PatientID != nil {
		s := a.PatientID.String()
		resp.PatientID = &s
	}
	if a.AddressID != nil {
		s := a.AddressID.String()
		resp.AddressID = &s
	}
	return resp
}

// CreateAppointment handles POST /appointments
func (h *Handler) CreateAppointment(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := identitydomain.TenantIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusBadRequest, "missing_tenant", "No tenant in context")
		return
	}

	var req createAppointmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_body", "Invalid request body")
		return
	}

	profID, err := uuid.Parse(req.ProfessionalID)
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_professional_id", "Invalid professional_id")
		return
	}
	startAt, err := time.Parse(time.RFC3339, req.StartAt)
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_start_at", "start_at must be RFC3339")
		return
	}
	endAt, err := time.Parse(time.RFC3339, req.EndAt)
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_end_at", "end_at must be RFC3339")
		return
	}

	a := &domain.Appointment{
		TenantID:       tenantID,
		ProfessionalID: profID,
		StartAt:        startAt,
		EndAt:          endAt,
		ServiceMode:    domain.ServiceMode(req.ServiceMode),
		Source:         domain.AppointmentSource(req.Source),
		Notes:          req.Notes,
	}

	if req.PatientID != nil {
		pid, err := uuid.Parse(*req.PatientID)
		if err != nil {
			server.RenderError(w, r, http.StatusBadRequest, "invalid_patient_id", "Invalid patient_id")
			return
		}
		a.PatientID = &pid
	}
	if req.AddressID != nil {
		aid, err := uuid.Parse(*req.AddressID)
		if err != nil {
			server.RenderError(w, r, http.StatusBadRequest, "invalid_address_id", "Invalid address_id")
			return
		}
		a.AddressID = &aid
	}
	if a.Source == "" {
		a.Source = domain.SourceProfessional
	}

	actorID, _ := identitydomain.UserIDFromContext(r.Context())
	if err := h.uc.CreateAppointment(r.Context(), a, &actorID); err != nil {
		if errors.Is(err, domain.ErrConflict) {
			server.RenderError(w, r, http.StatusConflict, "time_conflict", "Time slot conflicts with an existing appointment")
			return
		}
		server.RenderError(w, r, http.StatusInternalServerError, "create_appointment_failed", "Failed to create appointment")
		return
	}

	entry := audit.NewEntry(
		audit.WithTenantID(tenantID),
		audit.WithActor(actorID, audit.ScopeTenant),
		audit.WithEntity("appointment", a.ID),
		audit.WithAction("appointment_created"),
		audit.WithIPAddress(r.RemoteAddr),
	)
	_ = h.auditSvc.Write(r.Context(), h.pool, entry)

	server.RenderJSON(w, http.StatusCreated, toApptResponse(a))
}

// ListAppointments handles GET /appointments?professional_id=...&from=...&to=...&status=...
func (h *Handler) ListAppointments(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := identitydomain.TenantIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusBadRequest, "missing_tenant", "No tenant in context")
		return
	}

	var profID *uuid.UUID
	if p := r.URL.Query().Get("professional_id"); p != "" {
		id, err := uuid.Parse(p)
		if err != nil {
			server.RenderError(w, r, http.StatusBadRequest, "invalid_professional_id", "Invalid professional_id")
			return
		}
		profID = &id
	}

	var from, to *time.Time
	if f := r.URL.Query().Get("from"); f != "" {
		t, err := time.Parse(time.RFC3339, f)
		if err != nil {
			server.RenderError(w, r, http.StatusBadRequest, "invalid_from", "from must be RFC3339")
			return
		}
		from = &t
	}
	if t := r.URL.Query().Get("to"); t != "" {
		parsed, err := time.Parse(time.RFC3339, t)
		if err != nil {
			server.RenderError(w, r, http.StatusBadRequest, "invalid_to", "to must be RFC3339")
			return
		}
		to = &parsed
	}

	var status *domain.AppointmentStatus
	if s := r.URL.Query().Get("status"); s != "" {
		st := domain.AppointmentStatus(s)
		status = &st
	}

	appts, err := h.uc.ListAppointments(r.Context(), tenantID, profID, from, to, status)
	if err != nil {
		server.RenderError(w, r, http.StatusInternalServerError, "list_appointments_failed", "Failed to list appointments")
		return
	}

	resp := make([]appointmentResponse, len(appts))
	for i := range appts {
		resp[i] = toApptResponse(&appts[i])
	}
	server.RenderJSON(w, http.StatusOK, resp)
}

// GetAppointment handles GET /appointments/{id}
func (h *Handler) GetAppointment(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := identitydomain.TenantIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusBadRequest, "missing_tenant", "No tenant in context")
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_id", "Invalid appointment ID")
		return
	}

	a, err := h.uc.GetAppointmentByID(r.Context(), tenantID, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			server.RenderError(w, r, http.StatusNotFound, "not_found", "Appointment not found")
			return
		}
		server.RenderError(w, r, http.StatusInternalServerError, "get_appointment_failed", "Failed to get appointment")
		return
	}
	server.RenderJSON(w, http.StatusOK, toApptResponse(a))
}

// UpdateAppointmentStatus handles PATCH /appointments/{id}/status
type updateStatusRequest struct {
	Status             string `json:"status"`
	CancellationReason string `json:"cancellation_reason"`
}

func (h *Handler) UpdateAppointmentStatus(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := identitydomain.TenantIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusBadRequest, "missing_tenant", "No tenant in context")
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_id", "Invalid appointment ID")
		return
	}

	var req updateStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Status == "" {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_body", "status is required")
		return
	}

	actorID, _ := identitydomain.UserIDFromContext(r.Context())
	if err := h.uc.UpdateAppointmentStatus(r.Context(), tenantID, id, domain.AppointmentStatus(req.Status), req.CancellationReason, &actorID); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			server.RenderError(w, r, http.StatusNotFound, "not_found", "Appointment not found")
			return
		}
		if errors.Is(err, domain.ErrInvalidTransition) {
			server.RenderError(w, r, http.StatusUnprocessableEntity, "invalid_transition", "Invalid status transition")
			return
		}
		server.RenderError(w, r, http.StatusInternalServerError, "update_status_failed", "Failed to update status")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// RescheduleAppointment handles PATCH /appointments/{id}/reschedule
type rescheduleRequest struct {
	StartAt string `json:"start_at"`
	EndAt   string `json:"end_at"`
}

func (h *Handler) RescheduleAppointment(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := identitydomain.TenantIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusBadRequest, "missing_tenant", "No tenant in context")
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_id", "Invalid appointment ID")
		return
	}

	var req rescheduleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_body", "Invalid request body")
		return
	}

	newStart, err := time.Parse(time.RFC3339, req.StartAt)
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_start_at", "start_at must be RFC3339")
		return
	}
	newEnd, err := time.Parse(time.RFC3339, req.EndAt)
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_end_at", "end_at must be RFC3339")
		return
	}

	actorID, _ := identitydomain.UserIDFromContext(r.Context())
	if err := h.uc.RescheduleAppointment(r.Context(), tenantID, id, newStart, newEnd, &actorID); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			server.RenderError(w, r, http.StatusNotFound, "not_found", "Appointment not found")
			return
		}
		if errors.Is(err, domain.ErrConflict) {
			server.RenderError(w, r, http.StatusConflict, "time_conflict", "New time conflicts with an existing appointment")
			return
		}
		if errors.Is(err, domain.ErrInvalidTransition) {
			server.RenderError(w, r, http.StatusUnprocessableEntity, "invalid_state", "Appointment cannot be rescheduled in its current state")
			return
		}
		server.RenderError(w, r, http.StatusInternalServerError, "reschedule_failed", "Failed to reschedule")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
