package service

import (
	"strings"
	"user/constant"
	"user/pkg/util"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type createUserReq struct {
	Name        string  `json:"name"     validate:"required"`
	Password    string  `json:"password" validate:"required,min=8"`
	Dob         string  `json:"dob"`
	Address     string  `json:"address"`
	Description string  `json:"description"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
}

type modifyUserReq struct {
	Id float64 `json:"id" validate:"required"`
}

type addFriendReq struct {
	Uid float64 `json:"uid" validate:"required"`
	Fri float64 `json:"fri" validate:"required"`
}

func validateReq(req interface{}) error {
	if err := validate.Struct(req); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			for _, e := range ve {
				field := strings.ToLower(e.Field())
				return util.NewCodeError(constant.ERROR_PARAM_ERR, "param "+field+" not set")
			}
		}
		return util.NewCodeError(constant.ERROR_PARAM_ERR, err.Error())
	}
	return nil
}
