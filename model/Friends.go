package model

import "user/pkg/util"

type Friends struct {
	Uid        int           `json:"uid" gorm:"primaryKey"`
	FriendID   int           `json:"friend_id" gorm:"column:fri;primaryKey"`
	CreateTime util.JsonTime `json:"create_at"`
}

func (Friends) TableName() string {
	return "friends"
}

type RetListFriends struct {
	FriUid     int           `json:"fri_uid"`
	FriName    string        `json:"fri_name"`
	CreateTime util.JsonTime `json:"create_at"`
}

type RetNearbyFriendsList struct {
	FriUid     int           `json:"fri_uid"`
	FriName    string        `json:"fri_name"`
	CreateTime util.JsonTime `json:"create_at"`
	Latitude   float64       `json:"latitude"`
	Longitude  float64       `json:"longitude"`
	LocGeohash string        `json:"loc_geohash"`
}
