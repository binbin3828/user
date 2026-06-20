package service

import (
	"encoding/json"
	"strconv"
	"user/constant"

	"github.com/gin-gonic/gin"
)

func parsePagination(c *gin.Context) (page, pageSize int) {
	page = 1
	pageSize = 20

	if p := c.Query("page"); p != "" {
		if n, err := strconv.Atoi(p); err == nil && n > 0 {
			page = n
		}
	}
	if ps := c.Query("page_size"); ps != "" {
		if n, err := strconv.Atoi(ps); err == nil && n > 0 {
			pageSize = n
			if pageSize > 100 {
				pageSize = 100
			}
		}
	}
	return
}

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

	callerID := s.currentUserID(c)
	if callerID != uid {
		s.returnError(c, constant.ERROR_PERMISSION_DENIED, "permission denied")
		return
	}

	info, err := s.UserDao.FindUser(c.Request.Context(), uid)
	if err != nil {
		s.returnError(c, constant.ERROR_MYSQL_ERR, sanitizeErr(err).Error())
		return
	}
	geohashStr := info.LocGeohash

	precision := 6
	if p := c.Query("precision"); p != "" {
		n, err := strconv.Atoi(p)
		if err != nil || n < 1 || n > 12 {
			s.returnError(c, constant.ERROR_PARAM_ERR, "param precision must be 1-12")
			return
		}
		precision = n
	}
	if len(geohashStr) < precision {
		precision = len(geohashStr)
	}
	likeSubStr := geohashStr[:precision]

	page, pageSize := parsePagination(c)
	offset := (page - 1) * pageSize

	total, err := s.FriendsDao.CountNearbyFriend(c.Request.Context(), uid, likeSubStr)
	if err != nil {
		s.returnError(c, constant.ERROR_MYSQL_ERR, sanitizeErr(err).Error())
		return
	}

	list, err := s.FriendsDao.GetNearbyFriend(c.Request.Context(), uid, likeSubStr, pageSize, offset)
	if err != nil {
		s.returnError(c, constant.ERROR_MYSQL_ERR, sanitizeErr(err).Error())
		return
	}

	s.returnSuccess(c, gin.H{
		"uid":       uid,
		"list":      list,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
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

	callerID := s.currentUserID(c)
	if callerID != uid {
		s.returnError(c, constant.ERROR_PERMISSION_DENIED, "permission denied")
		return
	}

	_, err = s.UserDao.FindUser(c.Request.Context(), uid)
	if err != nil {
		s.returnError(c, constant.ERROR_MYSQL_ERR, sanitizeErr(err).Error())
		return
	}

	page, pageSize := parsePagination(c)
	offset := (page - 1) * pageSize

	total, err := s.FriendsDao.CountFriendsList(c.Request.Context(), uid)
	if err != nil {
		s.returnError(c, constant.ERROR_MYSQL_ERR, sanitizeErr(err).Error())
		return
	}

	list, err := s.FriendsDao.GetFriendsList(c.Request.Context(), uid, pageSize, offset)
	if err != nil {
		s.returnError(c, constant.ERROR_MYSQL_ERR, sanitizeErr(err).Error())
		return
	}

	s.returnSuccess(c, gin.H{
		"uid":       uid,
		"list":      list,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
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
		s.returnError(c, constant.ERROR_PARAM_ERR, err.Error())
		return
	}

	uid := int(req.Uid)

	callerID := s.currentUserID(c)
	if callerID != uid {
		s.returnError(c, constant.ERROR_PERMISSION_DENIED, "permission denied")
		return
	}

	friendID := int(req.Fri)

	_, err := s.UserDao.FindUser(c.Request.Context(), uid)
	if err != nil {
		s.returnError(c, constant.ERROR_PARAM_ERR, "invalid friend request")
		return
	}

	_, err = s.UserDao.FindUser(c.Request.Context(), friendID)
	if err != nil {
		s.returnError(c, constant.ERROR_PARAM_ERR, "invalid friend request")
		return
	}

	err = s.FriendsDao.AddFriend(c.Request.Context(), uid, friendID)
	if err != nil {
		s.returnError(c, constant.ERROR_MYSQL_ERR, sanitizeErr(err).Error())
		return
	}

	s.returnSuccess(c, gin.H{
		"uid":       uid,
		"friend_id": friendID,
	})
}
