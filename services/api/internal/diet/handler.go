package diet

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"nutrometra/api/internal/diet/domain"
	"nutrometra/api/internal/diet/usecase"
	identitydomain "nutrometra/api/internal/identity/domain"
	"nutrometra/api/internal/platform/audit"
	"nutrometra/api/internal/platform/server"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Handler exposes HTTP endpoints for diet operations.
type Handler struct {
	uc       *usecase.Usecase
	pool     *pgxpool.Pool
	auditSvc *audit.Service
}

// NewHandler creates a diet Handler.
func NewHandler(uc *usecase.Usecase, pool *pgxpool.Pool, auditSvc *audit.Service) *Handler {
	return &Handler{uc: uc, pool: pool, auditSvc: auditSvc}
}

// --- Request / Response types ---

type createDietRequest struct {
	ProfessionalID string  `json:"professional_id"`
	Title          string  `json:"title"`
	Objective      string  `json:"objective"`
	ValidFrom      *string `json:"valid_from"`
	ValidUntil     *string `json:"valid_until"`
}

type updateDietRequest struct {
	Title      string  `json:"title"`
	Objective  string  `json:"objective"`
	ValidFrom  *string `json:"valid_from"`
	ValidUntil *string `json:"valid_until"`
}

type mealRequest struct {
	MealName  string `json:"meal_name"`
	MealOrder int    `json:"meal_order"`
	Notes     string `json:"notes"`
}

type mealItemRequest struct {
	FoodItemID         string   `json:"food_item_id"`
	QuantityValue      float64  `json:"quantity_value"`
	QuantityUnit       string   `json:"quantity_unit"`
	HouseholdMeasureID *string  `json:"household_measure_id"`
	AmountDescription  string   `json:"amount_description"`
	PreparationNotes   string   `json:"preparation_notes"`
	SortOrder          int      `json:"sort_order"`
}

type substitutionRequest struct {
	SubstituteFoodItemID string   `json:"substitute_food_item_id"`
	QuantityValue        *float64 `json:"quantity_value"`
	QuantityUnit         *string  `json:"quantity_unit"`
	Notes                string   `json:"notes"`
	SortOrder            int      `json:"sort_order"`
}

type dietResponse struct {
	ID                string         `json:"id"`
	PatientID         string         `json:"patient_id"`
	ProfessionalID    string         `json:"professional_id"`
	Title             string         `json:"title"`
	Objective         string         `json:"objective,omitempty"`
	Status            string         `json:"status"`
	VersionNumber     int            `json:"version_number"`
	PreviousVersionID *string        `json:"previous_version_id,omitempty"`
	PublishedAt       *string        `json:"published_at,omitempty"`
	ValidFrom         *string        `json:"valid_from,omitempty"`
	ValidUntil        *string        `json:"valid_until,omitempty"`
	CreatedAt         string         `json:"created_at"`
	UpdatedAt         string         `json:"updated_at"`
	Meals             []mealResponse `json:"meals,omitempty"`
}

type mealResponse struct {
	ID        string             `json:"id"`
	DietID    string             `json:"diet_id"`
	MealName  string             `json:"meal_name"`
	MealOrder int                `json:"meal_order"`
	Notes     string             `json:"notes,omitempty"`
	Items     []mealItemResponse `json:"items,omitempty"`
}

type mealItemResponse struct {
	ID                 string                 `json:"id"`
	DietMealID         string                 `json:"diet_meal_id"`
	FoodItemID         string                 `json:"food_item_id"`
	QuantityValue      float64                `json:"quantity_value"`
	QuantityUnit       string                 `json:"quantity_unit"`
	HouseholdMeasureID *string                `json:"household_measure_id,omitempty"`
	AmountDescription  string                 `json:"amount_description,omitempty"`
	PreparationNotes   string                 `json:"preparation_notes,omitempty"`
	SortOrder          int                    `json:"sort_order"`
	Substitutions      []substitutionResponse `json:"substitutions,omitempty"`
}

type substitutionResponse struct {
	ID                   string   `json:"id"`
	DietMealItemID       string   `json:"diet_meal_item_id"`
	SubstituteFoodItemID string   `json:"substitute_food_item_id"`
	QuantityValue        *float64 `json:"quantity_value,omitempty"`
	QuantityUnit         *string  `json:"quantity_unit,omitempty"`
	Notes                string   `json:"notes,omitempty"`
	SortOrder            int      `json:"sort_order"`
}

