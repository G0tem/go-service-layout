package v1

import (
	"errors"
	"net/http"

	"github.com/G0tem/go-service-layout/internal/product/dto"
	"github.com/G0tem/go-service-layout/internal/product/entity"
	"github.com/G0tem/go-service-layout/pkg/logger"
	"github.com/G0tem/go-service-layout/pkg/otel/tracer"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary UpdateApple
// @Description Update Apple
// @Tags apple
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Param product body dto.UpdateAppleInput true "Apple update payload"
// @Success 200 {object} map[string]string "status: Apple"
// @Failure 400 {object} map[string]string "invalid payload"
// @Failure 401 {object} map[string]string "unauthorized"
// @Failure 500 {object} map[string]string "internal error"
// @Router /apples/update_apple/{id} [put]
func (h *Handlers) UpdateApple(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "http/v1 UpdateApple")
	defer span.End()

	log := logger.Ctx(c)

	id := c.Param("id")

	idUUID, err := uuid.Parse(id)
	if err != nil {
		log.Error().Err(err).Str("id_raw", id).Msg("uuid.Parse")
		c.JSON(http.StatusBadRequest, gin.H{"error": "uuid validate error"})
		return
	}

	input := dto.UpdateAppleInput{}

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

	output, err := h.uc.UpdateApple(ctx, idUUID, input)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrNotFound):
			log.Error().Err(err).Msg("uc.UpdateApple: not found")
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
	}

	c.JSON(http.StatusOK, output)
}
