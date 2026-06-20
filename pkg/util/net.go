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

func ReturnError(w http.ResponseWriter, log logger.Logger, errCode int, msg string) {
	errMsg := ErrMsg{Code: errCode, Msg: msg}
	json.NewEncoder(w).Encode(errMsg)
	bstr, _ := json.Marshal(errMsg)
	log.Error("return err:", string(bstr))
}

func ReturnSucc(w http.ResponseWriter, log logger.Logger, data interface{}) {
	succMsg := SuccMsg{Code: 0, Data: data}
	json.NewEncoder(w).Encode(succMsg)
	bstr, _ := json.Marshal(succMsg)
	log.Info("return succ:", string(bstr))
}
