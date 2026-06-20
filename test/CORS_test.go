package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"user/service"

	"github.com/gin-gonic/gin"
)

func TestCORS_SetsHeaders(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/user/1", nil)
	c.Request.Header.Set("Origin", "https://example.com")

	service.CORS()(c)

	if w.Header().Get("Access-Control-Allow-Origin") == "" {
		t.Error("expected Access-Control-Allow-Origin header")
	}
	if w.Header().Get("Access-Control-Allow-Methods") == "" {
		t.Error("expected Access-Control-Allow-Methods header")
	}
	if w.Header().Get("Access-Control-Allow-Headers") == "" {
		t.Error("expected Access-Control-Allow-Headers header")
	}
}

func TestCORS_HandlesPreflight(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("OPTIONS", "/user/1", nil)

	service.CORS()(c)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d", w.Code)
	}
}

func TestCORS_NonPreflightPassesThrough(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/user/1", nil)

	service.CORS()(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}
