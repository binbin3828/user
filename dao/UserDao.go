/*
 * @Autor: Bobby
 * @Description: User dao to do SQL option
 * @Date: 2022-06-06 17:55:22
 * @LastEditTime: 2022-06-09 15:15:52
 * @FilePath: \user\dao\UserDao.go
 */

package dao

import (
	"time"
	"user/model"
	"user/pkg/logger"
	"user/pkg/util"

	"github.com/jinzhu/gorm"
)

// IUserDao 用户数据访问接口
type IUserDao interface {
	CreateUser(user *model.User) error
	FindUser(id int) (*model.User, error)
	DeleteUser(uid int) error
	UpdateUser(uid int, modifyArr map[string]interface{}) error
}

// 编译期检查 UserDao 是否实现 IUserDao
var _ IUserDao = (*UserDao)(nil)

type UserDao struct {
	db     *gorm.DB
	logger logger.Logger
}

func NewUserDao(db *gorm.DB, log logger.Logger) *UserDao {
	return &UserDao{db: db, logger: log}
}

func (T *UserDao) CreateUser(user *model.User) error {
	user.CreateAt = util.JsonTime(time.Now())
	err := T.db.Table("user").Create(user).Error
	if err != nil {
		return err
	}
	return nil
}

func (T *UserDao) FindUser(id int) (*model.User, error) {
	var user model.User
	err := T.db.
		Table("user").
		Select("*").
		Where("id=?", id).
		First(&user).Error
	if err != nil {
		return &user, err
	}
	return &user, nil
}

func (T *UserDao) DeleteUser(uid int) error {
	err := T.db.
		Table("user").
		Where("id=?", uid).
		Delete(&model.User{}).Error
	return err
}

func (T *UserDao) UpdateUser(uid int, modifyArr map[string]interface{}) error {
	err := T.db.
		Table("user").
		Where("id=?", uid).
		Updates(modifyArr).Error
	return err
}
