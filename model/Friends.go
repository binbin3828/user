package model

import "user/pkg/util"

type Friends struct {
	Uid        int           `json:"uid"`
	FriendID   int           `json:"friend_id"`
	CreateTime util.JsonTime `json:"create_at"`
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
