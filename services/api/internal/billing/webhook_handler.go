package billing

import (
	"io"
	"log/slog"
	"net/http"

	"nutrometra/api/internal/billing/provider"
	"nutrometra/api/internal/billing/usecase"
	"nutrometra/api/internal/platform/server"
)

// WebhookHandler handles incoming billing provider webhooks.
type WebhookHandler struct {
	provider  provider.BillingProvider
	processor *usecase.WebhookProcessor
}

// NewWebhookHandler creates a webhook handler.
func NewWebhookHandler(bp provider.BillingProvider, processor *usecase.WebhookProcessor) *WebhookHandler {
	return &WebhookHandler{provider: bp, processor: processor}
}

// HandleStripeWebhook processes incoming Stripe webhook events.
// POST /webhooks/stripe — public route, no auth. Validated via HMAC signature.
func (h *WebhookHandler) HandleStripeWebhook(w http.ResponseWriter, r *http.Request) {
	const maxBodySize = 1 << 16 // 64KB
	body, err := io.ReadAll(io.LimitReader(r.Body, maxBodySize))
	if err != nil {
		server.RenderError(w, r, http.StatusBadRequest, "read_body_failed", "Failed to read request body")
		return
	}

	sigHeader := r.Header.Get("Stripe-Signature")
	if sigHeader == "" {
		server.RenderError(w, r, http.StatusBadRequest, "missing_signature", "Stripe-Signature header required")
		return
	}

	evt, err := h.provider.ConstructWebhookEvent(body, sigHeader)
	if err != nil {
		slog.ErrorContext(r.Context(), "webhook signature verification failed", "error", err)
		server.RenderError(w, r, http.StatusBadRequest, "invalid_signature", "Invalid webhook signature")
		return
	}

	// Process the event. Always return 200 to Stripe to prevent retries.
	if err := h.processor.ProcessEvent(r.Context(), evt.ID, evt.Type, evt.Data); err != nil {
		slog.ErrorContext(r.Context(), "webhook processing error", "event_id", evt.ID, "type", evt.Type, "error", err)
	}

	server.RenderJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
