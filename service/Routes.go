package service

import (
	"github.com/gin-gonic/gin"
)

func NewRouter(s *Service) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(s.LoggerMiddleware())

	r.GET("/user/:uid", s.GetUser)
	r.POST("/user", s.CreateUser)
	r.PUT("/user", s.ModifyUser)
	r.DELETE("/user/:uid", s.DeleteUser)

	r.POST("/friends", s.AddFriend)
	r.GET("/friends/:uid", s.GetFriendsList)
	r.GET("/nearbyfriends/:uid", s.GetNearbyFriend)

	return r
}