// --- Converters ---

func toDietResponse(d *domain.Diet) dietResponse {
	resp := dietResponse{
		ID:             d.ID.String(),
		PatientID:      d.PatientID.String(),
		ProfessionalID: d.ProfessionalID.String(),
		Title:          d.Title,
		Objective:      d.Objective,
		Status:         string(d.Status),
		VersionNumber:  d.VersionNumber,
		CreatedAt:      d.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      d.UpdatedAt.Format(time.RFC3339),
	}
	if d.PreviousVersionID != nil {
		s := d.PreviousVersionID.String()
		resp.PreviousVersionID = &s
	}
	if d.PublishedAt != nil {
		s := d.PublishedAt.Format(time.RFC3339)
		resp.PublishedAt = &s
	}
	if d.ValidFrom != nil {
		s := d.ValidFrom.Format("2006-01-02")
		resp.ValidFrom = &s
	}
	if d.ValidUntil != nil {
		s := d.ValidUntil.Format("2006-01-02")
		resp.ValidUntil = &s
	}
	if len(d.Meals) > 0 {
		resp.Meals = make([]mealResponse, len(d.Meals))
		for i := range d.Meals {
			resp.Meals[i] = toMealResponse(&d.Meals[i])
		}
	}
	return resp
}

func toDietListResponse(d *domain.Diet) dietResponse {
	resp := dietResponse{
		ID:             d.ID.String(),
		PatientID:      d.PatientID.String(),
		ProfessionalID: d.ProfessionalID.String(),
		Title:          d.Title,
		Objective:      d.Objective,
		Status:         string(d.Status),
		VersionNumber:  d.VersionNumber,
		CreatedAt:      d.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      d.UpdatedAt.Format(time.RFC3339),
	}
	if d.PreviousVersionID != nil {
		s := d.PreviousVersionID.String()
		resp.PreviousVersionID = &s
	}
	if d.PublishedAt != nil {
		s := d.PublishedAt.Format(time.RFC3339)
		resp.PublishedAt = &s
	}
	if d.ValidFrom != nil {
		s := d.ValidFrom.Format("2006-01-02")
		resp.ValidFrom = &s
	}
	if d.ValidUntil != nil {
		s := d.ValidUntil.Format("2006-01-02")
		resp.ValidUntil = &s
	}
	return resp
}

func toMealResponse(m *domain.DietMeal) mealResponse {
	resp := mealResponse{
		ID:        m.ID.String(),
		DietID:    m.DietID.String(),
		MealName:  m.MealName,
		MealOrder: m.MealOrder,
		Notes:     m.Notes,
	}
	if len(m.Items) > 0 {
		resp.Items = make([]mealItemResponse, len(m.Items))
		for i := range m.Items {
			resp.Items[i] = toMealItemResponse(&m.Items[i])
		}
	}
	return resp
}

func toMealItemResponse(item *domain.DietMealItem) mealItemResponse {
	resp := mealItemResponse{
		ID:                item.ID.String(),
		DietMealID:        item.DietMealID.String(),
		FoodItemID:        item.FoodItemID.String(),
		QuantityValue:     item.QuantityValue,
		QuantityUnit:      item.QuantityUnit,
		AmountDescription: item.AmountDescription,
		PreparationNotes:  item.PreparationNotes,
		SortOrder:         item.SortOrder,
	}
	if item.HouseholdMeasureID != nil {
		s := item.HouseholdMeasureID.String()
		resp.HouseholdMeasureID = &s
	}
	if len(item.Substitutions) > 0 {
		resp.Substitutions = make([]substitutionResponse, len(item.Substitutions))
		for i := range item.Substitutions {
			resp.Substitutions[i] = toSubstitutionResponse(&item.Substitutions[i])
		}
	}
	return resp
}

func toSubstitutionResponse(sub *domain.DietSubstitution) substitutionResponse {
	return substitutionResponse{
		ID:                   sub.ID.String(),
		DietMealItemID:       sub.DietMealItemID.String(),
		SubstituteFoodItemID: sub.SubstituteFoodItemID.String(),
		QuantityValue:        sub.QuantityValue,
		QuantityUnit:         sub.QuantityUnit,
		Notes:                sub.Notes,
		SortOrder:            sub.SortOrder,
	}
}

// --- Handlers ---

