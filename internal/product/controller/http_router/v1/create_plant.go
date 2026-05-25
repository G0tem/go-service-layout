package v1

import (
	"net/http"

	"github.com/G0tem/go-service-layout/pkg/otel/tracer"
	"github.com/gin-gonic/gin"
)

// @Summary CreatePlant
// @Description Create CreatePlant
// @Tags plant
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param peoduct body dto.CreatePlantInput true "Plant creation payload"
// @Success 200 {object} map[string]string "status: Plant"
// @Failure 400 {object} map[string]string "invalid payload"
// @Failure 401 {object} map[string]string "unauthorized"
// @Failure 500 {object} map[string]string "internal error"
// @Router /plant/create_plant [post]
func (h *Handlers) CreatePlant(c *gin.Context) {
	_, span := tracer.Start(c.Request.Context(), "http/v1 CreatePlant")
	defer span.End()

	c.JSON(http.StatusOK, "ok")
}
