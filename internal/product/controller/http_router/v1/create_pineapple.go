package v1

import (
	"net/http"

	"github.com/G0tem/go-service-layout/pkg/otel/tracer"
	"github.com/gin-gonic/gin"
)

// @Summary CreatePineapple
// @Description Create Pineapple
// @Tags pineapple
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param peoduct body dto.CreatePineAppleInput true "Pineapple creation payload"
// @Success 200 {object} map[string]string "status: Pineapple"
// @Failure 400 {object} map[string]string "invalid payload"
// @Failure 401 {object} map[string]string "unauthorized"
// @Failure 500 {object} map[string]string "internal error"
// @Router /pineapple/create_pineapple [post]
func (h *Handlers) CreatePineapple(c *gin.Context) {
	_, span := tracer.Start(c.Request.Context(), "http/v1 CreatePineapple")
	defer span.End()

	c.JSON(http.StatusOK, "ok")
}
