package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"GoShort/internal/models"
	"GoShort/pkg/config"
	"GoShort/pkg/logger"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestMain(m *testing.M) {
	logger.Init() // Ensure logger is initialized so calls to logger.Error won't panic
	code := m.Run()
	os.Exit(code)
}

func TestLoginHandler_Success(t *testing.T) {
	os.Setenv("JWT_SECRET", "testsecret")
	config.Load()

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open in-memory DB: %v", err)
	}

	if err := db.AutoMigrate(&models.User{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte("testpass"), bcrypt.DefaultCost)
	user := models.User{
		Username: "testuser",
		Password: string(hashed),
		Role:     "user",
	}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	body, _ := json.Marshal(LoginRequest{Username: "testuser", Password: "testpass"})
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	w := httptest.NewRecorder()

	LoginHandler(db)(w, req)
	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var lr LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&lr); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if lr.User.Username != "testuser" || lr.User.Role != "user" || lr.Token == "" {
		t.Fatal("unexpected login response")
	}

	var updated models.User
	if err := db.Where("username = ?", "testuser").First(&updated).Error; err != nil {
		t.Fatalf("query error: %v", err)
	}
	if updated.LastLogin.Before(time.Now().Add(-1 * time.Minute)) {
		t.Fatal("last login not updated")
	}
}

func TestLoginHandler_InvalidRequest(t *testing.T) {
	os.Setenv("JWT_SECRET", "testsecret")
	config.Load()

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("db open error: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader([]byte("invalid-json")))
	w := httptest.NewRecorder()

	LoginHandler(db)(w, req)
	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestLoginHandler_UserNotFound(t *testing.T) {
	os.Setenv("JWT_SECRET", "testsecret")
	config.Load()

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("db open error: %v", err)
	}

	body, _ := json.Marshal(LoginRequest{Username: "notexisting", Password: "fake"})
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	w := httptest.NewRecorder()

	LoginHandler(db)(w, req)
	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp.StatusCode)
	}
}

func TestLoginHandler_WrongPassword(t *testing.T) {
	os.Setenv("JWT_SECRET", "testsecret")
	config.Load()

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("db open error: %v", err)
	}

	if err := db.AutoMigrate(&models.User{}); err != nil {
		t.Fatalf("migrate error: %v", err)
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte("realpass"), bcrypt.DefaultCost)
	user := models.User{
		Username: "testuser",
		Password: string(hashed),
		Role:     "admin",
	}
	db.Create(&user)

	body, _ := json.Marshal(LoginRequest{Username: "testuser", Password: "wrongpass"})
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	w := httptest.NewRecorder()

	LoginHandler(db)(w, req)
	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp.StatusCode)
	}
}