/*
 * @Autor: Bobby
 * @Description: service 层依赖注入容器
 * @Date: 2022-06-20
 * @FilePath: \user\service\Service.go
 */

package service

import "user/dao"

// Service 持有各 DAO 接口，通过构造函数注入依赖
type Service struct {
	UserDao    dao.IUserDao
	FriendsDao dao.IFriendsDao
}

// NewService 创建 Service 实例，注入 DAO 依赖
func NewService(userDao dao.IUserDao, friendsDao dao.IFriendsDao) *Service {
	return &Service{
		UserDao:    userDao,
		FriendsDao: friendsDao,
	}
}
