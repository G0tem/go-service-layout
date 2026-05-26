package v1

import (
	"net/http"

	"github.com/G0tem/go-service-layout/internal/product/dto"
	"github.com/G0tem/go-service-layout/pkg/logger"
	"github.com/G0tem/go-service-layout/pkg/otel/tracer"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary DeletePlant
// @Description Delete Plant
// @Tags plant
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID"
// @Success 200 {object} map[string]string "status: Plant delete"
// @Failure 400 {object} map[string]string "invalid payload"
// @Failure 401 {object} map[string]string "unauthorized"
// @Failure 500 {object} map[string]string "internal error"
// @Router /plant/delete_plant/{id} [delete]
func (h *Handlers) DeletePlant(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "http/v1 DeletePlant")
	defer span.End()

	log := logger.Ctx(c)

	var input dto.DeletePlantInput
	var err error

	id := c.Param("id")

	input.ID, err = uuid.Parse(id)
	if err != nil {
		log.Error().Err(err).Str("id_raw", id).Msg("uuid.Parse")
		c.JSON(http.StatusBadRequest, gin.H{"error": "uuid validate error"})
		return
	}

	err = h.uc.DeletePlant(ctx, input)
	if err != nil {
		log.Error().Err(err).Msg("uc.DeletePlant")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Plant delete"})
}
