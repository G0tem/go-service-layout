package v1

import (
	"errors"
	"net/http"

	"github.com/G0tem/go-service-layout/internal/apple/dto"
	"github.com/G0tem/go-service-layout/internal/apple/entity"
	"github.com/G0tem/go-service-layout/pkg/otel/tracer"
	"github.com/G0tem/go-service-layout/pkg/render"
	"github.com/go-chi/chi/v5"
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
func (h *Handlers) GetApple(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "http/v1 GetApple")
	defer span.End()

	var (
		input dto.GetAppleInput
		err   error
	)

	id := chi.URLParam(r, "id")

	input.ID, err = uuid.Parse(id)
	if err != nil {
		log.Error().Err(err).Msg("uuid.Parse")
		http.Error(w, "uuid validate error", http.StatusBadRequest)

		return
	}

	output, err := h.uc.GetApple(ctx, input)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrNotFound):
			log.Error().Err(err).Msg("uc.CreateApple: not found")
			http.Error(w, "not found", http.StatusNotFound)

			return

		case errors.Is(err, entity.ErrUUIDInvalid), errors.Is(err, entity.ErrStatusInvalid):
			log.Error().Err(err).Msg("uc.CreateApple: validate error")
			http.Error(w, "validate error", http.StatusBadRequest)

			return

		default:
			log.Error().Err(err).Msg("uc.CreateApple: internal error")
			http.Error(w, "internal error", http.StatusInternalServerError)

			return
		}
	}

	render.JSON(w, output)
}
