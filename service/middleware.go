package service

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"user/constant"

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
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
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

func (s *Service) AuditLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		userID, _ := c.Get("user_id")
		uidStr := ""
		if id, ok := userID.(int); ok {
			uidStr = strconv.Itoa(id)
		}

		reqID := getReqID(c.Request.Context())
		status := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path

		s.Logger.Infof("[audit] request_id=%s user_id=%s method=%s path=%s status=%d",
			reqID, uidStr, method, path, status)
	}
}
