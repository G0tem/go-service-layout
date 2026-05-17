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

func (u *UseCase) UpdateApple(ctx context.Context, id uuid.UUID, input dto.UpdateAppleInput) (dto.UpdateAppleOutput, error) {
	ctx, span := tracer.Start(ctx, "usecase UpdateApple")
	defer span.End()

	var output dto.UpdateAppleOutput

	a := entity.Apple{
		ID:     id,
		Name:   input.Name,
		Status: entity.StatusUpdate,
	}

	log := log.Ctx(ctx)

	log.Info().
		Str("apple_id", a.ID.String()).
		Str("name", a.Name).
		Msg("UpdateApple started")

	err := u.postgres.UpdateApple(ctx, a)
	if err != nil {
		return output, fmt.Errorf("u.postgres.UpdateApple: %w", err)
	}

	err = u.redis.PutApple(ctx, a)
	if err != nil {
		return output, fmt.Errorf("u.redis.PutApple: %w", err)
	}

	output.ID = a.ID
	output.Name = a.Name
	output.Status = a.Status

	log.Info().
		Str("apple_id", a.ID.String()).
		Str("name", a.Name).
		Msg("UpdateApple completed")

	return output, nil
}