// CreateDiet handles POST /patients/{patient_id}/diets
func (h *Handler) CreateDiet(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := identitydomain.TenantIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusBadRequest, "missing_tenant", "No tenant in context")
		return
	}

	patientID, err := uuid.Parse(chi.URLParam(r, "patient_id"))
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_patient_id", "Invalid patient ID")
		return
	}

	var req createDietRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_body", "Invalid request body")
		return
	}

	profID, err := uuid.Parse(req.ProfessionalID)
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_professional_id", "Invalid professional_id")
		return
	}

	d := &domain.Diet{
		TenantID:       tenantID,
		PatientID:      patientID,
		ProfessionalID: profID,
		Title:          req.Title,
		Objective:      req.Objective,
	}

	if req.ValidFrom != nil {
		t, err := time.Parse("2006-01-02", *req.ValidFrom)
		if err != nil {
			server.RenderError(w, r, http.StatusBadRequest, "invalid_valid_from", "Invalid valid_from date")
			return
		}
		d.ValidFrom = &t
	}
	if req.ValidUntil != nil {
		t, err := time.Parse("2006-01-02", *req.ValidUntil)
		if err != nil {
			server.RenderError(w, r, http.StatusBadRequest, "invalid_valid_until", "Invalid valid_until date")
			return
		}
		d.ValidUntil = &t
	}

	if err := h.uc.CreateDiet(r.Context(), d); err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "create_diet_failed", err.Error())
		return
	}

	actorID, _ := identitydomain.UserIDFromContext(r.Context())
	entry := audit.NewEntry(
		audit.WithTenantID(tenantID),
		audit.WithActor(actorID, audit.ScopeTenant),
		audit.WithEntity("diet", d.ID),
		audit.WithAction("diet_created"),
		audit.WithIPAddress(r.RemoteAddr),
	)
	_ = h.auditSvc.Write(r.Context(), h.pool, entry)

	server.RenderJSON(w, http.StatusCreated, toDietResponse(d))
}

// ListByPatient handles GET /patients/{patient_id}/diets
func (h *Handler) ListByPatient(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := identitydomain.TenantIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusBadRequest, "missing_tenant", "No tenant in context")
		return
	}

	patientID, err := uuid.Parse(chi.URLParam(r, "patient_id"))
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_patient_id", "Invalid patient ID")
		return
	}

	diets, err := h.uc.ListByPatient(r.Context(), tenantID, patientID)
	if err != nil {
		server.RenderError(w, r, http.StatusInternalServerError, "list_diets_failed", "Failed to list diets")
		return
	}

	resp := make([]dietResponse, len(diets))
	for i := range diets {
		resp[i] = toDietListResponse(&diets[i])
	}
	server.RenderJSON(w, http.StatusOK, resp)
}

// GetByID handles GET /diets/{id}
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := identitydomain.TenantIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusBadRequest, "missing_tenant", "No tenant in context")
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_id", "Invalid diet ID")
		return
	}

	d, err := h.uc.GetByID(r.Context(), tenantID, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			server.RenderError(w, r, http.StatusNotFound, "not_found", "Diet not found")
			return
		}
		server.RenderError(w, r, http.StatusInternalServerError, "get_diet_failed", "Failed to get diet")
		return
	}
	server.RenderJSON(w, http.StatusOK, toDietResponse(d))
}

// Update handles PUT /diets/{id}
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := identitydomain.TenantIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusBadRequest, "missing_tenant", "No tenant in context")
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_id", "Invalid diet ID")
		return
	}

	var req updateDietRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_body", "Invalid request body")
		return
	}

	d := &domain.Diet{
		ID:        id,
		TenantID:  tenantID,
		Title:     req.Title,
		Objective: req.Objective,
	}

	if req.ValidFrom != nil {
		t, err := time.Parse("2006-01-02", *req.ValidFrom)
		if err != nil {
			server.RenderError(w, r, http.StatusBadRequest, "invalid_valid_from", "Invalid valid_from date")
			return
		}
		d.ValidFrom = &t
	}
	if req.ValidUntil != nil {
		t, err := time.Parse("2006-01-02", *req.ValidUntil)
		if err != nil {
			server.RenderError(w, r, http.StatusBadRequest, "invalid_valid_until", "Invalid valid_until date")
			return
		}
		d.ValidUntil = &t
	}

	if err := h.uc.UpdateDiet(r.Context(), d); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			server.RenderError(w, r, http.StatusNotFound, "not_found", "Diet not found")
			return
		}
		if errors.Is(err, domain.ErrDietNotDraft) {
			server.RenderError(w, r, http.StatusConflict, "not_draft", "Only draft diets can be updated")
			return
		}
		server.RenderError(w, r, http.StatusBadRequest, "update_diet_failed", err.Error())
		return
	}

	actorID, _ := identitydomain.UserIDFromContext(r.Context())
	entry := audit.NewEntry(
		audit.WithTenantID(tenantID),
		audit.WithActor(actorID, audit.ScopeTenant),
		audit.WithEntity("diet", id),
		audit.WithAction("diet_updated"),
		audit.WithIPAddress(r.RemoteAddr),
	)
	_ = h.auditSvc.Write(r.Context(), h.pool, entry)

	server.RenderJSON(w, http.StatusOK, toDietResponse(d))
}

