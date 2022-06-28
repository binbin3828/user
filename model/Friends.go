/*
 * @Autor: Bobby
 * @Description: friends model
 * @Date: 2022-06-09 15:19:24
 * @LastEditTime: 2022-06-10 11:26:55
 * @FilePath: \user\model\Friends.go
 */

package model

import "user/pkg/util"

type Friends struct {
	Uid        int           `json:"uid"`       //my uid
	Fri        int           `json:"name"`      //friend uid
	CreateTime util.JsonTime `json:"create_at"` //relation create time
}

type RetListFriends struct {
	FriUid     int           `json:"fri_uid"`   //friend uid
	FriName    string        `json:"fri_name"`  //friends nick name
	CreateTime util.JsonTime `json:"create_at"` //relation create time
}

type RetNearbyFriendsList struct {
	FriUid     int           `json:"fri_uid"`     //friend uid
	FriName    string        `json:"fri_name"`    //friends nick name
	CreateTime util.JsonTime `json:"create_at"`   //relation create time
	Latitude   float64       `json:"latitude"`    // user loction latitude
	Longitude  float64       `json:"longitude"`   // user loction latitude
	LocGeohash string        `json:"loc_geohash"` // user loction geo hash value
}
