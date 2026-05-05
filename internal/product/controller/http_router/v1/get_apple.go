package v1

import (
	"errors"
	"net/http"

	"github.com/G0tem/go-service-layout/internal/product/dto"
	"github.com/G0tem/go-service-layout/internal/product/entity"
	"github.com/G0tem/go-service-layout/pkg/otel/tracer"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// @Summary GetApple
// @Description Get Apple
// @Tags apple
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID"
// @Success 200 {object} map[string]string "status: Apple"
// @Failure 400 {object} map[string]string "invalid payload"
// @Failure 401 {object} map[string]string "unauthorized"
// @Failure 500 {object} map[string]string "internal error"
// @Router /apples/get_apple/{id} [get]
func (h *Handlers) GetApple(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "http/v1 GetApple")
	defer span.End()

	var (
		input dto.GetAppleInput
		err   error
	)

	id := c.Param("id")

	input.ID, err = uuid.Parse(id)
	if err != nil {
		log.Error().Err(err).Msg("uuid.Parse")
		c.JSON(http.StatusBadRequest, gin.H{"error": "uuid validate error"})
		return
	}

	output, err := h.uc.GetApple(ctx, input)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrNotFound):
			log.Error().Err(err).Msg("uc.GetApple: not found")
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return

		case errors.Is(err, entity.ErrUUIDInvalid), errors.Is(err, entity.ErrStatusInvalid):
			log.Error().Err(err).Msg("uc.GetApple: validate error")
			c.JSON(http.StatusBadRequest, gin.H{"error": "validate error"})
			return

		default:
			log.Error().Err(err).Msg("uc.GetApple: internal error")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}
	}

	c.JSON(http.StatusOK, output)
}
