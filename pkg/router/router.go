package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func New() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	r.GET("/live", probe)
	r.GET("/ready", probe)

	return r
}

func probe(c *gin.Context) {
	c.Status(http.StatusOK)
}
