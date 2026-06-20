package service

import (
	"encoding/json"
	"errors"
	"strconv"
	"user/constant"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (s *Service) SendFriendRequest(c *gin.Context) {
	reqBody, _ := c.GetRawData()

	var req sendFriendRequestReq
	if err := json.Unmarshal(reqBody, &req); err != nil {
		s.returnError(c, constant.ERROR_PARAM_ERR, err.Error())
		return
	}
	if err := validateReq(req); err != nil {
		s.returnError(c, constant.ERROR_PARAM_ERR, err.Error())
		return
	}

	fromUID := s.currentUserID(c)
	toUID := req.ToUID

	if fromUID == toUID {
		s.returnError(c, constant.ERROR_PARAM_ERR, "cannot friend yourself")
		return
	}

	ctx := c.Request.Context()

	blocked, err := s.BlacklistDao.IsBlocked(ctx, toUID, fromUID)
	if err != nil {
		s.returnError(c, constant.ERROR_MYSQL_ERR, sanitizeErr(err).Error())
		return
	}
	if blocked {
		s.returnError(c, constant.ERROR_BLOCKED, "you have been blocked by this user")
		return
	}

	_, err = s.UserDao.FindUser(ctx, fromUID)
	if err != nil {
		s.returnError(c, constant.ERROR_PARAM_ERR, "invalid request")
		return
	}

	_, err = s.UserDao.FindUser(ctx, toUID)
	if err != nil {
		s.returnError(c, constant.ERROR_PARAM_ERR, "target user not found")
		return
	}

	alreadyFriends, err := s.FriendRequestDao.AreAlreadyFriends(ctx, fromUID, toUID)
	if err != nil {
		s.returnError(c, constant.ERROR_MYSQL_ERR, sanitizeErr(err).Error())
		return
	}
	if alreadyFriends {
		s.returnError(c, constant.ERROR_ALREADY_FRIENDS, "already friends")
		return
	}

	hasPending, err := s.FriendRequestDao.HasPendingRequest(ctx, fromUID, toUID)
	if err != nil {
		s.returnError(c, constant.ERROR_MYSQL_ERR, sanitizeErr(err).Error())
		return
	}
	if hasPending {
		s.returnError(c, constant.ERROR_FRIEND_REQUEST_EXISTS, "friend request already sent")
		return
	}

	reversePending, err := s.FriendRequestDao.HasPendingRequest(ctx, toUID, fromUID)
	if err != nil {
		s.returnError(c, constant.ERROR_MYSQL_ERR, sanitizeErr(err).Error())
		return
	}
	if reversePending {
		err = s.FriendsDao.AddFriend(ctx, fromUID, toUID)
		if err != nil {
			s.returnError(c, constant.ERROR_MYSQL_ERR, sanitizeErr(err).Error())
			return
		}
		s.returnSuccess(c, gin.H{
			"from_uid": fromUID,
			"to_uid":   toUID,
			"status":   "accepted",
			"message":  "friend request already pending, auto-accepted",
		})
		return
	}

	request, err := s.FriendRequestDao.CreateRequest(ctx, fromUID, toUID)
	if err != nil {
		s.returnError(c, constant.ERROR_MYSQL_ERR, sanitizeErr(err).Error())
		return
	}

	friendRequestsSent.Inc()

	s.returnSuccess(c, gin.H{
		"id":       request.Id,
		"from_uid": request.FromUID,
		"to_uid":   request.ToUID,
		"status":   request.Status,
		"message":  "friend request sent",
	})
}

func (s *Service) GetIncomingFriendRequests(c *gin.Context) {
	toUID := s.currentUserID(c)

	page, pageSize := parsePagination(c)
	offset := (page - 1) * pageSize

	status := c.DefaultQuery("status", "pending")

	list, total, err := s.FriendRequestDao.GetIncomingRequests(c.Request.Context(), toUID, status, pageSize, offset)
	if err != nil {
		s.returnError(c, constant.ERROR_MYSQL_ERR, sanitizeErr(err).Error())
		return
	}

	s.returnPaginated(c, list, total, page, pageSize)
}

func (s *Service) GetOutgoingFriendRequests(c *gin.Context) {
	fromUID := s.currentUserID(c)

	page, pageSize := parsePagination(c)
	offset := (page - 1) * pageSize

	status := c.DefaultQuery("status", "pending")

	list, total, err := s.FriendRequestDao.GetOutgoingRequests(c.Request.Context(), fromUID, status, pageSize, offset)
	if err != nil {
		s.returnError(c, constant.ERROR_MYSQL_ERR, sanitizeErr(err).Error())
		return
	}

	s.returnPaginated(c, list, total, page, pageSize)
}

func (s *Service) AcceptFriendRequest(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		s.returnError(c, constant.ERROR_PARAM_ERR, "param id not set")
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		s.returnError(c, constant.ERROR_PARAM_ERR, err.Error())
		return
	}

	ctx := c.Request.Context()

	req, err := s.FriendRequestDao.GetRequestByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.returnError(c, constant.ERROR_REQUEST_NOT_FOUND, "friend request not found")
			return
		}
		s.returnError(c, constant.ERROR_MYSQL_ERR, sanitizeErr(err).Error())
		return
	}

	currentUID := s.currentUserID(c)
	if req.ToUID != currentUID {
		s.returnError(c, constant.ERROR_PERMISSION_DENIED, "permission denied")
		return
	}

	if req.Status != "pending" {
		s.returnError(c, constant.ERROR_REQUEST_NOT_PENDING, "request is no longer pending")
		return
	}

	err = s.FriendsDao.AddFriend(ctx, req.FromUID, req.ToUID)
	if err != nil {
		s.returnError(c, constant.ERROR_MYSQL_ERR, sanitizeErr(err).Error())
		return
	}

	err = s.FriendRequestDao.UpdateRequestStatus(ctx, req.Id, "accepted")
	if err != nil {
		s.returnError(c, constant.ERROR_MYSQL_ERR, sanitizeErr(err).Error())
		return
	}

	friendAdditions.Inc()
	friendRequestsAccepted.Inc()

	s.returnSuccess(c, gin.H{
		"id":       req.Id,
		"from_uid": req.FromUID,
		"to_uid":   req.ToUID,
		"status":   "accepted",
	})
}

func (s *Service) RejectFriendRequest(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		s.returnError(c, constant.ERROR_PARAM_ERR, "param id not set")
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		s.returnError(c, constant.ERROR_PARAM_ERR, err.Error())
		return
	}

	ctx := c.Request.Context()

	req, err := s.FriendRequestDao.GetRequestByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.returnError(c, constant.ERROR_REQUEST_NOT_FOUND, "friend request not found")
			return
		}
		s.returnError(c, constant.ERROR_MYSQL_ERR, sanitizeErr(err).Error())
		return
	}

	currentUID := s.currentUserID(c)
	if req.ToUID != currentUID {
		s.returnError(c, constant.ERROR_PERMISSION_DENIED, "permission denied")
		return
	}

	if req.Status != "pending" {
		s.returnError(c, constant.ERROR_REQUEST_NOT_PENDING, "request is no longer pending")
		return
	}

	err = s.FriendRequestDao.UpdateRequestStatus(ctx, req.Id, "rejected")
	if err != nil {
		s.returnError(c, constant.ERROR_MYSQL_ERR, sanitizeErr(err).Error())
		return
	}

	s.returnSuccess(c, gin.H{
		"id":       req.Id,
		"from_uid": req.FromUID,
		"to_uid":   req.ToUID,
		"status":   "rejected",
	})
}
