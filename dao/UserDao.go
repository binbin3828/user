package dao

import (
	"context"
	"time"
	"user/model"
	"user/pkg/logger"
	"user/pkg/util"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

var daoUserTracer = otel.Tracer("dao.user")

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
	ctx, span := daoUserTracer.Start(ctx, "UserDao.CreateUser", trace.WithAttributes(attribute.String("table", "user")))
	defer span.End()

	user.CreateAt = util.JsonTime(time.Now())
	return T.db.WithContext(ctx).Table("user").Create(user).Error
}

func (T *UserDao) FindUser(ctx context.Context, id int) (*model.User, error) {
	ctx, span := daoUserTracer.Start(ctx, "UserDao.FindUser", trace.WithAttributes(attribute.Int("user.id", id)))
	defer span.End()

	var user model.User
	err := T.db.WithContext(ctx).Table("user").Where("id=?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (T *UserDao) FindUserByName(ctx context.Context, name string) (*model.User, error) {
	ctx, span := daoUserTracer.Start(ctx, "UserDao.FindUserByName", trace.WithAttributes(attribute.String("user.name", name)))
	defer span.End()

	var user model.User
	err := T.db.WithContext(ctx).Table("user").Where("name=?", name).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (T *UserDao) DeleteUser(ctx context.Context, uid int) error {
	ctx, span := daoUserTracer.Start(ctx, "UserDao.DeleteUser", trace.WithAttributes(attribute.Int("user.id", uid)))
	defer span.End()

	return T.db.WithContext(ctx).Table("user").Where("id=?", uid).Delete(&model.User{}).Error
}

func (T *UserDao) UpdateUser(ctx context.Context, uid int, modifyArr map[string]interface{}) error {
	ctx, span := daoUserTracer.Start(ctx, "UserDao.UpdateUser", trace.WithAttributes(attribute.Int("user.id", uid)))
	defer span.End()

	return T.db.WithContext(ctx).Table("user").Where("id=?", uid).Updates(modifyArr).Error
}
