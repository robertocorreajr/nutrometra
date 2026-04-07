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
	"nutrometra/api/internal/identity"
	identityoidc "nutrometra/api/internal/identity/oidc"
	"nutrometra/api/internal/platform/audit"
	"nutrometra/api/internal/platform/config"
	"nutrometra/api/internal/platform/db"
	"nutrometra/api/internal/platform/logger"
	"nutrometra/api/internal/platform/observability"
	apiredis "nutrometra/api/internal/platform/redis"
	"nutrometra/api/internal/platform/server"
	"nutrometra/api/internal/rbac"
	rbacrepo "nutrometra/api/internal/rbac/repository"
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

	// --- Services ---

	auditSvc := audit.NewService()
	tenancyUC := tenancyuc.NewTenantUsecase(tenancyRepo)
	entitlementSvc := billinguc.NewEntitlementService(billingRepo)

	// --- Handlers ---

	tenancyHandler := tenancy.NewHandler(tenancyUC)
	rbacHandler := rbac.NewHandler(rbacRepo, pool, auditSvc)
	billingHandler := billing.NewHandler(billingRepo, entitlementSvc, pool, auditSvc)
	healthChecker := observability.NewHealthChecker(pool, redisClient, cfg.Zitadel.Issuer)

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
