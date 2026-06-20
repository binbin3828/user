package dao

import (
	"time"
	"user/model"
	"user/pkg/logger"
	"user/pkg/util"

	"gorm.io/gorm"
)

type IFriendsDao interface {
	AddFriend(uid, fri int) error
	GetFriendsList(uid int) ([]*model.RetListFriends, error)
	GetNearbyFriend(uid int, subStr string) ([]*model.RetNearbyFriendsList, error)
}

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
	err := T.db.Raw("select user.id as fri_uid,user.name as fri_name, friends.create_time from friends,user where friends.uid = ? AND friends.fri = user.id", uid).Scan(&list).Error
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
	if err := tx.Table("friends").Create(&friends1).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Table("friends").Create(&friends2).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
