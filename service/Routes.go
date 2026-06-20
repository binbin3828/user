package service

import (
	"time"

	"github.com/gin-gonic/gin"
)

func NewRouter(s *Service) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(RequestID())
	r.Use(CORS())
	r.Use(SecurityHeaders())
	r.Use(TracingMiddleware("user-service"))
	r.Use(MetricsMiddleware())
	r.Use(s.LoggerMiddleware())
	r.Use(MaxBodySize())
	r.Use(RequestTimeout(30 * time.Second))
	r.Use(s.AuditLog())

	loginLimiter := NewRateLimiter(10, time.Minute)

	r.GET("/healthz", s.Healthz)
	r.GET("/readyz", s.Readyz)
	r.GET("/metrics", MetricsHandler())
	r.POST("/auth/login", RateLimitMiddleware(loginLimiter), s.Login)

	auth := r.Group("")
	auth.Use(s.AuthRequired())
	{
		auth.GET("/user/:uid", s.GetUser)
		auth.POST("/user", s.CreateUser)
		auth.PUT("/user", s.ModifyUser)
		auth.DELETE("/user/:uid", s.DeleteUser)

		auth.POST("/friends", s.AddFriend)
		auth.GET("/friends/:uid", s.GetFriendsList)
		auth.GET("/nearbyfriends/:uid", s.GetNearbyFriend)
	}

	return r
}
