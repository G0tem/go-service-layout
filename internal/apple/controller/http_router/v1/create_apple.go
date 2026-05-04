package v1

import (
	"errors"
	"net/http"

	"github.com/G0tem/go-service-layout/internal/apple/dto"
	"github.com/G0tem/go-service-layout/internal/apple/entity"
	"github.com/G0tem/go-service-layout/pkg/otel/tracer"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/attribute"
)

// @Summary CreateApple
// @Description Creates Apple
// @Tags apple
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param peoduct body dto.CreateAppleInput true "Apple creation payload"
// @Success 201 {object} map[string]string "status: Apple created"
// @Failure 400 {object} map[string]string "invalid payload"
// @Failure 401 {object} map[string]string "unauthorized"
// @Failure 500 {object} map[string]string "internal error"
// @Router /apples/create_apple [post]
func (h *Handlers) CreateApple(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "http/v1 CreateApple")
	defer span.End()

	input := dto.CreateAppleInput{}

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

	output, err := h.uc.CreateApple(ctx, input)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrNotFound):
			log.Error().Err(err).Msg("uc.CreateApple: not found")
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return

		case errors.Is(err, entity.ErrUUIDInvalid), errors.Is(err, entity.ErrStatusInvalid):
			log.Error().Err(err).Msg("uc.CreateApple: validate error")
			c.JSON(http.StatusBadRequest, gin.H{"error": "validate error"})
			return

		default:
			log.Error().Err(err).Msg("uc.CreateApple: internal error")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			span.SetAttributes(attribute.KeyValue{Key: "error", Value: attribute.StringValue(err.Error())})
			tracer.SetStatus(span, err)
			return
		}
	}

	c.JSON(http.StatusCreated, output)
}
