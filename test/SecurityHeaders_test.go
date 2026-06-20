package test

import (
	"net/http/httptest"
	"testing"

	"user/service"

	"github.com/gin-gonic/gin"
)

func TestSecurityHeaders_SetsAll(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/user/1", nil)

	service.SecurityHeaders()(c)

	headers := map[string]string{
		"X-Content-Type-Options":  "nosniff",
		"X-Frame-Options":         "DENY",
		"X-XSS-Protection":        "1; mode=block",
		"Referrer-Policy":         "no-referrer",
		"Content-Security-Policy": "default-src 'none'",
	}

	for name, expected := range headers {
		actual := w.Header().Get(name)
		if actual == "" {
			t.Errorf("missing header: %s", name)
		} else if actual != expected {
			t.Errorf("header %s: expected '%s', got '%s'", name, expected, actual)
		}
	}
}

func TestSecurityHeaders_HSTSNotSetOnHTTP(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/user/1", nil)

	service.SecurityHeaders()(c)

	if v := w.Header().Get("Strict-Transport-Security"); v != "" {
		t.Errorf("expected no HSTS on HTTP, got '%s'", v)
	}
}
