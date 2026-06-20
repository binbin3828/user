package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"user/constant"
)

func logWithTrace(s *Service, ctx context.Context, format string, args ...interface{}) {
	traceID := traceIDFromContext(ctx)
	spanID := spanIDFromContext(ctx)
	if traceID != "" {
		format = "[trace_id=" + traceID + "] [span_id=" + spanID + "] " + format
	}
	s.Logger.Infof(format, args...)
}

func (s *Service) LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqID := getReqID(c.Request.Context())
		ctx := c.Request.Context()
		logWithTrace(s, ctx, "[%s] request begin: Method: %v, request url: %s", reqID, c.Request.Method, c.Request.Host+c.Request.RequestURI)

		if c.Request.Body != nil {
			body, _ := c.GetRawData()
			logWithTrace(s, ctx, "[%s] request body: %s", reqID, body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		}

		startTime := time.Now()
		c.Next()
		logWithTrace(s, ctx, "[%s] exec end: api: %v, execute time: %vms", reqID, c.Request.RequestURI, time.Since(startTime).Milliseconds())
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
	ctx := c.Request.Context()
	traceID := traceIDFromContext(ctx)
	spanID := spanIDFromContext(ctx)
	logFormat := "[%s] api: %s, code: %d, msg: %s"
	if traceID != "" {
		logFormat = "[trace_id=" + traceID + "] [span_id=" + spanID + "] " + logFormat
	}
	s.Logger.Errorf(logFormat, reqID, c.Request.RequestURI, code, msg)
	c.JSON(httpStatusFromCode(code), gin.H{"code": code, "msg": msg})
}

func (s *Service) returnErrorf(c *gin.Context, code int, format string, a ...interface{}) {
	s.returnError(c, code, fmt.Sprintf(format, a...))
}

func (s *Service) returnSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": data})
}

func (s *Service) returnPaginated(c *gin.Context, data interface{}, total int64, page, pageSize int) {
	totalPages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		totalPages++
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": data,
		"pagination": gin.H{
			"total":       total,
			"page":        page,
			"page_size":   pageSize,
			"total_pages": totalPages,
		},
	})
}

// @Summary      Liveness check
// @Description  Returns OK if the service is alive
// @Tags         System
// @Produce      json
// @Success      200  {object}  util.SuccMsg
// @Router       /healthz [get]
func (s *Service) Healthz(c *gin.Context) {
	s.returnSuccess(c, "ok")
}

// @Summary      Readiness check
// @Description  Returns OK if the service can reach the database
// @Tags         System
// @Produce      json
// @Success      200  {object}  util.SuccMsg
// @Failure      503  {object}  util.ErrMsg
// @Router       /readyz [get]
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
