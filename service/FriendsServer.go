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
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"user/dao"

	"github.com/gorilla/mux"
)

func GetNearbyFriend(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	if _, ok := vars["uid"]; !ok {
		return nil, errors.New("vars param uid not set")
	}
	uid, err := strconv.Atoi(vars["uid"])
	if err != nil {
		return nil, err
	}

	//check uid exist
	userDao := dao.UserDao{}
	info, err := userDao.FindUser(uid)
	if err != nil {
		return nil, err
	}
	geohashStr := info.LocGeohash

	//附近，参数n代表Geohash，精确的位数，也就是大概距离；n=6时候，大概为附近1千米
	likeSubStr := geohashStr[0:6]

	friendsDao := dao.FriendsDao{}
	list, err := friendsDao.GetNearbyFriend(uid, likeSubStr)
	if err != nil {
		return nil, err
	}

	data := make(map[string]interface{})
	data["uid"] = uid
	data["list"] = list
	return data, nil
}

func GetFriendsList(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	if _, ok := vars["uid"]; !ok {
		return nil, errors.New("vars param uid not set")
	}
	uid, err := strconv.Atoi(vars["uid"])
	if err != nil {
		return nil, err
	}

	//check uid exist
	userDao := dao.UserDao{}
	_, err = userDao.FindUser(uid)
	if err != nil {
		return nil, err
	}
	friendsDao := dao.FriendsDao{}
	list, err := friendsDao.GetFriendsList(uid)
	if err != nil {
		return nil, err
	}
	data := make(map[string]interface{})
	data["uid"] = uid
	data["list"] = list
	return data, nil
}

func AddFriend(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	data := make(map[string]interface{})
	json.Unmarshal(reqBody, &data)

	uidft, ok := data["uid"].(float64)
	if !ok {
		return nil, errors.New("param uid not set")
	}
	uid := int(uidft)

	//check uid exist
	userDao := dao.UserDao{}
	_, err := userDao.FindUser(uid)
	if err != nil {
		return nil, err
	}

	friFt, ok := data["fri"].(float64)
	if !ok {
		return nil, errors.New("param fri not set")
	}
	fri := int(friFt)

	//check fri exist
	_, err = userDao.FindUser(fri)
	if err != nil {
		return nil, err
	}

	friendsDao := dao.FriendsDao{}
	err = friendsDao.AddFriend(uid, fri)
	if err != nil {
		return nil, err
	}

	data = make(map[string]interface{})
	data["uid"] = uid
	data["fri"] = fri
	return data, nil
}