// Delete handles DELETE /diets/{id}
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := identitydomain.TenantIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusBadRequest, "missing_tenant", "No tenant in context")
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_id", "Invalid diet ID")
		return
	}

	if err := h.uc.DeleteDiet(r.Context(), tenantID, id); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			server.RenderError(w, r, http.StatusNotFound, "not_found", "Diet not found")
			return
		}
		if errors.Is(err, domain.ErrDietNotDraft) {
			server.RenderError(w, r, http.StatusConflict, "not_draft", "Only draft diets can be deleted")
			return
		}
		server.RenderError(w, r, http.StatusInternalServerError, "delete_diet_failed", "Failed to delete diet")
		return
	}

	actorID, _ := identitydomain.UserIDFromContext(r.Context())
	entry := audit.NewEntry(
		audit.WithTenantID(tenantID),
		audit.WithActor(actorID, audit.ScopeTenant),
		audit.WithEntity("diet", id),
		audit.WithAction("diet_deleted"),
		audit.WithIPAddress(r.RemoteAddr),
	)
	_ = h.auditSvc.Write(r.Context(), h.pool, entry)

	w.WriteHeader(http.StatusNoContent)
}

// Publish handles POST /diets/{id}/publish
func (h *Handler) Publish(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := identitydomain.TenantIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusBadRequest, "missing_tenant", "No tenant in context")
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_id", "Invalid diet ID")
		return
	}

	userID, _ := identitydomain.UserIDFromContext(r.Context())

	if err := h.uc.PublishDiet(r.Context(), tenantID, id, userID); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			server.RenderError(w, r, http.StatusNotFound, "not_found", "Diet not found")
			return
		}
		if errors.Is(err, domain.ErrAlreadyPublished) {
			server.RenderError(w, r, http.StatusConflict, "already_published", "Diet is already published")
			return
		}
		if errors.Is(err, domain.ErrDietEmpty) {
			server.RenderError(w, r, http.StatusUnprocessableEntity, "diet_empty", "Diet must have at least one meal with one item to publish")
			return
		}
		server.RenderError(w, r, http.StatusInternalServerError, "publish_failed", "Failed to publish diet")
		return
	}

	entry := audit.NewEntry(
		audit.WithTenantID(tenantID),
		audit.WithActor(userID, audit.ScopeTenant),
		audit.WithEntity("diet", id),
		audit.WithAction("diet_published"),
		audit.WithIPAddress(r.RemoteAddr),
	)
	_ = h.auditSvc.Write(r.Context(), h.pool, entry)

	w.WriteHeader(http.StatusNoContent)
}

// Archive handles POST /diets/{id}/archive
func (h *Handler) Archive(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := identitydomain.TenantIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusBadRequest, "missing_tenant", "No tenant in context")
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_id", "Invalid diet ID")
		return
	}

	if err := h.uc.ArchiveDiet(r.Context(), tenantID, id); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			server.RenderError(w, r, http.StatusNotFound, "not_found", "Diet not found")
			return
		}
		if errors.Is(err, domain.ErrNotPublished) {
			server.RenderError(w, r, http.StatusConflict, "not_published", "Only published diets can be archived")
			return
		}
		server.RenderError(w, r, http.StatusInternalServerError, "archive_failed", "Failed to archive diet")
		return
	}

	actorID, _ := identitydomain.UserIDFromContext(r.Context())
	entry := audit.NewEntry(
		audit.WithTenantID(tenantID),
		audit.WithActor(actorID, audit.ScopeTenant),
		audit.WithEntity("diet", id),
		audit.WithAction("diet_archived"),
		audit.WithIPAddress(r.RemoteAddr),
	)
	_ = h.auditSvc.Write(r.Context(), h.pool, entry)

	w.WriteHeader(http.StatusNoContent)
}

