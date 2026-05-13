package router

import (
	"net/http"

	"github.com/G0tem/go-service-layout/pkg/logger"
	"github.com/G0tem/go-service-layout/pkg/prometheus"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func New() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(otelgin.Middleware("go-service-layout"))
	r.Use(logger.GinTraceLogger())
	r.Use(prometheus.PrometheusMiddleware())

	r.GET("/live", probe)
	r.GET("/ready", probe)

	// Health check
	r.GET("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	// Swagger UI (доступен на /swagger/index.html)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,
		ginSwagger.URL("/swagger/doc.json"),     // URL для загрузки swagger.json
		ginSwagger.DefaultModelsExpandDepth(-1), // скрыть модели по умолчанию
	))

	// Endpoint для Prometheus
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	return r
}

func probe(c *gin.Context) {
	c.Status(http.StatusOK)
}
