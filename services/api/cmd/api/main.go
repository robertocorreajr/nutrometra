package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"nutrometra/api/internal/billing"
	billingrepo "nutrometra/api/internal/billing/repository"
	billinguc "nutrometra/api/internal/billing/usecase"
	"nutrometra/api/internal/bioimpedance"
	biorepo "nutrometra/api/internal/bioimpedance/repository"
	biouc "nutrometra/api/internal/bioimpedance/usecase"
	"nutrometra/api/internal/catalog"
	catrepo "nutrometra/api/internal/catalog/repository"
	catuc "nutrometra/api/internal/catalog/usecase"
	"nutrometra/api/internal/clinical"
	clinrepo "nutrometra/api/internal/clinical/repository"
	clinuc "nutrometra/api/internal/clinical/usecase"
	"nutrometra/api/internal/identity"
	identityoidc "nutrometra/api/internal/identity/oidc"
	"nutrometra/api/internal/patient"
	patrepo "nutrometra/api/internal/patient/repository"
	patuc "nutrometra/api/internal/patient/usecase"
	"nutrometra/api/internal/platform/audit"
	"nutrometra/api/internal/platform/config"
	"nutrometra/api/internal/platform/db"
	"nutrometra/api/internal/platform/logger"
	"nutrometra/api/internal/platform/observability"
	apiredis "nutrometra/api/internal/platform/redis"
	"nutrometra/api/internal/platform/server"
	"nutrometra/api/internal/professional"
	profrepo "nutrometra/api/internal/professional/repository"
	profuc "nutrometra/api/internal/professional/usecase"
	"nutrometra/api/internal/rbac"
	rbacrepo "nutrometra/api/internal/rbac/repository"
	"nutrometra/api/internal/scheduling"
	schedrepo "nutrometra/api/internal/scheduling/repository"
	scheduc "nutrometra/api/internal/scheduling/usecase"
	"nutrometra/api/internal/tenancy"
	tenancyrepo "nutrometra/api/internal/tenancy/repository"
	tenancyuc "nutrometra/api/internal/tenancy/usecase"

	"github.com/go-chi/chi/v5"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	log := logger.New(cfg.API.Env, cfg.API.LogLevel)
	slog.SetDefault(log)

	ctx := context.Background()

	// --- Infrastructure ---

	pool, err := db.New(ctx, cfg.Postgres)
	if err != nil {
		log.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()
	log.Info("database connected")

	redisClient, err := apiredis.New(ctx, cfg.Redis)
	if err != nil {
		log.Error("failed to connect to redis", "error", err)
		os.Exit(1)
	}
	defer redisClient.Close()
	log.Info("redis connected")

	// --- OIDC Validator ---

	validator := identityoidc.NewValidator(cfg.Zitadel.Issuer, cfg.Zitadel.ClientID)
	if err := validator.Init(ctx); err != nil {
		log.Error("failed to initialize OIDC validator", "error", err)
		os.Exit(1)
	}
	defer validator.Close()
	log.Info("OIDC validator initialized")

	// --- Repositories ---

	tenancyRepo := tenancyrepo.NewPostgresRepository(pool)
	rbacRepo := rbacrepo.New(pool)
	billingRepo := billingrepo.New(pool)
	profRepo := profrepo.New(pool)
	schedRepo := schedrepo.New(pool)
	patRepo := patrepo.New(pool)
	clinRepo := clinrepo.New(pool)
	bioRepo := biorepo.New(pool)
	catRepo := catrepo.New(pool)

	// --- Services ---

	auditSvc := audit.NewService()
	tenancyUC := tenancyuc.NewTenantUsecase(tenancyRepo)
	entitlementSvc := billinguc.NewEntitlementService(billingRepo)
	entitlementAdapter := billinguc.NewEntitlementAdapter(entitlementSvc)

	// Phase 2 usecases
	profUC := profuc.New(profRepo, entitlementAdapter)
	schedUC := scheduc.New(schedRepo)
	patUC := patuc.New(patRepo, entitlementAdapter)
	clinUC := clinuc.New(clinRepo)
	bioUC := biouc.New(bioRepo)
	catUC := catuc.New(catRepo)

	// --- Handlers ---

	tenancyHandler := tenancy.NewHandler(tenancyUC)
	rbacHandler := rbac.NewHandler(rbacRepo, pool, auditSvc)
	billingHandler := billing.NewHandler(billingRepo, entitlementSvc, pool, auditSvc)
	healthChecker := observability.NewHealthChecker(pool, redisClient, cfg.Zitadel.Issuer)

	// Phase 2 handlers
	profHandler := professional.NewHandler(profUC, pool, auditSvc)
	schedHandler := scheduling.NewHandler(schedUC, pool, auditSvc)
	patHandler := patient.NewHandler(patUC, pool, auditSvc)
	clinHandler := clinical.NewHandler(clinUC, pool, auditSvc)
	bioHandler := bioimpedance.NewHandler(bioUC, pool, auditSvc)
	catHandler := catalog.NewHandler(catUC, pool, auditSvc)

	// --- Middlewares ---

	authMW := identity.AuthMiddleware(validator)
	resolverMW := identity.UserResolverMiddleware(pool)
	tenantMW := tenancy.TenantMiddleware(tenancyUC)

	// --- Router ---

	srv := server.New(cfg.API.Port)
	r := srv.Router()

	// Public routes
	r.Get("/health", healthChecker.Health)
	r.Get("/ready", healthChecker.Ready)
	r.Get("/plans", billingHandler.ListPlans)

	// Authenticated routes
	r.Group(func(r chi.Router) {
		r.Use(authMW, resolverMW)

		r.Get("/auth/me", identity.MeHandler)

		// Invite activation (auth only, no tenant required)
		r.Post("/invites/activate", patHandler.ActivatePortalAccess)

		// Tenant-scoped routes (require X-Tenant-ID header + membership)
		r.Group(func(r chi.Router) {
			r.Use(tenantMW)

			r.Get("/tenants/current", tenancyHandler.GetCurrent)

			// Billing
			r.Get("/subscription", billingHandler.GetSubscription)
			r.Post("/subscription/trial", billingHandler.ActivateTrial)
			r.Get("/entitlements", billingHandler.GetEntitlements)

			// RBAC
			r.Get("/roles", rbacHandler.ListRoles)
			r.Post("/members/{member_id}/roles", rbacHandler.AssignRole)
			r.Delete("/members/{member_id}/roles/{role_id}", rbacHandler.RevokeRole)

			// --- Phase 2: Clinical Operations ---

			// Professionals
			r.Route("/professionals", func(r chi.Router) {
				r.With(rbac.RequirePermission("patients:read", rbacRepo)).Get("/", profHandler.List)
				r.With(rbac.RequirePermission("tenant:manage", rbacRepo)).Post("/", profHandler.Create)
				r.Route("/{id}", func(r chi.Router) {
					r.With(rbac.RequirePermission("patients:read", rbacRepo)).Get("/", profHandler.GetByID)
					r.With(rbac.RequirePermission("tenant:manage", rbacRepo)).Put("/", profHandler.Update)

					// Addresses
					r.With(rbac.RequirePermission("tenant:manage", rbacRepo)).Post("/addresses", profHandler.CreateAddress)
					r.With(rbac.RequirePermission("schedule:read", rbacRepo)).Get("/addresses", profHandler.ListAddresses)
					r.With(rbac.RequirePermission("tenant:manage", rbacRepo)).Put("/addresses/{addr_id}", profHandler.UpdateAddress)

					// Service Modes
					r.With(rbac.RequirePermission("tenant:manage", rbacRepo)).Put("/service-modes", profHandler.SetServiceMode)
					r.With(rbac.RequirePermission("schedule:read", rbacRepo)).Get("/service-modes", profHandler.ListServiceModes)
				})
			})

			// Scheduling — availability & blocks under /professionals/{prof_id}
			r.Route("/professionals/{prof_id}/availability", func(r chi.Router) {
				r.With(rbac.RequirePermission("schedule:manage", rbacRepo)).Post("/", schedHandler.CreateAvailabilityRule)
				r.With(rbac.RequirePermission("schedule:read", rbacRepo)).Get("/", schedHandler.ListAvailabilityRules)
			})
			r.With(rbac.RequirePermission("schedule:read", rbacRepo)).Get("/professionals/{prof_id}/slots", schedHandler.GetAvailableSlots)
			r.Route("/professionals/{prof_id}/blocks", func(r chi.Router) {
				r.With(rbac.RequirePermission("schedule:manage", rbacRepo)).Post("/", schedHandler.CreateBlock)
				r.With(rbac.RequirePermission("schedule:read", rbacRepo)).Get("/", schedHandler.ListBlocks)
			})

			// Scheduling — standalone availability/block resources
			r.With(rbac.RequirePermission("schedule:manage", rbacRepo)).Put("/availability/{id}", schedHandler.UpdateAvailabilityRule)
			r.With(rbac.RequirePermission("schedule:manage", rbacRepo)).Delete("/availability/{id}", schedHandler.DeleteAvailabilityRule)
			r.With(rbac.RequirePermission("schedule:manage", rbacRepo)).Delete("/blocks/{id}", schedHandler.DeleteBlock)

			// Appointments
			r.Route("/appointments", func(r chi.Router) {
				r.With(rbac.RequirePermission("schedule:manage", rbacRepo)).Post("/", schedHandler.CreateAppointment)
				r.With(rbac.RequirePermission("appointments:read", rbacRepo)).Get("/", schedHandler.ListAppointments)
				r.Route("/{id}", func(r chi.Router) {
					r.With(rbac.RequirePermission("appointments:read", rbacRepo)).Get("/", schedHandler.GetAppointment)
					r.With(rbac.RequirePermission("schedule:manage", rbacRepo)).Patch("/status", schedHandler.UpdateAppointmentStatus)
					r.With(rbac.RequirePermission("schedule:manage", rbacRepo)).Patch("/reschedule", schedHandler.RescheduleAppointment)
				})
			})

			// Patients
			r.Route("/patients", func(r chi.Router) {
				r.With(rbac.RequirePermission("patients:write", rbacRepo)).Post("/", patHandler.Create)
				r.With(rbac.RequirePermission("patients:read", rbacRepo)).Get("/", patHandler.List)
				r.Route("/{patient_id}", func(r chi.Router) {
					r.With(rbac.RequirePermission("patients:read", rbacRepo)).Get("/", patHandler.GetByID)
					r.With(rbac.RequirePermission("patients:write", rbacRepo)).Put("/", patHandler.Update)
					r.With(rbac.RequirePermission("patients:write", rbacRepo)).Post("/invites", patHandler.GenerateInvite)

					// Profiles
					r.With(rbac.RequirePermission("patients:write", rbacRepo)).Post("/profiles", patHandler.UpsertProfile)
					r.With(rbac.RequirePermission("patients:read", rbacRepo)).Get("/profiles", patHandler.GetProfile)

					// Clinical — anamneses
					r.With(rbac.RequirePermission("clinical:write", rbacRepo)).Post("/anamneses", clinHandler.CreateAnamnesis)
					r.With(rbac.RequirePermission("clinical:read", rbacRepo)).Get("/anamneses", clinHandler.ListAnamneses)

					// Clinical — progress notes
					r.With(rbac.RequirePermission("clinical:write", rbacRepo)).Post("/notes", clinHandler.CreateProgressNote)
					r.With(rbac.RequirePermission("clinical:read", rbacRepo)).Get("/notes", clinHandler.ListProgressNotes)

					// Clinical — attachments
					r.With(rbac.RequirePermission("clinical:write", rbacRepo)).Post("/attachments", clinHandler.CreateAttachment)
					r.With(rbac.RequirePermission("clinical:read", rbacRepo)).Get("/attachments", clinHandler.ListAttachments)

					// Bioimpedance — measurements
					r.With(rbac.RequirePermission("clinical:read", rbacRepo)).Get("/measurements", bioHandler.ListByPatient)
				})
			})

			// Clinical — standalone anamnesis resources
			r.Route("/anamneses/{id}", func(r chi.Router) {
				r.With(rbac.RequirePermission("clinical:read", rbacRepo)).Get("/", clinHandler.GetAnamnesis)
				r.With(rbac.RequirePermission("clinical:write", rbacRepo)).Put("/", clinHandler.UpdateAnamnesis)
				r.With(rbac.RequirePermission("clinical:write", rbacRepo)).Post("/finalize", clinHandler.FinalizeAnamnesis)
			})

			// Bioimpedance — measurements
			r.Route("/measurements", func(r chi.Router) {
				r.With(rbac.RequirePermission("clinical:write", rbacRepo)).Post("/", bioHandler.Create)
				r.Route("/{id}", func(r chi.Router) {
					r.With(rbac.RequirePermission("clinical:read", rbacRepo)).Get("/", bioHandler.GetByID)
					r.With(rbac.RequirePermission("clinical:write", rbacRepo)).Post("/publish", bioHandler.Publish)
				})
			})

			// Food Catalog
			r.Route("/foods", func(r chi.Router) {
				r.With(rbac.RequirePermission("diet:read", rbacRepo)).Get("/", catHandler.List)
				r.With(rbac.RequirePermission("diet:read", rbacRepo)).Get("/groups", catHandler.ListFoodGroups)
				r.With(rbac.RequirePermission("clinical:write", rbacRepo)).Post("/", catHandler.Create)
				r.Route("/{id}", func(r chi.Router) {
					r.With(rbac.RequirePermission("diet:read", rbacRepo)).Get("/", catHandler.GetByID)
					r.With(rbac.RequirePermission("clinical:write", rbacRepo)).Put("/", catHandler.Update)
					r.With(rbac.RequirePermission("clinical:write", rbacRepo)).Delete("/", catHandler.Delete)
				})
			})
		})
	})

	// --- Start server ---

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Info("API starting", "port", cfg.API.Port)
		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			log.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	<-quit
	log.Info("shutting down...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error("shutdown error", "error", err)
	}
	log.Info("server stopped")
}
