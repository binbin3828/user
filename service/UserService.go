/*
 * @Autor: Bobby
 * @Description: User API service
 * @Date: 2022-06-06 11:02:06
 * @LastEditTime: 2022-06-16 14:41:16
 * @FilePath: \user\service\UserService.go
 */
package service

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"user/dao"
	"user/model"
	"user/pkg/logger"

	"github.com/mmcloughlin/geohash"

	"github.com/gorilla/mux"
)

func GetUser(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	if _, ok := vars["uid"]; !ok {
		return nil, errors.New("vars param uid not set")
	}
	uid, err := strconv.Atoi(vars["uid"])
	if err != nil {
		return nil, err
	}
	userDao := dao.UserDao{}
	userInfo, err := userDao.FindUser(uid)
	if err != nil {
		return nil, err
	}
	return userInfo, nil
}

func CreateUser(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	logger.SugarLogger.Infof("request body: %s", reqBody)
	data := make(map[string]interface{})
	json.Unmarshal(reqBody, &data)

	var user model.User
	tmp, ok := data["name"].(string)
	if !ok {
		return nil, errors.New("param name not set")
	}
	user.Name = tmp

	if tmp, ok := data["dob"].(string); ok {
		user.Dob = tmp
	}
	if tmp, ok := data["address"].(string); ok {
		user.Address = tmp
	}
	if tmp, ok := data["description"].(string); ok {
		user.Description = tmp
	}

	//calc geohash string
	tmpLatitude, ok1 := data["latitude"].(float64)
	if ok1 {
		user.Latitude = tmpLatitude
	}
	tmpLongitude, ok2 := data["longitude"].(float64)
	if ok2 {
		user.Longitude = tmpLongitude
	}
	if ok1 && ok2 && tmpLatitude >= 0 && tmpLongitude >= 0 {
		//当geohash base32编码长度为8时，精度在19米左右，而当编码长度为9时，精度在2米左右，编码长度需要根据数据情况进行选择
		hash_base32 := geohash.EncodeWithPrecision(tmpLatitude, tmpLongitude, 8)
		user.LocGeohash = hash_base32
	}

	userDao := dao.UserDao{}
	err := userDao.CreateUser(&user)
	if err != nil {
		return nil, err
	}
	uid := user.Id

	//find new user and return succ
	info, err := userDao.FindUser(uid)
	if err != nil {
		return nil, err
	}
	return info, nil
}

func DeleteUser(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	uid, err := strconv.Atoi(vars["uid"])
	if err != nil {
		return nil, err
	}
	userDao := dao.UserDao{}
	err = userDao.DeleteUser(uid)
	if err != nil {
		return nil, err
	}
	data := "delete succ"
	return data, nil
}

func ModifyUser(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	data := make(map[string]interface{})
	json.Unmarshal(reqBody, &data)

	if _, ok := data["id"].(float64); !ok {
		return nil, errors.New("user id is must param")
	}

	uid := int(data["id"].(float64))

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

	//calc geohash string
	tmpLatitude, ok1 := data["latitude"].(float64)
	if ok1 {
		modifyArr["latitude"] = tmpLatitude
	}
	tmpLongitude, ok2 := data["longitude"].(float64)
	if ok2 {
		modifyArr["longitude"] = tmpLongitude
	}
	if ok1 && ok2 && tmpLatitude >= 0 && tmpLongitude >= 0 {
		//当geohash base32编码长度为8时，精度在19米左右，而当编码长度为9时，精度在2米左右，编码长度需要根据数据情况进行选择
		hash_base32 := geohash.EncodeWithPrecision(tmpLatitude, tmpLongitude, 8)
		modifyArr["loc_geohash"] = hash_base32
	}

	userDao := dao.UserDao{}
	err := userDao.UpdateUser(uid, modifyArr)
	if err != nil {
		return nil, err
	}
	//find new player and return
	info, err := userDao.FindUser(uid)
	if err != nil {
		return nil, err
	}
	return info, nil
}
