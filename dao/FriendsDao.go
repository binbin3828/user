/*
 * @Autor: Bobby
 * @Description: friends dao
 * @Date: 2022-06-09 15:20:11
 * @LastEditTime: 2022-06-10 11:30:34
 * @FilePath: \user\dao\FriendsDao.go
 */

package dao

import (
	"time"
	"user/model"
	"user/pkg/logger"
	"user/pkg/util"

	"github.com/jinzhu/gorm"
)

// IFriendsDao 好友数据访问接口
type IFriendsDao interface {
	AddFriend(uid, fri int) error
	GetFriendsList(uid int) ([]*model.RetListFriends, error)
	GetNearbyFriend(uid int, subStr string) ([]*model.RetNearbyFriendsList, error)
}

// 编译期检查 FriendsDao 是否实现 IFriendsDao
var _ IFriendsDao = (*FriendsDao)(nil)

type FriendsDao struct {
	db     *gorm.DB
	logger logger.Logger
}

func NewFriendsDao(db *gorm.DB, log logger.Logger) *FriendsDao {
	return &FriendsDao{db: db, logger: log}
}

func (T *FriendsDao) GetNearbyFriend(uid int, subStr string) ([]*model.RetNearbyFriendsList, error) {
	var list []*model.RetNearbyFriendsList
	T.logger.Debug("substr:", subStr)
	subStr = subStr + "%"
	err := T.db.Raw("select user.id as fri_uid,user.name as fri_name, user.latitude, user.longitude, user.loc_geohash, friends.create_time from friends,user where friends.uid = ? AND friends.fri = user.id AND user.loc_geohash LIKE ?", uid, subStr).Scan(&list).Error
	if err != nil {
		return list, err
	}
	return list, nil
}

func (T *FriendsDao) GetFriendsList(uid int) ([]*model.RetListFriends, error) {
	var list []*model.RetListFriends
	sql := "select user.id as fri_uid,user.name as fri_name, friends.create_time from friends,user where friends.uid = ? AND friends.fri = user.id"
	err := T.db.Raw(sql, uid).Scan(&list).Error
	if err != nil {
		return list, err
	}
	return list, nil
}

func (T *FriendsDao) AddFriend(uid, fri int) error {
	friends1 := model.Friends{
		Uid:        uid,
		Fri:        fri,
		CreateTime: util.JsonTime(time.Now()),
	}
	friends2 := model.Friends{
		Uid:        fri,
		Fri:        uid,
		CreateTime: util.JsonTime(time.Now()),
	}

	tx := T.db.Begin()
	err := tx.Table("friends").Create(&friends1).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Table("friends").Create(&friends2).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
