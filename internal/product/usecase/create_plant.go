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

func (u *UseCase) CreatePlant(ctx context.Context, input dto.CreatePlantInput) (dto.CreatePlantOutput, error) {
	ctx, span := tracer.Start(ctx, "usecase AddPlant")
	defer span.End()

	var output dto.CreatePlantOutput

	// ctx, err := transaction.Begin(ctx)
	// if err != nil {
	// 	return output, fmt.Errorf("transaction.Begin: %w", err)
	// }

	// defer transaction.Rollback(ctx)

	a := entity.Plant{
		ID:     uuid.New(),
		Name:   entity.Name(input.Name),
		Status: entity.NewStatus("success"),
	}

	log := log.Ctx(ctx)

	log.Info().
		Str("apple_id", a.ID.String()).
		Str("name", string(a.Name)).
		Msg("CreateApple started")

	err := u.postgres.CreatePlant(ctx, a)
	if err != nil {
		return output, fmt.Errorf("u.postgres.CreatePlant: %w", err)
	}

	err = u.redis.PutPlant(ctx, a)
	if err != nil {
		return output, fmt.Errorf("u.redis.PutPlant: %w", err)
	}

	event := entity.CreateEvent{
		ID:   a.ID,
		Name: input.Name,
	}

	err = u.kafka.CreateEvent(ctx, event)
	if err != nil {
		return output, fmt.Errorf("u.kafka.CreateEvent: %w", err)
	}

	// err = transaction.Commit(ctx)
	// if err != nil {
	// 	return output, fmt.Errorf("transaction.Commit: %w", err)
	// }

	output.ID = a.ID

	log.Info().
		Str("plant_id", a.ID.String()).
		Msg("CreatePlant completed")

	return output, nil
}
