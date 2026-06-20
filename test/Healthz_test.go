package test

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestHealthz_ReturnsOk(t *testing.T) {
	svc, _, _ := newTestService()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/healthz", nil)

	svc.Healthz(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if int(resp["code"].(float64)) != 0 {
		t.Fatalf("expected code=0, got %d", int(resp["code"].(float64)))
	}
	if resp["data"] != "ok" {
		t.Errorf("expected data='ok', got '%v'", resp["data"])
	}
}

func TestReadyz_Ready(t *testing.T) {
	svc, userDao, _ := newTestService()
	userDao.Users[0] = nil

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/readyz", nil)

	svc.Readyz(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if int(resp["code"].(float64)) != 0 {
		t.Fatalf("expected code=0, got %d: %v", int(resp["code"].(float64)), resp["msg"])
	}
}
