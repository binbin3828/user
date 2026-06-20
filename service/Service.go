/*
 * @Autor: Bobby
 * @Description: service 层依赖注入容器
 * @Date: 2022-06-20
 * @FilePath: \user\service\Service.go
 */

package service

import (
	"user/dao"
	"user/pkg/logger"
	"user/pkg/mailer"
	"user/pkg/ratelimit"
)

// Service 持有各 DAO 接口和日志，通过构造函数注入依赖
type Service struct {
	Logger           logger.Logger
	UserDao          dao.IUserDao
	FriendsDao       dao.IFriendsDao
	FriendRequestDao  dao.IFriendRequestDao
	BlacklistDao      dao.IBlacklistDao
	PasswordResetDao  dao.IPasswordResetDao
	Mailer            mailer.Mailer
	RateLimiter       ratelimit.Limiter
}

// NewService 创建 Service 实例，注入依赖
func NewService(log logger.Logger, userDao dao.IUserDao, friendsDao dao.IFriendsDao, friendRequestDao dao.IFriendRequestDao, blacklistDao dao.IBlacklistDao, passwordResetDao dao.IPasswordResetDao, mailer mailer.Mailer, rl ratelimit.Limiter) *Service {
	return &Service{
		Logger:           log,
		UserDao:          userDao,
		FriendsDao:       friendsDao,
		FriendRequestDao: friendRequestDao,
		BlacklistDao:     blacklistDao,
		PasswordResetDao: passwordResetDao,
		Mailer:           mailer,
		RateLimiter:      rl,
	}
}
