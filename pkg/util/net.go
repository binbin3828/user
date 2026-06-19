/*
 * @Autor: Bobby
 * @Description: Success message, error message return definition
 * @Date: 2022-06-06 21:39:53
 * @LastEditTime: 2022-06-10 15:57:15
 * @FilePath: \user\pkg\util\net.go
 */

package util

import (
	"encoding/json"
	"fmt"
	"net/http"
	"user/pkg/logger"
)

type CodeError struct {
	Code int
	Msg  string
}

func (e *CodeError) Error() string {
	return e.Msg
}

func NewCodeError(code int, msg string) error {
	return &CodeError{Code: code, Msg: msg}
}

func NewCodeErrorf(code int, format string, a ...interface{}) error {
	return &CodeError{Code: code, Msg: fmt.Sprintf(format, a...)}
}

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
	bstr, _ := json.Marshal(errMsg)
	logger.SugarLogger.Error("return err:", string(bstr))
}

func ReturnSucc(w http.ResponseWriter, data interface{}) {
	var succMsg SuccMsg
	succMsg.Code = 0
	succMsg.Data = data
	json.NewEncoder(w).Encode(succMsg)
	bstr, _ := json.Marshal(succMsg)
	logger.SugarLogger.Info("return succ:", string(bstr))
}
