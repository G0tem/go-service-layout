package router

import (
	"net/http"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/gin-gonic/gin"
)

func New() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(prometheus.PrometheusMiddleware())

	r.GET("/live", probe)
	r.GET("/ready", probe)

	// Endpoint для Prometheus
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	return r
}

func probe(c *gin.Context) {
	c.Status(http.StatusOK)
}
