package main

import (
	"app/internal/api"
	"app/internal/config"
	"app/internal/infrastructure/database"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	requestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)
	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)
)

func init() {
	prometheus.MustRegister(requestCount)
	prometheus.MustRegister(requestDuration)
}

func main() {
	cfg := config.NewConfig()
	logger := config.NewSugaredLogger()
	db := database.New(cfg.Database, logger)

	r := gin.Default()
	r.Use(monitoringMiddleware)
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	api.Run(r, db, logger)
}

func monitoringMiddleware(c *gin.Context) {
	path := c.FullPath()
	method := c.Request.Method

	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		requestDuration.WithLabelValues(method, path).Observe(v)
	}))
	defer timer.ObserveDuration()

	c.Next()

	status := http.StatusText(c.Writer.Status())
	requestCount.WithLabelValues(method, path, status).Inc()
}
