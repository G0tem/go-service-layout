package v1

import (
	"net/http"

	"github.com/G0tem/go-service-layout/internal/product/dto"
	"github.com/G0tem/go-service-layout/pkg/logger"
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
	ctx, span := tracer.Start(c.Request.Context(), "http/v1 CreatePlant")
	defer span.End()

	log := logger.Ctx(c)

	input := dto.CreatePlantInput{}

	err := c.ShouldBindJSON(&input)
	if err != nil {
		log.Error().Err(err).Msg("c.ShouldBindJSON")
		c.JSON(http.StatusBadRequest, gin.H{"error": "json error"})
		return
	}

	err = input.Validate()
	if err != nil {
		log.Error().Err(err).Msg("input.Validate")
		c.JSON(http.StatusBadRequest, gin.H{"error": "validate error"})
		return
	}

	output, err := h.uc.CreatePlant(ctx, input)
	if err != nil {
		log.Error().Err(err).Msg("uc.CreatePlant: internal error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusCreated, output)
}
