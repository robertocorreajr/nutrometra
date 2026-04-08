package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	aihandler "nutrometra/api/internal/ai"
	aiprovider "nutrometra/api/internal/ai/provider"
	airepo "nutrometra/api/internal/ai/repository"
	aiuc "nutrometra/api/internal/ai/usecase"
	"nutrometra/api/internal/backoffice"
	"nutrometra/api/internal/billing"
	"nutrometra/api/internal/billing/provider"
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
	"nutrometra/api/internal/diet"
	dietrepo "nutrometra/api/internal/diet/repository"
	dietuc "nutrometra/api/internal/diet/usecase"
	"nutrometra/api/internal/document"
	docrepo "nutrometra/api/internal/document/repository"
	docuc "nutrometra/api/internal/document/usecase"
	"nutrometra/api/internal/export"
	exportrepo "nutrometra/api/internal/export/repository"
	exportuc "nutrometra/api/internal/export/usecase"
	"nutrometra/api/internal/identity"
	identityoidc "nutrometra/api/internal/identity/oidc"
	googleint "nutrometra/api/internal/integrations/google"
	"nutrometra/api/internal/patient"
	patrepo "nutrometra/api/internal/patient/repository"
	patuc "nutrometra/api/internal/patient/usecase"
	"nutrometra/api/internal/platform/audit"
	"nutrometra/api/internal/platform/cache"
	"nutrometra/api/internal/platform/config"
	"nutrometra/api/internal/platform/db"
	"nutrometra/api/internal/platform/logger"
	"nutrometra/api/internal/platform/observability"
	"nutrometra/api/internal/platform/pdfgen"
	"nutrometra/api/internal/platform/queue"
	apiredis "nutrometra/api/internal/platform/redis"
	"nutrometra/api/internal/platform/server"
	"nutrometra/api/internal/platform/worker"
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

	// --- Cache ---
	redisCache := cache.NewRedisCache(redisClient)

	// --- Persistent Job Queue ---
	jobQueue := queue.NewPostgresQueue(pool)

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
	dietRepo := dietrepo.New(pool)
	docRepo := docrepo.New(pool)
	exportRepo := exportrepo.New(pool)

	// Phase 4 repositories
	webhookRepo := billingrepo.NewWebhookRepository(pool)
	subRepo := billingrepo.NewSubscriptionRepository(pool)
	invoiceRepo := billingrepo.NewInvoiceRepository(pool)
	paymentRepo := billingrepo.NewPaymentRepository(pool)
	overrideRepo := billingrepo.NewOverrideRepository(pool)
	boRepo := backoffice.NewRepository(pool)

	// --- Services ---

	auditSvc := audit.NewService()
	tenancyUC := tenancyuc.NewTenantUsecase(tenancyRepo)
	entitlementSvc := billinguc.NewEntitlementService(billingRepo)
	cachedEntitlementSvc := billinguc.NewCachedEntitlementService(entitlementSvc, redisCache)
	entitlementAdapter := billinguc.NewEntitlementAdapter(entitlementSvc)

	// Phase 2 usecases
	profUC := profuc.New(profRepo, entitlementAdapter)
	schedUC := scheduc.New(schedRepo)
	patUC := patuc.New(patRepo, entitlementAdapter)
	clinUC := clinuc.New(clinRepo)
	bioUC := biouc.New(bioRepo)
	catUC := catuc.New(catRepo)

	// Phase 3 usecases
	pdfGen := pdfgen.New()
	bgWorker := worker.New(100)
	dietUC := dietuc.New(dietRepo, pool)
	docUC := docuc.New(docRepo, pool)
	exportUC := exportuc.New(exportRepo, pool, entitlementAdapter, pdfGen, bgWorker)

	// Phase 4: Billing provider + usecases
	var billingProvider provider.BillingProvider
	if cfg.Stripe.SecretKey != "" {
		billingProvider = provider.NewStripeProvider(cfg.Stripe.SecretKey, cfg.Stripe.WebhookSecret)
		log.Info("Stripe billing provider initialized")
	} else {
		log.Warn("STRIPE_SECRET_KEY not set — billing provider disabled")
	}

	webhookProcessor := billinguc.NewWebhookProcessor(pool, webhookRepo, subRepo, invoiceRepo, paymentRepo, auditSvc)
	subUC := billinguc.NewSubscriptionUsecase(pool, billingProvider, billingRepo, subRepo, auditSvc)
	overrideUC := billinguc.NewOverrideUsecase(overrideRepo, pool, auditSvc)

	// --- Handlers ---

	tenancyHandler := tenancy.NewHandler(tenancyUC)
	rbacHandler := rbac.NewHandler(rbacRepo, pool, auditSvc)
	billingHandler := billing.NewHandler(billingRepo, entitlementSvc, cachedEntitlementSvc, pool, auditSvc)
	healthChecker := observability.NewHealthChecker(pool, redisClient, cfg.Zitadel.Issuer)

	// Phase 2 handlers
	profHandler := professional.NewHandler(profUC, pool, auditSvc)
	schedHandler := scheduling.NewHandler(schedUC, pool, auditSvc)
	patHandler := patient.NewHandler(patUC, pool, auditSvc)
	clinHandler := clinical.NewHandler(clinUC, pool, auditSvc)
	bioHandler := bioimpedance.NewHandler(bioUC, pool, auditSvc)
	catHandler := catalog.NewHandler(catUC, pool, auditSvc)

	// Phase 3 handlers
	dietHandler := diet.NewHandler(dietUC, pool, auditSvc)
	docHandler := document.NewHandler(docUC, pool, auditSvc)
	exportHandler := export.NewHandler(exportUC, pool, auditSvc)

	// Phase 4 handlers
	var webhookHandler *billing.WebhookHandler
	if billingProvider != nil {
		webhookHandler = billing.NewWebhookHandler(billingProvider, webhookProcessor)
	}
	subHandler := billing.NewSubscriptionHandler(subUC, invoiceRepo, paymentRepo)
	boHandler := backoffice.NewHandler(boRepo, subUC, overrideUC, invoiceRepo, paymentRepo)

	// Phase 5: Google Calendar integration
	var googleHandler *googleint.Handler
	if cfg.Google.ClientID != "" && cfg.Google.EncryptionKey != "" {
		googleCalProvider := googleint.NewGoogleCalendarProvider()
		googleRepo := googleint.NewRepository(pool, cfg.Google.EncryptionKey)
		googleHandler = googleint.NewHandler(pool, googleRepo, googleCalProvider, auditSvc, cfg.Google, jobQueue)
		jobQueue.RegisterHandler("calendar_sync", googleint.NewSyncHandler(googleRepo, googleCalProvider, cfg.Google.EncryptionKey))
		log.Info("Google Calendar integration initialized")
	} else {
		log.Warn("GOOGLE_CLIENT_ID or GOOGLE_ENCRYPTION_KEY not set — Google Calendar disabled")
	}

	// Phase 5: AI Assistive
	var aiHandler *aihandler.Handler
	if cfg.AI.AnthropicAPIKey != "" {
		claudeProvider := aiprovider.NewClaudeProvider(cfg.AI.AnthropicAPIKey, cfg.AI.Model, cfg.AI.MaxTokens)
		aiRepo := airepo.New(pool)
		stubCtxProvider := &aiuc.StubPatientContextProvider{}
		aiSuggestionSvc := aiuc.NewSuggestionService(aiRepo, claudeProvider, stubCtxProvider, jobQueue)
		aiHandler = aihandler.NewHandler(aiSuggestionSvc, pool, auditSvc)
		generateHandler := aiuc.NewGenerateHandler(aiRepo, claudeProvider)
		jobQueue.RegisterHandler("ai_generate", generateHandler.Handle)
		log.Info("AI assistive module initialized")
	} else {
		log.Warn("ANTHROPIC_API_KEY not set — AI module disabled")
	}

	// --- Middlewares ---

	authMW := identity.AuthMiddleware(validator)
	resolverMW := identity.UserResolverMiddleware(pool)
	tenantMW := tenancy.TenantMiddleware(tenancyUC)
	backofficeMW := backoffice.BackofficeMiddleware(boRepo)

	// --- Router ---

	srv := server.New(cfg.API.Port)
	r := srv.Router()

	// Public routes
	r.Get("/health", healthChecker.Health)
	r.Get("/ready", healthChecker.Ready)
	r.Get("/plans", billingHandler.ListPlans)

	// Stripe webhook — public, validated via HMAC signature
	if webhookHandler != nil {
		r.Post("/webhooks/stripe", webhookHandler.HandleStripeWebhook)
	}

	// Authenticated routes
	r.Group(func(r chi.Router) {
		r.Use(authMW, resolverMW)

		r.Get("/auth/me", identity.MeHandler)
		r.Get("/auth/me/tenants", identity.MeTenantsHandler(tenancyRepo))

		// Patient portal (auth only, no tenant required)
		r.Post("/invites/activate", patHandler.ActivatePortalAccess)
		r.Get("/patients/me", patHandler.GetMe)

		// Tenant-scoped routes (require X-Tenant-ID header + membership)
		r.Group(func(r chi.Router) {
			r.Use(tenantMW)

			r.Get("/tenants/current", tenancyHandler.GetCurrent)

			// Billing
			r.Get("/subscription", billingHandler.GetSubscription)
			r.Post("/subscription/trial", billingHandler.ActivateTrial)
			r.Get("/entitlements", billingHandler.GetEntitlements)

			// Phase 4 — Subscription self-service
			r.With(rbac.RequirePermission("subscription:manage", rbacRepo)).Post("/subscription/checkout", subHandler.Checkout)
			r.With(rbac.RequirePermission("subscription:manage", rbacRepo)).Patch("/subscription/plan", subHandler.ChangePlan)
			r.With(rbac.RequirePermission("subscription:manage", rbacRepo)).Post("/subscription/cancel", subHandler.CancelSubscription)
			r.Get("/invoices", subHandler.ListInvoices)
			r.Get("/invoices/{id}", subHandler.GetInvoice)
			r.Get("/payments", subHandler.ListPayments)

			// RBAC
			r.Get("/roles", rbacHandler.ListRoles)
			r.Post("/members/{member_id}/roles", rbacHandler.AssignRole)
			r.Delete("/members/{member_id}/roles/{role_id}", rbacHandler.RevokeRole)

			// --- Phase 2: Clinical Operations ---

			// Professionals
			r.Route("/professionals", func(r chi.Router) {
				r.With(rbac.RequirePermission("patients:read", rbacRepo)).Get("/", profHandler.List)
				r.With(rbac.RequirePermission("tenant:manage", rbacRepo)).Post("/", profHandler.Create)
				r.With(rbac.RequirePermission("patients:read", rbacRepo)).Get("/me", profHandler.GetMe)
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

					// Phase 3 — Diets
					r.Route("/diets", func(r chi.Router) {
						r.With(rbac.RequirePermission("diet:write", rbacRepo)).Post("/", dietHandler.CreateDiet)
						r.With(rbac.RequirePermission("diet:read", rbacRepo)).Get("/", dietHandler.ListByPatient)
					})

					// Phase 3 — Documents
					r.Route("/documents", func(r chi.Router) {
						r.With(rbac.RequirePermission("document:write", rbacRepo)).Post("/", docHandler.Create)
						r.With(rbac.RequirePermission("clinical:read", rbacRepo)).Get("/", docHandler.ListByPatient)
					})
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

			// --- Phase 3: Prescription & Publishing ---

			// Diets — standalone resources
			r.Route("/diets/{id}", func(r chi.Router) {
				r.With(rbac.RequirePermission("diet:read", rbacRepo)).Get("/", dietHandler.GetByID)
				r.With(rbac.RequirePermission("diet:write", rbacRepo)).Put("/", dietHandler.Update)
				r.With(rbac.RequirePermission("diet:write", rbacRepo)).Delete("/", dietHandler.Delete)
				r.With(rbac.RequirePermission("diet:publish", rbacRepo)).Post("/publish", dietHandler.Publish)
				r.With(rbac.RequirePermission("diet:write", rbacRepo)).Post("/archive", dietHandler.Archive)
				r.With(rbac.RequirePermission("diet:write", rbacRepo)).Post("/new-version", dietHandler.NewVersion)

				// Meals
				r.Route("/meals", func(r chi.Router) {
					r.With(rbac.RequirePermission("diet:write", rbacRepo)).Post("/", dietHandler.AddMeal)
					r.With(rbac.RequirePermission("diet:write", rbacRepo)).Put("/{meal_id}", dietHandler.UpdateMeal)
					r.With(rbac.RequirePermission("diet:write", rbacRepo)).Delete("/{meal_id}", dietHandler.RemoveMeal)

					// Items
					r.With(rbac.RequirePermission("diet:write", rbacRepo)).Post("/{meal_id}/items", dietHandler.AddMealItem)
				})
			})

			// Diet items — standalone
			r.With(rbac.RequirePermission("diet:write", rbacRepo)).Put("/diet-items/{item_id}", dietHandler.UpdateMealItem)
			r.With(rbac.RequirePermission("diet:write", rbacRepo)).Delete("/diet-items/{item_id}", dietHandler.RemoveMealItem)
			r.With(rbac.RequirePermission("diet:write", rbacRepo)).Post("/diet-items/{item_id}/substitutions", dietHandler.AddSubstitution)
			r.With(rbac.RequirePermission("diet:write", rbacRepo)).Delete("/diet-substitutions/{sub_id}", dietHandler.RemoveSubstitution)

			// Documents — standalone resources
			r.Route("/documents/{id}", func(r chi.Router) {
				r.With(rbac.RequirePermission("clinical:read", rbacRepo)).Get("/", docHandler.GetByID)
				r.With(rbac.RequirePermission("document:write", rbacRepo)).Put("/", docHandler.Update)
				r.With(rbac.RequirePermission("document:write", rbacRepo)).Post("/finalize", docHandler.Finalize)
				r.With(rbac.RequirePermission("document:publish", rbacRepo)).Post("/publish", docHandler.Publish)
				r.With(rbac.RequirePermission("document:write", rbacRepo)).Post("/new-version", docHandler.NewVersion)
				r.With(rbac.RequirePermission("clinical:read", rbacRepo)).Get("/versions", docHandler.ListVersions)
			})

			// Exports
			r.Route("/exports", func(r chi.Router) {
				r.With(rbac.RequirePermission("export:pdf", rbacRepo)).Post("/", exportHandler.RequestExport)
				r.With(rbac.RequirePermission("export:pdf", rbacRepo)).Get("/", exportHandler.List)
				r.Route("/{id}", func(r chi.Router) {
					r.With(rbac.RequirePermission("export:pdf", rbacRepo)).Get("/", exportHandler.GetByID)
					r.With(rbac.RequirePermission("export:pdf", rbacRepo)).Get("/download", exportHandler.Download)
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

			// Phase 5 — Google Calendar integration
			if googleHandler != nil {
				r.Route("/integrations/google", func(r chi.Router) {
					r.With(rbac.RequirePermission("schedule:manage", rbacRepo)).Get("/authorize", googleHandler.Authorize)
					r.With(rbac.RequirePermission("schedule:manage", rbacRepo)).Get("/callback", googleHandler.Callback)
					r.With(rbac.RequirePermission("schedule:manage", rbacRepo)).Post("/disconnect", googleHandler.Disconnect)
					r.With(rbac.RequirePermission("schedule:manage", rbacRepo)).Get("/status", googleHandler.Status)
				})
			}

			// Phase 5 — AI assistive
			if aiHandler != nil {
				r.Route("/ai/suggestions", func(r chi.Router) {
					r.With(rbac.RequirePermission("ai:suggest", rbacRepo)).Post("/", aiHandler.Create)
					r.With(rbac.RequirePermission("ai:read", rbacRepo)).Get("/", aiHandler.List)
					r.Route("/{id}", func(r chi.Router) {
						r.With(rbac.RequirePermission("ai:read", rbacRepo)).Get("/", aiHandler.GetByID)
						r.With(rbac.RequirePermission("ai:suggest", rbacRepo)).Post("/accept", aiHandler.Accept)
						r.With(rbac.RequirePermission("ai:suggest", rbacRepo)).Post("/reject", aiHandler.Reject)
					})
				})
			}
		})

		// --- Phase 4: Backoffice routes ---
		r.Route("/backoffice", func(r chi.Router) {
			r.Use(backofficeMW)

			// Tenants
			r.With(backoffice.RequireBackofficePermission("tenants:manage", boRepo)).Get("/tenants", boHandler.ListTenants)

			r.Route("/tenants/{id}", func(r chi.Router) {
				r.With(backoffice.RequireBackofficePermission("tenants:manage", boRepo)).Get("/", boHandler.GetTenantDetail)

				// Subscription management
				r.With(backoffice.RequireBackofficePermission("billing:manage", boRepo)).Get("/subscription", boHandler.GetSubscription)
				r.With(backoffice.RequireBackofficePermission("billing:manage", boRepo)).Patch("/subscription/plan", boHandler.ChangePlan)
				r.With(backoffice.RequireBackofficePermission("billing:manage", boRepo)).Post("/subscription/cancel", boHandler.CancelSubscription)
				r.With(backoffice.RequireBackofficePermission("billing:manage", boRepo)).Post("/subscription/reactivate", boHandler.ReactivateSubscription)

				// Invoices & payments
				r.With(backoffice.RequireBackofficePermission("billing:manage", boRepo)).Get("/invoices", boHandler.ListInvoices)
				r.With(backoffice.RequireBackofficePermission("billing:manage", boRepo)).Get("/payments", boHandler.ListPayments)

				// Overrides
				r.With(backoffice.RequireBackofficePermission("overrides:manage", boRepo)).Get("/overrides", boHandler.ListOverrides)
				r.With(backoffice.RequireBackofficePermission("overrides:manage", boRepo)).Post("/overrides", boHandler.CreateOverride)
				r.With(backoffice.RequireBackofficePermission("overrides:manage", boRepo)).Put("/overrides/{feature_key}", boHandler.UpdateOverride)
				r.With(backoffice.RequireBackofficePermission("overrides:manage", boRepo)).Delete("/overrides/{feature_key}", boHandler.DeleteOverride)

				// Audit trail
				r.With(backoffice.RequireBackofficePermission("support:manage", boRepo)).Get("/audit", boHandler.ListAuditLogs)
			})
		})
	})

	// --- Start worker & server ---

	jobQueue.Start(ctx)
	bgWorker.Start(ctx)

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
	jobQueue.Stop()
	bgWorker.Stop()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error("shutdown error", "error", err)
	}
	log.Info("server stopped")
}
