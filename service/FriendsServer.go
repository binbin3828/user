package service

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"user/constant"
	"user/pkg/util"

	"github.com/go-chi/chi/v5"
)

func (s *Service) GetNearbyFriend(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	ctx := r.Context()
	uidStr := chi.URLParam(r, "uid")
	if uidStr == "" {
		return nil, util.NewCodeError(constant.ERROR_PARAM_ERR, "param uid not set")
	}
	uid, err := strconv.Atoi(uidStr)
	if err != nil {
		return nil, err
	}

	info, err := s.UserDao.FindUser(ctx, uid)
	if err != nil {
		return nil, err
	}
	geohashStr := info.LocGeohash

	likeSubStr := geohashStr[0:6]

	list, err := s.FriendsDao.GetNearbyFriend(ctx, uid, likeSubStr)
	if err != nil {
		return nil, err
	}

	data := make(map[string]interface{})
	data["uid"] = uid
	data["list"] = list
	return data, nil
}

func (s *Service) GetFriendsList(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	ctx := r.Context()
	uidStr := chi.URLParam(r, "uid")
	if uidStr == "" {
		return nil, util.NewCodeError(constant.ERROR_PARAM_ERR, "param uid not set")
	}
	uid, err := strconv.Atoi(uidStr)
	if err != nil {
		return nil, err
	}

	_, err = s.UserDao.FindUser(ctx, uid)
	if err != nil {
		return nil, err
	}
	list, err := s.FriendsDao.GetFriendsList(ctx, uid)
	if err != nil {
		return nil, err
	}
	data := make(map[string]interface{})
	data["uid"] = uid
	data["list"] = list
	return data, nil
}

func (s *Service) AddFriend(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	ctx := r.Context()
	reqBody, _ := io.ReadAll(r.Body)

	var req addFriendReq
	if err := json.Unmarshal(reqBody, &req); err != nil {
		return nil, util.NewCodeError(constant.ERROR_PARAM_ERR, err.Error())
	}
	if err := validateReq(req); err != nil {
		return nil, err
	}

	uid := int(req.Uid)
	friendID := int(req.Fri)

	_, err := s.UserDao.FindUser(ctx, uid)
	if err != nil {
		return nil, err
	}

	_, err = s.UserDao.FindUser(ctx, friendID)
	if err != nil {
		return nil, err
	}

	err = s.FriendsDao.AddFriend(ctx, uid, friendID)
	if err != nil {
		return nil, err
	}

	data := make(map[string]interface{})
	data["uid"] = uid
	data["friend_id"] = friendID
	return data, nil
}
