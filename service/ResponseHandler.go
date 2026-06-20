package service

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (s *Service) LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		s.Logger.Infof("request begin: Method: %v, request url: %s", c.Request.Method, c.Request.Host+c.Request.RequestURI)

		if c.Request.Body != nil {
			body, _ := c.GetRawData()
			s.Logger.Infof("request body: %s", body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		}

		startTime := time.Now()
		c.Next()
		s.Logger.Infof("exec end: api: %v, execute time: %vms", c.Request.RequestURI, time.Since(startTime).Milliseconds())
	}
}

func (s *Service) returnError(c *gin.Context, code int, msg string) {
	s.Logger.Errorf("api: %s, code: %d, msg: %s", c.Request.RequestURI, code, msg)
	c.JSON(http.StatusOK, gin.H{"code": code, "msg": msg})
}

func (s *Service) returnErrorf(c *gin.Context, code int, format string, a ...interface{}) {
	s.returnError(c, code, fmt.Sprintf(format, a...))
}

func (s *Service) returnSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": data})
}