// NewVersion handles POST /diets/{id}/new-version
func (h *Handler) NewVersion(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := identitydomain.TenantIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusBadRequest, "missing_tenant", "No tenant in context")
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_id", "Invalid diet ID")
		return
	}

	newDiet, err := h.uc.CreateNewVersion(r.Context(), tenantID, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			server.RenderError(w, r, http.StatusNotFound, "not_found", "Diet not found")
			return
		}
		server.RenderError(w, r, http.StatusInternalServerError, "new_version_failed", "Failed to create new version")
		return
	}

	actorID, _ := identitydomain.UserIDFromContext(r.Context())
	entry := audit.NewEntry(
		audit.WithTenantID(tenantID),
		audit.WithActor(actorID, audit.ScopeTenant),
		audit.WithEntity("diet", newDiet.ID),
		audit.WithAction("diet_version_created"),
		audit.WithIPAddress(r.RemoteAddr),
	)
	_ = h.auditSvc.Write(r.Context(), h.pool, entry)

	server.RenderJSON(w, http.StatusCreated, toDietResponse(newDiet))
}

// --- Meal handlers ---

// AddMeal handles POST /diets/{id}/meals
func (h *Handler) AddMeal(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := identitydomain.TenantIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusBadRequest, "missing_tenant", "No tenant in context")
		return
	}

	dietID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_id", "Invalid diet ID")
		return
	}

	var req mealRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_body", "Invalid request body")
		return
	}

	m := &domain.DietMeal{
		TenantID:  tenantID,
		DietID:    dietID,
		MealName:  req.MealName,
		MealOrder: req.MealOrder,
		Notes:     req.Notes,
	}

	if err := h.uc.AddMeal(r.Context(), m); err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "add_meal_failed", err.Error())
		return
	}

	actorID, _ := identitydomain.UserIDFromContext(r.Context())
	entry := audit.NewEntry(
		audit.WithTenantID(tenantID),
		audit.WithActor(actorID, audit.ScopeTenant),
		audit.WithEntity("diet_meal", m.ID),
		audit.WithAction("meal_added"),
		audit.WithIPAddress(r.RemoteAddr),
	)
	_ = h.auditSvc.Write(r.Context(), h.pool, entry)

	server.RenderJSON(w, http.StatusCreated, toMealResponse(m))
}

// UpdateMeal handles PUT /diets/{id}/meals/{meal_id}
func (h *Handler) UpdateMeal(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := identitydomain.TenantIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusBadRequest, "missing_tenant", "No tenant in context")
		return
	}

	dietID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_id", "Invalid diet ID")
		return
	}

	mealID, err := uuid.Parse(chi.URLParam(r, "meal_id"))
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_meal_id", "Invalid meal ID")
		return
	}

	var req mealRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_body", "Invalid request body")
		return
	}

	m := &domain.DietMeal{
		ID:        mealID,
		TenantID:  tenantID,
		DietID:    dietID,
		MealName:  req.MealName,
		MealOrder: req.MealOrder,
		Notes:     req.Notes,
	}

	if err := h.uc.UpdateMeal(r.Context(), m); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			server.RenderError(w, r, http.StatusNotFound, "not_found", "Meal not found")
			return
		}
		server.RenderError(w, r, http.StatusBadRequest, "update_meal_failed", err.Error())
		return
	}

	actorID, _ := identitydomain.UserIDFromContext(r.Context())
	entry := audit.NewEntry(
		audit.WithTenantID(tenantID),
		audit.WithActor(actorID, audit.ScopeTenant),
		audit.WithEntity("diet_meal", mealID),
		audit.WithAction("meal_updated"),
	)
	_ = h.auditSvc.Write(r.Context(), h.pool, entry)

	server.RenderJSON(w, http.StatusOK, toMealResponse(m))
}

