package prometheus

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	httpDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
		},
		[]string{"method", "endpoint"},
	)
)

func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		duration := time.Since(start).Seconds()
		status := fmt.Sprintf("%d", c.Writer.Status())

		httpRequestsTotal.WithLabelValues(
			c.Request.Method,
			c.FullPath(),
			status,
		).Inc()

		httpDuration.WithLabelValues(
			c.Request.Method,
			c.FullPath(),
		).Observe(duration)
	}
}
