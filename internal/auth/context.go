package auth

import (
	"context"
	"errors"
)

type contextKey string

const claimsContextKey contextKey = "claims"

// AddClaimsToContext adds JWT claims to the request context
func AddClaimsToContext(ctx context.Context, claims *Claims) context.Context {
	return context.WithValue(ctx, claimsContextKey, claims)
}

// GetClaimsFromContext retrieves JWT claims from the request context
func GetClaimsFromContext(ctx context.Context) (*Claims, error) {
	claims, ok := ctx.Value(claimsContextKey).(*Claims)
	if !ok {
		return nil, errors.New("no claims found in context")
	}
	return claims, nil
}