// RemoveMeal handles DELETE /diets/{id}/meals/{meal_id}
func (h *Handler) RemoveMeal(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := identitydomain.TenantIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusBadRequest, "missing_tenant", "No tenant in context")
		return
	}

	mealID, err := uuid.Parse(chi.URLParam(r, "meal_id"))
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_meal_id", "Invalid meal ID")
		return
	}

	if err := h.uc.RemoveMeal(r.Context(), tenantID, mealID); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			server.RenderError(w, r, http.StatusNotFound, "not_found", "Meal not found")
			return
		}
		server.RenderError(w, r, http.StatusInternalServerError, "remove_meal_failed", "Failed to remove meal")
		return
	}

	actorID, _ := identitydomain.UserIDFromContext(r.Context())
	entry := audit.NewEntry(
		audit.WithTenantID(tenantID),
		audit.WithActor(actorID, audit.ScopeTenant),
		audit.WithEntity("diet_meal", mealID),
		audit.WithAction("meal_removed"),
		audit.WithIPAddress(r.RemoteAddr),
	)
	_ = h.auditSvc.Write(r.Context(), h.pool, entry)

	w.WriteHeader(http.StatusNoContent)
}

// --- Meal Item handlers ---

// AddMealItem handles POST /diets/{id}/meals/{meal_id}/items
func (h *Handler) AddMealItem(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := identitydomain.TenantIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusBadRequest, "missing_tenant", "No tenant in context")
		return
	}

	mealID, err := uuid.Parse(chi.URLParam(r, "meal_id"))
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_meal_id", "Invalid meal ID")
		return
	}

	var req mealItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_body", "Invalid request body")
		return
	}

	foodID, err := uuid.Parse(req.FoodItemID)
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_food_item_id", "Invalid food_item_id")
		return
	}

	item := &domain.DietMealItem{
		TenantID:          tenantID,
		DietMealID:        mealID,
		FoodItemID:        foodID,
		QuantityValue:     req.QuantityValue,
		QuantityUnit:      req.QuantityUnit,
		AmountDescription: req.AmountDescription,
		PreparationNotes:  req.PreparationNotes,
		SortOrder:         req.SortOrder,
	}

	if req.HouseholdMeasureID != nil {
		hmID, err := uuid.Parse(*req.HouseholdMeasureID)
		if err != nil {
			server.RenderError(w, r, http.StatusBadRequest, "invalid_household_measure_id", "Invalid household_measure_id")
			return
		}
		item.HouseholdMeasureID = &hmID
	}

	if err := h.uc.AddMealItem(r.Context(), item); err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "add_item_failed", err.Error())
		return
	}

	actorID, _ := identitydomain.UserIDFromContext(r.Context())
	entry := audit.NewEntry(
		audit.WithTenantID(tenantID),
		audit.WithActor(actorID, audit.ScopeTenant),
		audit.WithEntity("diet_meal_item", item.ID),
		audit.WithAction("meal_item_added"),
		audit.WithIPAddress(r.RemoteAddr),
	)
	_ = h.auditSvc.Write(r.Context(), h.pool, entry)

	server.RenderJSON(w, http.StatusCreated, toMealItemResponse(item))
}

// UpdateMealItem handles PUT /diet-items/{item_id}
func (h *Handler) UpdateMealItem(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := identitydomain.TenantIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusBadRequest, "missing_tenant", "No tenant in context")
		return
	}

	itemID, err := uuid.Parse(chi.URLParam(r, "item_id"))
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_item_id", "Invalid item ID")
		return
	}

	var req mealItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_body", "Invalid request body")
		return
	}

	foodID, err := uuid.Parse(req.FoodItemID)
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_food_item_id", "Invalid food_item_id")
		return
	}

	item := &domain.DietMealItem{
		ID:                itemID,
		TenantID:          tenantID,
		FoodItemID:        foodID,
		QuantityValue:     req.QuantityValue,
		QuantityUnit:      req.QuantityUnit,
		AmountDescription: req.AmountDescription,
		PreparationNotes:  req.PreparationNotes,
		SortOrder:         req.SortOrder,
	}

	if req.HouseholdMeasureID != nil {
		hmID, err := uuid.Parse(*req.HouseholdMeasureID)
		if err != nil {
			server.RenderError(w, r, http.StatusBadRequest, "invalid_household_measure_id", "Invalid household_measure_id")
			return
		}
		item.HouseholdMeasureID = &hmID
	}

	if err := h.uc.UpdateMealItem(r.Context(), item); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			server.RenderError(w, r, http.StatusNotFound, "not_found", "Meal item not found")
			return
		}
		server.RenderError(w, r, http.StatusBadRequest, "update_item_failed", err.Error())
		return
	}

	actorID, _ := identitydomain.UserIDFromContext(r.Context())
	entry := audit.NewEntry(
		audit.WithTenantID(tenantID),
		audit.WithActor(actorID, audit.ScopeTenant),
		audit.WithEntity("diet_meal_item", itemID),
		audit.WithAction("meal_item_updated"),
	)
	_ = h.auditSvc.Write(r.Context(), h.pool, entry)

	server.RenderJSON(w, http.StatusOK, toMealItemResponse(item))
}

