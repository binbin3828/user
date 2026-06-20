package service

import (
	"github.com/gin-gonic/gin"
)

func NewRouter(s *Service) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(CORS())
	r.Use(s.LoggerMiddleware())
	r.Use(MaxBodySize())

	r.POST("/auth/login", s.Login)

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
