package service

import (
	"encoding/json"
	"strconv"
	"user/constant"

	"github.com/gin-gonic/gin"
)

type blockUserReq struct {
	BlockedUID int `json:"blocked_uid" validate:"required"`
}

func (s *Service) BlockUser(c *gin.Context) {
	reqBody, _ := c.GetRawData()

	var req blockUserReq
	if err := json.Unmarshal(reqBody, &req); err != nil {
		s.returnError(c, constant.ERROR_PARAM_ERR, err.Error())
		return
	}
	if err := validateReq(req); err != nil {
		s.returnError(c, constant.ERROR_PARAM_ERR, err.Error())
		return
	}

	uid := s.currentUserID(c)
	if uid == req.BlockedUID {
		s.returnError(c, constant.ERROR_PARAM_ERR, "cannot block yourself")
		return
	}

	ctx := c.Request.Context()

	_, err := s.UserDao.FindUser(ctx, req.BlockedUID)
	if err != nil {
		s.returnError(c, constant.ERROR_PARAM_ERR, "target user not found")
		return
	}

	alreadyBlocked, err := s.BlacklistDao.IsBlocked(ctx, uid, req.BlockedUID)
	if err != nil {
		s.returnError(c, constant.ERROR_MYSQL_ERR, sanitizeErr(err).Error())
		return
	}
	if alreadyBlocked {
		s.returnError(c, constant.ERROR_PARAM_ERR, "user already blocked")
		return
	}

	err = s.BlacklistDao.Block(ctx, uid, req.BlockedUID)
	if err != nil {
		s.returnError(c, constant.ERROR_MYSQL_ERR, sanitizeErr(err).Error())
		return
	}

	s.returnSuccess(c, gin.H{
		"uid":         uid,
		"blocked_uid": req.BlockedUID,
		"status":      "blocked",
	})
}

func (s *Service) UnblockUser(c *gin.Context) {
	targetStr := c.Param("uid")
	if targetStr == "" {
		s.returnError(c, constant.ERROR_PARAM_ERR, "param uid not set")
		return
	}
	blockedUID, err := strconv.Atoi(targetStr)
	if err != nil {
		s.returnError(c, constant.ERROR_PARAM_ERR, err.Error())
		return
	}

	uid := s.currentUserID(c)

	err = s.BlacklistDao.Unblock(c.Request.Context(), uid, blockedUID)
	if err != nil {
		s.returnError(c, constant.ERROR_MYSQL_ERR, sanitizeErr(err).Error())
		return
	}

	s.returnSuccess(c, gin.H{
		"uid":         uid,
		"blocked_uid": blockedUID,
		"status":      "unblocked",
	})
}

func (s *Service) GetBlockedList(c *gin.Context) {
	uid := s.currentUserID(c)

	page, pageSize := parsePagination(c)
	offset := (page - 1) * pageSize

	list, total, err := s.BlacklistDao.GetBlockedList(c.Request.Context(), uid, pageSize, offset)
	if err != nil {
		s.returnError(c, constant.ERROR_MYSQL_ERR, sanitizeErr(err).Error())
		return
	}

	s.returnPaginated(c, list, total, page, pageSize)
}
