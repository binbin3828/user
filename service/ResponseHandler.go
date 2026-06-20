package service

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"user/constant"
)

func (s *Service) LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqID := getReqID(c.Request.Context())
		s.Logger.Infof("[%s] request begin: Method: %v, request url: %s", reqID, c.Request.Method, c.Request.Host+c.Request.RequestURI)

		if c.Request.Body != nil {
			body, _ := c.GetRawData()
			s.Logger.Infof("[%s] request body: %s", reqID, body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		}

		startTime := time.Now()
		c.Next()
		s.Logger.Infof("[%s] exec end: api: %v, execute time: %vms", reqID, c.Request.RequestURI, time.Since(startTime).Milliseconds())
	}
}

func httpStatusFromCode(code int) int {
	switch code {
	case constant.ERROR_PARAM_ERR:
		return http.StatusBadRequest
	case constant.ERROR_AUTH_FAIL:
		return http.StatusUnauthorized
	case constant.ERROR_PERMISSION_DENIED:
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}

func (s *Service) returnError(c *gin.Context, code int, msg string) {
	reqID := getReqID(c.Request.Context())
	s.Logger.Errorf("[%s] api: %s, code: %d, msg: %s", reqID, c.Request.RequestURI, code, msg)
	c.JSON(httpStatusFromCode(code), gin.H{"code": code, "msg": msg})
}

func (s *Service) returnErrorf(c *gin.Context, code int, format string, a ...interface{}) {
	s.returnError(c, code, fmt.Sprintf(format, a...))
}

func (s *Service) returnSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": data})
}

func (s *Service) Healthz(c *gin.Context) {
	s.returnSuccess(c, "ok")
}

func (s *Service) Readyz(c *gin.Context) {
	_, err := s.UserDao.FindUser(c.Request.Context(), 0)
	if err == nil {
		s.returnSuccess(c, nil)
		return
	}
	c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
		"code": -1,
		"msg":  "not ready",
	})
}
