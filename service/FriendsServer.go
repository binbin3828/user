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
	uidStr := chi.URLParam(r, "uid")
	if uidStr == "" {
		return nil, util.NewCodeError(constant.ERROR_PARAM_ERR, "param uid not set")
	}
	uid, err := strconv.Atoi(uidStr)
	if err != nil {
		return nil, err
	}

	info, err := s.UserDao.FindUser(uid)
	if err != nil {
		return nil, err
	}
	geohashStr := info.LocGeohash

	// n=6 时约 1km 精度
	likeSubStr := geohashStr[0:6]

	list, err := s.FriendsDao.GetNearbyFriend(uid, likeSubStr)
	if err != nil {
		return nil, err
	}

	data := make(map[string]interface{})
	data["uid"] = uid
	data["list"] = list
	return data, nil
}

func (s *Service) GetFriendsList(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	uidStr := chi.URLParam(r, "uid")
	if uidStr == "" {
		return nil, util.NewCodeError(constant.ERROR_PARAM_ERR, "param uid not set")
	}
	uid, err := strconv.Atoi(uidStr)
	if err != nil {
		return nil, err
	}

	_, err = s.UserDao.FindUser(uid)
	if err != nil {
		return nil, err
	}
	list, err := s.FriendsDao.GetFriendsList(uid)
	if err != nil {
		return nil, err
	}
	data := make(map[string]interface{})
	data["uid"] = uid
	data["list"] = list
	return data, nil
}

func (s *Service) AddFriend(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	reqBody, _ := io.ReadAll(r.Body)
	data := make(map[string]interface{})
	json.Unmarshal(reqBody, &data)

	uidft, ok := data["uid"].(float64)
	if !ok {
		return nil, util.NewCodeError(constant.ERROR_PARAM_ERR, "param uid not set")
	}
	uid := int(uidft)

	_, err := s.UserDao.FindUser(uid)
	if err != nil {
		return nil, err
	}

	friFt, ok := data["fri"].(float64)
	if !ok {
		return nil, util.NewCodeError(constant.ERROR_PARAM_ERR, "param fri not set")
	}
	fri := int(friFt)

	_, err = s.UserDao.FindUser(fri)
	if err != nil {
		return nil, err
	}

	err = s.FriendsDao.AddFriend(uid, fri)
	if err != nil {
		return nil, err
	}

	data = make(map[string]interface{})
	data["uid"] = uid
	data["fri"] = fri
	return data, nil
}
