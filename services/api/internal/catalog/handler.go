package catalog

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"nutrometra/api/internal/catalog/domain"
	"nutrometra/api/internal/catalog/usecase"

	identitydomain "nutrometra/api/internal/identity/domain"
	"nutrometra/api/internal/platform/audit"
	"nutrometra/api/internal/platform/server"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Handler exposes HTTP endpoints for the food catalog.
type Handler struct {
	uc       *usecase.Usecase
	pool     *pgxpool.Pool
	auditSvc *audit.Service
}

// NewHandler creates a catalog Handler.
func NewHandler(uc *usecase.Usecase, pool *pgxpool.Pool, auditSvc *audit.Service) *Handler {
	return &Handler{uc: uc, pool: pool, auditSvc: auditSvc}
}

type foodItemResponse struct {
	ID            string              `json:"id"`
	TenantID      *string             `json:"tenant_id,omitempty"`
	Name          string              `json:"name"`
	FoodGroup     string              `json:"food_group"`
	Brand         string              `json:"brand,omitempty"`
	Barcode       string              `json:"barcode,omitempty"`
	ServingSizeG  float64             `json:"serving_size_g"`
	ServingLabel  string              `json:"serving_label"`
	Source        string              `json:"source"`
	Active        bool                `json:"active"`
	NutritionFacts *nutritionResponse `json:"nutrition_facts,omitempty"`
	HouseholdMeasures []measureResponse `json:"household_measures,omitempty"`
	CreatedAt     string              `json:"created_at"`
	UpdatedAt     string              `json:"updated_at"`
}

type nutritionResponse struct {
	CaloriesKcal  *float64 `json:"calories_kcal,omitempty"`
	ProteinG      *float64 `json:"protein_g,omitempty"`
	CarbsG        *float64 `json:"carbs_g,omitempty"`
	FiberG        *float64 `json:"fiber_g,omitempty"`
	SugarG        *float64 `json:"sugar_g,omitempty"`
	TotalFatG     *float64 `json:"total_fat_g,omitempty"`
	SaturatedFatG *float64 `json:"saturated_fat_g,omitempty"`
	TransFatG     *float64 `json:"trans_fat_g,omitempty"`
	CholesterolMg *float64 `json:"cholesterol_mg,omitempty"`
	SodiumMg      *float64 `json:"sodium_mg,omitempty"`
	PotassiumMg   *float64 `json:"potassium_mg,omitempty"`
	CalciumMg     *float64 `json:"calcium_mg,omitempty"`
	IronMg        *float64 `json:"iron_mg,omitempty"`
}

type measureResponse struct {
	ID    string  `json:"id"`
	Label string  `json:"label"`
	Grams float64 `json:"grams"`
}

func toFoodResponse(f *domain.FoodItem) foodItemResponse {
	resp := foodItemResponse{
		ID:           f.ID.String(),
		Name:         f.Name,
		FoodGroup:    f.FoodGroup,
		Brand:        f.Brand,
		Barcode:      f.Barcode,
		ServingSizeG: f.ServingSizeG,
		ServingLabel: f.ServingLabel,
		Source:       string(f.Source),
		Active:       f.Active,
		CreatedAt:    f.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    f.UpdatedAt.Format(time.RFC3339),
	}
	if f.TenantID != nil {
		s := f.TenantID.String()
		resp.TenantID = &s
	}
	if f.NutritionFacts != nil {
		resp.NutritionFacts = &nutritionResponse{
			CaloriesKcal:  f.NutritionFacts.CaloriesKcal,
			ProteinG:      f.NutritionFacts.ProteinG,
			CarbsG:        f.NutritionFacts.CarbsG,
			FiberG:        f.NutritionFacts.FiberG,
			SugarG:        f.NutritionFacts.SugarG,
			TotalFatG:     f.NutritionFacts.TotalFatG,
			SaturatedFatG: f.NutritionFacts.SaturatedFatG,
			TransFatG:     f.NutritionFacts.TransFatG,
			CholesterolMg: f.NutritionFacts.CholesterolMg,
			SodiumMg:      f.NutritionFacts.SodiumMg,
			PotassiumMg:   f.NutritionFacts.PotassiumMg,
			CalciumMg:     f.NutritionFacts.CalciumMg,
			IronMg:        f.NutritionFacts.IronMg,
		}
	}
	if len(f.HouseholdMeasures) > 0 {
		resp.HouseholdMeasures = make([]measureResponse, len(f.HouseholdMeasures))
		for i, m := range f.HouseholdMeasures {
			resp.HouseholdMeasures[i] = measureResponse{
				ID:    m.ID.String(),
				Label: m.Label,
				Grams: m.Grams,
			}
		}
	}
	return resp
}

