package service

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	loginAttempts = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "login_attempts_total",
			Help: "Total number of login attempts",
		},
		[]string{"status"},
	)

	userCreations = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "user_creations_total",
			Help: "Total number of user registrations",
		},
	)

	friendRequestsSent = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "friend_requests_sent_total",
			Help: "Total number of friend requests sent",
		},
	)

	friendRequestsAccepted = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "friend_requests_accepted_total",
			Help: "Total number of friend requests accepted",
		},
	)

	friendAdditions = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "friend_additions_total",
			Help: "Total number of friend additions",
		},
	)

	reqCnt = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	reqDur = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: []float64{.01, .05, .1, .5, 1, 3, 5, 10, 30},
		},
		[]string{"method", "path"},
	)

	reqInFlight = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_requests_in_flight",
			Help: "Current number of in-flight HTTP requests",
		},
	)
)

func init() {
	prometheus.MustRegister(reqCnt, reqDur, reqInFlight, loginAttempts, userCreations, friendAdditions, friendRequestsSent, friendRequestsAccepted)
	prometheus.Register(collectors.NewGoCollector())
}

func MetricsHandler() gin.HandlerFunc {
	h := promhttp.Handler()
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.FullPath()
		if path == "" {
			path = "unknown"
		}

		reqInFlight.Inc()
		start := time.Now()

		c.Next()

		status := strconv.Itoa(c.Writer.Status())
		reqCnt.WithLabelValues(c.Request.Method, path, status).Inc()
		reqDur.WithLabelValues(c.Request.Method, path).Observe(time.Since(start).Seconds())
		reqInFlight.Dec()
	}
}
