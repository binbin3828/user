/*
 * @Autor: Bobby
 * @Description: Success message, error message return definition
 * @Date: 2022-06-06 21:39:53
 * @LastEditTime: 2022-06-07 22:01:25
 * @FilePath: \User\pkg\util\net.go
 */

package util

import (
	"encoding/json"
	"net/http"
)

type ErrMsg struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type SuccMsg struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

func ReturnError(w http.ResponseWriter, errCode int, msg string) {
	var errMsg ErrMsg
	errMsg.Code = errCode
	errMsg.Msg = msg
	json.NewEncoder(w).Encode(errMsg)
}

func ReturnSucc(w http.ResponseWriter, data interface{}) {
	var succMsg SuccMsg
	succMsg.Code = 0
	succMsg.Data = data
	json.NewEncoder(w).Encode(succMsg)
}
