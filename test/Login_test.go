package test

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
	"user/model"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func TestLogin_Success(t *testing.T) {
	svc, userDao, _ := newTestService()
	hash, _ := bcrypt.GenerateFromPassword([]byte("testpass"), bcrypt.DefaultCost)
	userDao.Users[1] = &model.User{Id: 1, Name: "testuser", Password: string(hash)}

	body := `{"name":"testuser","password":"testpass"}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/auth/login", strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	svc.Login(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if int(resp["code"].(float64)) != 0 {
		t.Fatalf("expected code=0, got %d: %s", int(resp["code"].(float64)), resp["msg"])
	}
	data := resp["data"].(map[string]interface{})
	if data["token"] == "" {
		t.Error("expected non-empty token")
	}
	if int(data["user_id"].(float64)) != 1 {
		t.Errorf("expected user_id=1, got %v", data["user_id"])
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	svc, userDao, _ := newTestService()
	hash, _ := bcrypt.GenerateFromPassword([]byte("testpass"), bcrypt.DefaultCost)
	userDao.Users[1] = &model.User{Id: 1, Name: "testuser", Password: string(hash)}

	body := `{"name":"testuser","password":"wrongpass"}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/auth/login", strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	svc.Login(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if int(resp["code"].(float64)) != -3 {
		t.Fatalf("expected code=-3, got %d", int(resp["code"].(float64)))
	}
	if resp["msg"] != "invalid credentials" {
		t.Errorf("expected 'invalid credentials', got '%v'", resp["msg"])
	}
}

func TestLogin_UserNotFound(t *testing.T) {
	svc, _, _ := newTestService()

	body := `{"name":"nonexistent","password":"testpass"}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/auth/login", strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	svc.Login(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if int(resp["code"].(float64)) != -3 {
		t.Fatalf("expected code=-3, got %d", int(resp["code"].(float64)))
	}
	if resp["msg"] != "invalid credentials" {
		t.Errorf("expected 'invalid credentials', got '%v'", resp["msg"])
	}
}

func TestLogin_MissingName(t *testing.T) {
	svc, _, _ := newTestService()

	body := `{"password":"testpass"}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/auth/login", strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	svc.Login(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if int(resp["code"].(float64)) != -1 {
		t.Fatalf("expected code=-1, got %d", int(resp["code"].(float64)))
	}
}

func TestLogin_MissingPassword(t *testing.T) {
	svc, _, _ := newTestService()

	body := `{"name":"testuser"}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/auth/login", strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	svc.Login(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if int(resp["code"].(float64)) != -1 {
		t.Fatalf("expected code=-1, got %d", int(resp["code"].(float64)))
	}
}
