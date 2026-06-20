package service

import (
	"encoding/json"
	"strconv"
	"user/constant"
	"user/model"
	"user/pkg/util"

	"github.com/gin-gonic/gin"
	"github.com/mmcloughlin/geohash"
)

func (s *Service) GetUser(c *gin.Context) {
	uidStr := c.Param("uid")
	if uidStr == "" {
		s.returnError(c, constant.ERROR_PARAM_ERR, "param uid not set")
		return
	}
	uid, err := strconv.Atoi(uidStr)
	if err != nil {
		s.returnError(c, constant.ERROR_PARAM_ERR, err.Error())
		return
	}
	userInfo, err := s.UserDao.FindUser(c.Request.Context(), uid)
	if err != nil {
		code := constant.ERROR_MYSQL_ERR
		if ce, ok := err.(*util.CodeError); ok {
			code = ce.Code
		}
		s.returnError(c, code, err.Error())
		return
	}
	s.returnSuccess(c, userInfo)
}

func (s *Service) CreateUser(c *gin.Context) {
	reqBody, _ := c.GetRawData()

	var req createUserReq
	if err := json.Unmarshal(reqBody, &req); err != nil {
		s.returnError(c, constant.ERROR_PARAM_ERR, err.Error())
		return
	}
	if err := validateReq(req); err != nil {
		code := constant.ERROR_PARAM_ERR
		if ce, ok := err.(*util.CodeError); ok {
			code = ce.Code
		}
		s.returnError(c, code, err.Error())
		return
	}

	user := model.User{
		Name:        req.Name,
		Dob:         req.Dob,
		Address:     req.Address,
		Description: req.Description,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
	}

	if req.Latitude >= 0 && req.Longitude >= 0 {
		hash_base32 := geohash.EncodeWithPrecision(req.Latitude, req.Longitude, 8)
		user.LocGeohash = hash_base32
	}

	err := s.UserDao.CreateUser(c.Request.Context(), &user)
	if err != nil {
		code := constant.ERROR_MYSQL_ERR
		if ce, ok := err.(*util.CodeError); ok {
			code = ce.Code
		}
		s.returnError(c, code, err.Error())
		return
	}

	info, err := s.UserDao.FindUser(c.Request.Context(), user.Id)
	if err != nil {
		code := constant.ERROR_MYSQL_ERR
		if ce, ok := err.(*util.CodeError); ok {
			code = ce.Code
		}
		s.returnError(c, code, err.Error())
		return
	}
	s.returnSuccess(c, info)
}

func (s *Service) DeleteUser(c *gin.Context) {
	uidStr := c.Param("uid")
	if uidStr == "" {
		s.returnError(c, constant.ERROR_PARAM_ERR, "param uid not set")
		return
	}
	uid, err := strconv.Atoi(uidStr)
	if err != nil {
		s.returnError(c, constant.ERROR_PARAM_ERR, err.Error())
		return
	}
	err = s.UserDao.DeleteUser(c.Request.Context(), uid)
	if err != nil {
		code := constant.ERROR_MYSQL_ERR
		if ce, ok := err.(*util.CodeError); ok {
			code = ce.Code
		}
		s.returnError(c, code, err.Error())
		return
	}
	s.returnSuccess(c, "delete succ")
}

func (s *Service) ModifyUser(c *gin.Context) {
	reqBody, _ := c.GetRawData()

	var req modifyUserReq
	if err := json.Unmarshal(reqBody, &req); err != nil {
		s.returnError(c, constant.ERROR_PARAM_ERR, err.Error())
		return
	}
	if err := validateReq(req); err != nil {
		code := constant.ERROR_PARAM_ERR
		if ce, ok := err.(*util.CodeError); ok {
			code = ce.Code
		}
		s.returnError(c, code, err.Error())
		return
	}
	uid := int(req.Id)

	data := make(map[string]interface{})
	json.Unmarshal(reqBody, &data)

	modifyArr := make(map[string]interface{})
	if tmp, ok := data["name"].(string); ok {
		modifyArr["name"] = tmp
	}
	if tmp, ok := data["dob"].(string); ok {
		modifyArr["dob"] = tmp
	}
	if tmp, ok := data["address"].(string); ok {
		modifyArr["address"] = tmp
	}
	if tmp, ok := data["description"].(string); ok {
		modifyArr["description"] = tmp
	}

	tmpLatitude, ok1 := data["latitude"].(float64)
	if ok1 {
		modifyArr["latitude"] = tmpLatitude
	}
	tmpLongitude, ok2 := data["longitude"].(float64)
	if ok2 {
		modifyArr["longitude"] = tmpLongitude
	}
	if ok1 && ok2 && tmpLatitude >= 0 && tmpLongitude >= 0 {
		hash_base32 := geohash.EncodeWithPrecision(tmpLatitude, tmpLongitude, 8)
		modifyArr["loc_geohash"] = hash_base32
	}

	err := s.UserDao.UpdateUser(c.Request.Context(), uid, modifyArr)
	if err != nil {
		code := constant.ERROR_MYSQL_ERR
		if ce, ok := err.(*util.CodeError); ok {
			code = ce.Code
		}
		s.returnError(c, code, err.Error())
		return
	}
	info, err := s.UserDao.FindUser(c.Request.Context(), uid)
	if err != nil {
		code := constant.ERROR_MYSQL_ERR
		if ce, ok := err.(*util.CodeError); ok {
			code = ce.Code
		}
		s.returnError(c, code, err.Error())
		return
	}
	s.returnSuccess(c, info)
}
