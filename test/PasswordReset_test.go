package test

import (
	"encoding/json"
	"errors"
	"net/http/httptest"
	"strings"
	"testing"
	"user/model"

	"github.com/gin-gonic/gin"
)

func TestForgotPassword_Success(t *testing.T) {
	svc, userDao, _, _, _, resetDao := newTestService()
	userDao.Users[1] = &model.User{Id: 1, Name: "alice", Email: "alice@example.com", Password: "hash"}

	body := `{"email":"alice@example.com"}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/auth/forgot-password", strings.NewReader(body))

	svc.ForgotPassword(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	code := int(resp["code"].(float64))
	if code != 0 {
		t.Fatalf("expected code=0, got %d: %s", code, resp["msg"])
	}
	data := resp["data"].(map[string]interface{})
	if data["token"] == "" {
		t.Error("expected token in response")
	}
	if len(resetDao.Tokens) != 1 {
		t.Errorf("expected 1 token created, got %d", len(resetDao.Tokens))
	}
}

func TestForgotPassword_EmailNotFound(t *testing.T) {
	svc, _, _, _, _, _ := newTestService()

	body := `{"email":"unknown@example.com"}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/auth/forgot-password", strings.NewReader(body))

	svc.ForgotPassword(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	code := int(resp["code"].(float64))
	if code != 0 {
		t.Fatalf("expected code=0 (no info leak), got %d", code)
	}
	data := resp["data"].(map[string]interface{})
	if msg, ok := data["message"].(string); !ok || msg == "" {
		t.Error("expected friendly message")
	}
}

func TestForgotPassword_MissingEmail(t *testing.T) {
	svc, _, _, _, _, _ := newTestService()

	body := `{}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/auth/forgot-password", strings.NewReader(body))

	svc.ForgotPassword(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if int(resp["code"].(float64)) != -1 {
		t.Errorf("expected code=-1, got %v", resp["code"])
	}
}

func TestForgotPassword_InvalidEmail(t *testing.T) {
	svc, _, _, _, _, _ := newTestService()

	body := `{"email":"not-an-email"}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/auth/forgot-password", strings.NewReader(body))

	svc.ForgotPassword(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if int(resp["code"].(float64)) != -1 {
		t.Errorf("expected code=-1, got %v", resp["code"])
	}
}

func TestResetPassword_Success(t *testing.T) {
	svc, userDao, _, _, _, resetDao := newTestService()
	userDao.Users[1] = &model.User{Id: 1, Name: "alice", Email: "alice@example.com", Password: "oldhash"}
	tok, _ := resetDao.CreateToken(nil, 1)

	body := `{"token":"` + tok.Token + `","new_password":"newpassword123"}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/auth/reset-password", strings.NewReader(body))

	svc.ResetPassword(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	code := int(resp["code"].(float64))
	if code != 0 {
		t.Fatalf("expected code=0, got %d: %s", code, resp["msg"])
	}
	if userDao.Users[1].Password == "oldhash" {
		t.Error("password was not updated")
	}
	if !tok.Used {
		t.Error("token was not marked as used")
	}
}

func TestResetPassword_InvalidToken(t *testing.T) {
	svc, _, _, _, _, _ := newTestService()

	body := `{"token":"invalid-token","new_password":"newpassword123"}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/auth/reset-password", strings.NewReader(body))

	svc.ResetPassword(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if int(resp["code"].(float64)) != -1 || resp["msg"] != "invalid or expired token" {
		t.Errorf("got code=%v msg=%v, want code=-1 msg='invalid or expired token'", resp["code"], resp["msg"])
	}
}

func TestResetPassword_ExpiredToken(t *testing.T) {
	svc, userDao, _, _, _, resetDao := newTestService()
	userDao.Users[1] = &model.User{Id: 1, Name: "alice", Password: "old"}
	tok, _ := resetDao.CreateToken(nil, 1)
	tok.Used = true

	body := `{"token":"` + tok.Token + `","new_password":"newpassword123"}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/auth/reset-password", strings.NewReader(body))

	svc.ResetPassword(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if int(resp["code"].(float64)) != -1 {
		t.Errorf("expected error, got success")
	}
}

func TestResetPassword_MissingFields(t *testing.T) {
	svc, _, _, _, _, _ := newTestService()

	body := `{"token":"xxx"}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/auth/reset-password", strings.NewReader(body))

	svc.ResetPassword(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if int(resp["code"].(float64)) != -1 {
		t.Errorf("expected code=-1, got %v", resp["code"])
	}
}

func TestResetPassword_ShortPassword(t *testing.T) {
	svc, _, _, _, _, _ := newTestService()

	body := `{"token":"token","new_password":"short"}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/auth/reset-password", strings.NewReader(body))

	svc.ResetPassword(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if int(resp["code"].(float64)) != -1 {
		t.Errorf("expected code=-1, got %v", resp["code"])
	}
}

func TestResetPassword_DAOCreateError(t *testing.T) {
	svc, userDao, _, _, _, resetDao := newTestService()
	userDao.Users[1] = &model.User{Id: 1, Name: "alice", Email: "alice@example.com"}
	resetDao.CreateErr = errors.New("db error")

	body := `{"email":"alice@example.com"}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/auth/forgot-password", strings.NewReader(body))

	svc.ForgotPassword(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if int(resp["code"].(float64)) == 0 {
		t.Fatal("expected error, got success")
	}
}
