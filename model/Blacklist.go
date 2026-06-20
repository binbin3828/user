package model

import "user/pkg/util"

type Blacklist struct {
	Uid        int           `json:"uid" gorm:"primaryKey"`
	BlockedUID int           `json:"blocked_uid" gorm:"column:blocked_uid;primaryKey"`
	CreatedAt  util.JsonTime `json:"created_at"`
}

func (Blacklist) TableName() string {
	return "blacklist"
}
