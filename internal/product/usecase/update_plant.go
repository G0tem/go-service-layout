package usecase

import (
	"context"
	"fmt"

	"github.com/G0tem/go-service-layout/internal/product/dto"
	"github.com/G0tem/go-service-layout/internal/product/entity"
	"github.com/G0tem/go-service-layout/pkg/otel/tracer"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func (u *UseCase) UpdatePlant(ctx context.Context, id uuid.UUID, input dto.UpdatePlantInput) (dto.UpdatePlantOutput, error) {
        ctx, span := tracer.Start(ctx, "usecase UpdatePlant")
        defer span.End()

        var output dto.UpdatePlantOutput

        pl := entity.Plant{
                ID:     id,
                Name:   entity.Name(input.Name),
                Status: entity.NewStatus("update"),
        }

        log := log.Ctx(ctx)

        log.Info().
                Str("plant_id", pl.ID.String()).
                Str("name", string(pl.Name)).
                Msg("UpdatePlant started")

        err := u.postgres.UpdatePlant(ctx, pl)
        if err != nil {
                return output, fmt.Errorf("u.postgres.UpdatePlant: %w", err)
        }

        err = u.redis.PutPlant(ctx, pl)
        if err != nil {
                return output, fmt.Errorf("u.redis.PutPlant: %w", err)
        }

        output.ID = pl.ID
        output.Name = string(pl.Name)
        output.Status = int(pl.Status)

        log.Info().
                Str("plant_id", pl.ID.String()).
                Str("name", string(pl.Name)).
                Msg("UpdatePlant completed")

        return output, nil
}