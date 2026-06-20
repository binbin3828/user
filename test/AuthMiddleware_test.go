package test

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func testToken(userID int) string {
	claims := jwt.MapClaims{
		"user_id": float64(userID),
		"role":    "user",
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, _ := token.SignedString([]byte("change-me-in-production-use-a-long-random-string"))
	return s
}

func TestAuthRequired_MissingHeader(t *testing.T) {
	svc, _, _ := newTestService()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/user/1", nil)

	handler := svc.AuthRequired()
	handler(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	code := int(resp["code"].(float64))
	if code != -3 {
		t.Fatalf("expected code=-3, got %d: %v", code, resp["msg"])
	}
	if resp["msg"] != "authorization required" {
		t.Errorf("expected 'authorization required', got '%v'", resp["msg"])
	}
}

func TestAuthRequired_InvalidToken(t *testing.T) {
	svc, _, _ := newTestService()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/user/1", nil)
	c.Request.Header.Set("Authorization", "Bearer this.is.definitely.not.a.valid.jwt.token")

	handler := svc.AuthRequired()
	handler(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	code := int(resp["code"].(float64))
	if code != -3 {
		t.Fatalf("expected code=-3, got %d: %v", code, resp["msg"])
	}
	if resp["msg"] != "invalid token" {
		t.Errorf("expected 'invalid token', got '%v'", resp["msg"])
	}
}

func TestAuthRequired_InvalidScheme(t *testing.T) {
	svc, _, _ := newTestService()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/user/1", nil)
	c.Request.Header.Set("Authorization", "Basic dXNlcjpwYXNz")

	handler := svc.AuthRequired()
	handler(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	code := int(resp["code"].(float64))
	if code != -3 {
		t.Fatalf("expected code=-3, got %d: %v", code, resp["msg"])
	}
	if resp["msg"] != "authorization required" {
		t.Errorf("expected 'authorization required', got '%v'", resp["msg"])
	}
}

func TestAuthRequired_EmptyBearerToken(t *testing.T) {
	svc, _, _ := newTestService()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/user/1", nil)
	c.Request.Header.Set("Authorization", "Bearer ")

	handler := svc.AuthRequired()
	handler(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	code := int(resp["code"].(float64))
	if code != -3 {
		t.Fatalf("expected code=-3, got %d: %v", code, resp["msg"])
	}
}

func TestAuthRequired_SetsUserIDAndRole(t *testing.T) {
	svc, _, _ := newTestService()
	token := testToken(42)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/user/1", nil)
	c.Request.Header.Set("Authorization", "Bearer "+token)

	handler := svc.AuthRequired()
	handler(c)

	if c.IsAborted() {
		t.Fatal("expected request to proceed, got aborted")
	}
	uid, exists := c.Get("user_id")
	if !exists {
		t.Fatal("expected user_id in context")
	}
	if uid.(int) != 42 {
		t.Errorf("expected user_id=42, got %v", uid)
	}
	role, exists := c.Get("role")
	if !exists {
		t.Fatal("expected role in context")
	}
	if role.(string) != "user" {
		t.Errorf("expected role='user', got '%v'", role)
	}
}
