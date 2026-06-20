package model

import (
	"user/pkg/util"
)

type User struct {
	Id          int           `json:"id"`
	Name        string        `json:"name"`
	Password    string        `json:"-"`
	Dob         string        `json:"dob"`
	Address     string        `json:"address"`
	Description string        `json:"description"`
	CreateAt    util.JsonTime `json:"create_at"`
	Latitude    float64       `json:"latitude"`
	Longitude   float64       `json:"longitude"`
	LocGeohash  string        `json:"loc_geohash"`
}
