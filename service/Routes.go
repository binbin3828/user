package service

import (
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "user/docs"
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

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/v1")
	{
		v1.POST("/auth/login", RateLimitMiddleware(loginLimiter), s.Login)

		auth := v1.Group("")
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
	}

	return r
}