// List handles GET /foods?q=...&group=...&limit=...
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := identitydomain.TenantIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusBadRequest, "missing_tenant", "No tenant in context")
		return
	}

	query := r.URL.Query().Get("q")
	group := r.URL.Query().Get("group")
	limit := 50
	if l := r.URL.Query().Get("limit"); l != "" {
		if n, err := strconv.Atoi(l); err == nil && n > 0 {
			limit = n
		}
	}

	var items []domain.FoodItem
	var err error
	if query != "" {
		items, err = h.uc.Search(r.Context(), &tenantID, query, limit)
	} else {
		items, err = h.uc.List(r.Context(), &tenantID, group, limit)
	}

	if err != nil {
		server.RenderError(w, r, http.StatusInternalServerError, "list_failed", "Failed to list food items")
		return
	}

	resp := make([]foodItemResponse, len(items))
	for i := range items {
		resp[i] = toFoodResponse(&items[i])
	}
	server.RenderJSON(w, http.StatusOK, resp)
}

// GetByID handles GET /foods/{id}
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_id", "Invalid food item ID")
		return
	}

	f, err := h.uc.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			server.RenderError(w, r, http.StatusNotFound, "not_found", "Food item not found")
			return
		}
		server.RenderError(w, r, http.StatusInternalServerError, "get_failed", "Failed to get food item")
		return
	}
	server.RenderJSON(w, http.StatusOK, toFoodResponse(f))
}

type createFoodRequest struct {
	Name           string           `json:"name"`
	FoodGroup      string           `json:"food_group"`
	Brand          string           `json:"brand"`
	Barcode        string           `json:"barcode"`
	ServingSizeG   float64          `json:"serving_size_g"`
	ServingLabel   string           `json:"serving_label"`
	NutritionFacts *nutritionRequest `json:"nutrition_facts"`
}

type nutritionRequest struct {
	CaloriesKcal  *float64 `json:"calories_kcal"`
	ProteinG      *float64 `json:"protein_g"`
	CarbsG        *float64 `json:"carbs_g"`
	FiberG        *float64 `json:"fiber_g"`
	SugarG        *float64 `json:"sugar_g"`
	TotalFatG     *float64 `json:"total_fat_g"`
	SaturatedFatG *float64 `json:"saturated_fat_g"`
	TransFatG     *float64 `json:"trans_fat_g"`
	CholesterolMg *float64 `json:"cholesterol_mg"`
	SodiumMg      *float64 `json:"sodium_mg"`
	PotassiumMg   *float64 `json:"potassium_mg"`
	CalciumMg     *float64 `json:"calcium_mg"`
	IronMg        *float64 `json:"iron_mg"`
}

// Create handles POST /foods
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := identitydomain.TenantIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusBadRequest, "missing_tenant", "No tenant in context")
		return
	}

	var req createFoodRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_body", "Invalid request body")
		return
	}

	f := &domain.FoodItem{
		TenantID:     &tenantID,
		Name:         req.Name,
		FoodGroup:    req.FoodGroup,
		Brand:        req.Brand,
		Barcode:      req.Barcode,
		ServingSizeG: req.ServingSizeG,
		ServingLabel: req.ServingLabel,
	}
	if f.ServingLabel == "" {
		f.ServingLabel = "100g"
	}

	var nf *domain.NutritionFacts
	if req.NutritionFacts != nil {
		nf = &domain.NutritionFacts{
			CaloriesKcal:  req.NutritionFacts.CaloriesKcal,
			ProteinG:      req.NutritionFacts.ProteinG,
			CarbsG:        req.NutritionFacts.CarbsG,
			FiberG:        req.NutritionFacts.FiberG,
			SugarG:        req.NutritionFacts.SugarG,
			TotalFatG:     req.NutritionFacts.TotalFatG,
			SaturatedFatG: req.NutritionFacts.SaturatedFatG,
			TransFatG:     req.NutritionFacts.TransFatG,
			CholesterolMg: req.NutritionFacts.CholesterolMg,
			SodiumMg:      req.NutritionFacts.SodiumMg,
			PotassiumMg:   req.NutritionFacts.PotassiumMg,
			CalciumMg:     req.NutritionFacts.CalciumMg,
			IronMg:        req.NutritionFacts.IronMg,
		}
	}

	if err := h.uc.Create(r.Context(), f, nf); err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "create_failed", err.Error())
		return
	}

	actorID, _ := identitydomain.UserIDFromContext(r.Context())
	entry := audit.NewEntry(
		audit.WithTenantID(tenantID),
		audit.WithActor(actorID, audit.ScopeTenant),
		audit.WithEntity("food_item", f.ID),
		audit.WithAction("food_item_created"),
		audit.WithIPAddress(r.RemoteAddr),
	)
	_ = h.auditSvc.Write(r.Context(), h.pool, entry)

	server.RenderJSON(w, http.StatusCreated, toFoodResponse(f))
}

