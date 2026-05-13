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

func (u *UseCase) CreateApple(ctx context.Context, input dto.CreateAppleInput) (dto.CreateAppleOutput, error) {
	ctx, span := tracer.Start(ctx, "usecase CreateApple")
	defer span.End()

	var output dto.CreateAppleOutput

	a := entity.Apple{
		ID:     uuid.New(),
		Name:   input.Name,
		Status: entity.StatusNew,
	}

	log := log.Ctx(ctx)

	log.Info().
		Str("apple_id", a.ID.String()).
		Str("name", a.Name).
		Msg("CreateApple started")

	err := u.postgres.CreateApple(ctx, a)
	if err != nil {
		return output, fmt.Errorf("u.postgres.CreateApple: %w", err)
	}

	err = u.redis.PutApple(ctx, a)
	if err != nil {
		return output, fmt.Errorf("u.redis.PutApple: %w", err)
	}

	event := entity.CreateEvent{
		ID:   a.ID,
		Name: input.Name,
	}

	err = u.kafka.CreateEvent(ctx, event)
	if err != nil {
		return output, fmt.Errorf("u.kafka.CreateEvent: %w", err)
	}

	output.ID = a.ID

	log.Info().
		Str("apple_id", a.ID.String()).
		Msg("CreateApple completed")

	return output, nil
}
