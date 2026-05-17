package v1

import (
	"net/http"

	"github.com/G0tem/go-service-layout/pkg/otel/tracer"
	"github.com/gin-gonic/gin"
)

// @Summary GetPineapple
// @Description Get Pineapple
// @Tags pineapple
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID"
// @Success 200 {object} map[string]string "status: Pineapple"
// @Failure 400 {object} map[string]string "invalid payload"
// @Failure 401 {object} map[string]string "unauthorized"
// @Failure 500 {object} map[string]string "internal error"
// @Router /pineapple/get_pineapple/{id} [get]
func (h *Handlers) GetPineapple(c *gin.Context) {
	_, span := tracer.Start(c.Request.Context(), "http/v1 GetPineapple")
	defer span.End()

	c.JSON(http.StatusOK, "ok")
}
