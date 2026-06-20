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
	"io/ioutil"
	"net/http"
	"strconv"
	"user/constant"
	"user/model"
	"user/pkg/logger"
	"user/pkg/util"

	"github.com/mmcloughlin/geohash"

	"github.com/gorilla/mux"
)

func (s *Service) GetUser(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	if _, ok := vars["uid"]; !ok {
		return nil, util.NewCodeError(constant.ERROR_PARAM_ERR, "vars param uid not set")
	}
	uid, err := strconv.Atoi(vars["uid"])
	if err != nil {
		return nil, err
	}
	userInfo, err := s.UserDao.FindUser(uid)
	if err != nil {
		return nil, err
	}
	return userInfo, nil
}

func (s *Service) CreateUser(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	logger.SugarLogger.Infof("request body: %s", reqBody)
	data := make(map[string]interface{})
	json.Unmarshal(reqBody, &data)

	var user model.User
	tmp, ok := data["name"].(string)
	if !ok {
		return nil, util.NewCodeError(constant.ERROR_PARAM_ERR, "param name not set")
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

	err := s.UserDao.CreateUser(&user)
	if err != nil {
		return nil, err
	}
	uid := user.Id

	//find new user and return succ
	info, err := s.UserDao.FindUser(uid)
	if err != nil {
		return nil, err
	}
	return info, nil
}

func (s *Service) DeleteUser(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	uid, err := strconv.Atoi(vars["uid"])
	if err != nil {
		return nil, err
	}
	err = s.UserDao.DeleteUser(uid)
	if err != nil {
		return nil, err
	}
	data := "delete succ"
	return data, nil
}

func (s *Service) ModifyUser(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	data := make(map[string]interface{})
	json.Unmarshal(reqBody, &data)

	if _, ok := data["id"].(float64); !ok {
		return nil, util.NewCodeError(constant.ERROR_PARAM_ERR, "user id is must param")
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

	err := s.UserDao.UpdateUser(uid, modifyArr)
	if err != nil {
		return nil, err
	}
	//find new player and return
	info, err := s.UserDao.FindUser(uid)
	if err != nil {
		return nil, err
	}
	return info, nil
}
