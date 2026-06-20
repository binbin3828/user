package model

import "user/pkg/util"

type FriendRequest struct {
	Id        int           `json:"id" gorm:"primaryKey;autoIncrement"`
	FromUID   int           `json:"from_uid"`
	ToUID     int           `json:"to_uid"`
	Status    string        `json:"status"` // pending, accepted, rejected
	CreatedAt util.JsonTime `json:"created_at"`
	UpdatedAt util.JsonTime `json:"updated_at"`
}

func (FriendRequest) TableName() string {
	return "friend_requests"
}
