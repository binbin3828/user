package dao

import (
	"time"
	"user/model"
	"user/pkg/logger"
	"user/pkg/util"

	"gorm.io/gorm"
)

type IUserDao interface {
	CreateUser(user *model.User) error
	FindUser(id int) (*model.User, error)
	DeleteUser(uid int) error
	UpdateUser(uid int, modifyArr map[string]interface{}) error
}

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
	return T.db.Table("user").Create(user).Error
}

func (T *UserDao) FindUser(id int) (*model.User, error) {
	var user model.User
	err := T.db.Table("user").Where("id=?", id).First(&user).Error
	if err != nil {
		return &user, err
	}
	return &user, nil
}

func (T *UserDao) DeleteUser(uid int) error {
	return T.db.Table("user").Where("id=?", uid).Delete(&model.User{}).Error
}

func (T *UserDao) UpdateUser(uid int, modifyArr map[string]interface{}) error {
	return T.db.Table("user").Where("id=?", uid).Updates(modifyArr).Error
}
