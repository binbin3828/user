/*
 * @Autor: Bobby
 * @Description: Friends
 * @Date: 2022-06-09 15:08:07
 * @LastEditTime: 2022-06-16 14:26:38
 * @FilePath: \user\service\FriendsServer.go
 */

package service

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"user/constant"
	"user/pkg/util"

	"github.com/gorilla/mux"
)

func (s *Service) GetNearbyFriend(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	if _, ok := vars["uid"]; !ok {
		return nil, util.NewCodeError(constant.ERROR_PARAM_ERR, "vars param uid not set")
	}
	uid, err := strconv.Atoi(vars["uid"])
	if err != nil {
		return nil, err
	}

	//check uid exist
	info, err := s.UserDao.FindUser(uid)
	if err != nil {
		return nil, err
	}
	geohashStr := info.LocGeohash

	//附近，参数n代表Geohash，精确的位数，也就是大概距离；n=6时候，大概为附近1千米
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
	vars := mux.Vars(r)
	if _, ok := vars["uid"]; !ok {
		return nil, util.NewCodeError(constant.ERROR_PARAM_ERR, "vars param uid not set")
	}
	uid, err := strconv.Atoi(vars["uid"])
	if err != nil {
		return nil, err
	}

	//check uid exist
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
	reqBody, _ := ioutil.ReadAll(r.Body)
	data := make(map[string]interface{})
	json.Unmarshal(reqBody, &data)

	uidft, ok := data["uid"].(float64)
	if !ok {
		return nil, util.NewCodeError(constant.ERROR_PARAM_ERR, "param uid not set")
	}
	uid := int(uidft)

	//check uid exist
	_, err := s.UserDao.FindUser(uid)
	if err != nil {
		return nil, err
	}

	friFt, ok := data["fri"].(float64)
	if !ok {
		return nil, util.NewCodeError(constant.ERROR_PARAM_ERR, "param fri not set")
	}
	fri := int(friFt)

	//check fri exist
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
