package auth

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"GoShort/internal/models"
	"GoShort/pkg/config"
)

func setupMiddlewareTests() {
	os.Setenv("JWT_SECRET", "testsecret")
	config.Load()
}

func init() {
	setupMiddlewareTests()
}

func TestAuthMiddleware_NoToken(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	AuthMiddleware(dummyHandler)(w, req)
	if w.Result().StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Result().StatusCode)
	}
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer invalidtoken")
	w := httptest.NewRecorder()

	AuthMiddleware(dummyHandler)(w, req)
	if w.Result().StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Result().StatusCode)
	}
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	user := models.User{ID: 1, Username: "testuser", Role: "user"}
	token, _ := GenerateToken(user)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	AuthMiddleware(dummyHandler)(w, req)
	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Result().StatusCode)
	}
}

func TestAdminOnlyMiddleware_AsAdmin(t *testing.T) {
	user := models.User{ID: 2, Username: "adminuser", Role: "admin"}
	token, _ := GenerateToken(user)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	AdminOnlyMiddleware(dummyHandler)(w, req)
	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Result().StatusCode)
	}
}

func TestAdminOnlyMiddleware_AsUser(t *testing.T) {
	user := models.User{ID: 3, Username: "normaluser", Role: "user"}
	token, _ := GenerateToken(user)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	AdminOnlyMiddleware(dummyHandler)(w, req)
	if w.Result().StatusCode != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", w.Result().StatusCode)
	}
}

func dummyHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "OK")
}
