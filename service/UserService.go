/*
 * @Autor: Bobby
 * @Description: User API service
 * @Date: 2022-06-06 11:02:06
 * @LastEditTime: 2022-06-08 14:59:09
 * @FilePath: \user\service\UserService.go
 */
package service

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"user/constant"
	"user/dao"
	"user/model"
	"user/pkg/util"

	"github.com/gorilla/mux"
)

func GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if _, ok := vars["uid"]; !ok {
		util.ReturnError(w, constant.ERROR_PARAM_ERR, "vars param uid not set")
		return
	}
	uid, err := strconv.Atoi(vars["uid"])
	if err != nil {
		util.ReturnError(w, constant.ERROR_PARAM_ERR, err.Error())
		return
	}
	userDao := dao.UserDao{}
	userInfo, err := userDao.FindUser(uid)
	if err != nil {
		util.ReturnError(w, -1, err.Error())
		return
	}
	util.ReturnSucc(w, userInfo)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	data := make(map[string]interface{})
	json.Unmarshal(reqBody, &data)

	var user model.User
	if tmp, ok := data["name"].(string); ok {
		user.Name = tmp
	}
	if tmp, ok := data["dob"].(string); ok {
		user.Dob = tmp
	}
	if tmp, ok := data["address"].(string); ok {
		user.Address = tmp
	}
	if tmp, ok := data["description"].(string); ok {
		user.Description = tmp
	}
	userDao := dao.UserDao{}
	err := userDao.CreateUser(&user)
	if err != nil {
		util.ReturnError(w, -1, err.Error())
		return
	}
	uid := user.Id

	//查询出来新建的用户并且返回
	info, err := userDao.FindUser(uid)
	if err != nil {
		util.ReturnError(w, -1, err.Error())
		return
	}
	util.ReturnSucc(w, info)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid, err := strconv.Atoi(vars["uid"])
	if err != nil {
		util.ReturnError(w, constant.ERROR_PARAM_ERR, err.Error())
		return
	}
	userDao := dao.UserDao{}
	err = userDao.DeleteUser(uid)
	if err != nil {
		util.ReturnError(w, constant.ERROR_MYSQL_ERR, err.Error())
		return
	}
	util.ReturnSucc(w, "删除成功")
}

func ModifyUser(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	data := make(map[string]interface{})
	json.Unmarshal(reqBody, &data)

	if _, ok := data["id"].(float64); !ok {
		util.ReturnError(w, constant.ERROR_PARAM_ERR, "user id is must param")
		return
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

	userDao := dao.UserDao{}
	err := userDao.UpdateUser(uid, modifyArr)
	if err != nil {
		util.ReturnError(w, constant.ERROR_MYSQL_ERR, err.Error())
		return
	}
	//查询出来新建的用户并且返回
	info, err := userDao.FindUser(uid)
	if err != nil {
		util.ReturnError(w, -1, err.Error())
		return
	}
	util.ReturnSucc(w, info)
}
