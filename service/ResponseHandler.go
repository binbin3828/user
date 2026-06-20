package service

import (
	"io/ioutil"
	"net/http"
	"runtime/debug"
	"time"
	"user/constant"
	"user/pkg/util"
)

type handler func(w http.ResponseWriter, req *http.Request) (data interface{}, err error)

func (s *Service) responseHandler(h handler) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		defer func() {
			if errRecover := recover(); errRecover != nil {
				s.Logger.Error("api: ", req.RequestURI, " errRecover: ", errRecover)
				s.Logger.Errorf("panic:%s\n", string(debug.Stack()))
			}
		}()

		s.Logger.Infof("request begin: Method: %v, request url: %s", req.Method, req.Host+req.RequestURI)
		reqBody, _ := ioutil.ReadAll(req.Body)
		s.Logger.Infof("request body: %s", reqBody)
		startTime := time.Now().UnixNano() / 1e6
		data, err := h(w, req)
		s.Logger.Infof("exec end: api: %v, execute time : %vms,", req.RequestURI, time.Now().UnixNano()/1e6-startTime)
		if err != nil {
			code := constant.ERROR_MYSQL_ERR
			if ce, ok := err.(*util.CodeError); ok {
				code = ce.Code
			}
			util.ReturnError(w, s.Logger, code, err.Error())
			return
		}
		util.ReturnSucc(w, s.Logger, data)
	}
}
