package oidc

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"nutrometra/api/internal/identity/domain"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
)

// sentinel errors for the validator.
var (
	ErrNotInitialized = errors.New("validator not initialized")
	ErrTokenInvalid   = errors.New("token validation failed")
)

// Validator validates OIDC JWTs issued by Zitadel.
// Call Init() before using Validate().
// Use ValidateRaw() for testing without JWKS (no signature verification).
type Validator struct {
	issuer   string
	audience string
	jwks     keyfunc.Keyfunc
}

// NewValidator creates a new Validator. It does NOT make any network calls.
// You must call Init() before using Validate().
func NewValidator(issuer, audience string) *Validator {
	return &Validator{
		issuer:   strings.TrimRight(issuer, "/"),
		audience: audience,
	}
}

// Init fetches the JWKS from the Zitadel issuer. It retries up to 10 times
// with 3s delay between attempts. Respects context cancellation.
func (v *Validator) Init(ctx context.Context) error {
	jwksURL := v.issuer + "/oauth/v2/keys"

	var lastErr error
	const maxRetries = 10
	const retryDelay = 3 * time.Second

	for attempt := range maxRetries {
		select {
		case <-ctx.Done():
			return fmt.Errorf("%w: context cancelled during JWKS init", ErrTokenInvalid)
		default:
		}

		jwks, err := keyfunc.NewDefaultCtx(ctx, []string{jwksURL})
		if err == nil {
			v.jwks = jwks
			return nil
		}

		lastErr = err

		if attempt < maxRetries-1 {
			select {
			case <-ctx.Done():
				return fmt.Errorf("%w: context cancelled during JWKS init", ErrTokenInvalid)
			case <-time.After(retryDelay):
			}
		}
	}

	return fmt.Errorf("failed to fetch JWKS after %d attempts: %w", maxRetries, lastErr)
}

// Validate validates a JWT token using the JWKS keys fetched during Init().
//
// Security properties enforced:
//   - Only RS256 signing algorithm is accepted.
//   - Issuer must match exactly.
//   - Audience must contain the configured client ID.
//   - Expiry (exp) is required and validated.
//   - Issued-at (iat) is validated.
//   - Signature is verified against the JWKS keys.
//
// Returns opaque errors that do not leak key material or internal details.
func (v *Validator) Validate(ctx context.Context, rawToken string) (*domain.Claims, error) {
	if v.jwks == nil {
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
	token, err := parser.ParseWithClaims(rawToken, &claims, v.jwks.KeyfuncCtx(ctx))
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

// ValidateRaw parses a JWT token WITHOUT verifying the signature.
// This method exists for unit tests that cannot reach a real Zitadel instance.
// It still validates token structure, expiry, issuer, and audience.
//
// DO NOT use this method in production code paths.
func (v *Validator) ValidateRaw(_ context.Context, rawToken string) (*domain.Claims, error) {
	parser := jwt.NewParser(
		jwt.WithExpirationRequired(),
		jwt.WithIssuer(v.issuer),
		jwt.WithAudience(v.audience),
		jwt.WithIssuedAt(),
	)

	var claims zitadelClaims
	token, _, err := parser.ParseUnverified(rawToken, &claims)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrTokenInvalid, sanitizeError(err))
	}

	// ParseUnverified does not run the claims validator, so we must do it manually.
	validator := jwt.NewValidator(
		jwt.WithExpirationRequired(),
		jwt.WithIssuer(v.issuer),
		jwt.WithAudience(v.audience),
		jwt.WithIssuedAt(),
	)
	if err := validator.Validate(token.Claims); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrTokenInvalid, sanitizeError(err))
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

// zitadelClaims represents the JWT claims structure from Zitadel tokens.
type zitadelClaims struct {
	jwt.RegisteredClaims
	Email string `json:"email"`
	Name  string `json:"name"`
}

// sanitizeError strips internal details from JWT errors to prevent
// leaking key material or internal architecture details.
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
