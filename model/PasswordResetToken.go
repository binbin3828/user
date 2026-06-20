package model

import "user/pkg/util"

type PasswordResetToken struct {
	Id        int           `json:"id" gorm:"primaryKey;autoIncrement"`
	UID       int           `json:"uid"`
	Token     string        `json:"token"`
	ExpiresAt util.JsonTime `json:"expires_at"`
	Used      bool          `json:"used"`
	CreatedAt util.JsonTime `json:"created_at"`
}

func (PasswordResetToken) TableName() string {
	return "password_reset_tokens"
}
