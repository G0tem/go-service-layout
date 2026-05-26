package v1

import (
	"net/http"

	"github.com/G0tem/go-service-layout/internal/product/dto"
	"github.com/G0tem/go-service-layout/pkg/logger"
	"github.com/G0tem/go-service-layout/pkg/otel/tracer"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary UpdatePlant
// @Description Update Plant
// @Tags plant
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID"
// @Param product body dto.UpdatePlantInput true "Plant update payload"
// @Success 200 {object} map[string]string "status: Plant"
// @Failure 400 {object} map[string]string "invalid payload"
// @Failure 401 {object} map[string]string "unauthorized"
// @Failure 500 {object} map[string]string "internal error"
// @Router /plant/update_plant/{id} [put]
func (h *Handlers) UpdatePlant(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "http/v1 UpdatePlant")
	defer span.End()

	log := logger.Ctx(c)

	id := c.Param("id")

	idUUID, err := uuid.Parse(id)
	if err != nil {
		log.Error().Err(err).Str("id_raw", id).Msg("uuid.Parse")
		c.JSON(http.StatusBadRequest, gin.H{"error": "uuid validate error"})
		return
	}

	input := dto.UpdatePlantInput{}

	err = c.ShouldBindJSON(&input)
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

	output, err := h.uc.UpdatePlant(ctx, idUUID, input)
	if err != nil {
		log.Error().Err(err).Msg("uc.UpdatePlant")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, output)
}
