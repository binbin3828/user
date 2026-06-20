package dao

import (
	"context"
	"time"
	"user/model"
	"user/pkg/logger"
	"user/pkg/util"

	"gorm.io/gorm"
)

type IFriendsDao interface {
	AddFriend(ctx context.Context, uid, friendID int) error
	GetFriendsList(ctx context.Context, uid int, limit, offset int) ([]*model.RetListFriends, error)
	CountFriendsList(ctx context.Context, uid int) (int64, error)
	GetNearbyFriend(ctx context.Context, uid int, subStr string, limit, offset int) ([]*model.RetNearbyFriendsList, error)
	CountNearbyFriend(ctx context.Context, uid int, subStr string) (int64, error)
}

var _ IFriendsDao = (*FriendsDao)(nil)

type FriendsDao struct {
	db     *gorm.DB
	logger logger.Logger
}

func NewFriendsDao(db *gorm.DB, log logger.Logger) *FriendsDao {
	return &FriendsDao{db: db, logger: log}
}

func (T *FriendsDao) GetNearbyFriend(ctx context.Context, uid int, subStr string, limit, offset int) ([]*model.RetNearbyFriendsList, error) {
	var list []*model.RetNearbyFriendsList
	T.logger.Debug("substr:", subStr)
	subStr = subStr + "%"
	err := T.db.WithContext(ctx).Raw("select user.id as fri_uid,user.name as fri_name, user.latitude, user.longitude, user.loc_geohash, friends.create_time from friends,user where friends.uid = ? AND friends.fri = user.id AND user.loc_geohash LIKE ? LIMIT ? OFFSET ?", uid, subStr, limit, offset).Scan(&list).Error
	if err != nil {
		return list, err
	}
	return list, nil
}

func (T *FriendsDao) CountNearbyFriend(ctx context.Context, uid int, subStr string) (int64, error) {
	var total int64
	subStr = subStr + "%"
	err := T.db.WithContext(ctx).Raw("select COUNT(*) from friends,user where friends.uid = ? AND friends.fri = user.id AND user.loc_geohash LIKE ?", uid, subStr).Scan(&total).Error
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (T *FriendsDao) GetFriendsList(ctx context.Context, uid int, limit, offset int) ([]*model.RetListFriends, error) {
	var list []*model.RetListFriends
	err := T.db.WithContext(ctx).Raw("select user.id as fri_uid,user.name as fri_name, friends.create_time from friends,user where friends.uid = ? AND friends.fri = user.id LIMIT ? OFFSET ?", uid, limit, offset).Scan(&list).Error
	if err != nil {
		return list, err
	}
	return list, nil
}

func (T *FriendsDao) CountFriendsList(ctx context.Context, uid int) (int64, error) {
	var total int64
	err := T.db.WithContext(ctx).Raw("select COUNT(*) from friends where uid = ?", uid).Scan(&total).Error
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (T *FriendsDao) AddFriend(ctx context.Context, uid, friendID int) error {
	friends1 := model.Friends{
		Uid:        uid,
		FriendID:   friendID,
		CreateTime: util.JsonTime(time.Now()),
	}
	friends2 := model.Friends{
		Uid:        friendID,
		FriendID:   uid,
		CreateTime: util.JsonTime(time.Now()),
	}

	tx := T.db.WithContext(ctx).Begin()
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
