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
	"user/pkg/dbconn"
	"user/pkg/logger"
	"user/pkg/util"
)

type FriendsDao struct {
}

func (T *FriendsDao) GetNearbyFriend(uid int, subStr string) ([]*model.RetNearbyFriendsList, error) {
	var list []*model.RetNearbyFriendsList
	logger.SugarLogger.Debug("substr:", subStr)
	subStr = subStr + "%"
	err := dbconn.GetMysql().Raw("select user.id as fri_uid,user.name as fri_name, user.latitude, user.longitude, user.loc_geohash, friends.create_time from friends,user where friends.uid = ? AND friends.fri = user.id AND user.loc_geohash LIKE ?", uid, subStr).Scan(&list).Error
	if err != nil {
		return list, err
	}
	return list, nil
}

func (T *FriendsDao) GetFriendsList(uid int) ([]*model.RetListFriends, error) {
	var list []*model.RetListFriends
	sql := "select user.id as fri_uid,user.name as fri_name, friends.create_time from friends,user where friends.uid = ? AND friends.fri = user.id"
	err := dbconn.GetMysql().Raw(sql, uid).Scan(&list).Error
	if err != nil {
		return list, err
	}
	return list, nil
}

func (T *FriendsDao) AddFriend(uid, fri int) error {

	// Be friends with each other
	// insert two records
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

	// use transaction
	tx := dbconn.GetMysql().Begin()
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
