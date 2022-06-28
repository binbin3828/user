/*
 * @Autor: Bobby
 * @Description: User model
 * @Date: 2022-06-06 14:54:26
 * @LastEditTime: 2022-06-09 18:00:40
 * @FilePath: \user\model\User.go
 */

package model

import (
	"user/pkg/util"
)

type User struct {
	Id          int           `json:"id"`          // user ID
	Name        string        `json:"name"`        // user name
	Dob         string        `json:"dob"`         // date of birth
	Address     string        `json:"address"`     // user address
	Description string        `json:"description"` // user description
	CreateAt    util.JsonTime `json:"create_at"`   // user created date
	Latitude    float64       `json:"latitude"`    // user loction latitude
	Longitude   float64       `json:"longitude"`   // user loction latitude
	LocGeohash  string        `json:"loc_geohash"` // user loction geo hash value
}
