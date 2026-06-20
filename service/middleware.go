package service

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"user/constant"
	"user/pkg/config"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ctxKey string

const reqIDKey ctxKey = "request_id"

func getReqID(ctx context.Context) string {
	id, _ := ctx.Value(reqIDKey).(string)
	return id
}

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader("X-Request-ID")
		if id == "" {
			id = uuid.NewString()
		}
		c.Set(string(reqIDKey), id)
		c.Header("X-Request-ID", id)
		c.Next()
	}
}

func (s *Service) AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			s.returnError(c, constant.ERROR_AUTH_FAIL, "authorization required")
			c.Abort()
			return
		}
		tokenStr := strings.TrimPrefix(auth, "Bearer ")
		claims, err := s.ParseToken(tokenStr)
		if err != nil {
			s.returnError(c, constant.ERROR_AUTH_FAIL, "invalid token")
			c.Abort()
			return
		}
		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}

func CORS() gin.HandlerFunc {
	origins := "*"
	if v := config.Get("config.cors.allowedOrigins"); v != nil {
		if s, ok := v.(string); ok && s != "" {
			origins = s
		}
	}

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origins == "*" {
			c.Header("Access-Control-Allow-Origin", "*")
		} else {
			for _, allowed := range strings.Split(origins, ",") {
				allowed = strings.TrimSpace(allowed)
				if allowed == origin || allowed == "*" {
					c.Header("Access-Control-Allow-Origin", origin)
					break
				}
			}
		}
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

const maxBodySize = 1 << 20

func MaxBodySize() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.ContentLength > maxBodySize {
			c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, gin.H{
				"code": constant.ERROR_PARAM_ERR,
				"msg":  "request body too large",
			})
			return
		}
		c.Next()
	}
}

func RequestTimeout(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "no-referrer")
		c.Header("Content-Security-Policy", "default-src 'none'")
		if c.Request.TLS != nil {
			c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}
		c.Next()
	}
}

type RateLimiter struct {
	mu       sync.Mutex
	visitors map[string]*visitor
	rate     int
	window   time.Duration
}

type visitor struct {
	count    int
	lastSeen time.Time
}

func NewRateLimiter(rate int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*visitor),
		rate:     rate,
		window:   window,
	}
	go func() {
		ticker := time.NewTicker(window)
		defer ticker.Stop()
		for range ticker.C {
			rl.mu.Lock()
			for ip, v := range rl.visitors {
				if time.Since(v.lastSeen) > window {
					delete(rl.visitors, ip)
				}
			}
			rl.mu.Unlock()
		}
	}()
	return rl
}

func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, ok := rl.visitors[ip]
	if !ok {
		rl.visitors[ip] = &visitor{count: 1, lastSeen: time.Now()}
		return true
	}

	if time.Since(v.lastSeen) > rl.window {
		v.count = 1
		v.lastSeen = time.Now()
		return true
	}

	v.lastSeen = time.Now()
	v.count++
	return v.count <= rl.rate
}

func RateLimitMiddleware(rl *RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !rl.Allow(ip) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"code": constant.ERROR_PARAM_ERR,
				"msg":  "too many requests",
			})
			return
		}
		c.Next()
	}
}

func (s *Service) AuditLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		userID, _ := c.Get("user_id")
		uidStr := ""
		if id, ok := userID.(int); ok {
			uidStr = strconv.Itoa(id)
		}

		reqID := getReqID(c.Request.Context())
		ctx := c.Request.Context()
		traceID := traceIDFromContext(ctx)
		spanID := spanIDFromContext(ctx)
		status := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path

		auditFormat := "[audit] request_id=%s user_id=%s method=%s path=%s status=%d"
		if traceID != "" {
			auditFormat = "[trace_id=" + traceID + "] [span_id=" + spanID + "] " + auditFormat
		}
		s.Logger.Infof(auditFormat, reqID, uidStr, method, path, status)
	}
}
