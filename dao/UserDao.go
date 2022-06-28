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
	"user/pkg/dbconn"
	"user/pkg/util"
)

type UserDao struct {
}

func (T *UserDao) CreateUser(user *model.User) error {
	user.CreateAt = util.JsonTime(time.Now())
	err := dbconn.GetMysql().Table("user").Create(user).Error
	if err != nil {
		return err
	}
	return nil
}

func (T *UserDao) FindUser(id int) (*model.User, error) {
	var user model.User
	err := dbconn.GetMysql().
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
	err := dbconn.GetMysql().
		Table("user").
		Where("id=?", uid).
		Delete(&model.User{}).Error
	return err
}

func (T *UserDao) UpdateUser(uid int, modifyArr map[string]interface{}) error {
	err := dbconn.GetMysql().
		Table("user").
		Where("id=?", uid).
		Updates(modifyArr).Error
	return err
}
