package dao

import (
	"context"
	"time"
	"user/model"
	"user/pkg/logger"
	"user/pkg/util"

	"gorm.io/gorm"
)

type IUserDao interface {
	CreateUser(ctx context.Context, user *model.User) error
	FindUser(ctx context.Context, id int) (*model.User, error)
	FindUserByName(ctx context.Context, name string) (*model.User, error)
	DeleteUser(ctx context.Context, uid int) error
	UpdateUser(ctx context.Context, uid int, modifyArr map[string]interface{}) error
}

var _ IUserDao = (*UserDao)(nil)

type UserDao struct {
	db     *gorm.DB
	logger logger.Logger
}

func NewUserDao(db *gorm.DB, log logger.Logger) *UserDao {
	return &UserDao{db: db, logger: log}
}

func (T *UserDao) CreateUser(ctx context.Context, user *model.User) error {
	user.CreateAt = util.JsonTime(time.Now())
	return T.db.WithContext(ctx).Table("user").Create(user).Error
}

func (T *UserDao) FindUser(ctx context.Context, id int) (*model.User, error) {
	var user model.User
	err := T.db.WithContext(ctx).Table("user").Where("id=?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (T *UserDao) FindUserByName(ctx context.Context, name string) (*model.User, error) {
	var user model.User
	err := T.db.WithContext(ctx).Table("user").Where("name=?", name).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (T *UserDao) DeleteUser(ctx context.Context, uid int) error {
	return T.db.WithContext(ctx).Table("user").Where("id=?", uid).Delete(&model.User{}).Error
}

func (T *UserDao) UpdateUser(ctx context.Context, uid int, modifyArr map[string]interface{}) error {
	return T.db.WithContext(ctx).Table("user").Where("id=?", uid).Updates(modifyArr).Error
}
