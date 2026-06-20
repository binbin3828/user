package test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"user/pkg/ratelimit"
	"user/service"

	"github.com/gin-gonic/gin"
)

func TestRateLimiter_AllowsWithinLimit(t *testing.T) {
	rl := ratelimit.NewMemoryLimiter(5, time.Minute)

	for i := 0; i < 5; i++ {
		if !rl.Allow("1.2.3.4") {
			t.Fatalf("iteration %d: expected allow, got block", i)
		}
	}
}

func TestRateLimiter_BlocksExceeded(t *testing.T) {
	rl := ratelimit.NewMemoryLimiter(2, time.Minute)

	if !rl.Allow("1.2.3.4") {
		t.Fatal("expected allow")
	}
	if !rl.Allow("1.2.3.4") {
		t.Fatal("expected allow")
	}
	if rl.Allow("1.2.3.4") {
		t.Fatal("expected block after limit exceeded")
	}
}

func TestRateLimiter_DifferentIPsIndependent(t *testing.T) {
	rl := ratelimit.NewMemoryLimiter(1, time.Minute)

	if !rl.Allow("1.2.3.4") {
		t.Fatal("expected allow for ip A")
	}
	if rl.Allow("1.2.3.4") {
		t.Fatal("expected block for ip A after limit")
	}
	if !rl.Allow("5.6.7.8") {
		t.Fatal("expected allow for ip B")
	}
}

func TestRateLimiter_ResetsAfterWindow(t *testing.T) {
	rl := ratelimit.NewMemoryLimiter(1, 50*time.Millisecond)

	if !rl.Allow("1.2.3.4") {
		t.Fatal("expected allow")
	}
	if rl.Allow("1.2.3.4") {
		t.Fatal("expected block")
	}

	time.Sleep(60 * time.Millisecond)

	if !rl.Allow("1.2.3.4") {
		t.Fatal("expected allow after window reset")
	}
}

func TestRateLimitMiddleware_DoesNotCrash(t *testing.T) {
	rl := ratelimit.NewMemoryLimiter(100, time.Minute)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)

	service.RateLimitMiddleware(rl)(c)

	if w.Code == http.StatusTooManyRequests {
		t.Fatal("unexpected rate limit with high threshold")
	}
}
