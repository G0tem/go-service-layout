package usecase

import (
	"context"
	"fmt"

	"github.com/G0tem/go-service-layout/internal/product/dto"
	"github.com/G0tem/go-service-layout/pkg/otel/tracer"
	"github.com/rs/zerolog/log"
)

func (u *UseCase) DeletePlant(ctx context.Context, input dto.DeletePlantInput) error {
	ctx, span := tracer.Start(ctx, "usecase DeletePlant")
	defer span.End()

	log := log.Ctx(ctx)

	log.Info().
		Str("plant_id", input.ID.String()).
		Msg("DeletePlant started")

	err := u.postgres.DeletePlant(ctx, input.ID)
	if err != nil {
		return fmt.Errorf("u.postgres.DeletePlant: %w", err)
	}

	log.Info().
		Str("plant_id", input.ID.String()).
		Msg("DeletePlant completed")

	return nil
}
