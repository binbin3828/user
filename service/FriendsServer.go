package service

import (
	"encoding/json"
	"strconv"
	"user/constant"
	"user/pkg/util"

	"github.com/gin-gonic/gin"
)

func (s *Service) GetNearbyFriend(c *gin.Context) {
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

	info, err := s.UserDao.FindUser(c.Request.Context(), uid)
	if err != nil {
		code := constant.ERROR_MYSQL_ERR
		if ce, ok := err.(*util.CodeError); ok {
			code = ce.Code
		}
		s.returnError(c, code, err.Error())
		return
	}
	geohashStr := info.LocGeohash
	likeSubStr := geohashStr[0:6]

	list, err := s.FriendsDao.GetNearbyFriend(c.Request.Context(), uid, likeSubStr)
	if err != nil {
		code := constant.ERROR_MYSQL_ERR
		if ce, ok := err.(*util.CodeError); ok {
			code = ce.Code
		}
		s.returnError(c, code, err.Error())
		return
	}

	s.returnSuccess(c, gin.H{
		"uid":  uid,
		"list": list,
	})
}

func (s *Service) GetFriendsList(c *gin.Context) {
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

	_, err = s.UserDao.FindUser(c.Request.Context(), uid)
	if err != nil {
		code := constant.ERROR_MYSQL_ERR
		if ce, ok := err.(*util.CodeError); ok {
			code = ce.Code
		}
		s.returnError(c, code, err.Error())
		return
	}
	list, err := s.FriendsDao.GetFriendsList(c.Request.Context(), uid)
	if err != nil {
		code := constant.ERROR_MYSQL_ERR
		if ce, ok := err.(*util.CodeError); ok {
			code = ce.Code
		}
		s.returnError(c, code, err.Error())
		return
	}
	s.returnSuccess(c, gin.H{
		"uid":  uid,
		"list": list,
	})
}

func (s *Service) AddFriend(c *gin.Context) {
	reqBody, _ := c.GetRawData()

	var req addFriendReq
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

	uid := int(req.Uid)
	friendID := int(req.Fri)

	_, err := s.UserDao.FindUser(c.Request.Context(), uid)
	if err != nil {
		code := constant.ERROR_MYSQL_ERR
		if ce, ok := err.(*util.CodeError); ok {
			code = ce.Code
		}
		s.returnError(c, code, err.Error())
		return
	}

	_, err = s.UserDao.FindUser(c.Request.Context(), friendID)
	if err != nil {
		code := constant.ERROR_MYSQL_ERR
		if ce, ok := err.(*util.CodeError); ok {
			code = ce.Code
		}
		s.returnError(c, code, err.Error())
		return
	}

	err = s.FriendsDao.AddFriend(c.Request.Context(), uid, friendID)
	if err != nil {
		code := constant.ERROR_MYSQL_ERR
		if ce, ok := err.(*util.CodeError); ok {
			code = ce.Code
		}
		s.returnError(c, code, err.Error())
		return
	}

	s.returnSuccess(c, gin.H{
		"uid":       uid,
		"friend_id": friendID,
	})
}