// Update handles PUT /foods/{id}
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := identitydomain.TenantIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusBadRequest, "missing_tenant", "No tenant in context")
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_id", "Invalid food item ID")
		return
	}

	var req createFoodRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_body", "Invalid request body")
		return
	}

	f := &domain.FoodItem{
		ID:           id,
		TenantID:     &tenantID,
		Name:         req.Name,
		FoodGroup:    req.FoodGroup,
		Brand:        req.Brand,
		Barcode:      req.Barcode,
		ServingSizeG: req.ServingSizeG,
		ServingLabel: req.ServingLabel,
	}

	var nf *domain.NutritionFacts
	if req.NutritionFacts != nil {
		nf = &domain.NutritionFacts{
			CaloriesKcal:  req.NutritionFacts.CaloriesKcal,
			ProteinG:      req.NutritionFacts.ProteinG,
			CarbsG:        req.NutritionFacts.CarbsG,
			FiberG:        req.NutritionFacts.FiberG,
			SugarG:        req.NutritionFacts.SugarG,
			TotalFatG:     req.NutritionFacts.TotalFatG,
			SaturatedFatG: req.NutritionFacts.SaturatedFatG,
			TransFatG:     req.NutritionFacts.TransFatG,
			CholesterolMg: req.NutritionFacts.CholesterolMg,
			SodiumMg:      req.NutritionFacts.SodiumMg,
			PotassiumMg:   req.NutritionFacts.PotassiumMg,
			CalciumMg:     req.NutritionFacts.CalciumMg,
			IronMg:        req.NutritionFacts.IronMg,
		}
	}

	if err := h.uc.Update(r.Context(), f, nf); err != nil {
		if errors.Is(err, domain.ErrNotEditable) {
			server.RenderError(w, r, http.StatusForbidden, "not_editable", "System food items cannot be edited")
			return
		}
		server.RenderError(w, r, http.StatusBadRequest, "update_failed", err.Error())
		return
	}

	server.RenderJSON(w, http.StatusOK, toFoodResponse(f))
}

// Delete handles DELETE /foods/{id}
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := identitydomain.TenantIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusBadRequest, "missing_tenant", "No tenant in context")
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "invalid_id", "Invalid food item ID")
		return
	}

	if err := h.uc.Delete(r.Context(), id); err != nil {
		if errors.Is(err, domain.ErrNotEditable) {
			server.RenderError(w, r, http.StatusForbidden, "not_editable", "System food items cannot be deleted")
			return
		}
		server.RenderError(w, r, http.StatusInternalServerError, "delete_failed", "Failed to delete food item")
		return
	}

	actorID, _ := identitydomain.UserIDFromContext(r.Context())
	entry := audit.NewEntry(
		audit.WithTenantID(tenantID),
		audit.WithActor(actorID, audit.ScopeTenant),
		audit.WithEntity("food_item", id),
		audit.WithAction("food_item_deleted"),
		audit.WithIPAddress(r.RemoteAddr),
	)
	_ = h.auditSvc.Write(r.Context(), h.pool, entry)

	w.WriteHeader(http.StatusNoContent)
}

// ListFoodGroups handles GET /foods/groups
func (h *Handler) ListFoodGroups(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := identitydomain.TenantIDFromContext(r.Context())
	if !ok {
		server.RenderError(w, r, http.StatusBadRequest, "missing_tenant", "No tenant in context")
		return
	}

	groups, err := h.uc.ListFoodGroups(r.Context(), &tenantID)
	if err != nil {
		server.RenderError(w, r, http.StatusInternalServerError, "list_groups_failed", "Failed to list food groups")
		return
	}
	server.RenderJSON(w, http.StatusOK, groups)
}
