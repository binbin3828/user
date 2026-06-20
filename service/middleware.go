package service

import (
	"net/http"
	"strings"

	"user/constant"

	"github.com/gin-gonic/gin"
)

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
