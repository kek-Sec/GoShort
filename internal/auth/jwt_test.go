package auth

import (
	"net/http/httptest"
	"testing"

	"GoShort/internal/models"
	"GoShort/pkg/config"
)

func TestGenerateToken(t *testing.T) {
	config.Load()
	user := models.User{ID: 1, Username: "tester", Role: "admin"}
	token, err := GenerateToken(user)
	if err != nil {
		t.Fatal(err)
	}
	if token == "" {
		t.Fatal("got empty token")
	}
}

func TestVerifyToken(t *testing.T) {
	config.Load()
	user := models.User{ID: 2, Username: "tester2", Role: "user"}
	token, _ := GenerateToken(user)
	claims, err := VerifyToken(token)
	if err != nil {
		t.Fatal(err)
	}
	if claims.Username != user.Username || claims.Role != user.Role {
		t.Fatal("claims do not match user")
	}
}

func TestExtractTokenFromRequest(t *testing.T) {
	tokenString := "dummy-token"
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	got, err := ExtractTokenFromRequest(req)
	if err != nil {
		t.Fatal(err)
	}
	if got != tokenString {
		t.Fatalf("wanted %s, got %s", tokenString, got)
	}
}

func TestGetUserFromToken(t *testing.T) {
	config.Load()
	user := models.User{ID: 3, Username: "tester3", Role: "admin"}
	token, _ := GenerateToken(user)
	gotUser, err := GetUserFromToken(token)
	if err != nil {
		t.Fatal(err)
	}
	if gotUser.Username != user.Username || gotUser.Role != user.Role {
		t.Fatal("user from token does not match")
	}
}
