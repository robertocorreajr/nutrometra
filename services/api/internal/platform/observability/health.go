package observability

import (
	"context"
	"net/http"
	"time"

	"nutrometra/api/internal/platform/server"

	"github.com/jackc/pgx/v5/pgxpool"
	goredis "github.com/redis/go-redis/v9"
)

// HealthChecker provides health and readiness probes.
type HealthChecker struct {
	db            *pgxpool.Pool
	redis         goredis.UniversalClient
	zitadelIssuer string
}

// NewHealthChecker creates a HealthChecker.
func NewHealthChecker(db *pgxpool.Pool, redis goredis.UniversalClient, zitadelIssuer string) *HealthChecker {
	return &HealthChecker{db: db, redis: redis, zitadelIssuer: zitadelIssuer}
}

type serviceStatus struct {
	Status  string `json:"status"`
	Latency string `json:"latency,omitempty"`
}

type healthResponse struct {
	Status   string                   `json:"status"`
	Services map[string]serviceStatus `json:"services"`
}

// Health checks connectivity with PostgreSQL, Redis and Zitadel.
// GET /health
func (h *HealthChecker) Health(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	services := map[string]serviceStatus{}
	overall := "ok"

	// PostgreSQL
	start := time.Now()
	if err := h.db.Ping(ctx); err != nil {
		services["postgres"] = serviceStatus{Status: "fail"}
		overall = "degraded"
	} else {
		services["postgres"] = serviceStatus{Status: "ok", Latency: time.Since(start).String()}
	}

	// Redis
	start = time.Now()
	if err := h.redis.Ping(ctx).Err(); err != nil {
		services["redis"] = serviceStatus{Status: "fail"}
		overall = "degraded"
	} else {
		services["redis"] = serviceStatus{Status: "ok", Latency: time.Since(start).String()}
	}

	// Zitadel
	start = time.Now()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, h.zitadelIssuer+"/debug/healthz", nil)
	if err != nil {
		services["zitadel"] = serviceStatus{Status: "fail"}
		overall = "degraded"
	} else {
		resp, err := http.DefaultClient.Do(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			services["zitadel"] = serviceStatus{Status: "fail"}
			overall = "degraded"
			if resp != nil {
				resp.Body.Close()
			}
		} else {
			resp.Body.Close()
			services["zitadel"] = serviceStatus{Status: "ok", Latency: time.Since(start).String()}
		}
	}

	status := http.StatusOK
	if overall != "ok" {
		status = http.StatusServiceUnavailable
	}

	server.RenderJSON(w, status, healthResponse{Status: overall, Services: services})
}

// Ready returns 200 only if PostgreSQL and Redis are healthy.
// GET /ready
func (h *HealthChecker) Ready(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	if err := h.db.Ping(ctx); err != nil {
		server.RenderError(w, r, http.StatusServiceUnavailable, "not_ready", "Database not ready")
		return
	}
	if err := h.redis.Ping(ctx).Err(); err != nil {
		server.RenderError(w, r, http.StatusServiceUnavailable, "not_ready", "Redis not ready")
		return
	}
	server.RenderJSON(w, http.StatusOK, map[string]string{"status": "ready"})
}
