package service

import (
	"io/ioutil"
	"net/http"
	"runtime/debug"
	"time"
	"user/constant"
	"user/pkg/logger"
	"user/pkg/util"
)

type handler func(w http.ResponseWriter, req *http.Request) (data interface{}, err error)

func responseHandler(h handler) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		defer func() {
			if errRecover := recover(); errRecover != nil {
				logger.SugarLogger.Error("api: ", req.RequestURI, " errRecover: ", errRecover)
				logger.SugarLogger.Errorf("panic:%s\n", string(debug.Stack()))
			}
		}()

		logger.SugarLogger.Infof("request begin: Method: %v, request url: %s", req.Method, req.Host+req.RequestURI)
		reqBody, _ := ioutil.ReadAll(req.Body)
		logger.SugarLogger.Infof("request body: %s", reqBody)
		startTime := time.Now().UnixNano() / 1e6
		data, err := h(w, req)
		logger.SugarLogger.Infof("exec end: api: %v, execute time : %vms,", req.RequestURI, time.Now().UnixNano()/1e6-startTime)
		if err != nil {
			code := constant.ERROR_MYSQL_ERR
			if ce, ok := err.(*util.CodeError); ok {
				code = ce.Code
			}
			util.ReturnError(w, code, err.Error())
			return
		}
		util.ReturnSucc(w, data)
	}
}
