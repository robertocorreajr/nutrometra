package oidc

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	"nutrometra/api/internal/identity/domain"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
)

// Sentinel errors.
var (
	ErrNotInitialized = errors.New("validator not initialized")
	ErrTokenInvalid   = errors.New("token validation failed")
	ErrInitFailed     = errors.New("validator initialization failed")
)

// Validator validates OIDC JWTs issued by Zitadel via JWKS.
// Call Init() before using Validate().
type Validator struct {
	issuer     string
	audience   string
	mu         sync.RWMutex // protects jwks and cancelJWKS
	jwks       keyfunc.Keyfunc
	cancelJWKS context.CancelFunc
}

// NewValidator creates a Validator. No network calls at construction.
// Call Init(ctx) before using Validate.
func NewValidator(issuer, audience string) *Validator {
	return &Validator{
		issuer:   strings.TrimRight(issuer, "/"),
		audience: audience,
	}
}

// Init fetches the JWKS from {issuer}/oauth/v2/keys and starts the
// background refresh goroutine. Retries up to 10 times with 3s delay.
// Respects context cancellation.
func (v *Validator) Init(ctx context.Context) error {
	u, err := url.Parse(v.issuer)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return fmt.Errorf("%w: invalid issuer URL %q (must be absolute with scheme and host)", ErrInitFailed, v.issuer)
	}

	jwksURL := v.issuer + "/oauth/v2/keys"

	const maxRetries = 10
	const retryDelay = 3 * time.Second

	var lastErr error
	for attempt := range maxRetries {
		select {
		case <-ctx.Done():
			return fmt.Errorf("%w: context cancelled", ErrInitFailed)
		default:
		}

		// Use a long-lived background context for the refresh goroutine lifecycle.
		// Cancellation is owned by Close().
		jwksCtx, cancel := context.WithCancel(context.Background())
		jwks, err := keyfunc.NewDefaultCtx(jwksCtx, []string{jwksURL})
		if err == nil {
			v.mu.Lock()
			if v.cancelJWKS != nil {
				v.cancelJWKS() // stop any previous refresh goroutine
			}
			v.jwks = jwks
			v.cancelJWKS = cancel
			v.mu.Unlock()
			return nil
		}
		cancel() // clean up goroutine if init failed
		lastErr = err

		if attempt < maxRetries-1 {
			select {
			case <-ctx.Done():
				return fmt.Errorf("%w: context cancelled", ErrInitFailed)
			case <-time.After(retryDelay):
			}
		}
	}

	// lastErr is intentionally not wrapped to avoid leaking network/TLS details into logs.
	_ = lastErr
	return fmt.Errorf("%w: failed to fetch JWKS after %d attempts", ErrInitFailed, maxRetries)
}

// Close stops the background JWKS refresh goroutine. Call during graceful shutdown.
func (v *Validator) Close() {
	v.mu.Lock()
	defer v.mu.Unlock()
	if v.cancelJWKS != nil {
		v.cancelJWKS()
		v.cancelJWKS = nil
		v.jwks = nil
	}
}

// Validate validates a JWT token against the JWKS fetched during Init.
//
// Security properties enforced:
//   - Only RS256 signing algorithm accepted (prevents algorithm confusion attacks).
//   - Issuer must match exactly.
//   - Audience must contain the configured client ID.
//   - Expiry (exp) is required and validated.
//   - Issued-at (iat) is validated.
//   - Signature verified against JWKS keys.
//
// Returns opaque errors that do not expose key material or internal details.
func (v *Validator) Validate(ctx context.Context, rawToken string) (*domain.Claims, error) {
	v.mu.RLock()
	jwks := v.jwks
	v.mu.RUnlock()

	if jwks == nil {
		return nil, fmt.Errorf("%w: call Init() first", ErrNotInitialized)
	}

	parser := jwt.NewParser(
		jwt.WithValidMethods([]string{"RS256"}),
		jwt.WithIssuer(v.issuer),
		jwt.WithAudience(v.audience),
		jwt.WithExpirationRequired(),
		jwt.WithIssuedAt(),
	)

	var claims zitadelClaims
	token, err := parser.ParseWithClaims(rawToken, &claims, jwks.KeyfuncCtx(ctx))
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrTokenInvalid, sanitizeError(err))
	}

	if !token.Valid {
		return nil, ErrTokenInvalid
	}

	if claims.Subject == "" {
		return nil, fmt.Errorf("%w: missing subject claim", ErrTokenInvalid)
	}

	return &domain.Claims{
		Subject: claims.Subject,
		Email:   claims.Email,
		Name:    claims.Name,
	}, nil
}

// zitadelClaims maps the JWT fields emitted by Zitadel.
type zitadelClaims struct {
	jwt.RegisteredClaims
	Email string `json:"email"`
	Name  string `json:"name"`
}

// sanitizeError maps JWT library errors to opaque messages, preventing
// leakage of key material or internal implementation details.
func sanitizeError(err error) string {
	switch {
	case errors.Is(err, jwt.ErrTokenMalformed):
		return "token is malformed"
	case errors.Is(err, jwt.ErrTokenExpired):
		return "token is expired"
	case errors.Is(err, jwt.ErrTokenNotValidYet):
		return "token is not valid yet"
	case errors.Is(err, jwt.ErrTokenSignatureInvalid):
		return "token signature is invalid"
	case errors.Is(err, jwt.ErrTokenInvalidAudience):
		return "token has invalid audience"
	case errors.Is(err, jwt.ErrTokenInvalidIssuer):
		return "token has invalid issuer"
	case errors.Is(err, jwt.ErrTokenUnverifiable):
		return "token is unverifiable"
	case errors.Is(err, jwt.ErrTokenInvalidClaims):
		return "token has invalid claims"
	case errors.Is(err, jwt.ErrTokenRequiredClaimMissing):
		return "token is missing required claim"
	default:
		return "token validation error"
	}
}
