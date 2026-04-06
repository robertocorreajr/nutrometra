// Package oidc — test-only helpers.
// This file is compiled ONLY during `go test` and never in the production binary.
package oidc

import (
	"context"
	"fmt"

	"nutrometra/api/internal/identity/domain"

	"github.com/golang-jwt/jwt/v5"
)

// ValidateRaw parses a JWT without verifying the signature.
// It is defined here (package oidc, _test.go) so it is never compiled into
// the production binary. Tests in package oidc_test can call it because it is
// an exported method added to *Validator only during test builds.
//
// It still validates token structure, expiry, issuer, and audience.
// Accept both RS256 and HS256 to allow test tokens signed with symmetric keys.
func (v *Validator) ValidateRaw(_ context.Context, rawToken string) (*domain.Claims, error) {
	parser := jwt.NewParser(
		jwt.WithValidMethods([]string{"RS256", "HS256"}),
		jwt.WithExpirationRequired(),
		jwt.WithIssuedAt(),
	)

	var claims zitadelClaims
	token, _, err := parser.ParseUnverified(rawToken, &claims)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrTokenInvalid, sanitizeError(err))
	}

	// ParseUnverified does not run the claims validator — do it manually.
	claimsValidator := jwt.NewValidator(
		jwt.WithExpirationRequired(),
		jwt.WithIssuer(v.issuer),
		jwt.WithAudience(v.audience),
		jwt.WithIssuedAt(),
	)
	if err := claimsValidator.Validate(token.Claims); err != nil {
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
