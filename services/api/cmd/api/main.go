package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"nutrometra/api/internal/platform/config"
	"nutrometra/api/internal/platform/db"
	"nutrometra/api/internal/platform/logger"
	apiredis "nutrometra/api/internal/platform/redis"
	"nutrometra/api/internal/platform/server"
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

	srv := server.New(cfg.API.Port)

	// Temporary health endpoint (full observability added in Task 17)
	srv.Router().Get("/health", func(w http.ResponseWriter, r *http.Request) {
		server.RenderJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	// TODO: register module routes (Tasks 13-18)

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
