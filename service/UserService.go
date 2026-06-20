package service

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"user/constant"
	"user/model"
	"user/pkg/util"

	"github.com/go-chi/chi/v5"
	"github.com/mmcloughlin/geohash"
)

func (s *Service) GetUser(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	uidStr := chi.URLParam(r, "uid")
	if uidStr == "" {
		return nil, util.NewCodeError(constant.ERROR_PARAM_ERR, "param uid not set")
	}
	uid, err := strconv.Atoi(uidStr)
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
	reqBody, _ := io.ReadAll(r.Body)
	s.Logger.Infof("request body: %s", reqBody)
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

	tmpLatitude, ok1 := data["latitude"].(float64)
	if ok1 {
		user.Latitude = tmpLatitude
	}
	tmpLongitude, ok2 := data["longitude"].(float64)
	if ok2 {
		user.Longitude = tmpLongitude
	}
	if ok1 && ok2 && tmpLatitude >= 0 && tmpLongitude >= 0 {
		hash_base32 := geohash.EncodeWithPrecision(tmpLatitude, tmpLongitude, 8)
		user.LocGeohash = hash_base32
	}

	err := s.UserDao.CreateUser(&user)
	if err != nil {
		return nil, err
	}

	info, err := s.UserDao.FindUser(user.Id)
	if err != nil {
		return nil, err
	}
	return info, nil
}

func (s *Service) DeleteUser(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	uidStr := chi.URLParam(r, "uid")
	if uidStr == "" {
		return nil, util.NewCodeError(constant.ERROR_PARAM_ERR, "param uid not set")
	}
	uid, err := strconv.Atoi(uidStr)
	if err != nil {
		return nil, err
	}
	err = s.UserDao.DeleteUser(uid)
	if err != nil {
		return nil, err
	}
	return "delete succ", nil
}

func (s *Service) ModifyUser(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	reqBody, _ := io.ReadAll(r.Body)
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

	tmpLatitude, ok1 := data["latitude"].(float64)
	if ok1 {
		modifyArr["latitude"] = tmpLatitude
	}
	tmpLongitude, ok2 := data["longitude"].(float64)
	if ok2 {
		modifyArr["longitude"] = tmpLongitude
	}
	if ok1 && ok2 && tmpLatitude >= 0 && tmpLongitude >= 0 {
		hash_base32 := geohash.EncodeWithPrecision(tmpLatitude, tmpLongitude, 8)
		modifyArr["loc_geohash"] = hash_base32
	}

	err := s.UserDao.UpdateUser(uid, modifyArr)
	if err != nil {
		return nil, err
	}
	info, err := s.UserDao.FindUser(uid)
	if err != nil {
		return nil, err
	}
	return info, nil
}