// RemoveMealItem handles DELETE /diet-items/{item_id}
func (h *Handler) RemoveMealItem(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := identitydomain.TenantIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusBadRequest, "missing_tenant", "No tenant in context")
		return
	}

	itemID, err := uuid.Parse(chi.URLParam(r, "item_id"))
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_item_id", "Invalid item ID")
		return
	}

	if err := h.uc.RemoveMealItem(r.Context(), tenantID, itemID); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			server.RenderError(w, r, http.StatusNotFound, "not_found", "Meal item not found")
			return
		}
		server.RenderError(w, r, http.StatusInternalServerError, "remove_item_failed", "Failed to remove meal item")
		return
	}

	actorID, _ := identitydomain.UserIDFromContext(r.Context())
	entry := audit.NewEntry(
		audit.WithTenantID(tenantID),
		audit.WithActor(actorID, audit.ScopeTenant),
		audit.WithEntity("diet_meal_item", itemID),
		audit.WithAction("meal_item_removed"),
		audit.WithIPAddress(r.RemoteAddr),
	)
	_ = h.auditSvc.Write(r.Context(), h.pool, entry)

	w.WriteHeader(http.StatusNoContent)
}

// --- Substitution handlers ---

// AddSubstitution handles POST /diet-items/{item_id}/substitutions
func (h *Handler) AddSubstitution(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := identitydomain.TenantIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusBadRequest, "missing_tenant", "No tenant in context")
		return
	}

	itemID, err := uuid.Parse(chi.URLParam(r, "item_id"))
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_item_id", "Invalid item ID")
		return
	}

	var req substitutionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_body", "Invalid request body")
		return
	}

	subFoodID, err := uuid.Parse(req.SubstituteFoodItemID)
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_substitute_food_item_id", "Invalid substitute_food_item_id")
		return
	}

	sub := &domain.DietSubstitution{
		TenantID:             tenantID,
		DietMealItemID:       itemID,
		SubstituteFoodItemID: subFoodID,
		QuantityValue:        req.QuantityValue,
		QuantityUnit:         req.QuantityUnit,
		Notes:                req.Notes,
		SortOrder:            req.SortOrder,
	}

	if err := h.uc.AddSubstitution(r.Context(), sub); err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "add_substitution_failed", err.Error())
		return
	}

	actorID, _ := identitydomain.UserIDFromContext(r.Context())
	entry := audit.NewEntry(
		audit.WithTenantID(tenantID),
		audit.WithActor(actorID, audit.ScopeTenant),
		audit.WithEntity("diet_substitution", sub.ID),
		audit.WithAction("substitution_added"),
		audit.WithIPAddress(r.RemoteAddr),
	)
	_ = h.auditSvc.Write(r.Context(), h.pool, entry)

	server.RenderJSON(w, http.StatusCreated, toSubstitutionResponse(sub))
}

// RemoveSubstitution handles DELETE /diet-substitutions/{sub_id}
func (h *Handler) RemoveSubstitution(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := identitydomain.TenantIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusBadRequest, "missing_tenant", "No tenant in context")
		return
	}

	subID, err := uuid.Parse(chi.URLParam(r, "sub_id"))
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_sub_id", "Invalid substitution ID")
		return
	}

	if err := h.uc.RemoveSubstitution(r.Context(), tenantID, subID); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			server.RenderError(w, r, http.StatusNotFound, "not_found", "Substitution not found")
			return
		}
		server.RenderError(w, r, http.StatusInternalServerError, "remove_substitution_failed", "Failed to remove substitution")
		return
	}

	actorID, _ := identitydomain.UserIDFromContext(r.Context())
	entry := audit.NewEntry(
		audit.WithTenantID(tenantID),
		audit.WithActor(actorID, audit.ScopeTenant),
		audit.WithEntity("diet_substitution", subID),
		audit.WithAction("substitution_removed"),
		audit.WithIPAddress(r.RemoteAddr),
	)
	_ = h.auditSvc.Write(r.Context(), h.pool, entry)

	w.WriteHeader(http.StatusNoContent)
}
