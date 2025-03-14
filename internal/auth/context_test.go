package auth

import (
	"context"
	"testing"
)

func TestAddClaimsToContextAndGetClaimsFromContext(t *testing.T) {
	want := &Claims{
		Username: "testuser",
		Role:     "admin",
	}
	ctx := AddClaimsToContext(context.Background(), want)

	got, err := GetClaimsFromContext(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Username != want.Username || got.Role != want.Role {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

func TestGetClaimsFromContext_NoClaims(t *testing.T) {
	_, err := GetClaimsFromContext(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
