package v1

import (
	"net/http"

	"github.com/G0tem/go-service-layout/internal/product/dto"
	"github.com/G0tem/go-service-layout/pkg/logger"
	"github.com/G0tem/go-service-layout/pkg/otel/tracer"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary DeleteApple
// @Description Delete Apple
// @Tags apple
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Success 200 {object} map[string]string "status: Apple delete"
// @Failure 400 {object} map[string]string "invalid payload"
// @Failure 401 {object} map[string]string "unauthorized"
// @Failure 500 {object} map[string]string "internal error"
// @Router /apples/delete_apple/{id} [delete]
func (h *Handlers) DeleteApple(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "http/v1 DeleteApple")
	defer span.End()

	log := logger.Ctx(c)

	var (
		input dto.DeleteAppleInput
		err   error
	)

	id := c.Param("id")

	input.ID, err = uuid.Parse(id)
	if err != nil {
		log.Error().Err(err).Str("id_raw", id).Msg("uuid.Parse")
		c.JSON(http.StatusBadRequest, gin.H{"error": "uuid validate error"})
		return
	}

	err = h.uc.DeleteApple(ctx, input)
	if err != nil {
		log.Error().Err(err).Msg("uc.DeleteApple")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Apple delete"})
}
